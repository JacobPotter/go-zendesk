package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type TriggerCategory struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Position  int64     `json:"position"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TriggerListOptions is options for GetTriggers
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#list-triggers
type TriggerCategoryListOptions struct {
	CursorPagination
	Sort string `url:"sort,omitempty"`
}

// TriggerAPI an interface containing all trigger related methods
type TriggerCategoryAPI interface {
	GetTriggerCategories(ctx context.Context, opts *TriggerCategoryListOptions) ([]TriggerCategory, CursorPaginationMeta, error)
	CreateTriggerCategory(ctx context.Context, triggerCategory TriggerCategory) (TriggerCategory, error)
	GetTriggerCategory(ctx context.Context, id int64) (TriggerCategory, error)
	UpdateTriggerCategory(ctx context.Context, id int64, triggerCategory TriggerCategory) (TriggerCategory, error)
	DeleteTriggerCategory(ctx context.Context, id int64) error
}

// GetTriggerCategories fetch trigger category list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#getting-triggers
func (z *Client) GetTriggerCategories(ctx context.Context, opts *TriggerCategoryListOptions) ([]TriggerCategory, CursorPaginationMeta, error) {
	var data struct {
		TriggerCategories []TriggerCategory `json:"trigger_categories"`
		CursorPaginationMeta
	}

	if opts == nil {
		return []TriggerCategory{}, CursorPaginationMeta{}, &OptionsError{opts}
	}

	u, err := addOptions("/trigger_categories.json", opts)
	if err != nil {
		return []TriggerCategory{}, CursorPaginationMeta{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return []TriggerCategory{}, CursorPaginationMeta{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []TriggerCategory{}, CursorPaginationMeta{}, err
	}
	return data.TriggerCategories, data.CursorPaginationMeta, nil
}

// CreateTriggerCategories creates new trigger category
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#create-trigger
func (z *Client) CreateTriggerCategory(ctx context.Context, triggerCategory TriggerCategory) (TriggerCategory, error) {
	var data, result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}
	data.TriggerCategory = triggerCategory

	body, err := z.post(ctx, "/trigger_categories.json", data)
	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}
	return result.TriggerCategory, nil
}

// GetTrigger returns the specified trigger category
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#getting-triggers
func (z *Client) GetTriggerCategory(ctx context.Context, id int64) (TriggerCategory, error) {
	var result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/trigger_categories/%d.json", id))
	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}
	return result.TriggerCategory, nil
}

// UpdateTrigger updates the specified trigger category and returns the updated one
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#update-trigger
func (z *Client) UpdateTriggerCategory(ctx context.Context, id int64, triggerCategory TriggerCategory) (TriggerCategory, error) {
	var data, result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}

	data.TriggerCategory = triggerCategory
	body, err := z.put(ctx, fmt.Sprintf("/trigger_categories/%d.json", id), data)
	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}

	return result.TriggerCategory, nil
}

// DeleteTrigger deletes the specified trigger category
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#delete-trigger
func (z *Client) DeleteTriggerCategory(ctx context.Context, id int64) error {
	err := z.delete(ctx, fmt.Sprintf("/trigger_categories/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
