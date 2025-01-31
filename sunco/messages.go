package sunco

import (
	"context"
	"encoding/json"
	"fmt"
)

type MessagesAPI interface {
	ListMessages(ctx context.Context, conversationId string) (MessageResponse, error)
	PostMessage(ctx context.Context, message Message, conversationId string) (MessageResponse, error)
}

func (c *Client) ListMessages(ctx context.Context, conversationId string) (MessageResponse, error) {
	var response MessageResponse

	body, err := c.Get(ctx, fmt.Sprintf("/conversations/%s/messages", conversationId))
	if err != nil {
		return MessageResponse{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return MessageResponse{}, err
	}
	return response, nil
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
