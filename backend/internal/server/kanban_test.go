package server_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"taskflow-backend/internal/testsupport"
)

func makeTeam(t *testing.T, ts *httptest.Server, auth, name string) map[string]any {
	t.Helper()
	_, body := testsupport.PostJSONFull(t, ts, "/teams", auth, map[string]any{"name": name})
	return body
}

func TestCreateTask(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, _ := testsupport.SignupAndLogin(t, ts, "leader@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")

	resp, body := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/tasks", team["id"]), auth, map[string]any{"title": "테스트 태스크"})
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d body=%v", resp.StatusCode, body)
	}
	if body["status"] != "TODO" {
		t.Fatalf("expected TODO, got %v", body["status"])
	}
	if body["assignee_id"] != nil {
		t.Fatalf("expected nil assignee_id, got %v", body["assignee_id"])
	}
}

func TestStatusChange(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, _ := testsupport.SignupAndLogin(t, ts, "leader2@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")
	_, task := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/tasks", team["id"]), auth, map[string]any{"title": "이동할 태스크"})

	resp, body := testsupport.PatchJSON(t, ts, fmt.Sprintf("/tasks/%v/status", task["id"]), auth, map[string]any{"status": "DOING"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if body["status"] != "DOING" {
		t.Fatalf("expected DOING, got %v", body["status"])
	}
}

func TestDeleteByCreator(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, _ := testsupport.SignupAndLogin(t, ts, "leader3@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")
	_, task := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/tasks", team["id"]), auth, map[string]any{"title": "삭제될 태스크"})

	resp, _ := testsupport.DeleteReq(t, ts, fmt.Sprintf("/tasks/%v", task["id"]), auth)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDeleteByOwnerOnOtherMembersTask(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	ownerAuth, _ := testsupport.SignupAndLogin(t, ts, "owner@example.com", "password1")
	team := makeTeam(t, ts, ownerAuth, "Frontiers")

	memberAuth, _ := testsupport.SignupAndLogin(t, ts, "member@example.com", "password1")
	testsupport.PostJSONFull(t, ts, "/teams/join", memberAuth, map[string]any{"invite_code": team["invite_code"]})

	_, task := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/tasks", team["id"]), memberAuth, map[string]any{"title": "멤버 태스크"})

	resp, _ := testsupport.DeleteReq(t, ts, fmt.Sprintf("/tasks/%v", task["id"]), ownerAuth)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 (owner override), got %d", resp.StatusCode)
	}
}

func TestDeleteForbiddenForOtherMember(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	ownerAuth, _ := testsupport.SignupAndLogin(t, ts, "owner2@example.com", "password1")
	team := makeTeam(t, ts, ownerAuth, "Frontiers")

	memberAAuth, _ := testsupport.SignupAndLogin(t, ts, "membera@example.com", "password1")
	testsupport.PostJSONFull(t, ts, "/teams/join", memberAAuth, map[string]any{"invite_code": team["invite_code"]})

	memberBAuth, _ := testsupport.SignupAndLogin(t, ts, "memberb@example.com", "password1")
	testsupport.PostJSONFull(t, ts, "/teams/join", memberBAuth, map[string]any{"invite_code": team["invite_code"]})

	_, task := testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/tasks", team["id"]), memberAAuth, map[string]any{"title": "A의 태스크"})

	resp, body := testsupport.DeleteReq(t, ts, fmt.Sprintf("/tasks/%v", task["id"]), memberBAuth)
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d body=%v", resp.StatusCode, body)
	}
	errObj := body["error"].(map[string]any)
	if errObj["code"] != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %v", errObj["code"])
	}
}

func TestAssigneeFilter(t *testing.T) {
	ts := testsupport.NewTestServer(t)
	auth, user := testsupport.SignupAndLogin(t, ts, "filteruser@example.com", "password1")
	team := makeTeam(t, ts, auth, "Frontiers")

	testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/tasks", team["id"]), auth, map[string]any{"title": "내 태스크", "assignee_id": user["id"]})
	testsupport.PostJSONFull(t, ts, fmt.Sprintf("/teams/%v/tasks", team["id"]), auth, map[string]any{"title": "미할당 태스크"})

	_, meTasks := testsupport.GetJSONSlice(t, ts, fmt.Sprintf("/teams/%v/tasks?filter=me", team["id"]), auth)
	if len(meTasks) != 1 {
		t.Fatalf("expected 1 task for @me filter, got %d", len(meTasks))
	}

	_, unassigned := testsupport.GetJSONSlice(t, ts, fmt.Sprintf("/teams/%v/tasks?filter=unassigned", team["id"]), auth)
	if len(unassigned) != 1 {
		t.Fatalf("expected 1 unassigned task, got %d", len(unassigned))
	}
}
