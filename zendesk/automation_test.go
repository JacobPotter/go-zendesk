package zendesk

import (
	"errors"
	"github.com/JacobPotter/go-zendesk/client"
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAutomations(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "automations.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	automations, _, err := c.GetAutomations(ctx, &AutomationListOptions{})
	if err != nil {
		t.Fatalf("Failed to get automations: %s", err)
	}

	if len(automations) != 3 {
		t.Fatalf("expected length of automations is , but got %d", len(automations))
	}
}

func TestGetAutomationsWithNil(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "automations.json")
	c := NewTestClient(mockAPI)

	_, _, err := c.GetAutomations(ctx, nil)
	if err == nil {
		t.Fatal("expected an OptionsError, but no error")
	}

	var optionsError *client.OptionsError
	ok := errors.As(err, &optionsError)
	if !ok {
		t.Fatalf("unexpected error type: %v", err)
	}
}

func TestCreateAutomation(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "automations.json", http.StatusCreated)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateAutomation(ctx, Automation{})
	if err != nil {
		t.Fatalf("Failed to send request to create automation: %s", err)
	}
}

func TestGetAutomation(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "automation.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	trg, err := c.GetAutomation(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get automation: %s", err)
	}

	expectedID := int64(360017421099)
	if trg.ID != expectedID {
		t.Fatalf("Returned automation does not have the expected ID %d. Automation id is %d", expectedID, trg.ID)
	}
}

func TestGetAutomationFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.GetAutomation(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestUpdateAutomation(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "automations.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	trg, err := c.UpdateAutomation(ctx, 123, Automation{})
	if err != nil {
		t.Fatalf("Failed to get automation: %s", err)
	}

	expectedID := int64(360017421099)
	if trg.ID != expectedID {
		t.Fatalf("Returned automation does not have the expected ID %d. Automation id is %d", expectedID, trg.ID)
	}
}

func TestUpdateAutomationFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.UpdateAutomation(ctx, 1234, Automation{})
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestDeleteAutomation(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteAutomation(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete automation: %s", err)
	}
}

func TestDeleteAutomationFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteAutomation(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}
