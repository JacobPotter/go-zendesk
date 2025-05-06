package zendesk

import (
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetGroupsIterator(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "groups.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	ops := NewPaginationOptions()
	ops.PageSize = 10

	it := c.GetGroupsIterator(ctx, ops)

	expectedLength := 1
	groupsCount := 0
	for it.HasMore() {
		groups, err := it.GetNext()
		if len(groups) != expectedLength {
			t.Fatalf("expected length of groups is 1, but got %d", len(groups))
		}
		groupsCount += len(groups)
		if err != nil {
			t.Fatalf("Failed to get groups: %s", err)
		}
	}
	if groupsCount != 1 {
		t.Fatalf("expected length of groups is 1, but got %d", groupsCount)
	}
}

func TestGetGroups(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "groups.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	groups, _, err := c.GetGroups(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get groups: %s", err)
	}

	if len(groups) != 1 {
		t.Fatalf("expected length of groups is 1, but got %d", len(groups))
	}
}

func TestGetGroupsOBP(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "groups.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	groups, _, err := c.GetGroupsOBP(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get groups: %s", err)
	}

	if len(groups) != 1 {
		t.Fatalf("expected length of groups is 1, but got %d", len(groups))
	}
}

func TestGetGroupsCBP(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "groups.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	groups, _, err := c.GetGroupsCBP(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get groups: %s", err)
	}

	if len(groups) != 1 {
		t.Fatalf("expected length of groups is 1, but got %d", len(groups))
	}
}

func TestCreateGroup(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "groups.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateGroup(ctx, Group{})
	if err != nil {
		t.Fatalf("Failed to send request to create group: %s", err)
	}
}

func TestGetGroup(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "group.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	group, err := c.GetGroup(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get group: %s", err)
	}

	expectedID := int64(360002440594)
	if group.ID != expectedID {
		t.Fatalf("Returned group does not have the expected ID %d. Group ID is %d", expectedID, group.ID)
	}
}

func TestUpdateGroup(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "groups.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	updatedGroup, err := c.UpdateGroup(ctx, int64(1234), Group{})
	if err != nil {
		t.Fatalf("Failed to send request to create group: %s", err)
	}

	expectedID := int64(360002440594)
	if updatedGroup.ID != expectedID {
		t.Fatalf("Updated group %v did not have expected id %d", updatedGroup, expectedID)
	}
}

func TestDeleteGroup(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteGroup(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete group: %s", err)
	}
}
