package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	client2 "github.com/JacobPotter/go-zendesk/client"
	"time"
)

// Automation is zendesk automation JSON payload format
//
// ref: https://developer.zendesk.com/rest_api/docs/core/automations#json-format
type Automation struct {
	ID          int64      `json:"id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Active      bool       `json:"active"`
	Position    int64      `json:"position,omitempty"`
	Conditions  Conditions `json:"conditions"`
	Actions     []Action   `json:"actions"`
	RawTitle    string     `json:"raw_title,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	URL         string     `json:"url,omitempty"`
}

// AutomationListOptions is options for GetAutomations
//
// ref: https://developer.zendesk.com/rest_api/docs/support/automations#list-automations
type AutomationListOptions struct {
	PageOptions
	Active    bool   `url:"active"`
	SortBy    string `url:"sort_by,omitempty"`
	SortOrder string `url:"sort_order,omitempty"`
}

// AutomationAPI an interface containing all automation related methods
type AutomationAPI interface {
	GetAutomations(ctx context.Context, opts *AutomationListOptions) ([]Automation, Page, error)
	CreateAutomation(ctx context.Context, automation Automation) (Automation, error)
	GetAutomation(ctx context.Context, id int64) (Automation, error)
	UpdateAutomation(ctx context.Context, id int64, automation Automation) (Automation, error)
	DeleteAutomation(ctx context.Context, id int64) error
	GetAutomationsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[Automation]
	GetAutomationsOBP(ctx context.Context, opts *OBPOptions) ([]Automation, Page, error)
	GetAutomationsCBP(ctx context.Context, opts *CBPOptions) ([]Automation, client2.CursorPaginationMeta, error)
}

// GetAutomations fetch automation list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/automations#getting-automations
func (z *Client) GetAutomations(ctx context.Context, opts *AutomationListOptions) ([]Automation, Page, error) {
	var data struct {
		Automations []Automation `json:"automations"`
		Page
	}

	if opts == nil {
		return []Automation{}, Page{}, &client2.OptionsError{Opts: opts}
	}

	u, err := client2.AddOptions("/automations.json", opts)
	if err != nil {
		return []Automation{}, Page{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return []Automation{}, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []Automation{}, Page{}, err
	}

	return data.Automations, data.Page, nil
}

// CreateAutomation creates new automation
//
// ref: https://developer.zendesk.com/rest_api/docs/support/automations#create-automation
func (z *Client) CreateAutomation(ctx context.Context, automation Automation) (Automation, error) {
	var data, result struct {
		Automation Automation `json:"automation"`
	}

	data.Automation = automation
	body, err := z.Post(ctx, "/automations.json", data)

	if err != nil {
		return Automation{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Automation{}, err
	}

	return result.Automation, nil
}

// GetAutomation returns the specified automation
//
// ref: https://developer.zendesk.com/rest_api/docs/support/automations#getting-automations
func (z *Client) GetAutomation(ctx context.Context, id int64) (Automation, error) {
	var result struct {
		Automation Automation `json:"automation"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/automations/%d.json", id))
	if err != nil {
		return Automation{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Automation{}, err
	}

	return result.Automation, nil
}

// UpdateAutomation updates the specified automation and returns the updated one
//
// ref: https://developer.zendesk.com/rest_api/docs/support/automations#update-automation
func (z *Client) UpdateAutomation(ctx context.Context, id int64, automation Automation) (Automation, error) {
	var data, result struct {
		Automation Automation `json:"automation"`
	}

	data.Automation = automation
	body, err := z.Put(ctx, fmt.Sprintf("/automations/%d.json", id), data)

	if err != nil {
		return Automation{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Automation{}, err
	}

	return result.Automation, nil
}

// DeleteAutomation deletes the specified automation
//
// ref: https://developer.zendesk.com/rest_api/docs/support/automations#delete-automation
func (z *Client) DeleteAutomation(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/automations/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
