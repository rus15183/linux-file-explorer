// Linux File Explorer — a web file manager with a sequential copy/move queue.
package main

import (
	"crypto/subtle"
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

//go:embed web
var webFS embed.FS

var version = "dev"

// App holds the configuration shared by HTTP handlers and the job queue.
type App struct {
	root   string // real directory exposed as "/" in the UI
	syncOn bool   // flush data to disk during and after copy
	q      *Queue
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func main() {
	root := flag.String("root", env("LFE_ROOT", "/mnt"), "directory exposed in the web UI")
	listen := flag.String("listen", env("LFE_LISTEN", ":8080"), "listen address")
	user := flag.String("user", env("LFE_USER", ""), "basic auth user (optional)")
	pass := flag.String("password", env("LFE_PASSWORD", ""), "basic auth password (optional)")
	flag.Parse()

	abs, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}
	if r, err := filepath.EvalSymlinks(abs); err == nil {
		abs = r
	}
	if st, err := os.Stat(abs); err != nil || !st.IsDir() {
		log.Fatalf("root %q is not a directory", abs)
	}

	app := &App{root: abs, syncOn: env("LFE_SYNC", "1") != "0"}
	app.q = newQueue(app)
	go app.q.loop()

	mux := http.NewServeMux()
	app.routes(mux)
	static, _ := fs.Sub(webFS, "web")
	files := http.FileServer(http.FS(static))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		files.ServeHTTP(w, r)
	}))

	srv := &http.Server{
		Addr:              *listen,
		Handler:           guard(mux, *user, *pass),
		ReadHeaderTimeout: 15 * time.Second,
	}
	log.Printf("Linux File Explorer %s: root=%s listen=%s auth=%v", version, abs, *listen, *user != "")
	log.Fatal(srv.ListenAndServe())
}

// guard applies optional basic auth and blocks cross-site state-changing requests.
func guard(next http.Handler, user, pass string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user != "" || pass != "" {
			u, p, ok := r.BasicAuth()
			okU := subtle.ConstantTimeCompare([]byte(u), []byte(user)) == 1
			okP := subtle.ConstantTimeCompare([]byte(p), []byte(pass)) == 1
			if !ok || !okU || !okP {
				w.Header().Set("WWW-Authenticate", `Basic realm="Linux File Explorer", charset="UTF-8"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Header.Get("X-LFE") == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// resolve maps a virtual UI path ("/sda1/dir") to a real path under root.
func (a *App) resolve(p string) (string, error) {
	if strings.ContainsRune(p, 0) {
		return "", errors.New("invalid path")
	}
	return filepath.Join(a.root, filepath.FromSlash(path.Clean("/"+p))), nil
}

// virt maps a real path back to the virtual UI path.
func (a *App) virt(real string) string {
	rel, err := filepath.Rel(a.root, real)
	if err != nil || rel == "." {
		return "/"
	}
	return "/" + filepath.ToSlash(rel)
}
