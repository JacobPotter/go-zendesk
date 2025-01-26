package zendesk

import (
	"errors"
	"github.com/JacobPotter/go-zendesk/client"
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSLAPolicies(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "sla_policies.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	slaPolicies, _, err := c.GetSLAPolicies(ctx, &SLAPolicyListOptions{})
	if err != nil {
		t.Fatalf("Failed to get sla policies: %s", err)
	}

	if len(slaPolicies) != 3 {
		t.Fatalf("expected length of sla policies is , but got %d", len(slaPolicies))
	}
}

func TestGetSLAPoliciesWithNil(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "sla_policies.json")
	c := NewTestClient(mockAPI)

	_, _, err := c.GetSLAPolicies(ctx, nil)
	if err == nil {
		t.Fatal("expected an OptionsError, but no error")
	}

	var optionsError *client.OptionsError
	ok := errors.As(err, &optionsError)
	if !ok {
		t.Fatalf("unexpected error type: %v", err)
	}
}

func TestCreateSLAPolicy(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "sla_policies.json", http.StatusCreated)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	policy, err := c.CreateSLAPolicy(ctx, SLAPolicy{})
	if err != nil {
		t.Fatalf("Failed to send request to create sla policy: %s", err)
	}

	if len(policy.PolicyMetrics) == 0 {
		t.Fatal("Failed to set the policy metrics from the json response")
	}
}

func TestGetSLAPolicy(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "sla_policy.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	sla, err := c.GetSLAPolicy(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get sla policy: %s", err)
	}

	expectedID := int64(360000068060)
	if sla.ID != expectedID {
		t.Fatalf("Returned sla policy does not have the expected ID %d. Sla policy id is %d", expectedID, sla.ID)
	}
}

func TestGetSLAPolicyFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.GetSLAPolicy(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestUpdateSLAPolicy(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "sla_policies.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	sla, err := c.UpdateSLAPolicy(ctx, 123, SLAPolicy{})
	if err != nil {
		t.Fatalf("Failed to get sla policy: %s", err)
	}

	expectedID := int64(360000068060)
	if sla.ID != expectedID {
		t.Fatalf("Returned slaPolicy does not have the expected ID %d. Sla policy id is %d", expectedID, sla.ID)
	}
}

func TestUpdateSLAPolicyFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.UpdateSLAPolicy(ctx, 1234, SLAPolicy{})
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestDeleteSLAPolicy(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteSLAPolicy(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete sla policy: %s", err)
	}
}

func TestDeleteSLAPolicyFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteSLAPolicy(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}
