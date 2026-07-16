package models

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	TeamID       *int64    `db:"team_id" json:"team_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Team struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	InviteCode string    `db:"invite_code" json:"invite_code"`
	OwnerID    int64     `db:"owner_id" json:"owner_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Task struct {
	ID         int64     `db:"id" json:"id"`
	TeamID     int64     `db:"team_id" json:"team_id"`
	Title      string    `db:"title" json:"title"`
	Status     string    `db:"status" json:"status"`
	CreatorID  int64     `db:"creator_id" json:"creator_id"`
	AssigneeID *int64    `db:"assignee_id" json:"assignee_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Message struct {
	ID        int64     `db:"id" json:"id"`
	TeamID    int64     `db:"team_id" json:"team_id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
