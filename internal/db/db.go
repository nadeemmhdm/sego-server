package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	if _, err := db.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000;`); err != nil {
		db.Close()
		return nil, err
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 username TEXT NOT NULL UNIQUE,
 password_hash TEXT NOT NULL,
 role TEXT NOT NULL CHECK(role IN ('super_admin','admin','group_manager','user')),
 quota_bytes INTEGER NOT NULL DEFAULT 10737418240,
 reserved_bytes INTEGER NOT NULL DEFAULT 0,
 disabled_at TEXT,
 created_at TEXT NOT NULL
);`,
		`CREATE TABLE IF NOT EXISTS sessions (
 token_hash TEXT PRIMARY KEY,
 user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 expires_at TEXT NOT NULL,
 created_at TEXT NOT NULL
);`,
		`CREATE TABLE IF NOT EXISTS file_index (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 root_id TEXT NOT NULL,
 rel_path TEXT NOT NULL,
 owner_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 size_bytes INTEGER NOT NULL DEFAULT 0,
 is_dir INTEGER NOT NULL DEFAULT 0,
 created_at TEXT NOT NULL,
 updated_at TEXT NOT NULL,
 UNIQUE(root_id, rel_path)
);`,
		`CREATE TABLE IF NOT EXISTS shares (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 token_hash TEXT NOT NULL UNIQUE,
 owner_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 root_id TEXT NOT NULL,
 rel_path TEXT NOT NULL,
 allow_download INTEGER NOT NULL DEFAULT 1,
 expires_at TEXT,
 max_downloads INTEGER,
 download_count INTEGER NOT NULL DEFAULT 0,
 revoked_at TEXT,
 created_at TEXT NOT NULL
);`,
		`CREATE TABLE IF NOT EXISTS preview_tokens (
 token_hash TEXT PRIMARY KEY,
 owner_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 root_id TEXT NOT NULL,
 rel_path TEXT NOT NULL,
 expires_at TEXT NOT NULL,
 created_at TEXT NOT NULL
);`,
		`CREATE TABLE IF NOT EXISTS security_state (
 id INTEGER PRIMARY KEY CHECK(id=1),
 active INTEGER NOT NULL DEFAULT 0,
 activated_at TEXT,
 activated_by_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
 reason TEXT
);`,
		`INSERT OR IGNORE INTO security_state(id,active) VALUES(1,0);`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 action TEXT NOT NULL,
 actor_user_id INTEGER,
 actor_username TEXT NOT NULL DEFAULT '',
 actor_role TEXT NOT NULL DEFAULT '',
 ip TEXT NOT NULL DEFAULT '',
 user_agent TEXT NOT NULL DEFAULT '',
 target TEXT NOT NULL DEFAULT '',
 outcome TEXT NOT NULL DEFAULT 'success',
 details_json TEXT NOT NULL DEFAULT '{}',
 created_at TEXT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_file_owner_path ON file_index(owner_user_id, root_id, rel_path);`,
		`CREATE INDEX IF NOT EXISTS idx_session_expiry ON sessions(expires_at);`,
		`CREATE INDEX IF NOT EXISTS idx_share_hash ON shares(token_hash);`,
		`CREATE INDEX IF NOT EXISTS idx_preview_expiry ON preview_tokens(expires_at);`,
		`CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_logs(created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_audit_actor ON audit_logs(actor_user_id, created_at DESC);`,
	}
	for i, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
	}
	// Existing installations created before account disabling was introduced do
	// not have disabled_at. Add it idempotently without resetting user data.
	if err := ensureColumn(db, "users", "disabled_at", "TEXT"); err != nil {
		return fmt.Errorf("users disabled_at migration: %w", err)
	}
	return nil
}

func ensureColumn(db *sql.DB, table, column, definition string) error {
	rows, err := db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull, pk int
		var defaultValue any
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if name == column {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_, err = db.Exec(`ALTER TABLE ` + table + ` ADD COLUMN ` + column + ` ` + definition)
	return err
}

type User struct {
	ID                        int64
	Username, Role            string
	QuotaBytes, ReservedBytes int64
}

func AuthenticateSession(ctx context.Context, d *sql.DB, tokenHash string) (User, error) {
	var u User
	err := d.QueryRowContext(ctx, `SELECT u.id,u.username,u.role,u.quota_bytes,u.reserved_bytes
FROM sessions s JOIN users u ON u.id=s.user_id
WHERE s.token_hash=? AND s.expires_at>? AND u.disabled_at IS NULL`, tokenHash, time.Now().UTC().Format(time.RFC3339)).Scan(&u.ID, &u.Username, &u.Role, &u.QuotaBytes, &u.ReservedBytes)
	return u, err
}
