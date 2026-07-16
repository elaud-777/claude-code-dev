package testsupport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"

	"taskflow-backend/internal/app"
	"taskflow-backend/internal/config"
	"taskflow-backend/internal/server"
)

const testSchema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	team_id INTEGER,
	created_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS teams (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	invite_code TEXT UNIQUE NOT NULL,
	owner_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	team_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'TODO',
	creator_id INTEGER NOT NULL,
	assignee_id INTEGER,
	created_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	team_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);
`

// NewTestServer spins up an httptest.Server backed by a fresh in-memory SQLite DB.
func NewTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(testSchema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	a := &app.App{
		DB: db,
		Settings: config.Settings{
			DatabaseURL:    ":memory:",
			JWTSecret:      "test-secret",
			JWTExpireHours: 24,
			CORSOrigins:    []string{"*"},
		},
	}

	ts := httptest.NewServer(server.NewRouter(a))
	t.Cleanup(func() {
		ts.Close()
		db.Close()
	})
	return ts
}

// SignupAndLogin creates a user and returns an Authorization header value and the parsed user JSON.
func SignupAndLogin(t *testing.T, ts *httptest.Server, email, password string) (string, map[string]any) {
	t.Helper()
	PostJSON(t, ts, "/auth/signup", "", map[string]any{"email": email, "password": password})
	resp := PostJSON(t, ts, "/auth/login", "", map[string]any{"email": email, "password": password})
	token := resp["token"].(string)
	user := resp["user"].(map[string]any)
	return "Bearer " + token, user
}

func doRequest(t *testing.T, method, url, auth string, body any) *http.Response {
	t.Helper()
	var reader *strings.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = strings.NewReader(string(b))
	} else {
		reader = strings.NewReader("")
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func decodeBody(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return map[string]any{}
	}
	return out
}

func decodeBodySlice(t *testing.T, resp *http.Response) []any {
	t.Helper()
	defer resp.Body.Close()
	var out []any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out
}

func PostJSON(t *testing.T, ts *httptest.Server, path, auth string, body any) map[string]any {
	t.Helper()
	resp := doRequest(t, http.MethodPost, ts.URL+path, auth, body)
	return decodeBody(t, resp)
}

func PostJSONFull(t *testing.T, ts *httptest.Server, path, auth string, body any) (*http.Response, map[string]any) {
	t.Helper()
	resp := doRequest(t, http.MethodPost, ts.URL+path, auth, body)
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp, out
}

func GetJSON(t *testing.T, ts *httptest.Server, path, auth string) (*http.Response, map[string]any) {
	t.Helper()
	resp := doRequest(t, http.MethodGet, ts.URL+path, auth, nil)
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp, out
}

func GetJSONSlice(t *testing.T, ts *httptest.Server, path, auth string) (*http.Response, []any) {
	t.Helper()
	resp := doRequest(t, http.MethodGet, ts.URL+path, auth, nil)
	defer resp.Body.Close()
	var out []any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp, out
}

func PatchJSON(t *testing.T, ts *httptest.Server, path, auth string, body any) (*http.Response, map[string]any) {
	t.Helper()
	resp := doRequest(t, http.MethodPatch, ts.URL+path, auth, body)
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp, out
}

func DeleteReq(t *testing.T, ts *httptest.Server, path, auth string) (*http.Response, map[string]any) {
	t.Helper()
	resp := doRequest(t, http.MethodDelete, ts.URL+path, auth, nil)
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp, out
}
