package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTargets(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "targets.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	targets, _, err := c.GetTargets(ctx)
	if err != nil {
		t.Fatalf("Failed to get targets: %s", err)
	}

	if len(targets) != 2 {
		t.Fatalf("expected length of targets is , but got %d", len(targets))
	}
}

func TestGetTarget(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "target.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	target, err := c.GetTarget(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get targets: %s", err)
	}

	expectedID := int64(360000217439)
	if target.ID != expectedID {
		t.Fatalf("Returned target does not have the expected ID %d. Ticket id is %d", expectedID, target.ID)
	}
}

func TestCreateTarget(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "target.json", http.StatusCreated)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateTarget(ctx, Target{})
	if err != nil {
		t.Fatalf("Failed to send request to create target: %s", err)
	}
}

func TestUpdateTarget(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "target.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	updatedField, err := c.UpdateTarget(ctx, int64(1234), Target{})
	if err != nil {
		t.Fatalf("Failed to send request to create target: %s", err)
	}

	expectedID := int64(360000217439)
	if updatedField.ID != expectedID {
		t.Fatalf("Updated field %v did not have expected id %d", updatedField, expectedID)
	}
}

func TestDeleteTarget(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteTarget(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete target: %s", err)
	}
}
