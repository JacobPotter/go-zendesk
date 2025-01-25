package sunco

import (
	"context"
	"encoding/json"
)

type ConversationsAPI interface {
	CreateConversation(ctx context.Context, conversation Conversation) (Conversation, error)
}

func (c *Client) CreateConversation(ctx context.Context, conversation Conversation) (Conversation, error) {
	var response ConversationResponse

	body, err := c.Post(ctx, "/conversations", conversation)
	if err != nil {
		return Conversation{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Conversation{}, err
	}

	return response.Conversation, nil
}
