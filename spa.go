package spa

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

//go:embed all:frontend/.output/public
var spaFiles embed.FS

var distFS fs.FS

func init() {
	var err error
	distFS, err = fs.Sub(spaFiles, "frontend/.output/public")

	if err != nil {
		panic(err)
	}
}

var mimeTypes = map[string]string{
	".js":    "application/javascript",
	".css":   "text/css",
	".html":  "text/html; charset=utf-8",
	".json":  "application/json",
	".svg":   "image/svg+xml",
	".png":   "image/png",
	".ico":   "image/x-icon",
	".woff2": "font/woff2",
	".woff":  "font/woff",
}

func SPAHandler() http.Handler {
	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/")

		// serve static assets (files only, not directories)
		if isFile(distFS, p) {
			// Set correct MIME type for JS modules
			ext := path.Ext(p)
			if ct, ok := mimeTypes[ext]; ok {
				w.Header().Set("Content-Type", ct)
			}
			fileServer.ServeHTTP(w, r)
			return
		}

		// fallback to index.html for Vue Router
		indexHTML, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		defer indexHTML.Close()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", time.Time{}, indexHTML.(io.ReadSeeker))
	})
}

// isFile returns true only if the path exists AND is not a directory
func isFile(fsys fs.FS, name string) bool {
	if name == "" {
		return false
	}
	info, err := fs.Stat(fsys, name)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func fileExists(fsys fs.FS, name string) bool {
	_, err := fs.Stat(fsys, name)
	return err == nil
}

func fsModTime(fsys fs.FS, name string) time.Time {
	info, err := fs.Stat(fsys, name)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
