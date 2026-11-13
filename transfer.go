package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const (
	chunkSize  = 4 << 20  // bytes per kernel copy call: progress / cancel granularity
	syncWindow = 32 << 20 // write-back window for large files
	fsyncMin   = 16 << 20 // files at least this big are fdatasync'ed after copy
)

// Copy strategies, fastest first. The runner degrades when the kernel or
// filesystem pair does not support one.
const (
	modeCopyRange = iota // copy_file_range: in-kernel copy, reflink on btrfs/xfs
	modeSendfile         // sendfile: in-kernel copy through the page cache
	modeReadWrite        // plain read/write with a 1 MiB buffer
)

type runner struct {
	app *App
	job *Job
	ctx context.Context

	noRename bool // rename() failed with EXDEV: copy + delete instead
	mode     int  // current copy strategy
	tmpN     int
	wrote    bool // some data was actually copied (worth a final syncfs)
}

func (r *runner) addErr(realPath string, err error) {
	r.job.addErr(r.app.virt(realPath) + ": " + errMsg(err))
}

type item struct {
	src, dst string
	info     fs.FileInfo
	check    bool // destination exists and must be merged / overwritten
}

func (r *runner) runTransfer() error {
	j := r.job
	move := j.Op == "move"
	destDir, err := r.app.resolve(j.Dest)
	if err != nil {
		return err
	}
	if st, err := os.Stat(destDir); err != nil {
		return fmt.Errorf("destination: %s", errMsg(err))
	} else if !st.IsDir() {
		return errors.New("destination is not a directory")
	}
	sep := string(filepath.Separator)

	// Plan: resolve name conflicts and try an instant rename for moves.
	var items []item
	for _, vs := range j.Sources {
		if err := r.ctx.Err(); err != nil {
			return err
		}
		src, err := r.app.resolve(vs)
		if err != nil || src == r.app.root {
			j.addErr(vs + ": invalid source")
			continue
		}
		info, err := os.Lstat(src)
		if err != nil {
			r.addErr(src, err)
			continue
		}
		if info.IsDir() && (destDir == src || strings.HasPrefix(destDir+sep, src+sep)) {
			r.addErr(src, errors.New("cannot copy or move a folder into itself"))
			continue
		}
		dst := filepath.Join(destDir, filepath.Base(src))
		check := false
		if dst == src {
			if move {
				continue // moving onto itself: nothing to do
			}
			dst = uniqueName(dst, info.IsDir())
		} else if dinfo, err := os.Lstat(dst); err == nil {
			switch j.Conflict {
			case "skip":
				if !info.IsDir() || !dinfo.IsDir() {
					continue
				}
				check = true // folders are merged; existing files inside are skipped
			case "rename":
				dst = uniqueName(dst, info.IsDir())
			default:
				check = true
			}
		}
		mergeDir := check && info.IsDir()
		if move && !mergeDir {
			err := os.Rename(src, dst)
			if err == nil {
				continue // same filesystem: done instantly
			}
			if !errors.Is(err, syscall.EXDEV) {
				r.addErr(src, err)
				continue
			}
		}
		items = append(items, item{src, dst, info, check})
	}
	if len(items) == 0 {
		return nil
	}

	j.setPhase("scan")
	for _, it := range items {
		if err := r.scan(it.src, it.info, false); err != nil {
			return err
		}
	}

	j.setPhase(j.Op)
	for _, it := range items {
		r.mode = modeCopyRange
		r.noRename = move && !(it.check && it.info.IsDir())
		if err := r.transfer(it.src, it.dst, it.info, it.check); err != nil {
			return err
		}
	}

	if r.wrote && r.app.syncOn {
		// Make sure everything is on the disk before reporting success.
		j.setPhase("sync")
		if d, err := os.Open(destDir); err == nil {
			_ = unix.Syncfs(int(d.Fd()))
			d.Close()
		}
	}
	return nil
}

// scan adds the size of a tree to the job totals.
func (r *runner) scan(p string, info fs.FileInfo, countDirs bool) error {
	if err := r.ctx.Err(); err != nil {
		return err
	}
	j := r.job
	if !info.IsDir() {
		j.totalFiles.Add(1)
		if info.Mode().IsRegular() {
			j.totalBytes.Add(info.Size())
		}
		return nil
	}
	if countDirs {
		j.totalFiles.Add(1)
	}
	entries, _ := os.ReadDir(p) // errors are reported during the real pass
	for _, e := range entries {
		fi, err := e.Info()
		if err != nil {
			continue
		}
		if err := r.scan(filepath.Join(p, e.Name()), fi, countDirs); err != nil {
			return err
		}
	}
	return nil
}

// transfer copies (or moves) one tree. It returns an error only on cancel;
// per-file failures are recorded in the job and the walk continues.
func (r *runner) transfer(src, dst string, info fs.FileInfo, check bool) error {
	if err := r.ctx.Err(); err != nil {
		return err
	}
	j := r.job
	move := j.Op == "move"
	mode := info.Mode()

	if mode.IsDir() {
		existed := false
		if err := os.Mkdir(dst, 0o755); err != nil {
			if st, serr := os.Stat(dst); serr != nil || !st.IsDir() {
				r.addErr(src, err)
				return nil
			}
			existed = true
		}
		entries, err := os.ReadDir(src)
		if err != nil {
			r.addErr(src, err)
		}
		for _, e := range entries {
			child := filepath.Join(src, e.Name())
			fi, err := e.Info()
			if err != nil {
				r.addErr(child, err)
				continue
			}
			if err := r.transfer(child, filepath.Join(dst, e.Name()), fi, existed); err != nil {
				return err
			}
		}
		if !existed {
			copyMeta(dst, info)
		}
		if move {
			// The directory stays if something inside was skipped or failed.
			if err := os.Remove(src); err != nil && !errors.Is(err, syscall.ENOTEMPTY) && !errors.Is(err, syscall.EEXIST) {
				r.addErr(src, err)
			}
		}
		return nil
	}

	cur := r.app.virt(src)
	j.current.Store(&cur)
	var size int64
	if mode.IsRegular() {
		size = info.Size()
	}
	account := func(copied int64) {
		if size > copied {
			j.doneBytes.Add(size - copied)
		}
		j.doneFiles.Add(1)
	}

	if check {
		if _, err := os.Lstat(dst); err == nil {
			switch j.Conflict {
			case "skip":
				account(0)
				return nil
			case "rename":
				dst = uniqueName(dst, false)
			}
		}
	}

	if move && !r.noRename {
		err := os.Rename(src, dst)
		if err == nil {
			account(0)
			return nil
		}
		if !errors.Is(err, syscall.EXDEV) {
			r.addErr(src, err)
			account(0)
			return nil
		}
		r.noRename = true
	}

	switch {
	case mode.IsRegular():
		n, err := r.copyFile(src, dst, info)
		if err != nil {
			if cerr := r.ctx.Err(); cerr != nil {
				return cerr
			}
			r.addErr(src, err)
			account(n)
			return nil
		}
		account(size) // bytes were already counted while copying
	case mode&fs.ModeSymlink != 0:
		target, err := os.Readlink(src)
		if err == nil {
			if check {
				_ = os.Remove(dst)
			}
			err = os.Symlink(target, dst)
		}
		if err != nil {
			r.addErr(src, err)
			account(0)
			return nil
		}
		if st, ok := info.Sys().(*syscall.Stat_t); ok {
			_ = os.Lchown(dst, int(st.Uid), int(st.Gid))
		}
		account(0)
	default:
		r.addErr(src, errors.New("special file skipped"))
		account(0)
		return nil
	}

	if move {
		if err := os.Remove(src); err != nil {
			r.addErr(src, err)
		}
	}
	return nil
}

// copyFile copies one regular file through a temporary name, so an
// interrupted copy never leaves a truncated file under the final name.
func (r *runner) copyFile(src, dst string, info fs.FileInfo) (int64, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	r.tmpN++
	tmp := filepath.Join(filepath.Dir(dst), fmt.Sprintf(".lfe-%d-%d.part", r.job.ID, r.tmpN))
	out, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return 0, err
	}
	ok := false
	defer func() {
		if !ok {
			out.Close()
			os.Remove(tmp)
		}
	}()

	n, err := r.pump(in, out, info.Size())
	if err != nil {
		return n, err
	}
	r.wrote = true
	if r.app.syncOn && n >= fsyncMin {
		if err := unix.Fdatasync(int(out.Fd())); err != nil && err != unix.EINVAL && err != unix.ENOSYS {
			return n, err
		}
	}
	if err := out.Close(); err != nil {
		return n, err
	}
	copyMeta(tmp, info)
	if err := os.Rename(tmp, dst); err != nil {
		return n, err
	}
	ok = true
	return n, nil
}

// pump moves the data, preferring in-kernel zero-copy paths. It checks for
// cancellation and reports progress after every chunk.
func (r *runner) pump(in, out *os.File, size int64) (int64, error) {
	infd, outfd := int(in.Fd()), int(out.Fd())
	_ = unix.Fadvise(infd, 0, 0, unix.FADV_SEQUENTIAL)

	// For large files, push dirty pages to the disk as we go. This keeps
	// memory use flat on small devices and makes speed and progress real
	// instead of "fast into the page cache, then a long stall".
	throttle := r.app.syncOn && size >= 2*syncWindow
	var off, winStart, prevStart int64
	var buf []byte

	for {
		if err := r.ctx.Err(); err != nil {
			return off, err
		}
		var n int
		var err error
		switch r.mode {
		case modeCopyRange:
			n, err = unix.CopyFileRange(infd, nil, outfd, nil, chunkSize, 0)
			if err == unix.EINTR {
				continue
			}
			if err != nil || (n == 0 && off < size) {
				r.mode = modeSendfile // not supported here; real errors resurface below
				continue
			}
		case modeSendfile:
			n, err = unix.Sendfile(outfd, infd, nil, chunkSize)
			if err == unix.EINTR || err == unix.EAGAIN {
				continue
			}
			if err == unix.EINVAL || err == unix.ENOSYS || err == unix.EOPNOTSUPP || (err == nil && n == 0 && off < size) {
				r.mode = modeReadWrite
				continue
			}
			if err != nil {
				return off, err
			}
		default:
			if buf == nil {
				buf = make([]byte, 1<<20)
			}
			n, err = in.Read(buf)
			if n > 0 {
				if _, werr := out.Write(buf[:n]); werr != nil {
					return off, werr
				}
			}
			if err != nil && err != io.EOF {
				return off, err
			}
		}
		if n == 0 {
			return off, nil
		}
		off += int64(n)
		r.job.doneBytes.Add(int64(n))

		if throttle && off-winStart >= syncWindow {
			_ = unix.SyncFileRange(outfd, winStart, off-winStart, unix.SYNC_FILE_RANGE_WRITE)
			if winStart > prevStart {
				_ = unix.SyncFileRange(outfd, prevStart, winStart-prevStart,
					unix.SYNC_FILE_RANGE_WAIT_BEFORE|unix.SYNC_FILE_RANGE_WRITE|unix.SYNC_FILE_RANGE_WAIT_AFTER)
			}
			prevStart, winStart = winStart, off
		}
	}
}

// copyMeta copies owner, permissions and modification time (best effort:
// FAT/exFAT/NTFS do not support all of them).
func copyMeta(dst string, info fs.FileInfo) {
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		_ = os.Chown(dst, int(st.Uid), int(st.Gid))
	}
	_ = os.Chmod(dst, info.Mode()&(fs.ModePerm|fs.ModeSetuid|fs.ModeSetgid|fs.ModeSticky))
	_ = os.Chtimes(dst, time.Time{}, info.ModTime())
}

// uniqueName returns "name (2).ext", "name (3).ext"... that does not exist yet.
func uniqueName(p string, isDir bool) string {
	dir, base := filepath.Split(p)
	ext := ""
	if !isDir {
		if e := filepath.Ext(base); e != base {
			ext = e
		}
	}
	stem := strings.TrimSuffix(base, ext)
	for i := 2; ; i++ {
		c := filepath.Join(dir, fmt.Sprintf("%s (%d)%s", stem, i, ext))
		if _, err := os.Lstat(c); err != nil {
			return c
		}
	}
}

func (r *runner) runDelete() error {
	j := r.job
	type target struct {
		path string
		info fs.FileInfo
	}
	var targets []target
	j.setPhase("scan")
	for _, vs := range j.Sources {
		p, err := r.app.resolve(vs)
		if err != nil || p == r.app.root {
			j.addErr(vs + ": invalid path")
			continue
		}
		info, err := os.Lstat(p)
		if err != nil {
			r.addErr(p, err)
			continue
		}
		if err := r.scan(p, info, true); err != nil {
			return err
		}
		targets = append(targets, target{p, info})
	}
	j.setPhase("delete")
	for _, t := range targets {
		if err := r.remove(t.path, t.info); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) remove(p string, info fs.FileInfo) error {
	if err := r.ctx.Err(); err != nil {
		return err
	}
	j := r.job
	if info.IsDir() {
		entries, err := os.ReadDir(p)
		if err != nil {
			r.addErr(p, err)
		}
		for _, e := range entries {
			fi, err := e.Info()
			if err != nil {
				continue
			}
			if err := r.remove(filepath.Join(p, e.Name()), fi); err != nil {
				return err
			}
		}
	} else {
		cur := r.app.virt(p)
		j.current.Store(&cur)
		if info.Mode().IsRegular() {
			j.doneBytes.Add(info.Size())
		}
	}
	if err := os.Remove(p); err != nil && !errors.Is(err, fs.ErrNotExist) {
		r.addErr(p, err)
	}
	j.doneFiles.Add(1)
	return nil
}
