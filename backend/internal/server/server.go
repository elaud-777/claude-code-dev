package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"taskflow-backend/internal/app"
	"taskflow-backend/internal/handlers"
	mw "taskflow-backend/internal/middleware"
)

// NewRouter builds the full TaskFlow API route tree, shared by main.go and tests.
func NewRouter(a *app.App) http.Handler {
	r := chi.NewRouter()
	r.Use(mw.CORS(a.Settings.CORSOrigins))

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		mw.WriteJSON(w, 200, map[string]string{"status": "ok"})
	})

	r.Post("/auth/signup", mw.Wrap(handlers.Signup(a)))
	r.Post("/auth/login", mw.Wrap(handlers.Login(a)))

	r.Group(func(r chi.Router) {
		r.Use(mw.RequireAuth(a.DB, a.Settings.JWTSecret))

		r.Post("/auth/logout", mw.Wrap(handlers.Logout(a)))
		r.Get("/auth/me", mw.Wrap(handlers.Me(a)))

		r.Post("/teams", mw.Wrap(handlers.CreateTeam(a)))
		r.Post("/teams/join", mw.Wrap(handlers.JoinTeam(a)))
		r.Get("/teams/{teamId}", mw.Wrap(handlers.GetTeam(a)))
		r.Get("/teams/{teamId}/members", mw.Wrap(handlers.ListMembers(a)))
		r.Delete("/teams/{teamId}/leave", mw.Wrap(handlers.LeaveTeam(a)))

		r.Post("/teams/{teamId}/tasks", mw.Wrap(handlers.CreateTask(a)))
		r.Get("/teams/{teamId}/tasks", mw.Wrap(handlers.ListTasks(a)))
		r.Get("/tasks/{taskId}", mw.Wrap(handlers.GetTask(a)))
		r.Put("/tasks/{taskId}", mw.Wrap(handlers.UpdateTask(a)))
		r.Patch("/tasks/{taskId}/status", mw.Wrap(handlers.UpdateTaskStatus(a)))
		r.Delete("/tasks/{taskId}", mw.Wrap(handlers.DeleteTask(a)))

		r.Post("/teams/{teamId}/messages", mw.Wrap(handlers.SendMessage(a)))
		r.Get("/teams/{teamId}/messages", mw.Wrap(handlers.ListMessages(a)))
		r.Delete("/messages/{messageId}", mw.Wrap(handlers.DeleteMessage(a)))
	})

	return r
}
