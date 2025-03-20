package zendesk

import (
	"errors"
	"github.com/JacobPotter/go-zendesk/client"
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTriggers(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "triggers.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	triggers, _, err := c.GetTriggers(ctx, &TriggerListOptions{})
	if err != nil {
		t.Fatalf("Failed to get triggers: %s", err)
	}

	if len(triggers) != 8 {
		t.Fatalf("expected length of triggers is , but got %d", len(triggers))
	}
}

func TestGetTriggersWithNil(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "triggers.json")
	c := NewTestClient(mockAPI)

	_, _, err := c.GetTriggers(ctx, nil)
	if err == nil {
		t.Fatal("expected an OptionsError, but no error")
	}

	var optionsError *client.OptionsError
	ok := errors.As(err, &optionsError)
	if !ok {
		t.Fatalf("unexpected error type: %v", err)
	}
}

func TestCreateTrigger(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "triggers.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateTrigger(ctx, Trigger{})
	if err != nil {
		t.Fatalf("Failed to send request to create trigger: %s", err)
	}
}

func TestGetTrigger(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "trigger.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	trg, err := c.GetTrigger(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get trigger: %s", err)
	}

	expectedID := int64(360056295714)
	if trg.ID != expectedID {
		t.Fatalf("Returned trigger does not have the expected ID %d. Trigger id is %d", expectedID, trg.ID)
	}
}

func TestGetTriggerFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.GetTrigger(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestUpdateTrigger(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "triggers.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	trg, err := c.UpdateTrigger(ctx, 123, Trigger{})
	if err != nil {
		t.Fatalf("Failed to get trigger: %s", err)
	}

	expectedID := int64(360056295714)
	if trg.ID != expectedID {
		t.Fatalf("Returned trigger does not have the expected ID %d. Trigger id is %d", expectedID, trg.ID)
	}
}

func TestUpdateTriggerFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.UpdateTrigger(ctx, 1234, Trigger{})
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestDeleteTrigger(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteTrigger(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete trigger: %s", err)
	}
}

func TestDeleteTriggerFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteTrigger(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}
