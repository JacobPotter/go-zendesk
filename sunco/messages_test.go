package sunco

import (
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"testing"
)

func TestClient_PostMessage(t *testing.T) {
	t.Parallel()
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "messages.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)

	defer mockAPI.Close()

	testMessage := testhelper.MarshalMockData[Message](t, "message_body.json")

	_, err := c.PostMessage(ctx, testMessage, "123")
	if err != nil {
		t.Fatalf("Failed to post message: %v", err)
	}
}
