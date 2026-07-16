package server_test

import (
	"testing"

	"taskflow-backend/internal/testsupport"
)

func TestSignupSuccess(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	resp, body := testsupport.PostJSONFull(t, ts, "/auth/signup", "", map[string]any{
		"email": "new@example.com", "password": "password1",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d body=%v", resp.StatusCode, body)
	}
	if body["token"] == nil {
		t.Fatalf("expected token in response, got %v", body)
	}
	user := body["user"].(map[string]any)
	if user["team_id"] != nil {
		t.Fatalf("expected team_id nil, got %v", user["team_id"])
	}
}

func TestSignupDuplicateEmail(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	testsupport.PostJSONFull(t, ts, "/auth/signup", "", map[string]any{"email": "dup@example.com", "password": "password1"})
	resp, body := testsupport.PostJSONFull(t, ts, "/auth/signup", "", map[string]any{"email": "dup@example.com", "password": "password1"})
	if resp.StatusCode != 409 {
		t.Fatalf("expected 409, got %d", resp.StatusCode)
	}
	errObj := body["error"].(map[string]any)
	if errObj["code"] != "EMAIL_TAKEN" {
		t.Fatalf("expected EMAIL_TAKEN, got %v", errObj["code"])
	}
}

func TestSignupWeakPassword(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	resp, body := testsupport.PostJSONFull(t, ts, "/auth/signup", "", map[string]any{"email": "weak@example.com", "password": "short"})
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	errObj := body["error"].(map[string]any)
	if errObj["code"] != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %v", errObj["code"])
	}
}

func TestLoginSuccess(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	testsupport.PostJSONFull(t, ts, "/auth/signup", "", map[string]any{"email": "login@example.com", "password": "password1"})
	resp, body := testsupport.PostJSONFull(t, ts, "/auth/login", "", map[string]any{"email": "login@example.com", "password": "password1"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if body["token"] == nil {
		t.Fatalf("expected token, got %v", body)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	testsupport.PostJSONFull(t, ts, "/auth/signup", "", map[string]any{"email": "login2@example.com", "password": "password1"})
	resp, body := testsupport.PostJSONFull(t, ts, "/auth/login", "", map[string]any{"email": "login2@example.com", "password": "wrongpass"})
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	errObj := body["error"].(map[string]any)
	if errObj["code"] != "INVALID_CREDENTIALS" {
		t.Fatalf("expected INVALID_CREDENTIALS, got %v", errObj["code"])
	}
}

func TestMeRequiresAuth(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	resp, _ := testsupport.GetJSON(t, ts, "/auth/me", "")
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
