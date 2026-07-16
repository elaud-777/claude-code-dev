package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"taskflow-backend/internal/app"
	apierrors "taskflow-backend/internal/errors"
	mw "taskflow-backend/internal/middleware"
)

const maxMessageLength = 1000

type messageCreateRequest struct {
	Content string `json:"content"`
}

type messageOut struct {
	ID        int64     `json:"id"`
	TeamID    int64     `json:"team_id"`
	UserID    int64     `json:"user_id"`
	UserEmail string    `json:"user_email"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// SendMessage godoc
// @Summary 팀 채팅 메시지 전송
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} messageOut
// @Router /teams/{teamId}/messages [post]
func SendMessage(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		teamID, err := parseIDParam(chi.URLParam(r, "teamId"))
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, teamID, user); err != nil {
			return err
		}

		var req messageCreateRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}
		if len(req.Content) < 1 {
			return apierrors.ValidationError("메시지를 입력해주세요")
		}
		if len(req.Content) > maxMessageLength {
			return apierrors.TooLong(maxMessageLength, len(req.Content))
		}

		now := time.Now().UTC()
		res, err := a.DB.Exec(
			a.Q("INSERT INTO messages (team_id, user_id, content, created_at) VALUES (?, ?, ?, ?)"),
			teamID, user.ID, req.Content, now,
		)
		if err != nil {
			return err
		}
		msgID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		mw.WriteJSON(w, 201, messageOut{
			ID: msgID, TeamID: teamID, UserID: user.ID, UserEmail: user.Email, Content: req.Content, CreatedAt: now,
		})
		return nil
	}
}

type messageRow struct {
	ID        int64     `db:"id"`
	TeamID    int64     `db:"team_id"`
	UserID    int64     `db:"user_id"`
	UserEmail string    `db:"email"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

// ListMessages godoc
// @Summary 팀 채팅 메시지 조회 (폴링)
// @Tags messages
// @Security BearerAuth
// @Produce json
// @Param since query string false "RFC3339 timestamp"
// @Success 200 {array} messageOut
// @Router /teams/{teamId}/messages [get]
func ListMessages(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		teamID, err := parseIDParam(chi.URLParam(r, "teamId"))
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, teamID, user); err != nil {
			return err
		}

		since := r.URL.Query().Get("since")
		var rows []messageRow

		baseQuery := `
			SELECT m.id, m.team_id, m.user_id, u.email as email, m.content, m.created_at
			FROM messages m JOIN users u ON m.user_id = u.id
			WHERE m.team_id = ?`

		if since != "" {
			sinceTime, parseErr := time.Parse(time.RFC3339, since)
			if parseErr != nil {
				return apierrors.ValidationError("since 파라미터 형식이 올바르지 않습니다")
			}
			err = a.DB.Select(&rows, a.Q(baseQuery+" AND m.created_at > ? ORDER BY m.created_at ASC"), teamID, sinceTime)
		} else {
			err = a.DB.Select(&rows, a.Q(baseQuery+" ORDER BY m.created_at DESC LIMIT 50"), teamID)
			if err == nil {
				for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
					rows[i], rows[j] = rows[j], rows[i]
				}
			}
		}
		if err != nil {
			return err
		}

		out := make([]messageOut, 0, len(rows))
		for _, row := range rows {
			out = append(out, messageOut{
				ID: row.ID, TeamID: row.TeamID, UserID: row.UserID,
				UserEmail: row.UserEmail, Content: row.Content, CreatedAt: row.CreatedAt,
			})
		}
		mw.WriteJSON(w, 200, out)
		return nil
	}
}

// DeleteMessage godoc
// @Summary 메시지 삭제 (본인만)
// @Tags messages
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /messages/{messageId} [delete]
func DeleteMessage(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		msgID, err := parseIDParam(chi.URLParam(r, "messageId"))
		if err != nil {
			return err
		}

		var row struct {
			TeamID int64 `db:"team_id"`
			UserID int64 `db:"user_id"`
		}
		err = a.DB.Get(&row, a.Q("SELECT team_id, user_id FROM messages WHERE id = ?"), msgID)
		if err == sql.ErrNoRows {
			return apierrors.NotFound("")
		}
		if err != nil {
			return err
		}

		if _, err := requireTeamMember(a, row.TeamID, user); err != nil {
			return err
		}
		if row.UserID != user.ID {
			return apierrors.NotOwner()
		}

		if _, err := a.DB.Exec(a.Q("DELETE FROM messages WHERE id = ?"), msgID); err != nil {
			return err
		}

		mw.WriteJSON(w, 200, map[string]any{})
		return nil
	}
}
