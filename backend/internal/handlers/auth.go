package handlers

import (
	"net/http"
	"net/mail"
	"time"

	"taskflow-backend/internal/app"
	apierrors "taskflow-backend/internal/errors"
	mw "taskflow-backend/internal/middleware"
	"taskflow-backend/internal/models"
)

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userOut struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	TeamID *int64 `json:"team_id"`
}

type authResponse struct {
	Token string  `json:"token"`
	User  userOut `json:"user"`
}

func toUserOut(u *models.User) userOut {
	return userOut{ID: u.ID, Email: u.Email, TeamID: u.TeamID}
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Signup godoc
// @Summary 회원가입
// @Tags auth
// @Accept json
// @Produce json
// @Param body body signupRequest true "signup payload"
// @Success 201 {object} authResponse
// @Router /auth/signup [post]
func Signup(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var req signupRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}
		if !validateEmail(req.Email) {
			return apierrors.ValidationError("올바른 이메일 형식이 아닙니다")
		}
		if len(req.Password) < 8 {
			return apierrors.ValidationError("8자 이상 입력해주세요")
		}

		var count int
		if err := a.DB.Get(&count, a.Q("SELECT COUNT(*) FROM users WHERE email = ?"), req.Email); err != nil {
			return err
		}
		if count > 0 {
			return apierrors.EmailTaken()
		}

		hash, err := mw.HashPassword(req.Password)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		res, err := a.DB.Exec(
			a.Q("INSERT INTO users (email, password_hash, team_id, created_at) VALUES (?, ?, NULL, ?)"),
			req.Email, hash, now,
		)
		if err != nil {
			return err
		}
		userID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		token, err := mw.CreateAccessToken(userID, a.Settings.JWTSecret, a.Settings.JWTExpireHours)
		if err != nil {
			return err
		}

		mw.WriteJSON(w, 201, authResponse{
			Token: token,
			User:  userOut{ID: userID, Email: req.Email, TeamID: nil},
		})
		return nil
	}
}

// Login godoc
// @Summary 로그인
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "login payload"
// @Success 200 {object} authResponse
// @Router /auth/login [post]
func Login(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var req loginRequest
		if err := mw.DecodeJSON(r, &req); err != nil {
			return apierrors.ValidationError("")
		}

		var user models.User
		err := a.DB.Get(&user, a.Q("SELECT * FROM users WHERE email = ?"), req.Email)
		if err != nil {
			return apierrors.InvalidCredentials()
		}
		if !mw.VerifyPassword(req.Password, user.PasswordHash) {
			return apierrors.InvalidCredentials()
		}

		token, err := mw.CreateAccessToken(user.ID, a.Settings.JWTSecret, a.Settings.JWTExpireHours)
		if err != nil {
			return err
		}

		mw.WriteJSON(w, 200, authResponse{Token: token, User: toUserOut(&user)})
		return nil
	}
}

// Logout godoc
// @Summary 로그아웃 (stateless)
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /auth/logout [post]
func Logout(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		mw.WriteJSON(w, 200, map[string]any{})
		return nil
	}
}

// Me godoc
// @Summary 현재 사용자 조회
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} userOut
// @Router /auth/me [get]
func Me(a *app.App) mw.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := mw.CurrentUser(r)
		mw.WriteJSON(w, 200, toUserOut(user))
		return nil
	}
}
