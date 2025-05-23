package zendesk

import (
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTicketField_Validate(t *testing.T) {
	cases := []struct {
		testName    string
		ticketField TicketField
		shouldPass  bool
	}{
		{
			testName: "should validate ticket field",
			ticketField: TicketField{
				ID:                  0,
				URL:                 "",
				Type:                string(Text),
				Title:               "",
				RawTitle:            "",
				Description:         "",
				RawDescription:      "",
				Position:            0,
				Active:              false,
				Required:            false,
				CollapsedForAgents:  false,
				RegexpForValidation: "",
				TitleInPortal:       "",
				RawTitleInPortal:    "",
				VisibleInPortal:     false,
				EditableInPortal:    false,
				RequiredInPortal:    false,
				Tag:                 "",
				CreatedAt:           &time.Time{},
				UpdatedAt:           &time.Time{},
				SystemFieldOptions:  nil,
				CustomFieldOptions:  nil,
				SubTypeID:           0,
				Removable:           false,
				AgentDescription:    "",
			},
			shouldPass: true,
		},
		{
			testName: "should not validate ticket field when missing custom options",
			ticketField: TicketField{
				ID:                  0,
				URL:                 "",
				Type:                string(Tagger),
				Title:               "",
				RawTitle:            "",
				Description:         "",
				RawDescription:      "",
				Position:            0,
				Active:              false,
				Required:            false,
				CollapsedForAgents:  false,
				RegexpForValidation: "",
				TitleInPortal:       "",
				RawTitleInPortal:    "",
				VisibleInPortal:     false,
				EditableInPortal:    false,
				RequiredInPortal:    false,
				Tag:                 "",
				CreatedAt:           &time.Time{},
				UpdatedAt:           &time.Time{},
				SystemFieldOptions:  nil,
				CustomFieldOptions:  nil,
				SubTypeID:           0,
				Removable:           false,
				AgentDescription:    "",
			},
			shouldPass: false,
		},
		{
			testName: "should not validate ticket field with invalid field type",
			ticketField: TicketField{
				ID:                  0,
				URL:                 "",
				Type:                "blah",
				Title:               "",
				RawTitle:            "",
				Description:         "",
				RawDescription:      "",
				Position:            0,
				Active:              false,
				Required:            false,
				CollapsedForAgents:  false,
				RegexpForValidation: "",
				TitleInPortal:       "",
				RawTitleInPortal:    "",
				VisibleInPortal:     false,
				EditableInPortal:    false,
				RequiredInPortal:    false,
				Tag:                 "",
				CreatedAt:           &time.Time{},
				UpdatedAt:           &time.Time{},
				SystemFieldOptions:  nil,
				CustomFieldOptions:  nil,
				SubTypeID:           0,
				Removable:           false,
				AgentDescription:    "",
			},
			shouldPass: false,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := c.ticketField.Validate(); err != nil && c.shouldPass {
				t.Fatalf("error validating ticket field: %s", err)
			}
		})
	}
}

func TestGetTicketFields(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "ticket_fields.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	ticketFields, _, err := c.GetTicketFields(ctx)
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	if len(ticketFields) != 15 {
		t.Fatalf("expected length of ticket fields is , but got %d", len(ticketFields))
	}
}

func TestGetTicketField(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "ticket_field.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	ticketField, err := c.GetTicketField(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	expectedID := int64(360011737434)
	if ticketField.ID != expectedID {
		t.Fatalf("Returned ticket field does not have the expected ID %d. Ticket id is %d", expectedID, ticketField.ID)
	}
}

func TestCreateTicketField(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "ticket_fields.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateTicketField(ctx, TicketField{})
	if err != nil {
		t.Fatalf("Failed to send request to create ticket field: %s", err)
	}
}

func TestUpdateTicketField(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "ticket_field.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	updatedField, err := c.UpdateTicketField(ctx, int64(1234), TicketField{})
	if err != nil {
		t.Fatalf("Failed to send request to create ticket field: %s", err)
	}

	expectedID := int64(360011737434)
	if updatedField.ID != expectedID {
		t.Fatalf("Updated field %v did not have expected id %d", updatedField, expectedID)
	}
}

func TestDeleteTicketField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteTicketField(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete ticket field: %s", err)
	}
}
