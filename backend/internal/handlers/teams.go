package handlers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"

	"taskflow-backend/internal/app"
	apierrors "taskflow-backend/internal/errors"
	mw "taskflow-backend/internal/middleware"
	"taskflow-backend/internal/models"
)

var inviteCodeRe = regexp.MustCompile(`^[A-Z]{4}-[0-9]{4}$`)

type teamCreateRequest struct {
	Name string `json:"name"`
}

type teamJoinRequest struct {
	InviteCode string `json:"invite_code"`
}

type teamOut struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	InviteCode string    `json:"invite_code"`
	OwnerID    int64     `json:"owner_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type teamInfoOut struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	InviteCode  string `json:"invite_code"`
	MemberCount int    `json:"member_count"`
}

type memberOut struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	IsOwner   bool      `json:"is_owner"`
	CreatedAt time.Time `json:"created_at"`
}

func toTeamOut(t *models.Team) teamOut {
	return teamOut{ID: t.ID, Name: t.Name, InviteCode: t.InviteCode, OwnerID: t.OwnerID, CreatedAt: t.CreatedAt}
}

func generateInviteCode(a *app.App) (string, error) {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 20; i++ {
		code := make([]byte, 4)
		for j := range code {
			code[j] = letters[rand.Intn(len(letters))]
		}
		digits := rand.Intn(10000)
		candidate := fmt.Sprintf("%s-%04d", string(code), digits)

		var count int
		if err := a.DB.Get(&count, a.Q("SELECT COUNT(*) FROM teams WHERE invite_code = ?"), candidate); err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("초대코드 생성 실패: 재시도 한도 초과")
}

// CreateTeam godoc
// @Summary 팀 생성
// @Tags teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} teamOut
// @Router /teams [post]
func CreateTeam(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		var req teamCreateRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}
		if len(req.Name) < 1 || len(req.Name) > 30 {
			return apierrors.ValidationError("팀 이름은 1-30자여야 합니다")
		}

		code, err := generateInviteCode(a)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		res, err := a.DB.Exec(
			a.Q("INSERT INTO teams (name, invite_code, owner_id, created_at) VALUES (?, ?, ?, ?)"),
			req.Name, code, user.ID, now,
		)
		if err != nil {
			return err
		}
		teamID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		if _, err := a.DB.Exec(a.Q("UPDATE users SET team_id = ? WHERE id = ?"), teamID, user.ID); err != nil {
			return err
		}

		mw.WriteJSON(w, 201, teamOut{
			ID: teamID, Name: req.Name, InviteCode: code, OwnerID: user.ID, CreatedAt: now,
		})
		return nil
	}
}

// JoinTeam godoc
// @Summary 초대코드로 팀 합류
// @Tags teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} teamInfoOut
// @Router /teams/join [post]
func JoinTeam(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		var req teamJoinRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}
		if !inviteCodeRe.MatchString(req.InviteCode) {
			return apierrors.ValidationError("형식이 올바르지 않습니다")
		}
		if user.TeamID != nil {
			return apierrors.Conflict("")
		}

		var team models.Team
		err := a.DB.Get(&team, a.Q("SELECT * FROM teams WHERE invite_code = ?"), req.InviteCode)
		if err == sql.ErrNoRows {
			return apierrors.NotFound("해당 초대코드를 찾을 수 없습니다")
		}
		if err != nil {
			return err
		}

		if _, err := a.DB.Exec(a.Q("UPDATE users SET team_id = ? WHERE id = ?"), team.ID, user.ID); err != nil {
			return err
		}

		var memberCount int
		if err := a.DB.Get(&memberCount, a.Q("SELECT COUNT(*) FROM users WHERE team_id = ?"), team.ID); err != nil {
			return err
		}

		mw.WriteJSON(w, 200, teamInfoOut{ID: team.ID, Name: team.Name, InviteCode: team.InviteCode, MemberCount: memberCount})
		return nil
	}
}

// GetTeam godoc
// @Summary 팀 정보 조회
// @Tags teams
// @Security BearerAuth
// @Produce json
// @Success 200 {object} teamInfoOut
// @Router /teams/{teamId} [get]
func GetTeam(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		teamID, err := parseIDParam(chi.URLParam(r, "teamId"))
		if err != nil {
			return err
		}
		team, err := requireTeamMember(a, teamID, user)
		if err != nil {
			return err
		}

		var memberCount int
		if err := a.DB.Get(&memberCount, a.Q("SELECT COUNT(*) FROM users WHERE team_id = ?"), team.ID); err != nil {
			return err
		}

		mw.WriteJSON(w, 200, teamInfoOut{ID: team.ID, Name: team.Name, InviteCode: team.InviteCode, MemberCount: memberCount})
		return nil
	}
}

// ListMembers godoc
// @Summary 팀 멤버 목록
// @Tags teams
// @Security BearerAuth
// @Produce json
// @Success 200 {array} memberOut
// @Router /teams/{teamId}/members [get]
func ListMembers(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		teamID, err := parseIDParam(chi.URLParam(r, "teamId"))
		if err != nil {
			return err
		}
		team, err := requireTeamMember(a, teamID, user)
		if err != nil {
			return err
		}

		var members []models.User
		if err := a.DB.Select(&members, a.Q("SELECT * FROM users WHERE team_id = ?"), teamID); err != nil {
			return err
		}

		out := make([]memberOut, 0, len(members))
		for _, m := range members {
			out = append(out, memberOut{ID: m.ID, Email: m.Email, IsOwner: m.ID == team.OwnerID, CreatedAt: m.CreatedAt})
		}

		mw.WriteJSON(w, 200, out)
		return nil
	}
}

// LeaveTeam godoc
// @Summary 팀 나가기
// @Tags teams
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /teams/{teamId}/leave [delete]
func LeaveTeam(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		teamID, err := parseIDParam(chi.URLParam(r, "teamId"))
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, teamID, user); err != nil {
			return err
		}

		if _, err := a.DB.Exec(a.Q("UPDATE users SET team_id = NULL WHERE id = ?"), user.ID); err != nil {
			return err
		}

		mw.WriteJSON(w, 200, map[string]any{})
		return nil
	}
}
