package pbc

import (
	"os"
	"testing"
)

var (
	api        = New("http://localhost:8090")
	userUUID   = "testuserid12345"
	adminToken string
)

func assert[T comparable](t *testing.T, x, y T) {
	t.Helper()
	if x != y {
		t.Fatalf("[ASSERT-FAILED] - %v != %v", x, y)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestClient(t *testing.T) {
	t.Run("TestHealth", func(t *testing.T) {
		health, err := api.Health()
		t.Logf("health: %v - %v", health, err)
		assert(t, health.Code, 200)
		assert(t, health.Data.CanBackup, true)
	})
	t.Run("TestAdminAuth", func(t *testing.T) {
		resp, err := api.AdminAuth("testadmin@test.com", "testadmin1234")
		user := ResponseTo[ResponseAdminAuth](resp)
		t.Logf("user: %v, error: %v", user, err)
		assert(t, resp.StatusCode, 200)
		assert(t, user.Admin.Email, "testadmin@test.com")
		adminToken = user.Token
	})
	api.RecordDelete("users", userUUID, WithAuthorization(adminToken))
	t.Run("TestRecordCreate", func(t *testing.T) {
		resp, err := api.RecordCreate("users", WithData(map[string]any{
			"id":              userUUID,
			"username":        "testuser",
			"email":           "testuser@test.com",
			"password":        "testuser1234",
			"passwordConfirm": "testuser1234",
		}))
		t.Logf("resp: %v, error: %v", resp, err)
		assert(t, resp.StatusCode, 200)
	})
	var token string
	t.Run("TestUserAuth", func(t *testing.T) {
		resp, err := api.UserAuth("testuser@test.com", "testuser1234")
		user := ResponseTo[ResponseAuth](resp)
		t.Logf("user: %v, error: %v", user, err)
		assert(t, resp.StatusCode, 200)
		assert(t, user.Record.Email, "testuser@test.com")
		token = user.Token
	})
	t.Run("TestRecordUpdate", func(t *testing.T) {
		resp, err := api.RecordUpdate("users", userUUID, WithAuthorization(token), WithData(map[string]string{
			"name": "testuser2",
		}))
		user := ResponseTo[UserRecord](resp)
		t.Logf("user: %v, error: %v", user, err)
		assert(t, resp.StatusCode, 200)
		assert(t, user.Name, "testuser2")
	})
	t.Run("TestRecordGet", func(t *testing.T) {
		resp, err := api.RecordGet("users", userUUID, WithAuthorization(token))
		user := ResponseTo[UserRecord](resp)
		t.Logf("user: %v, error: %v", user, err)
		assert(t, resp.StatusCode, 200)
		assert(t, user.ID, userUUID)
		assert(t, user.Email, "testuser@test.com")
	})
	t.Run("TestRecordList", func(t *testing.T) {
		resp, err := api.RecordList("users", WithPage(1, 10, false), WithAuthorization(token))
		users := ResponseTo[Records[UserRecord]](resp)
		t.Logf("users: %v, error: %v", users, err)
		assert(t, resp.StatusCode, 200)
		assert(t, len(users.Items), 1)
	})
	t.Run("TestRecordDelete", func(t *testing.T) {
		resp, err := api.RecordDelete("users", userUUID, WithAuthorization(token))
		t.Logf("resp: %v, error: %v", resp, err)
		assert(t, resp.StatusCode, 204)
	})
}
