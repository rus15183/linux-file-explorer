package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func (a *App) routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/list", a.hList)
	mux.HandleFunc("GET /api/disks", a.hDisks)
	mux.HandleFunc("GET /api/download", a.hDownload)
	mux.HandleFunc("POST /api/mkdir", a.hMkdir)
	mux.HandleFunc("POST /api/rename", a.hRename)
	mux.HandleFunc("POST /api/check", a.hCheck)
	mux.HandleFunc("GET /api/jobs", a.hJobs)
	mux.HandleFunc("POST /api/jobs", a.hCreateJob)
	mux.HandleFunc("POST /api/jobs/clear", a.hClear)
	mux.HandleFunc("POST /api/jobs/{id}/cancel", a.hCancel)
	mux.HandleFunc("DELETE /api/jobs/{id}", a.hRemove)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, fs.ErrNotExist):
		status = http.StatusNotFound
	case errors.Is(err, fs.ErrPermission):
		status = http.StatusForbidden
	case errors.Is(err, fs.ErrExist):
		status = http.StatusConflict
	}
	writeJSON(w, status, map[string]string{"error": errMsg(err)})
}

func badRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}

func readJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<20)).Decode(v); err != nil {
		badRequest(w, "invalid request")
		return false
	}
	return true
}

// errMsg strips real filesystem paths from OS errors.
func errMsg(err error) string {
	var pe *fs.PathError
	if errors.As(err, &pe) {
		return pe.Err.Error()
	}
	var le *os.LinkError
	if errors.As(err, &le) {
		return le.Err.Error()
	}
	return err.Error()
}

func validName(n string) bool {
	return n != "" && n != "." && n != ".." && !strings.ContainsAny(n, "/\x00") && len(n) <= 255
}

type entry struct {
	Name  string `json:"name"`
	Dir   bool   `json:"dir"`
	Size  int64  `json:"size"`
	Mtime int64  `json:"mtime"`
	Link  bool   `json:"link,omitempty"`
}

func (a *App) hList(w http.ResponseWriter, r *http.Request) {
	dir, err := a.resolve(r.URL.Query().Get("path"))
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	des, err := os.ReadDir(dir)
	if err != nil {
		writeErr(w, err)
		return
	}
	out := make([]entry, 0, len(des))
	for _, de := range des {
		info, err := de.Info()
		if err != nil {
			continue
		}
		e := entry{Name: de.Name()}
		if info.Mode()&fs.ModeSymlink != 0 {
			e.Link = true
			if t, err := os.Stat(filepath.Join(dir, de.Name())); err == nil {
				info = t
			}
		}
		e.Dir = info.IsDir()
		e.Mtime = info.ModTime().UnixMilli()
		if !e.Dir {
			e.Size = info.Size()
		}
		out = append(out, e)
	}
	writeJSON(w, 200, map[string]any{"path": a.virt(dir), "entries": out})
}

type disk struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Mount  bool   `json:"mount"`
	FS     string `json:"fs,omitempty"`
	Device string `json:"device,omitempty"`
}

// mountTable returns mount point -> [fstype, source] from /proc/self/mountinfo.
func mountTable() map[string][2]string {
	m := map[string][2]string{}
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return m
	}
	defer f.Close()
	unesc := strings.NewReplacer(`\040`, " ", `\011`, "\t", `\012`, "\n", `\134`, `\`)
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	for sc.Scan() {
		left, right, ok := strings.Cut(sc.Text(), " - ")
		lf, rf := strings.Fields(left), strings.Fields(right)
		if !ok || len(lf) < 5 || len(rf) < 2 {
			continue
		}
		m[unesc.Replace(lf[4])] = [2]string{rf[0], unesc.Replace(rf[1])}
	}
	return m
}

func devOf(p string) (uint64, bool) {
	var st unix.Stat_t
	if err := unix.Stat(p, &st); err != nil {
		return 0, false
	}
	return uint64(st.Dev), true
}

func (a *App) diskInfo(name, real string, rootDev uint64, mounts map[string][2]string) disk {
	d := disk{Name: name, Path: a.virt(real)}
	var sf unix.Statfs_t
	if err := unix.Statfs(real, &sf); err == nil {
		d.Total = uint64(sf.Blocks) * uint64(sf.Bsize)
		d.Free = uint64(sf.Bavail) * uint64(sf.Bsize)
	}
	if mi, ok := mounts[real]; ok {
		d.Mount, d.FS, d.Device = true, mi[0], mi[1]
	} else if dev, ok := devOf(real); ok && dev != rootDev {
		d.Mount = true
	}
	return d
}

func (a *App) hDisks(w http.ResponseWriter, r *http.Request) {
	mounts := mountTable()
	rootDev, _ := devOf(a.root)
	disks := []disk{}
	des, _ := os.ReadDir(a.root)
	for _, de := range des {
		p := filepath.Join(a.root, de.Name())
		if st, err := os.Stat(p); err != nil || !st.IsDir() {
			continue
		}
		disks = append(disks, a.diskInfo(de.Name(), p, rootDev, mounts))
	}
	sort.Slice(disks, func(i, j int) bool { return disks[i].Name < disks[j].Name })
	writeJSON(w, 200, map[string]any{"disks": disks, "version": version})
}

func (a *App) hDownload(w http.ResponseWriter, r *http.Request) {
	p, err := a.resolve(r.URL.Query().Get("path"))
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	st, err := os.Stat(p)
	if err != nil {
		writeErr(w, err)
		return
	}
	if !st.Mode().IsRegular() {
		badRequest(w, "not a regular file")
		return
	}
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filepath.Base(p)}))
	http.ServeFile(w, r, p)
}

func (a *App) hMkdir(w http.ResponseWriter, r *http.Request) {
	var req struct{ Path, Name string }
	if !readJSON(w, r, &req) {
		return
	}
	dir, err := a.resolve(req.Path)
	if err != nil || !validName(req.Name) {
		badRequest(w, "invalid name")
		return
	}
	if err := os.Mkdir(filepath.Join(dir, req.Name), 0o755); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, map[string]bool{"ok": true})
}

func (a *App) hRename(w http.ResponseWriter, r *http.Request) {
	var req struct{ Path, Name string }
	if !readJSON(w, r, &req) {
		return
	}
	src, err := a.resolve(req.Path)
	if err != nil || src == a.root || !validName(req.Name) {
		badRequest(w, "invalid name")
		return
	}
	dst := filepath.Join(filepath.Dir(src), req.Name)
	if _, err := os.Lstat(dst); err == nil {
		writeErr(w, &fs.PathError{Op: "rename", Path: dst, Err: syscall.EEXIST})
		return
	}
	if err := os.Rename(src, dst); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, map[string]bool{"ok": true})
}

type jobReq struct {
	Op       string   `json:"op"`
	Sources  []string `json:"sources"`
	Dest     string   `json:"dest"`
	Conflict string   `json:"conflict"`
}

// hCheck reports which source names already exist in the destination.
func (a *App) hCheck(w http.ResponseWriter, r *http.Request) {
	var req jobReq
	if !readJSON(w, r, &req) {
		return
	}
	dest, err := a.resolve(req.Dest)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	conflicts := []string{}
	for _, s := range req.Sources {
		src, err := a.resolve(s)
		if err != nil {
			continue
		}
		dst := filepath.Join(dest, filepath.Base(src))
		if dst == src {
			continue
		}
		if _, err := os.Lstat(dst); err == nil {
			conflicts = append(conflicts, filepath.Base(src))
		}
	}
	writeJSON(w, 200, map[string]any{"conflicts": conflicts})
}

func (a *App) hCreateJob(w http.ResponseWriter, r *http.Request) {
	var req jobReq
	if !readJSON(w, r, &req) {
		return
	}
	if req.Op != "copy" && req.Op != "move" && req.Op != "delete" {
		badRequest(w, "invalid operation")
		return
	}
	if len(req.Sources) == 0 {
		badRequest(w, "nothing selected")
		return
	}
	switch req.Conflict {
	case "skip", "rename", "overwrite":
	default:
		req.Conflict = "rename"
	}
	srcs := make([]string, 0, len(req.Sources))
	for _, s := range req.Sources {
		p, err := a.resolve(s)
		if err != nil || p == a.root {
			badRequest(w, "invalid source path")
			return
		}
		srcs = append(srcs, a.virt(p))
	}
	dest := ""
	if req.Op != "delete" {
		p, err := a.resolve(req.Dest)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		if st, err := os.Stat(p); err != nil || !st.IsDir() {
			badRequest(w, "destination is not a directory")
			return
		}
		dest = a.virt(p)
	}
	j := a.q.add(req.Op, srcs, dest, req.Conflict)
	writeJSON(w, 200, j.view())
}

func (a *App) hJobs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"jobs": a.q.views()})
}

func (a *App) hCancel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	a.q.cancel(id)
	writeJSON(w, 200, map[string]bool{"ok": true})
}

func (a *App) hRemove(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	a.q.remove(id)
	writeJSON(w, 200, map[string]bool{"ok": true})
}

func (a *App) hClear(w http.ResponseWriter, r *http.Request) {
	a.q.remove(0)
	writeJSON(w, 200, map[string]bool{"ok": true})
}
