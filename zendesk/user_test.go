package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestUserRoleText(t *testing.T) {
	for key := UserRoleEndUser; key <= UserRoleAdmin; key++ {
		if text := UserRoleText(key); text == "" {
			t.Fatalf("key=%d is undefined", key)
		}
	}
}

func TestGetUsers(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "users.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	users, _, err := c.GetUsers(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get users: %s", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected length of users is 2, but got %d", len(users))
	}
}

func TestGetOrganizationUsers(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "users.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	users, _, err := c.GetOrganizationUsers(ctx, 1000006909040, nil)
	if err != nil {
		t.Fatalf("Failed to get users: %s", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected length of users is 2, but got %d", len(users))
	}
}

func TestGetManyUsers(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "users.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	users, _, err := c.GetManyUsers(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get many users: %s", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected length of many users is 2, but got %d", len(users))
	}
}

func TestSearchUsers(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "users.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	users, _, err := c.SearchUsers(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get many users: %s", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected length of many users is 2, but got %d", len(users))
	}
}

func TestGetUser(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodGet, "user.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	user, err := c.GetUser(ctx, 369531345753)
	if err != nil {
		t.Fatalf("Failed to get user: %s", err)
	}

	expectedID := int64(369531345753)
	if user.ID != expectedID {
		t.Fatalf("Returned user does not have the expected ID %d. User id is %d", expectedID, user.ID)
	}
}

func TestGetUserFailure(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodGet, "user.json", http.StatusInternalServerError)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.GetUser(ctx, 369531345753)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestGetUsersRolesEncodeCorrectly(t *testing.T) {
	expected := "role%5B%5D=admin&role%5B%5D=end-user"
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryString := r.URL.Query().Encode()
		if queryString != expected {
			t.Fatalf(`Did not get the expect query string: "%s". Was: "%s"`, expected, queryString)
		}
		_, err := w.Write(testhelper.ReadFixture(t, filepath.Join(http.MethodGet, "users.json")))
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	opts := UserListOptions{
		Roles: []string{
			"admin",
			"end-user",
		},
	}

	_, _, err := c.GetUsers(ctx, &opts)
	if err != nil {
		t.Fatalf("Failed to get users: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "users.json", http.StatusCreated)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	user, err := c.CreateUser(ctx, User{
		Email: "test@example.com",
		Name:  "testuser",
	})
	if err != nil {
		t.Fatalf("Failed to get valid response: %s", err)
	}
	if user.ID == 0 {
		t.Fatal("Failed to create user")
	}
}

func TestCreateOrUpdateUserCreated(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "users.json", http.StatusCreated)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	user, err := c.CreateOrUpdateUser(ctx, User{
		Email: "test@example.com",
		Name:  "testuser",
	})
	if err != nil {
		t.Fatalf("Failed to get valid response: %s", err)
	}
	if user.ID == 0 {
		t.Fatal("Failed to create or update user")
	}
}

func TestCreateOrUpdateUserUpdated(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "users.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	user, err := c.CreateOrUpdateUser(ctx, User{
		Email: "test@example.com",
		Name:  "testuser",
	})
	if err != nil {
		t.Fatalf("Failed to get valid response: %s", err)
	}
	if user.ID == 0 {
		t.Fatal("Failed to create or update user")
	}
}

func TestCreateOrUpdateUserFailure(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "users.json", http.StatusInternalServerError)

	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateOrUpdateUser(ctx, User{})
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestUpdateUser(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "user.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	user, err := c.UpdateUser(ctx, 369531345753, User{})
	if err != nil {
		t.Fatalf("Failed to update user: %s", err)
	}

	expectedID := int64(369531345753)
	if user.ID != expectedID {
		t.Fatalf("Returned user does not have the expected ID %d. User id is %d", expectedID, user.ID)
	}
}

func TestUpdateUserFailure(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "user.json", http.StatusInternalServerError)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.UpdateUser(ctx, 369531345753, User{})
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestGetUserRelated(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodGet, "user_related.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	userRelated, err := c.GetUserRelated(ctx, 369531345753)
	if err != nil {
		t.Fatalf("Failed to get user related information: %s", err)
	}

	expectedAssignedTickets := int64(5)
	if userRelated.AssignedTickets != expectedAssignedTickets {
		t.Fatalf("Returned user does not have the expected assigned tickets %d. It is %d", expectedAssignedTickets, userRelated.AssignedTickets)
	}
}
