package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JacobPotter/go-zendesk/internal/client"
	"time"
)

type TriggerCategory struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Position  int64     `json:"position"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TriggerCategoryListOptions is options for GetTriggers
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#list-triggers
type TriggerCategoryListOptions struct {
	client.CursorPagination
	Sort string `url:"sort,omitempty"`
}

// TriggerCategoryAPI an interface containing all trigger related methods
type TriggerCategoryAPI interface {
	GetTriggerCategories(ctx context.Context, opts *TriggerCategoryListOptions) ([]TriggerCategory, client.CursorPaginationMeta, error)
	CreateTriggerCategory(ctx context.Context, triggerCategory TriggerCategory) (TriggerCategory, error)
	GetTriggerCategory(ctx context.Context, id int64) (TriggerCategory, error)
	UpdateTriggerCategory(ctx context.Context, id int64, triggerCategory TriggerCategory) (TriggerCategory, error)
	DeleteTriggerCategory(ctx context.Context, id int64) error
}

// GetTriggerCategories fetch trigger category list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#getting-triggers
func (z *Client) GetTriggerCategories(ctx context.Context, opts *TriggerCategoryListOptions) ([]TriggerCategory, client.CursorPaginationMeta, error) {
	var data struct {
		TriggerCategories []TriggerCategory `json:"trigger_categories"`
		client.CursorPaginationMeta
	}

	if opts == nil {
		return []TriggerCategory{}, client.CursorPaginationMeta{}, &client.OptionsError{Opts: opts}
	}

	u, err := client.AddOptions("/trigger_categories.json", opts)
	if err != nil {
		return []TriggerCategory{}, client.CursorPaginationMeta{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return []TriggerCategory{}, client.CursorPaginationMeta{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []TriggerCategory{}, client.CursorPaginationMeta{}, err
	}
	return data.TriggerCategories, data.CursorPaginationMeta, nil
}

// CreateTriggerCategory creates new trigger category
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#create-trigger
func (z *Client) CreateTriggerCategory(ctx context.Context, triggerCategory TriggerCategory) (TriggerCategory, error) {
	var data, result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}
	data.TriggerCategory = triggerCategory

	body, err := z.Post(ctx, "/trigger_categories.json", data)
	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}
	return result.TriggerCategory, nil
}

// GetTriggerCategory returns the specified trigger category
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#getting-triggers
func (z *Client) GetTriggerCategory(ctx context.Context, id int64) (TriggerCategory, error) {
	var result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/trigger_categories/%d.json", id))
	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}
	return result.TriggerCategory, nil
}

// UpdateTriggerCategory updates the specified trigger category and returns the updated one
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#update-trigger
func (z *Client) UpdateTriggerCategory(ctx context.Context, id int64, triggerCategory TriggerCategory) (TriggerCategory, error) {
	var data, result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}

	data.TriggerCategory = triggerCategory
	body, err := z.Put(ctx, fmt.Sprintf("/trigger_categories/%d.json", id), data)
	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}

	return result.TriggerCategory, nil
}

// DeleteTriggerCategory deletes the specified trigger category
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#delete-trigger
func (z *Client) DeleteTriggerCategory(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/trigger_categories/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
