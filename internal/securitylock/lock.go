package securitylock

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type State struct {
	Active        bool   `json:"active"`
	ActivatedAt   string `json:"activated_at,omitempty"`
	ActivatedByID *int64 `json:"activated_by_id,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

func StateNow(ctx context.Context, d *sql.DB) (State, error) {
	var s State
	var active int
	var at sql.NullString
	var by sql.NullInt64
	var reason sql.NullString
	err := d.QueryRowContext(ctx, `SELECT active,activated_at,activated_by_user_id,reason FROM security_state WHERE id=1`).Scan(&active, &at, &by, &reason)
	if err != nil {
		return s, err
	}
	s.Active = active != 0
	if at.Valid {
		s.ActivatedAt = at.String
	}
	if by.Valid {
		v := by.Int64
		s.ActivatedByID = &v
	}
	if reason.Valid {
		s.Reason = reason.String
	}
	return s, nil
}

func IsActive(ctx context.Context, d *sql.DB) bool {
	s, err := StateNow(ctx, d)
	return err == nil && s.Active
}

type Actor struct {
	UserID    int64
	Username  string
	Role      string
	IP        string
	UserAgent string
}

func Activate(ctx context.Context, d *sql.DB, actor Actor, reason string) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := time.Now().UTC().Format(time.RFC3339)
	if reason == "" {
		reason = "Emergency SOS activated by administrator"
	}
	if _, err = tx.ExecContext(ctx, `UPDATE security_state SET active=1,activated_at=?,activated_by_user_id=?,reason=? WHERE id=1`, now, actor.UserID, reason); err != nil {
		return err
	}
	// All public links are permanently revoked/expired during an emergency. They
	// are intentionally not restored by unlock; owners must create new links.
	if _, err = tx.ExecContext(ctx, `UPDATE shares SET revoked_at=COALESCE(revoked_at,?),expires_at=? WHERE revoked_at IS NULL`, now, now); err != nil {
		return err
	}
	// Preview URLs are short lived but are also destroyed immediately.
	if _, err = tx.ExecContext(ctx, `DELETE FROM preview_tokens`); err != nil {
		return err
	}
	// Revoke every user/admin browser session, including the actor's session.
	if _, err = tx.ExecContext(ctx, `DELETE FROM sessions`); err != nil {
		return err
	}
	details, _ := json.Marshal(map[string]any{"reason": reason, "sessions_revoked": true, "share_links_revoked": true, "preview_tokens_revoked": true})
	if _, err = tx.ExecContext(ctx, `INSERT INTO audit_logs(action,actor_user_id,actor_username,actor_role,ip,user_agent,target,outcome,details_json,created_at)
VALUES(?,?,?,?,?,?,?,?,?,?)`, "emergency_lockdown_activated", actor.UserID, actor.Username, actor.Role, actor.IP, actor.UserAgent, "server", "success", string(details), now); err != nil {
		return err
	}
	return tx.Commit()
}

func DeactivateLocal(ctx context.Context, d *sql.DB, userID int64, username, role string) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err = tx.ExecContext(ctx, `UPDATE security_state SET active=0,activated_at=NULL,activated_by_user_id=NULL,reason=NULL WHERE id=1`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO audit_logs(action,actor_user_id,actor_username,actor_role,ip,user_agent,target,outcome,details_json,created_at)
VALUES(?,?,?,?,?,?,?,?,?,?)`, "emergency_lockdown_disabled", userID, username, role, "local-cli", "sego-cli", "server", "success", `{"source":"local-cli"}`, now); err != nil {
		return err
	}
	return tx.Commit()
}
