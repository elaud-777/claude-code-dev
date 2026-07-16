package handlers

import (
	"database/sql"
	"strconv"

	"taskflow-backend/internal/app"
	apierrors "taskflow-backend/internal/errors"
	"taskflow-backend/internal/models"
)

func parseIDParam(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, apierrors.NotFound("")
	}
	return id, nil
}

// requireTeamMember loads the team and ensures the current user belongs to it,
// mirroring the Python `require_team_member` dependency.
func requireTeamMember(a *app.App, teamID int64, user *models.User) (*models.Team, error) {
	var team models.Team
	err := a.DB.Get(&team, a.Q("SELECT * FROM teams WHERE id = ?"), teamID)
	if err == sql.ErrNoRows {
		return nil, apierrors.NotFound("팀을 찾을 수 없습니다")
	}
	if err != nil {
		return nil, err
	}
	if user.TeamID == nil || *user.TeamID != teamID {
		return nil, apierrors.Forbidden()
	}
	return &team, nil
}
