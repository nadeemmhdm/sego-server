package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nadeemmhdm/sego-server/internal/audit"
	"github.com/nadeemmhdm/sego-server/internal/auth"
	"github.com/nadeemmhdm/sego-server/internal/config"
	"github.com/nadeemmhdm/sego-server/internal/db"
	"github.com/nadeemmhdm/sego-server/internal/diagnostics"
	"github.com/nadeemmhdm/sego-server/internal/namespace"
	"github.com/nadeemmhdm/sego-server/internal/securitylock"
	"github.com/nadeemmhdm/sego-server/internal/share"
	"github.com/nadeemmhdm/sego-server/internal/storage"
)

type loginBucket struct {
	Attempts int
	Reset    time.Time
}

type Server struct {
	cfg               config.Config
	db                *sql.DB
	st                *storage.Manager
	log               *slog.Logger
	uploads           chan struct{}
	started           time.Time
	loginMu           sync.Mutex
	loginAttempts     map[string]loginBucket
	dummyPasswordHash string
}
type ctxKey int

const userKey ctxKey = 1

func New(cfg config.Config, d *sql.DB, st *storage.Manager, log *slog.Logger) *Server {
	dummy, _ := auth.HashPassword("sego-dummy-password-only")
	return &Server{cfg: cfg, db: d, st: st, log: log, uploads: make(chan struct{}, cfg.UploadConcurrency), started: time.Now(), loginAttempts: map[string]loginBucket{}, dummyPasswordHash: dummy}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(s.securityHeaders, s.requestLog)
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"ok": true, "name": "SEGO Server"})
	})
	r.Post("/api/auth/login", s.login)
	r.Get("/api/share/{token}/download", s.publicShareDownload)
	r.Get("/api/preview/{token}", s.previewByToken)
	r.Group(func(p chi.Router) {
		p.Use(s.requireAuth)
		p.Post("/api/auth/logout", s.logout)
		p.Get("/api/me", s.me)
		p.Get("/api/files", s.listFiles)
		p.Post("/api/files/mkdir", s.mkdir)
		p.Post("/api/files/upload", s.upload)
		p.Get("/api/files/download", s.download)
		p.Post("/api/files/preview-token", s.createPreviewToken)
		p.Post("/api/files/rename", s.rename)
		p.Delete("/api/files", s.trash)
		p.Post("/api/shares", s.createShare)
		p.Delete("/api/shares/{id}", s.revokeShare)
		p.Get("/api/diagnostics", s.diag)
		p.Get("/api/admin/security/status", s.adminSecurityStatus)
		p.Post("/api/admin/emergency/activate", s.emergencyActivate)
		p.Get("/api/admin/audit", s.auditList)
	})
	r.HandleFunc("/*", s.webUI)
	return r
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; object-src 'none'")
		next.ServeHTTP(w, r)
	})
}
func (s *Server) requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		path := r.URL.Path
		if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
			path = rc.RoutePattern()
		}
		// Route patterns prevent bearer-like share/preview tokens from entering logs.
		s.log.Info("http_request", "method", r.Method, "path", path, "duration_ms", time.Since(start).Milliseconds())
	})
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if securitylock.IsActive(r.Context(), s.db) {
			problem(w, 423, "emergency lockdown active")
			return
		}
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			problem(w, 401, "missing bearer token")
			return
		}
		u, err := db.AuthenticateSession(r.Context(), s.db, auth.HashToken(strings.TrimPrefix(h, "Bearer ")))
		if err != nil {
			problem(w, 401, "invalid or expired session")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, u)))
	})
}
func currentUser(r *http.Request) db.User { return r.Context().Value(userKey).(db.User) }

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if securitylock.IsActive(r.Context(), s.db) {
		problem(w, 423, "emergency lockdown active; local CLI recovery is required")
		return
	}
	if !s.allowLogin(r) {
		problem(w, 429, "too many login attempts")
		return
	}
	var in struct{ Username, Password string }
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&in); err != nil {
		problem(w, 400, "invalid request")
		return
	}
	var id int64
	var hash, role string
	var quota, reserved int64
	lookupErr := s.db.QueryRowContext(r.Context(), `SELECT id,password_hash,role,quota_bytes,reserved_bytes FROM users WHERE username=?`, in.Username).Scan(&id, &hash, &role, &quota, &reserved)
	verifyHash := hash
	if lookupErr != nil {
		verifyHash = s.dummyPasswordHash
	}
	validPassword := auth.VerifyPassword(verifyHash, in.Password)
	if lookupErr != nil || !validPassword {
		time.Sleep(250 * time.Millisecond)
		s.auditActor(r, "login_failed", nil, in.Username, "", "authentication", "denied", map[string]any{"reason": "invalid_credentials"})
		problem(w, 401, "invalid credentials")
		return
	}
	plain, tokenHash, err := auth.NewOpaqueToken()
	if err != nil {
		problem(w, 500, "session creation failed")
		return
	}
	exp := time.Now().UTC().Add(time.Duration(s.cfg.SessionHours) * time.Hour)
	_, err = s.db.Exec(`INSERT INTO sessions(token_hash,user_id,expires_at,created_at) VALUES(?,?,?,?)`, tokenHash, id, exp.Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		problem(w, 500, "session persistence failed")
		return
	}
	s.clearLogin(r)
	uid := id
	s.auditActor(r, "login_success", &uid, in.Username, role, "authentication", "success", nil)
	writeJSON(w, 200, map[string]any{"token": plain, "expires_at": exp, "user": map[string]any{"id": id, "username": in.Username, "role": role, "quota_bytes": quota, "reserved_bytes": reserved}})
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	s.auditUser(r, "logout", "authentication", "success", nil)
	plain := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	s.db.Exec(`DELETE FROM sessions WHERE token_hash=?`, auth.HashToken(plain))
	w.WriteHeader(204)
}
func (s *Server) me(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, currentUser(r)) }

func (s *Server) listFiles(w http.ResponseWriter, r *http.Request) {
	root := q(r, "root", "primary")
	path := r.URL.Query().Get("path")
	u := currentUser(r)
	internal := namespace.UserPath(u.ID, path)
	if path == "" {
		_ = s.st.Mkdir(root, internal)
	}
	items, err := s.st.List(root, internal)
	if err != nil {
		problem(w, 400, err.Error())
		return
	}
	for i := range items {
		items[i].Path = namespace.PublicPath(u.ID, items[i].Path)
	}
	writeJSON(w, 200, map[string]any{"root": root, "path": path, "items": items})
}
func (s *Server) mkdir(w http.ResponseWriter, r *http.Request) {
	var in struct{ Root, Path string }
	if decode(w, r, &in) != nil {
		return
	}
	if in.Root == "" {
		in.Root = "primary"
	}
	if err := s.st.Mkdir(in.Root, namespace.UserPath(currentUser(r).ID, in.Path)); err != nil {
		problem(w, 400, err.Error())
		return
	}
	s.auditUser(r, "folder_created", in.Path, "success", map[string]any{"root": in.Root})
	w.WriteHeader(204)
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	select {
	case s.uploads <- struct{}{}:
		defer func() { <-s.uploads }()
	default:
		problem(w, 429, "upload concurrency limit reached")
		return
	}
	u := currentUser(r)
	root := q(r, "root", "primary")
	path := r.URL.Query().Get("path")
	if path == "" {
		problem(w, 400, "path is required")
		return
	}
	cl := r.ContentLength
	if cl < 0 {
		cl = s.cfg.MaxUploadBytes
	}
	if cl > s.cfg.MaxUploadBytes {
		problem(w, 413, "upload too large")
		return
	}
	tx, err := s.db.BeginTx(r.Context(), nil)
	if err != nil {
		problem(w, 500, "quota check failed")
		return
	}
	var quota, reserved, used int64
	if err = tx.QueryRow(`SELECT quota_bytes,reserved_bytes FROM users WHERE id=?`, u.ID).Scan(&quota, &reserved); err == nil {
		err = tx.QueryRow(`SELECT COALESCE(SUM(size_bytes),0) FROM file_index WHERE owner_user_id=? AND is_dir=0`, u.ID).Scan(&used)
	}
	if err != nil {
		tx.Rollback()
		problem(w, 500, "quota check failed")
		return
	}
	if used+reserved+cl > quota {
		tx.Rollback()
		problem(w, 413, "quota exceeded")
		return
	}
	if _, err = tx.Exec(`UPDATE users SET reserved_bytes=reserved_bytes+? WHERE id=?`, cl, u.ID); err != nil {
		tx.Rollback()
		problem(w, 500, "quota reserve failed")
		return
	}
	if err = tx.Commit(); err != nil {
		problem(w, 500, "quota reserve failed")
		return
	}
	defer s.db.Exec(`UPDATE users SET reserved_bytes=MAX(reserved_bytes-?,0) WHERE id=?`, cl, u.ID)
	internalPath := namespace.UserPath(u.ID, path)
	n, err := s.st.StreamWrite(root, internalPath, http.MaxBytesReader(w, r.Body, s.cfg.MaxUploadBytes+1), s.cfg.MaxUploadBytes)
	if err != nil {
		problem(w, 400, err.Error())
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	_, err = s.db.Exec(`INSERT INTO file_index(root_id,rel_path,owner_user_id,size_bytes,is_dir,created_at,updated_at) VALUES(?,?,?,?,0,?,?) ON CONFLICT(root_id,rel_path) DO UPDATE SET owner_user_id=excluded.owner_user_id,size_bytes=excluded.size_bytes,updated_at=excluded.updated_at`, root, internalPath, u.ID, n, now, now)
	if err != nil {
		problem(w, 500, "metadata update failed")
		return
	}
	s.auditUser(r, "file_uploaded", path, "success", map[string]any{"root": root, "size": n})
	writeJSON(w, 201, map[string]any{"path": path, "size": n})
}

func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	root := q(r, "root", "primary")
	path := r.URL.Query().Get("path")
	f, info, err := s.st.Open(root, namespace.UserPath(currentUser(r).ID, path))
	if err != nil {
		problem(w, 404, err.Error())
		return
	}
	defer f.Close()
	s.auditUser(r, "file_downloaded", path, "success", map[string]any{"root": root, "size": info.Size()})
	w.Header().Set("Content-Type", contentTypeForDownload(info.Name()))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, sanitizeFilename(info.Name())))
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(w, f)
}
func (s *Server) rename(w http.ResponseWriter, r *http.Request) {
	var in struct{ Root, From, To string }
	if decode(w, r, &in) != nil {
		return
	}
	if in.Root == "" {
		in.Root = "primary"
	}
	u := currentUser(r)
	fromInternal, toInternal := namespace.UserPath(u.ID, in.From), namespace.UserPath(u.ID, in.To)
	if err := s.st.Rename(in.Root, fromInternal, toInternal); err != nil {
		problem(w, 400, err.Error())
		return
	}
	s.db.Exec(`UPDATE file_index SET rel_path=?,updated_at=? WHERE root_id=? AND rel_path=? AND owner_user_id=?`, toInternal, time.Now().UTC().Format(time.RFC3339), in.Root, fromInternal, u.ID)
	s.auditUser(r, "file_renamed", in.To, "success", map[string]any{"root": in.Root, "from": in.From, "to": in.To})
	w.WriteHeader(204)
}
func (s *Server) trash(w http.ResponseWriter, r *http.Request) {
	root := q(r, "root", "primary")
	path := r.URL.Query().Get("path")
	u := currentUser(r)
	internalPath := namespace.UserPath(u.ID, path)
	if err := s.st.Trash(root, internalPath); err != nil {
		problem(w, 400, err.Error())
		return
	}
	s.db.Exec(`DELETE FROM file_index WHERE root_id=? AND (rel_path=? OR rel_path LIKE ?) AND owner_user_id=?`, root, internalPath, internalPath+"/%", u.ID)
	s.auditUser(r, "item_moved_to_trash", path, "success", map[string]any{"root": root})
	w.WriteHeader(204)
}

func (s *Server) createShare(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Root, Path   string
		ExpiresHours int64 `json:"expires_hours"`
		MaxDownloads int64 `json:"max_downloads"`
	}
	if decode(w, r, &in) != nil {
		return
	}
	if in.Root == "" {
		in.Root = "primary"
	}
	u := currentUser(r)
	internalPath := namespace.UserPath(u.ID, in.Path)
	if _, err := s.st.Resolve(in.Root, internalPath, true); err != nil {
		problem(w, 400, "invalid file path")
		return
	}
	var exp *time.Time
	if in.ExpiresHours > 0 {
		t := time.Now().UTC().Add(time.Duration(in.ExpiresHours) * time.Hour)
		exp = &t
	}
	var max *int64
	if in.MaxDownloads > 0 {
		max = &in.MaxDownloads
	}
	token, err := share.Create(r.Context(), s.db, u.ID, in.Root, internalPath, exp, max)
	if err != nil {
		problem(w, 500, "share creation failed")
		return
	}
	s.auditUser(r, "share_link_created", in.Path, "success", map[string]any{"root": in.Root, "expires_hours": in.ExpiresHours, "max_downloads": in.MaxDownloads})
	writeJSON(w, 201, map[string]any{"token": token, "url": "/api/share/" + token + "/download"})
}
func (s *Server) revokeShare(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := share.Revoke(r.Context(), s.db, id, currentUser(r).ID); err != nil {
		problem(w, 404, err.Error())
		return
	}
	s.auditUser(r, "share_link_revoked", fmt.Sprintf("share:%d", id), "success", nil)
	w.WriteHeader(204)
}
func (s *Server) publicShareDownload(w http.ResponseWriter, r *http.Request) {
	if securitylock.IsActive(r.Context(), s.db) {
		problem(w, 423, "emergency lockdown active")
		return
	}
	rec, err := share.Resolve(r.Context(), s.db, chi.URLParam(r, "token"))
	if err != nil {
		problem(w, 404, "share unavailable")
		return
	}
	if !rec.AllowDownload {
		problem(w, 403, "download not allowed")
		return
	}
	f, info, err := s.st.Open(rec.RootID, rec.RelPath)
	if err != nil {
		problem(w, 404, "file unavailable")
		return
	}
	defer f.Close()
	if err := share.AcquireDownload(r.Context(), s.db, rec.ID); err != nil {
		problem(w, 500, "share accounting failed")
		return
	}
	_ = audit.Write(r.Context(), s.db, audit.Event{Action: "share_download", IP: clientIP(r), UserAgent: r.UserAgent(), Target: fmt.Sprintf("share:%d", rec.ID), Outcome: "success", Details: map[string]any{"size": info.Size()}})
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, sanitizeFilename(info.Name())))
	io.Copy(w, f)
}
func (s *Server) diag(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if u.Role != "super_admin" && u.Role != "admin" {
		problem(w, 403, "admin role required")
		return
	}
	writeJSON(w, 200, diagnostics.SnapshotNow(s.started))
}

func clientIP(r *http.Request) string {
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i > 0 {
		host = host[:i]
	}
	return host
}

func (s *Server) allowLogin(r *http.Request) bool {
	now := time.Now()
	key := clientIP(r)
	s.loginMu.Lock()
	defer s.loginMu.Unlock()
	b := s.loginAttempts[key]
	if b.Reset.IsZero() || now.After(b.Reset) {
		b = loginBucket{Reset: now.Add(time.Minute)}
	}
	if b.Attempts >= 5 {
		s.loginAttempts[key] = b
		return false
	}
	b.Attempts++
	s.loginAttempts[key] = b
	return true
}

func (s *Server) clearLogin(r *http.Request) {
	s.loginMu.Lock()
	delete(s.loginAttempts, clientIP(r))
	s.loginMu.Unlock()
}

func sanitizeFilename(s string) string {
	return strings.NewReplacer("\r", "_", "\n", "_", "\"", "_").Replace(s)
}
func q(r *http.Request, key, def string) string {
	if v := r.URL.Query().Get(key); v != "" {
		return v
	}
	return def
}
func decode(w http.ResponseWriter, r *http.Request, v any) error {
	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(v)
	if err != nil {
		problem(w, 400, "invalid JSON")
	}
	return err
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func problem(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg, "status": status})
}
