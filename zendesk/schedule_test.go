package zendesk

import (
	"net/http"
	"testing"
)

func TestClient_CreateSchedule(t *testing.T) {
	mockAPI := newMockAPI(http.MethodPost, "schedule.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := client.CreateSchedule(ctx, Schedule{})
	if err != nil {
		t.Fatalf("Error creating schedule: %v", err)
	}
}

func TestClient_GetSchedule(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "schedule.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := client.GetSchedule(ctx, 123)
	if err != nil {
		t.Fatalf("Error getting schedule: %v", err)
	}
}

func TestClient_UpdateSchedule(t *testing.T) {
	mockAPI := newMockAPI(http.MethodPut, "schedule.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := client.UpdateSchedule(ctx, 123, Schedule{})
	if err != nil {
		t.Fatalf("Error updating schedule: %v", err)
	}
}
