package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JacobPotter/go-zendesk/internal/client"
	"time"
)

// Trigger is zendesk trigger JSON payload format
//
// ref: https://developer.zendesk.com/rest_api/docs/core/triggers#json-format
type Trigger struct {
	ID          int64      `json:"id,omitempty"`
	Title       string     `json:"title"`
	Active      bool       `json:"active"`
	Position    int64      `json:"position,omitempty"`
	Conditions  Conditions `json:"conditions"`
	Actions     []Action   `json:"actions"`
	Description string     `json:"description,omitempty"`
	CategoryID  string     `json:"category_id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	URL         string     `json:"url,omitempty"`
}

// TriggerListOptions is options for GetTriggers
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#list-triggers
type TriggerListOptions struct {
	PageOptions
	Active     bool   `url:"active"`
	CategoryID string `url:"category_id,omitempty"`
	SortBy     string `url:"sort_by,omitempty"`
	SortOrder  string `url:"sort_order,omitempty"`
}

// TriggerAPI an interface containing all trigger related methods
type TriggerAPI interface {
	GetTriggers(ctx context.Context, opts *TriggerListOptions) ([]Trigger, Page, error)
	CreateTrigger(ctx context.Context, trigger Trigger) (Trigger, error)
	GetTrigger(ctx context.Context, id int64) (Trigger, error)
	UpdateTrigger(ctx context.Context, id int64, trigger Trigger) (Trigger, error)
	DeleteTrigger(ctx context.Context, id int64) error
	GetTriggersIterator(ctx context.Context, opts *PaginationOptions) *Iterator[Trigger]
	GetTriggersOBP(ctx context.Context, opts *OBPOptions) ([]Trigger, Page, error)
	GetTriggersCBP(ctx context.Context, opts *CBPOptions) ([]Trigger, client.CursorPaginationMeta, error)
}

// GetTriggers fetch trigger list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#getting-triggers
func (z *Client) GetTriggers(ctx context.Context, opts *TriggerListOptions) ([]Trigger, Page, error) {
	var data struct {
		Triggers []Trigger `json:"triggers"`
		Page
	}

	if opts == nil {
		return []Trigger{}, Page{}, &client.OptionsError{Opts: opts}
	}

	u, err := client.AddOptions("/triggers.json", opts)
	if err != nil {
		return []Trigger{}, Page{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return []Trigger{}, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []Trigger{}, Page{}, err
	}
	return data.Triggers, data.Page, nil
}

// CreateTrigger creates new trigger
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#create-trigger
func (z *Client) CreateTrigger(ctx context.Context, trigger Trigger) (Trigger, error) {
	var data, result struct {
		Trigger Trigger `json:"trigger"`
	}
	data.Trigger = trigger

	body, err := z.Post(ctx, "/triggers.json", data)
	if err != nil {
		return Trigger{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Trigger{}, err
	}
	return result.Trigger, nil
}

// GetTrigger returns the specified trigger
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#getting-triggers
func (z *Client) GetTrigger(ctx context.Context, id int64) (Trigger, error) {
	var result struct {
		Trigger Trigger `json:"trigger"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/triggers/%d.json", id))
	if err != nil {
		return Trigger{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Trigger{}, err
	}
	return result.Trigger, nil
}

// UpdateTrigger updates the specified trigger and returns the updated one
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#update-trigger
func (z *Client) UpdateTrigger(ctx context.Context, id int64, trigger Trigger) (Trigger, error) {
	var data, result struct {
		Trigger Trigger `json:"trigger"`
	}

	data.Trigger = trigger
	body, err := z.Put(ctx, fmt.Sprintf("/triggers/%d.json", id), data)
	if err != nil {
		return Trigger{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Trigger{}, err
	}

	return result.Trigger, nil
}

// DeleteTrigger deletes the specified trigger
//
// ref: https://developer.zendesk.com/rest_api/docs/support/triggers#delete-trigger
func (z *Client) DeleteTrigger(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/triggers/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
