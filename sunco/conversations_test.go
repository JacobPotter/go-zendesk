package sunco

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"testing"
)

func TestClient_CreateConversation(t *testing.T) {
	t.Parallel()
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "conversations.json", http.StatusOK)
	c := NewTestClient(mockAPI)

	defer mockAPI.Close()

	testConversation := testhelper.MarshalMockData[Conversation](t, "conversation_body.json")

	_, err := c.CreateConversation(ctx, testConversation)
	if err != nil {
		t.Fatalf("Failed to create conversation: %s", err)
	}
}
