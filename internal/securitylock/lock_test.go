package securitylock

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	segodb "github.com/nadeemmhdm/sego-server/internal/db"
)

func TestEmergencyLockdownRevokesAccess(t *testing.T) {
	d, err := segodb.Open(filepath.Join(t.TempDir(), "sego-test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	now := time.Now().UTC()
	res, err := d.Exec(`INSERT INTO users(username,password_hash,role,created_at) VALUES(?,?,?,?)`, "admin", "test-hash", "super_admin", now.Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}
	uid, _ := res.LastInsertId()
	_, err = d.Exec(`INSERT INTO sessions(token_hash,user_id,expires_at,created_at) VALUES(?,?,?,?)`, "session-hash", uid, now.Add(time.Hour).Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.Exec(`INSERT INTO shares(token_hash,owner_user_id,root_id,rel_path,created_at) VALUES(?,?,?,?,?)`, "share-hash", uid, "primary", "users/1/demo.txt", now.Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.Exec(`INSERT INTO preview_tokens(token_hash,owner_user_id,root_id,rel_path,expires_at,created_at) VALUES(?,?,?,?,?,?)`, "preview-hash", uid, "primary", "users/1/demo.txt", now.Add(time.Minute).Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}

	err = Activate(context.Background(), d, Actor{UserID: uid, Username: "admin", Role: "super_admin", IP: "127.0.0.1", UserAgent: "test"}, "security incident")
	if err != nil {
		t.Fatal(err)
	}

	state, err := StateNow(context.Background(), d)
	if err != nil || !state.Active {
		t.Fatalf("expected active lockdown, state=%+v err=%v", state, err)
	}
	assertCount(t, d, `SELECT COUNT(*) FROM sessions`, 0)
	assertCount(t, d, `SELECT COUNT(*) FROM preview_tokens`, 0)

	var revoked sql.NullString
	if err := d.QueryRow(`SELECT revoked_at FROM shares WHERE token_hash='share-hash'`).Scan(&revoked); err != nil {
		t.Fatal(err)
	}
	if !revoked.Valid || revoked.String == "" {
		t.Fatal("expected share link to be revoked")
	}
	assertCount(t, d, `SELECT COUNT(*) FROM audit_logs WHERE action='emergency_lockdown_activated'`, 1)

	if err := DeactivateLocal(context.Background(), d, uid, "admin", "super_admin"); err != nil {
		t.Fatal(err)
	}
	state, err = StateNow(context.Background(), d)
	if err != nil || state.Active {
		t.Fatalf("expected inactive lockdown after local recovery, state=%+v err=%v", state, err)
	}
	// Revoked share links intentionally stay revoked after recovery.
	if err := d.QueryRow(`SELECT revoked_at FROM shares WHERE token_hash='share-hash'`).Scan(&revoked); err != nil || !revoked.Valid {
		t.Fatalf("expected revoked link to stay revoked after recovery: %v", err)
	}
}

func assertCount(t *testing.T, d *sql.DB, query string, want int) {
	t.Helper()
	var got int
	if err := d.QueryRow(query).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("count for %q = %d, want %d", query, got, want)
	}
}
