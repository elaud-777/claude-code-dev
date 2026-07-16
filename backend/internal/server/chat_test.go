package server_test

import (
	"fmt"
	"strings"
	"testing"

	"taskflow-backend/internal/testsupport"
)

func TestSendAndListMessages(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, _ := testsupport.SignupAndLogin(t, ts, "chat1@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")

	resp, _ := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/messages", team["id"]), auth, map[string]any{"content": "안녕하세요"})
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	_, messages := testsupport.GetJSONSlice(t, ts, fmt.Sprintf("/teams/%v/messages", team["id"]), auth)
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
}

func TestMessageTooLong(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, _ := testsupport.SignupAndLogin(t, ts, "chat2@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")

	longContent := strings.Repeat("x", 1001)
	resp, body := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/messages", team["id"]), auth, map[string]any{"content": longContent})
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	errObj := body["error"].(map[string]any)
	if errObj["code"] != "TOO_LONG" {
		t.Fatalf("expected TOO_LONG, got %v", errObj["code"])
	}
}

func TestPollingSinceParameter(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, _ := testsupport.SignupAndLogin(t, ts, "chat3@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")

	_, first := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/messages", team["id"]), auth, map[string]any{"content": "첫 메시지"})

	_, none := testsupport.GetJSONSlice(t, ts, fmt.Sprintf("/teams/%v/messages?since=%v", team["id"], first["created_at"]), auth)
	if len(none) != 0 {
		t.Fatalf("expected 0 messages since first's own timestamp, got %d", len(none))
	}

	_, second := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/messages", team["id"]), auth, map[string]any{"content": "두번째 메시지"})
	_, some := testsupport.GetJSONSlice(t, ts, fmt.Sprintf("/teams/%v/messages?since=%v", team["id"], first["created_at"]), auth)
	if len(some) != 1 {
		t.Fatalf("expected 1 message after since, got %d", len(some))
	}
	got := some[0].(map[string]any)
	if fmt.Sprintf("%v", got["id"]) != fmt.Sprintf("%v", second["id"]) {
		t.Fatalf("expected second message id %v, got %v", second["id"], got["id"])
	}
}

func TestDeleteOwnMessage(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, _ := testsupport.SignupAndLogin(t, ts, "chat4@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")

	_, msg := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/messages", team["id"]), auth, map[string]any{"content": "삭제할 메시지"})

	resp, _ := testsupport.DeleteReq(t, ts, fmt.Sprintf("/messages/%v", msg["id"]), auth)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDeleteOthersMessageForbiddenEvenForOwner(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	ownerAuth, _ := testsupport.SignupAndLogin(t, ts, "chatowner@example.com", "password1")
	team := makeTeam(t, ts, ownerAuth, "Frontiers")

	memberAuth, _ := testsupport.SignupAndLogin(t, ts, "chatmember@example.com", "password1")
	testsupport.PostJSONFull(t, ts, "/teams/join", memberAuth, map[string]any{"invite_code": team["invite_code"]})

	_, msg := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/messages", team["id"]), memberAuth, map[string]any{"content": "멤버 메시지"})

	resp, body := testsupport.DeleteReq(t, ts, fmt.Sprintf("/messages/%v", msg["id"]), ownerAuth)
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
	errObj := body["error"].(map[string]any)
	if errObj["code"] != "NOT_OWNER" {
		t.Fatalf("expected NOT_OWNER, got %v", errObj["code"])
	}
}
