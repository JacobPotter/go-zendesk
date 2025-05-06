package zendesk

import (
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTicketForms(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "ticket_forms.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	ticketForms, _, err := c.GetTicketForms(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get ticket forms: %s", err)
	}

	if len(ticketForms) != 1 {
		t.Fatalf("expected length of ticket forms is , but got %d", len(ticketForms))
	}
}

func TestCreateTicketForm(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "ticket_form.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateTicketForm(ctx, TicketForm{})
	if err != nil {
		t.Fatalf("Failed to send request to create ticket form: %s", err)
	}
}

func TestDeleteTicketForm(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteTicketForm(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete ticket field: %s", err)
	}
}

func TestDeleteTicketFormFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteTicketForm(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestGetTicketForm(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "ticket_form.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	f, err := c.GetTicketForm(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	expectedID := int64(47)
	if f.ID != expectedID {
		t.Fatalf("Returned ticket form does not have the expected ID %d. Ticket id is %d", expectedID, f.ID)
	}
}

func TestGetTicketFormFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.GetTicketForm(ctx, 1234)
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestUpdateTicketForm(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "ticket_form.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	f, err := c.UpdateTicketForm(ctx, 123, TicketForm{})
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	expectedID := int64(47)
	if f.ID != expectedID {
		t.Fatalf("Returned ticket form does not have the expected ID %d. Ticket id is %d", expectedID, f.ID)
	}
}

func TestUpdateTicketFormFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	_, err := c.UpdateTicketForm(ctx, 1234, TicketForm{})
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}
