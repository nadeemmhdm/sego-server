package httpapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nadeemmhdm/sego-server/internal/audit"
	"github.com/nadeemmhdm/sego-server/internal/auth"
	"github.com/nadeemmhdm/sego-server/internal/db"
	"github.com/nadeemmhdm/sego-server/internal/namespace"
	"github.com/nadeemmhdm/sego-server/internal/securitylock"
)

func isAdmin(u db.User) bool { return u.Role == "super_admin" || u.Role == "admin" }

func (s *Server) auditUser(r *http.Request, action, target, outcome string, details any) {
	u := currentUser(r)
	uid := u.ID
	_ = audit.Write(r.Context(), s.db, audit.Event{
		Action: action, ActorUserID: &uid, ActorUsername: u.Username, ActorRole: u.Role,
		IP: clientIP(r), UserAgent: r.UserAgent(), Target: target, Outcome: outcome, Details: details,
	})
}

func (s *Server) auditActor(r *http.Request, action string, userID *int64, username, role, target, outcome string, details any) {
	_ = audit.Write(r.Context(), s.db, audit.Event{
		Action: action, ActorUserID: userID, ActorUsername: username, ActorRole: role,
		IP: clientIP(r), UserAgent: r.UserAgent(), Target: target, Outcome: outcome, Details: details,
	})
}

func (s *Server) adminSecurityStatus(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !isAdmin(u) {
		problem(w, 403, "admin role required")
		return
	}
	state, err := securitylock.StateNow(r.Context(), s.db)
	if err != nil {
		problem(w, 500, "security state unavailable")
		return
	}
	writeJSON(w, 200, state)
}

func (s *Server) emergencyActivate(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !isAdmin(u) {
		problem(w, 403, "admin role required")
		return
	}
	var in struct {
		Reason string `json:"reason"`
	}
	// Empty body is allowed so the UI can remain a true one-click emergency action.
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&in); err != nil {
			problem(w, 400, "invalid request")
			return
		}
	}
	err := securitylock.Activate(r.Context(), s.db, securitylock.Actor{
		UserID: u.ID, Username: u.Username, Role: u.Role, IP: clientIP(r), UserAgent: r.UserAgent(),
	}, strings.TrimSpace(in.Reason))
	if err != nil {
		problem(w, 500, "emergency lockdown failed")
		return
	}
	writeJSON(w, 200, map[string]any{
		"locked":  true,
		"message": "Emergency lockdown active. All sessions and preview tokens were revoked and all share links were expired. Remote logins remain blocked until local CLI recovery.",
	})
}

func (s *Server) auditList(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !isAdmin(u) {
		problem(w, 403, "admin role required")
		return
	}
	limit, _ := strconv.Atoi(q(r, "limit", "100"))
	offset, _ := strconv.Atoi(q(r, "offset", "0"))
	entries, err := audit.List(r.Context(), s.db, limit, offset)
	if err != nil {
		problem(w, 500, "audit log unavailable")
		return
	}
	writeJSON(w, 200, map[string]any{"items": entries, "limit": limit, "offset": offset})
}

type previewSpec struct {
	Kind, MIME string
	MaxBytes   int64
}

func previewType(name string) (previewSpec, bool) {
	ext := strings.ToLower(filepath.Ext(name))
	video := map[string]string{".mp4": "video/mp4", ".webm": "video/webm", ".ogv": "video/ogg", ".mov": "video/quicktime", ".m4v": "video/x-m4v"}
	audio := map[string]string{".mp3": "audio/mpeg", ".wav": "audio/wav", ".ogg": "audio/ogg", ".m4a": "audio/mp4", ".aac": "audio/aac", ".flac": "audio/flac"}
	images := map[string]string{".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png", ".webp": "image/webp", ".gif": "image/gif", ".bmp": "image/bmp", ".ico": "image/x-icon"}
	if m, ok := video[ext]; ok {
		return previewSpec{"video", m, 0}, true
	}
	if m, ok := audio[ext]; ok {
		return previewSpec{"audio", m, 0}, true
	}
	if m, ok := images[ext]; ok {
		return previewSpec{"image", m, 100 << 20}, true
	}
	if ext == ".pdf" {
		return previewSpec{"pdf", "application/pdf", 250 << 20}, true
	}
	textExt := map[string]bool{
		".txt": true, ".md": true, ".markdown": true, ".json": true, ".csv": true, ".log": true, ".xml": true, ".yaml": true, ".yml": true, ".toml": true, ".ini": true, ".conf": true,
		".go": true, ".py": true, ".js": true, ".jsx": true, ".ts": true, ".tsx": true, ".css": true, ".scss": true, ".html": false, ".htm": false, ".svg": false,
		".java": true, ".c": true, ".h": true, ".cpp": true, ".hpp": true, ".rs": true, ".sh": true, ".ps1": true, ".sql": true,
	}
	if allowed, exists := textExt[ext]; exists && allowed {
		return previewSpec{"text", "text/plain; charset=utf-8", 8 << 20}, true
	}
	return previewSpec{}, false
}

func (s *Server) createPreviewToken(w http.ResponseWriter, r *http.Request) {
	var in struct{ Root, Path string }
	if decode(w, r, &in) != nil {
		return
	}
	if in.Root == "" {
		in.Root = "primary"
	}
	u := currentUser(r)
	internalPath := namespace.UserPath(u.ID, in.Path)
	f, info, err := s.st.Open(in.Root, internalPath)
	if err != nil {
		problem(w, 404, "file unavailable")
		return
	}
	_ = f.Close()
	if info.IsDir() {
		problem(w, 400, "folders cannot be previewed")
		return
	}
	spec, ok := previewType(info.Name())
	if !ok {
		problem(w, 415, "this file type is not allowed for inline preview")
		return
	}
	if spec.MaxBytes > 0 && info.Size() > spec.MaxBytes {
		problem(w, 413, "file is too large for safe inline preview")
		return
	}
	plain, tokenHash, err := auth.NewOpaqueToken()
	if err != nil {
		problem(w, 500, "preview token creation failed")
		return
	}
	ttl := 5 * time.Minute
	if spec.Kind == "video" || spec.Kind == "audio" {
		// Media players issue additional Range requests while seeking. Keep the
		// token scoped and temporary, but long enough for normal playback.
		ttl = 30 * time.Minute
	}
	exp := time.Now().UTC().Add(ttl)
	_, err = s.db.ExecContext(r.Context(), `INSERT INTO preview_tokens(token_hash,owner_user_id,root_id,rel_path,expires_at,created_at) VALUES(?,?,?,?,?,?)`,
		tokenHash, u.ID, in.Root, internalPath, exp.Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		problem(w, 500, "preview token persistence failed")
		return
	}
	s.auditUser(r, "file_preview_token_created", in.Path, "success", map[string]any{"kind": spec.Kind, "expires_at": exp})
	writeJSON(w, 201, map[string]any{"url": "/api/preview/" + plain, "kind": spec.Kind, "mime": spec.MIME, "name": info.Name(), "size": info.Size(), "expires_at": exp})
}

func (s *Server) previewByToken(w http.ResponseWriter, r *http.Request) {
	if securitylock.IsActive(r.Context(), s.db) {
		problem(w, 423, "emergency lockdown active")
		return
	}
	tokenHash := auth.HashToken(chi.URLParam(r, "token"))
	var ownerID int64
	var rootID, relPath, username, role string
	err := s.db.QueryRowContext(r.Context(), `SELECT p.owner_user_id,p.root_id,p.rel_path,u.username,u.role
FROM preview_tokens p JOIN users u ON u.id=p.owner_user_id WHERE p.token_hash=? AND p.expires_at>?`, tokenHash, time.Now().UTC().Format(time.RFC3339)).Scan(&ownerID, &rootID, &relPath, &username, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			problem(w, 404, "preview unavailable")
			return
		}
		problem(w, 500, "preview unavailable")
		return
	}
	f, info, err := s.st.Open(rootID, relPath)
	if err != nil {
		problem(w, 404, "file unavailable")
		return
	}
	defer f.Close()
	spec, ok := previewType(info.Name())
	if !ok {
		problem(w, 415, "preview type is not allowed")
		return
	}
	if spec.MaxBytes > 0 && info.Size() > spec.MaxBytes {
		problem(w, 413, "file is too large for safe inline preview")
		return
	}
	w.Header().Set("Content-Type", spec.MIME)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename=%q`, sanitizeFilename(info.Name())))
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Accept-Ranges", "bytes")
	// Active HTML/SVG is never admitted by previewType. A restrictive policy is
	// still set on preview responses as a second layer for browser rendering.
	w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; media-src 'self'; style-src 'none'; script-src 'none'; sandbox")
	uid := ownerID
	_ = audit.Write(r.Context(), s.db, audit.Event{Action: "file_preview_opened", ActorUserID: &uid, ActorUsername: username, ActorRole: role, IP: clientIP(r), UserAgent: r.UserAgent(), Target: namespace.PublicPath(ownerID, relPath), Outcome: "success", Details: map[string]any{"kind": spec.Kind}})
	http.ServeContent(w, r, info.Name(), info.ModTime(), f)
}

func contentTypeForDownload(name string) string {
	if t := mime.TypeByExtension(filepath.Ext(name)); t != "" {
		return t
	}
	return "application/octet-stream"
}
