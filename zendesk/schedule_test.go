package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"testing"
)

func TestClient_CreateSchedule(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodPost, "schedule.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateSchedule(ctx, Schedule{})
	if err != nil {
		t.Fatalf("Error creating schedule: %v", err)
	}
}

func TestClient_GetSchedule(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "schedule.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.GetSchedule(ctx, 123)
	if err != nil {
		t.Fatalf("Error getting schedule: %v", err)
	}
}

func TestClient_UpdateSchedule(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodPut, "schedule.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.UpdateSchedule(ctx, 123, Schedule{})
	if err != nil {
		t.Fatalf("Error updating schedule: %v", err)
	}
}
