package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"taskflow-backend/internal/app"
	apierrors "taskflow-backend/internal/errors"
	mw "taskflow-backend/internal/middleware"
	"taskflow-backend/internal/models"
)

type taskCreateRequest struct {
	Title      string `json:"title"`
	AssigneeID *int64 `json:"assignee_id"`
}

type taskUpdateRequest struct {
	Title      *string `json:"title"`
	AssigneeID *int64  `json:"assignee_id"`
}

type taskStatusUpdateRequest struct {
	Status string `json:"status"`
}

type taskOut struct {
	ID         int64     `json:"id"`
	TeamID     int64     `json:"team_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	CreatorID  int64     `json:"creator_id"`
	AssigneeID *int64    `json:"assignee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func toTaskOut(t *models.Task) taskOut {
	return taskOut{
		ID: t.ID, TeamID: t.TeamID, Title: t.Title, Status: t.Status,
		CreatorID: t.CreatorID, AssigneeID: t.AssigneeID, CreatedAt: t.CreatedAt,
	}
}

func getTaskOr404(a *app.App, taskID int64) (*models.Task, error) {
	var task models.Task
	err := a.DB.Get(&task, a.Q("SELECT * FROM tasks WHERE id = ?"), taskID)
	if err == sql.ErrNoRows {
		return nil, apierrors.NotFound("")
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask godoc
// @Summary 태스크 생성
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} taskOut
// @Router /teams/{teamId}/tasks [post]
func CreateTask(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		teamID, err := parseIDParam(chi.URLParam(r, "teamId"))
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, teamID, user); err != nil {
			return err
		}

		var req taskCreateRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}
		if len(req.Title) < 1 || len(req.Title) > 100 {
			return apierrors.ValidationError("제목은 1-100자여야 합니다")
		}

		now := time.Now().UTC()
		res, err := a.DB.Exec(
			a.Q("INSERT INTO tasks (team_id, title, status, creator_id, assignee_id, created_at) VALUES (?, ?, 'TODO', ?, ?, ?)"),
			teamID, req.Title, user.ID, req.AssigneeID, now,
		)
		if err != nil {
			return err
		}
		taskID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		mw.WriteJSON(w, 201, taskOut{
			ID: taskID, TeamID: teamID, Title: req.Title, Status: "TODO",
			CreatorID: user.ID, AssigneeID: req.AssigneeID, CreatedAt: now,
		})
		return nil
	}
}

// ListTasks godoc
// @Summary 칸반 태스크 목록 조회
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param filter query string false "all|me|unassigned"
// @Success 200 {array} taskOut
// @Router /teams/{teamId}/tasks [get]
func ListTasks(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		teamID, err := parseIDParam(chi.URLParam(r, "teamId"))
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, teamID, user); err != nil {
			return err
		}

		filter := r.URL.Query().Get("filter")
		var tasks []models.Task
		switch filter {
		case "me":
			err = a.DB.Select(&tasks,
				a.Q("SELECT * FROM tasks WHERE team_id = ? AND assignee_id = ? ORDER BY created_at DESC"),
				teamID, user.ID)
		case "unassigned":
			err = a.DB.Select(&tasks,
				a.Q("SELECT * FROM tasks WHERE team_id = ? AND assignee_id IS NULL ORDER BY created_at DESC"),
				teamID)
		default:
			err = a.DB.Select(&tasks,
				a.Q("SELECT * FROM tasks WHERE team_id = ? ORDER BY created_at DESC"),
				teamID)
		}
		if err != nil {
			return err
		}

		out := make([]taskOut, 0, len(tasks))
		for _, t := range tasks {
			out = append(out, toTaskOut(&t))
		}
		mw.WriteJSON(w, 200, out)
		return nil
	}
}

// GetTask godoc
// @Summary 태스크 상세 조회
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Success 200 {object} taskOut
// @Router /tasks/{taskId} [get]
func GetTask(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		taskID, err := parseIDParam(chi.URLParam(r, "taskId"))
		if err != nil {
			return err
		}
		task, err := getTaskOr404(a, taskID)
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, task.TeamID, user); err != nil {
			return err
		}
		mw.WriteJSON(w, 200, toTaskOut(task))
		return nil
	}
}

// UpdateTask godoc
// @Summary 태스크 제목/담당자 수정
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} taskOut
// @Router /tasks/{taskId} [put]
func UpdateTask(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		taskID, err := parseIDParam(chi.URLParam(r, "taskId"))
		if err != nil {
			return err
		}
		task, err := getTaskOr404(a, taskID)
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, task.TeamID, user); err != nil {
			return err
		}

		var req taskUpdateRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}

		title := task.Title
		if req.Title != nil {
			title = *req.Title
		}

		if _, err := a.DB.Exec(a.Q("UPDATE tasks SET title = ?, assignee_id = ? WHERE id = ?"), title, req.AssigneeID, taskID); err != nil {
			return err
		}

		task.Title = title
		task.AssigneeID = req.AssigneeID
		mw.WriteJSON(w, 200, toTaskOut(task))
		return nil
	}
}

// UpdateTaskStatus godoc
// @Summary 태스크 상태 변경 (드래그)
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} taskOut
// @Router /tasks/{taskId}/status [patch]
func UpdateTaskStatus(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		taskID, err := parseIDParam(chi.URLParam(r, "taskId"))
		if err != nil {
			return err
		}
		task, err := getTaskOr404(a, taskID)
		if err != nil {
			return err
		}
		if _, err := requireTeamMember(a, task.TeamID, user); err != nil {
			return err
		}

		var req taskStatusUpdateRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}
		if req.Status != "TODO" && req.Status != "DOING" && req.Status != "DONE" {
			return apierrors.ValidationError("상태는 TODO/DOING/DONE 중 하나여야 합니다")
		}

		if _, err := a.DB.Exec(a.Q("UPDATE tasks SET status = ? WHERE id = ?"), req.Status, taskID); err != nil {
			return err
		}

		task.Status = req.Status
		mw.WriteJSON(w, 200, toTaskOut(task))
		return nil
	}
}

// DeleteTask godoc
// @Summary 태스크 삭제 (creator 또는 owner만)
// @Tags tasks
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tasks/{taskId} [delete]
func DeleteTask(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		taskID, err := parseIDParam(chi.URLParam(r, "taskId"))
		if err != nil {
			return err
		}
		task, err := getTaskOr404(a, taskID)
		if err != nil {
			return err
		}
		team, err := requireTeamMember(a, task.TeamID, user)
		if err != nil {
			return err
		}
		if task.CreatorID != user.ID && team.OwnerID != user.ID {
			return apierrors.Forbidden()
		}

		if _, err := a.DB.Exec(a.Q("DELETE FROM tasks WHERE id = ?"), taskID); err != nil {
			return err
		}

		mw.WriteJSON(w, 200, map[string]any{})
		return nil
	}
}
