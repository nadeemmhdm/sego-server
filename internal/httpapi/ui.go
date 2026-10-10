package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Server) webUI(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}
	root := s.cfg.WebRoot
	if root == "" {
		http.NotFound(w, r)
		return
	}
	rel := strings.TrimPrefix(filepath.Clean("/"+r.URL.Path), string(filepath.Separator))
	candidate := filepath.Join(root, rel)
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		http.ServeFile(w, r, candidate)
		return
	}
	index := filepath.Join(root, "index.html")
	if _, err := os.Stat(index); err != nil {
		problem(w, 503, "web dashboard is not installed")
		return
	}
	http.ServeFile(w, r, index)
}
