package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Event struct {
	Action        string
	ActorUserID   *int64
	ActorUsername string
	ActorRole     string
	IP            string
	UserAgent     string
	Target        string
	Outcome       string
	Details       any
}

type Entry struct {
	ID            int64           `json:"id"`
	Action        string          `json:"action"`
	ActorUserID   *int64          `json:"actor_user_id,omitempty"`
	ActorUsername string          `json:"actor_username,omitempty"`
	ActorRole     string          `json:"actor_role,omitempty"`
	IP            string          `json:"ip,omitempty"`
	UserAgent     string          `json:"user_agent,omitempty"`
	Target        string          `json:"target,omitempty"`
	Outcome       string          `json:"outcome"`
	Details       json.RawMessage `json:"details"`
	CreatedAt     string          `json:"created_at"`
}

func Write(ctx context.Context, d *sql.DB, e Event) error {
	details := []byte("{}")
	if e.Details != nil {
		if b, err := json.Marshal(e.Details); err == nil {
			details = b
		}
	}
	_, err := d.ExecContext(ctx, `INSERT INTO audit_logs(action,actor_user_id,actor_username,actor_role,ip,user_agent,target,outcome,details_json,created_at)
VALUES(?,?,?,?,?,?,?,?,?,?)`, e.Action, e.ActorUserID, e.ActorUsername, e.ActorRole, e.IP, e.UserAgent, e.Target, e.Outcome, string(details), time.Now().UTC().Format(time.RFC3339))
	return err
}

func List(ctx context.Context, d *sql.DB, limit, offset int) ([]Entry, error) {
	if limit < 1 || limit > 500 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := d.QueryContext(ctx, `SELECT id,action,actor_user_id,actor_username,actor_role,ip,user_agent,target,outcome,details_json,created_at
FROM audit_logs ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Entry, 0, limit)
	for rows.Next() {
		var e Entry
		var uid sql.NullInt64
		var details string
		if err := rows.Scan(&e.ID, &e.Action, &uid, &e.ActorUsername, &e.ActorRole, &e.IP, &e.UserAgent, &e.Target, &e.Outcome, &details, &e.CreatedAt); err != nil {
			return nil, err
		}
		if uid.Valid {
			v := uid.Int64
			e.ActorUserID = &v
		}
		e.Details = json.RawMessage(details)
		out = append(out, e)
	}
	return out, rows.Err()
}
