package main

import (
	"archive/zip"
	"bufio"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// Already compressed formats are stored as is: deflating them again only
// burns CPU, which matters on a router.
var zipStored = map[string]bool{}

func init() {
	for _, e := range strings.Fields("zip 7z rar gz tgz bz2 xz zst lz4 lzma cab iso jar apk " +
		"jpg jpeg png gif webp heic avif mp3 aac m4a ogg opus flac wma " +
		"mp4 m4v mkv avi mov wmv webm mpg mpeg ts m2ts flv 3gp " +
		"docx xlsx pptx odt ods odp epub pdf") {
		zipStored["."+e] = true
	}
}

type zipper struct {
	r    *runner
	zw   *zip.Writer
	skip string // the archive being written: never packed into itself
	loc  *time.Location
	buf  []byte
}

// runZip packs the job sources into one zip archive at j.Dest. The archive is
// written under a temporary name and renamed when complete.
func (r *runner) runZip() error {
	j := r.job
	dst, err := r.app.resolve(j.Dest)
	if err != nil {
		return err
	}
	type source struct {
		path string
		info fs.FileInfo
	}
	var srcs []source
	for _, vs := range j.Sources {
		p, err := r.app.resolve(vs)
		if err != nil || p == r.app.root {
			j.addErr(vs + ": invalid source")
			continue
		}
		info, err := os.Lstat(p)
		if err != nil {
			r.addErr(p, err)
			continue
		}
		srcs = append(srcs, source{p, info})
	}
	if len(srcs) == 0 {
		return errors.New("nothing to compress")
	}

	j.setPhase("scan")
	for _, s := range srcs {
		if err := r.scan(s.path, s.info, false); err != nil {
			return err
		}
	}

	j.setPhase("zip")
	tmp := filepath.Join(filepath.Dir(dst), fmt.Sprintf(".lfe-%d-zip.part", j.ID))
	out, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return errors.New(errMsg(err))
	}
	done := false
	defer func() {
		if !done {
			out.Close()
			os.Remove(tmp)
		}
	}()
	bw := bufio.NewWriterSize(out, 1<<20)
	z := &zipper{r: r, zw: zip.NewWriter(bw), skip: tmp, buf: make([]byte, 1<<20),
		loc: time.FixedZone("", j.TZ*60)}
	z.zw.RegisterCompressor(zip.Deflate, func(w io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(w, flate.BestSpeed)
	})
	for _, s := range srcs {
		if err := z.add(s.path, filepath.Base(s.path), s.info); err != nil {
			return err
		}
	}
	if err := z.zw.Close(); err != nil {
		return errors.New(errMsg(err))
	}
	if err := bw.Flush(); err != nil {
		return errors.New(errMsg(err))
	}
	if r.app.syncOn {
		j.setPhase("sync")
		_ = unix.Fdatasync(int(out.Fd()))
	}
	if err := out.Close(); err != nil {
		return errors.New(errMsg(err))
	}
	if _, err := os.Lstat(dst); err == nil {
		dst = uniqueName(dst, false) // the name was taken while the job waited
	}
	if err := os.Rename(tmp, dst); err != nil {
		return errors.New(errMsg(err))
	}
	done = true
	return nil
}

// add packs one tree. It returns an error on cancel or when the archive
// itself cannot be written; unreadable sources are recorded and skipped.
func (z *zipper) add(p, name string, info fs.FileInfo) error {
	r, j := z.r, z.r.job
	if err := r.check(); err != nil {
		return err
	}
	if p == z.skip {
		return nil
	}
	h, err := zip.FileInfoHeader(info)
	if err != nil {
		r.addErr(p, err)
		return nil
	}
	h.Name = name
	h.Modified = info.ModTime().In(z.loc)

	if info.IsDir() {
		h.Name += "/"
		h.Method = zip.Store
		if _, err := z.zw.CreateHeader(h); err != nil {
			return errors.New(errMsg(err))
		}
		entries, err := os.ReadDir(p)
		if err != nil {
			r.addErr(p, err)
		}
		for _, e := range entries {
			fi, err := e.Info()
			if err != nil {
				r.addErr(filepath.Join(p, e.Name()), err)
				continue
			}
			if err := z.add(filepath.Join(p, e.Name()), name+"/"+e.Name(), fi); err != nil {
				return err
			}
		}
		return nil
	}

	v := r.app.virt(p)
	j.current.Store(&v)
	defer j.doneFiles.Add(1)

	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		target, err := os.Readlink(p)
		if err != nil {
			r.addErr(p, err)
			return nil
		}
		h.Method = zip.Store
		w, err := z.zw.CreateHeader(h)
		if err == nil {
			_, err = w.Write([]byte(target))
		}
		if err != nil {
			return errors.New(errMsg(err))
		}
	case info.Mode().IsRegular():
		in, err := os.Open(p)
		if err != nil {
			r.addErr(p, err)
			j.doneBytes.Add(info.Size())
			return nil
		}
		defer in.Close()
		_ = unix.Fadvise(int(in.Fd()), 0, 0, unix.FADV_SEQUENTIAL)
		h.Method = zip.Deflate
		if info.Size() == 0 || zipStored[strings.ToLower(filepath.Ext(name))] {
			h.Method = zip.Store
		}
		w, err := z.zw.CreateHeader(h)
		if err != nil {
			return errors.New(errMsg(err))
		}
		for {
			if err := r.check(); err != nil {
				return err
			}
			n, rerr := in.Read(z.buf)
			if n > 0 {
				if _, err := w.Write(z.buf[:n]); err != nil {
					return errors.New(errMsg(err))
				}
				j.doneBytes.Add(int64(n))
			}
			if rerr == io.EOF {
				break
			}
			if rerr != nil {
				r.addErr(p, rerr) // the entry stays in the archive, truncated
				break
			}
		}
	default:
		r.addErr(p, errors.New("special file skipped"))
	}
	return nil
}
