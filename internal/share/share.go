package share

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/nadeemmhdm/sego-server/internal/auth"
)

type Record struct {
	ID              int64
	OwnerUserID     int64
	RootID, RelPath string
	AllowDownload   bool
	ExpiresAt       sql.NullString
	MaxDownloads    sql.NullInt64
	DownloadCount   int64
	RevokedAt       sql.NullString
}

func Create(ctx context.Context, d *sql.DB, owner int64, rootID, rel string, expires *time.Time, maxDownloads *int64) (string, error) {
	plain, hash, err := auth.NewOpaqueToken()
	if err != nil {
		return "", err
	}
	var exp any = nil
	if expires != nil {
		exp = expires.UTC().Format(time.RFC3339)
	}
	var max any = nil
	if maxDownloads != nil {
		max = *maxDownloads
	}
	_, err = d.ExecContext(ctx, `INSERT INTO shares(token_hash,owner_user_id,root_id,rel_path,allow_download,expires_at,max_downloads,created_at) VALUES(?,?,?,?,1,?,?,?)`, hash, owner, rootID, rel, exp, max, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return "", err
	}
	return plain, nil
}

func Resolve(ctx context.Context, d *sql.DB, plain string) (Record, error) {
	var r Record
	var allow int
	err := d.QueryRowContext(ctx, `SELECT id,owner_user_id,root_id,rel_path,allow_download,expires_at,max_downloads,download_count,revoked_at FROM shares WHERE token_hash=?`, auth.HashToken(plain)).Scan(&r.ID, &r.OwnerUserID, &r.RootID, &r.RelPath, &allow, &r.ExpiresAt, &r.MaxDownloads, &r.DownloadCount, &r.RevokedAt)
	if err != nil {
		return r, err
	}
	r.AllowDownload = allow == 1
	if r.RevokedAt.Valid {
		return r, errors.New("share revoked")
	}
	if r.ExpiresAt.Valid {
		t, e := time.Parse(time.RFC3339, r.ExpiresAt.String)
		if e != nil || time.Now().UTC().After(t) {
			return r, errors.New("share expired")
		}
	}
	if r.MaxDownloads.Valid && r.DownloadCount >= r.MaxDownloads.Int64 {
		return r, errors.New("download limit reached")
	}
	return r, nil
}

func AcquireDownload(ctx context.Context, d *sql.DB, id int64) error {
	res, err := d.ExecContext(ctx, `UPDATE shares SET download_count=download_count+1 WHERE id=? AND revoked_at IS NULL AND (max_downloads IS NULL OR download_count < max_downloads)`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("download limit reached or share revoked")
	}
	return nil
}
func Revoke(ctx context.Context, d *sql.DB, id, owner int64) error {
	res, e := d.ExecContext(ctx, `UPDATE shares SET revoked_at=? WHERE id=? AND owner_user_id=?`, time.Now().UTC().Format(time.RFC3339), id, owner)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("share not found")
	}
	return nil
}
