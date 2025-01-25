package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JacobPotter/go-zendesk/internal/client"
	"slices"
	"time"
)

// TicketForm ticket forms allow an admin to define a subset of ticket fields for display to both agents and end users.
// Accounts are limited to 300 ticket forms.
//
// TicketFieldIds is used to display all ticket fields which are in this ticket form.
// The products use the order of the ids to show the field values in the tickets.
//
// AgentConditions and EndUserConditions specify which ticket fields to show on the form based on a set of defined conditions
type TicketForm struct {
	Active             bool                     `json:"active"`
	AgentConditions    []ConditionalTicketField `json:"agent_conditions,omitempty"`
	CreatedAt          time.Time                `json:"created_at"`
	Default            bool                     `json:"default,omitempty"`
	DisplayName        string                   `json:"display_name"`
	EndUserConditions  []ConditionalTicketField `json:"end_user_conditions,omitempty"`
	EndUserVisible     bool                     `json:"end_user_visible,omitempty"`
	ID                 int64                    `json:"id"`
	InAllBrands        bool                     `json:"in_all_brands,omitempty"`
	Name               string                   `json:"name"`
	Position           int64                    `json:"position,omitempty"`
	RawDisplayName     string                   `json:"raw_display_name,omitempty"`
	RawName            string                   `json:"raw_name,omitempty"`
	RestrictedBrandIds []int64                  `json:"restricted_brand_ids,omitempty"`
	TicketFieldIds     []int64                  `json:"ticket_field_ids"`
	UpdatedAt          time.Time                `json:"updated_at"`
	Url                string                   `json:"url"`
}

// ConditionalTicketField condition which to display fields ParentFieldId is the ticket field the condition is for,
// where the condition is matching the value in Value.
// The Value will either be the tag value of a field or a case-sensitive match of a text field.
// Then, ChildFields is the set of fields to show when the condition is met.
// There is then a set of statuses to be required on, if it is required
type ConditionalTicketField struct {
	ParentFieldId int64        `json:"parent_field_id"`
	Value         string       `json:"value"`
	ChildFields   []ChildField `json:"child_fields"`
}

// ChildField is used to define when to show a field based on a condition from a parent field in a form.
// Has attributes IsRequired and RequiredOnStatuses, which indicate if the field to show is conditional
type ChildField struct {
	Id                 int64              `json:"id"`
	IsRequired         bool               `json:"is_required"`
	RequiredOnStatuses RequiredOnStatuses `json:"required_on_statuses"`
}

// RequiredOnStatuses is an object that defines how status requires a child field on a form.
// Valid types in Type attribute are ALL_STATUSES, NO_STATUSES, and SOME_STATUSES.
// Statuses will enclose the statuses the field is required when Type is SOME_STATUSES.
// Otherwise, it is empty
type RequiredOnStatuses struct {
	Statuses []string        `json:"statuses"`
	Type     RequirementType `json:"type"`
}

// Validate validates attributes in RequiredOnStatuses struct. Returns an error if validation fails
func (r RequiredOnStatuses) Validate() error {
	switch r.Type {
	case AllStatuses, NoStatuses:
		if len(r.Statuses) > 0 {
			return fmt.Errorf("error: no statuses when type is ALL_STATUSES or NO_STATUSES")
		}
	case SomeStatuses:
		if len(r.Statuses) == 0 {
			return fmt.Errorf("error: statuses required when type is SOME_STATUSES")
		}
		for _, status := range r.Statuses {
			if !slices.Contains(ValidRequirementStatuses, status) {
				return fmt.Errorf("error: status '%s' is not valid", status)
			}
		}
	}
	return nil
}

// ValidRequirementStatuses is an array of strings denoting which statuses are available to be used in RequiredOnStatuses struct
var ValidRequirementStatuses = []string{"new", "open", "pending", "hold", "solved"}

type RequirementType string

const (
	SomeStatuses RequirementType = "SOME_STATUSES"
	NoStatuses   RequirementType = "NO_STATUSES"
	AllStatuses  RequirementType = "ALL_STATUSES"
)

// ValidRequirementTypes denotes which status requirement types are valid
var ValidRequirementTypes = []RequirementType{
	SomeStatuses, AllStatuses, NoStatuses,
}

// Validate validates status requirement types
func (r RequirementType) Validate() error {
	if !slices.Contains(ValidRequirementTypes, r) {
		return fmt.Errorf("requirement type '%s' is not valid", r)
	}
	return nil
}

// TicketFormListOptions is options for GetTicketForms
//
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#available-parameters
type TicketFormListOptions struct {
	PageOptions
	Active            bool `url:"active"`
	EndUserVisible    bool `url:"end_user_visible,omitempty"`
	FallbackToDefault bool `url:"fallback_to_default,omitempty"`
	AssociatedToBrand bool `url:"associated_to_brand,omitempty"`
}

// TicketFormAPI an interface containing all ticket form related methods
type TicketFormAPI interface {
	GetTicketForms(ctx context.Context, options *TicketFormListOptions) ([]TicketForm, Page, error)
	CreateTicketForm(ctx context.Context, ticketForm TicketForm) (TicketForm, error)
	DeleteTicketForm(ctx context.Context, id int64) error
	UpdateTicketForm(ctx context.Context, id int64, form TicketForm) (TicketForm, error)
	GetTicketForm(ctx context.Context, id int64) (TicketForm, error)
	GetTicketFormsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[TicketForm]
	GetTicketFormsOBP(ctx context.Context, opts *OBPOptions) ([]TicketForm, Page, error)
	GetTicketFormsCBP(ctx context.Context, opts *CBPOptions) ([]TicketForm, client.CursorPaginationMeta, error)
}

// GetTicketForms fetches ticket forms
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#list-ticket-forms
func (z *Client) GetTicketForms(ctx context.Context, options *TicketFormListOptions) ([]TicketForm, Page, error) {
	var data struct {
		TicketForms []TicketForm `json:"ticket_forms"`
		Page
	}

	tmp := options
	if tmp == nil {
		tmp = &TicketFormListOptions{}
	}

	u, err := client.AddOptions("/ticket_forms.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return []TicketForm{}, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []TicketForm{}, Page{}, err
	}
	return data.TicketForms, data.Page, nil
}

// CreateTicketForm creates new ticket form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#create-ticket-forms
func (z *Client) CreateTicketForm(ctx context.Context, ticketForm TicketForm) (TicketForm, error) {
	var data, result struct {
		TicketForm TicketForm `json:"ticket_form"`
	}
	data.TicketForm = ticketForm

	body, err := z.Post(ctx, "/ticket_forms.json", data)
	if err != nil {
		return TicketForm{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TicketForm{}, err
	}
	return result.TicketForm, nil
}

// GetTicketForm returns the specified ticket form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#show-ticket-form
func (z *Client) GetTicketForm(ctx context.Context, id int64) (TicketForm, error) {
	var result struct {
		TicketForm TicketForm `json:"ticket_form"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/ticket_forms/%d.json", id))
	if err != nil {
		return TicketForm{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TicketForm{}, err
	}
	return result.TicketForm, nil
}

// UpdateTicketForm updates the specified ticket form and returns the updated form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#update-ticket-forms
func (z *Client) UpdateTicketForm(ctx context.Context, id int64, form TicketForm) (TicketForm, error) {
	var data, result struct {
		TicketForm TicketForm `json:"ticket_form"`
	}

	data.TicketForm = form
	body, err := z.Put(ctx, fmt.Sprintf("/ticket_forms/%d.json", id), data)
	if err != nil {
		return TicketForm{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TicketForm{}, err
	}

	return result.TicketForm, nil
}

// DeleteTicketForm deletes the specified ticket form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#delete-ticket-form
func (z *Client) DeleteTicketForm(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/ticket_forms/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
