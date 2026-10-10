package share

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	segodb "github.com/nadeemmhdm/sego-server/internal/db"
)

func testDB(t *testing.T) (context.Context, interface{ Close() error }, func() any) {
	t.Helper()
	ctx := context.Background()
	d, err := segodb.Open(filepath.Join(t.TempDir(), "sego.db"))
	if err != nil { t.Fatal(err) }
	_, err = d.Exec(`INSERT INTO users(username,password_hash,role,created_at) VALUES('owner','x','user',?)`, time.Now().UTC().Format(time.RFC3339))
	if err != nil { d.Close(); t.Fatal(err) }
	return ctx, d, func() any { return d }
}

func TestShareDownloadLimitAndRevoke(t *testing.T) {
	ctx := context.Background()
	d, err := segodb.Open(filepath.Join(t.TempDir(), "sego.db"))
	if err != nil { t.Fatal(err) }
	defer d.Close()
	res, err := d.Exec(`INSERT INTO users(username,password_hash,role,created_at) VALUES('owner','x','user',?)`, time.Now().UTC().Format(time.RFC3339))
	if err != nil { t.Fatal(err) }
	owner, err := res.LastInsertId(); if err != nil { t.Fatal(err) }
	exp := time.Now().UTC().Add(time.Hour)
	max := int64(1)
	token, err := Create(ctx, d, owner, "primary", "users/1/file.txt", &exp, &max)
	if err != nil { t.Fatal(err) }
	rec, err := Resolve(ctx, d, token)
	if err != nil { t.Fatalf("resolve active share: %v", err) }
	if rec.OwnerUserID != owner || rec.RootID != "primary" { t.Fatalf("unexpected share record: %+v", rec) }
	if err := AcquireDownload(ctx, d, rec.ID); err != nil { t.Fatalf("first download: %v", err) }
	if _, err := Resolve(ctx, d, token); err == nil { t.Fatal("expected download limit to block second resolve") }
	if err := Revoke(ctx, d, rec.ID, owner); err != nil { t.Fatalf("revoke: %v", err) }
	if _, err := Resolve(ctx, d, token); err == nil { t.Fatal("expected revoked share to stay unavailable") }
}

func TestExpiredShareIsRejected(t *testing.T) {
	ctx := context.Background()
	d, err := segodb.Open(filepath.Join(t.TempDir(), "sego.db"))
	if err != nil { t.Fatal(err) }
	defer d.Close()
	res, err := d.Exec(`INSERT INTO users(username,password_hash,role,created_at) VALUES('owner2','x','user',?)`, time.Now().UTC().Format(time.RFC3339))
	if err != nil { t.Fatal(err) }
	owner, _ := res.LastInsertId()
	exp := time.Now().UTC().Add(-time.Minute)
	token, err := Create(ctx, d, owner, "primary", "users/1/file.txt", &exp, nil)
	if err != nil { t.Fatal(err) }
	if _, err := Resolve(ctx, d, token); err == nil { t.Fatal("expected expired share to be rejected") }
}
