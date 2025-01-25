package sunco

import (
	"context"
	"encoding/json"
	"fmt"
)

type MessagesAPI interface {
	PostMessage(ctx context.Context, message Message, conversationId string) (MessageResponse, error)
}

func (c *Client) PostMessage(ctx context.Context, message Message, conversationId string) (MessageResponse, error) {
	var response MessageResponse

	body, err := c.Post(ctx, fmt.Sprintf("/conversations/%s/messages", conversationId), message)
	if err != nil {
		return MessageResponse{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return MessageResponse{}, err
	}
	return response, nil
}
