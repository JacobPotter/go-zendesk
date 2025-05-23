package zendesk

import (
	"context"
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateWebhook(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodPost, "webhooks.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	hook, err := c.CreateWebhook(context.Background(), Webhook{
		Authentication: &WebhookAuthentication{
			AddPosition: "header",
			Data: WebhookCredentials{
				HeaderName:  "Authentication",
				HeaderValue: "",
				Username:    "john_smith",
				Password:    "hello_123",
				Token:       "",
			},
			Type: "basic_auth",
		},
		Endpoint:      "https://example.com/status/200",
		HTTPMethod:    http.MethodGet,
		Name:          "Example Webhook",
		RequestFormat: "json",
		Status:        "active",
		Subscriptions: []string{"conditional_ticket_events"},
	})
	if err != nil {
		t.Fatalf("Failed to create webhook: %v", err)
	}

	if len(hook.Subscriptions) != 1 || hook.Authentication.AddPosition != "header" {
		t.Fatalf("Invalid response of webhook: %v", hook)
	}
}

func TestGetWebhook(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "webhook.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	hook, err := c.GetWebhook(ctx, "01EJFTSCC78X5V07NPY2MHR00M")
	if err != nil {
		t.Fatalf("Failed to get webhook: %s", err)
	}

	expectedID := "01EJFTSCC78X5V07NPY2MHR00M"
	if hook.ID != expectedID {
		t.Fatalf("Returned webhook does not have the expected ID %s. Webhook ID is %s", expectedID, hook.ID)
	}
}

func TestUpdateWebhook(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	err := c.UpdateWebhook(ctx, "01EJFTSCC78X5V07NPY2MHR00M", Webhook{})
	if err != nil {
		t.Fatalf("Failed to send request to create webhook: %s", err)
	}
}

func TestDeleteWebhook(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	err := c.DeleteWebhook(ctx, "01EJFTSCC78X5V07NPY2MHR00M")
	if err != nil {
		t.Fatalf("Failed to delete webhook: %s", err)
	}
}
