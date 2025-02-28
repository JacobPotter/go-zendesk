package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JacobPotter/go-zendesk/client"
	"slices"
	"time"
)

// TicketFieldSystemFieldOption is struct for value of `system_field_options`
type TicketFieldSystemFieldOption struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Position int64  `json:"position"`
	RawName  string `json:"raw_name"`
	URL      string `json:"url"`
	Value    string `json:"value"`
}

// TicketField is struct for ticket_field payload
type TicketField struct {
	ID                  int64                          `json:"id,omitempty"`
	URL                 string                         `json:"url,omitempty"`
	Type                string                         `json:"type"`
	Title               string                         `json:"title"`
	RawTitle            string                         `json:"raw_title,omitempty"`
	Description         string                         `json:"description,omitempty"`
	RawDescription      string                         `json:"raw_description,omitempty"`
	Position            int64                          `json:"position,omitempty"`
	Active              bool                           `json:"active"`
	Required            bool                           `json:"required"`
	CollapsedForAgents  bool                           `json:"collapsed_for_agents"`
	RegexpForValidation string                         `json:"regexp_for_validation,omitempty"`
	TitleInPortal       string                         `json:"title_in_portal,omitempty"`
	RawTitleInPortal    string                         `json:"raw_title_in_portal,omitempty"`
	VisibleInPortal     bool                           `json:"visible_in_portal"`
	EditableInPortal    bool                           `json:"editable_in_portal"`
	RequiredInPortal    bool                           `json:"required_in_portal"`
	Tag                 string                         `json:"tag,omitempty"`
	CreatedAt           *time.Time                     `json:"created_at,omitempty"`
	UpdatedAt           *time.Time                     `json:"updated_at,omitempty"`
	SystemFieldOptions  []TicketFieldSystemFieldOption `json:"system_field_options,omitempty"`
	CustomFieldOptions  []CustomFieldOption            `json:"custom_field_options,omitempty"`
	SubTypeID           int64                          `json:"sub_type_id,omitempty"`
	Removable           bool                           `json:"removable,omitempty"`
	AgentDescription    string                         `json:"agent_description,omitempty"`
}

func (f TicketField) Validate() error {
	if !slices.Contains(ValidTicketFieldsTypes, TicketFieldType(f.Type)) {
		return fmt.Errorf("ticket field type must be one of: %s", ValidTicketFieldsTypes)
	}

	if (f.Type == Multiselect.String() || f.Type == Tagger.String()) && len(f.CustomFieldOptions) < 1 {
		return fmt.Errorf("when ticket field is type tagger or multiselect, there must be at least one custom option")
	}

	return nil
}

type TicketFieldType string

func (t TicketFieldType) String() string {
	return string(t)
}

const (
	Subject           TicketFieldType = "subject"
	Description       TicketFieldType = "description"
	Status            TicketFieldType = "status"
	TicketType        TicketFieldType = "tickettype"
	Priority          TicketFieldType = "priority"
	GroupField        TicketFieldType = "group"
	Assignee          TicketFieldType = "assignee"
	CustomStatus      TicketFieldType = "custom_status"
	Text              TicketFieldType = "text"
	TextArea          TicketFieldType = "textarea"
	Checkbox          TicketFieldType = "checkbox"
	Date              TicketFieldType = "date"
	Integer           TicketFieldType = "integer"
	Decimal           TicketFieldType = "decimal"
	Regexp            TicketFieldType = "regexp"
	PartialCreditCard TicketFieldType = "partial_credit_card"
	Multiselect       TicketFieldType = "multiselect"
	Tagger            TicketFieldType = "tagger"
)

// TicketFieldsTypes is a slice containing all valid ticket field names.
type TicketFieldsTypes []TicketFieldType

func (t TicketFieldsTypes) StringSlice() []string {
	strings := make([]string, len(t))

	for i, fieldType := range t {
		strings[i] = fieldType.String()
	}

	return strings
}

var ValidTicketFieldsTypes = TicketFieldsTypes{
	Status,
	Description,
	Subject,
	TicketType,
	Priority,
	GroupField,
	Assignee,
	CustomStatus,
	Text,
	TextArea,
	Checkbox,
	Date,
	Integer,
	Decimal,
	Regexp,
	PartialCreditCard,
	Multiselect,
	Tagger,
}

// TicketFieldAPI an interface containing all the ticket field related zendesk methods
type TicketFieldAPI interface {
	GetTicketFields(ctx context.Context) ([]TicketField, Page, error)
	CreateTicketField(ctx context.Context, ticketField TicketField) (TicketField, error)
	GetTicketField(ctx context.Context, ticketID int64) (TicketField, error)
	UpdateTicketField(ctx context.Context, ticketID int64, field TicketField) (TicketField, error)
	DeleteTicketField(ctx context.Context, ticketID int64) error
	GetTicketFieldsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[TicketField]
	GetTicketFieldsOBP(ctx context.Context, opts *OBPOptions) ([]TicketField, Page, error)
	GetTicketFieldsCBP(ctx context.Context, opts *CBPOptions) ([]TicketField, client.CursorPaginationMeta, error)
}

// GetTicketFields fetches ticket field list
// ref: https://developer.zendesk.com/rest_api/docs/core/ticket_fields#list-ticket-fields
func (z *Client) GetTicketFields(ctx context.Context) ([]TicketField, Page, error) {
	var data struct {
		TicketFields []TicketField `json:"ticket_fields"`
		Page
	}

	body, err := z.Get(ctx, "/ticket_fields.json")
	if err != nil {
		return []TicketField{}, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []TicketField{}, Page{}, err
	}
	return data.TicketFields, data.Page, nil
}

// CreateTicketField creates new ticket field
// ref: https://developer.zendesk.com/rest_api/docs/core/ticket_fields#create-ticket-field
func (z *Client) CreateTicketField(ctx context.Context, ticketField TicketField) (TicketField, error) {
	var data, result struct {
		TicketField TicketField `json:"ticket_field"`
	}
	data.TicketField = ticketField

	body, err := z.Post(ctx, "/ticket_fields.json", data)
	if err != nil {
		return TicketField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TicketField{}, err
	}
	return result.TicketField, nil
}

// GetTicketField gets a specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_fields#show-ticket-field
func (z *Client) GetTicketField(ctx context.Context, ticketID int64) (TicketField, error) {
	var result struct {
		TicketField TicketField `json:"ticket_field"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/ticket_fields/%d.json", ticketID))

	if err != nil {
		return TicketField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TicketField{}, err
	}

	return result.TicketField, err
}

// UpdateTicketField updates a field with the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_fields#update-ticket-field
func (z *Client) UpdateTicketField(ctx context.Context, ticketID int64, field TicketField) (TicketField, error) {
	var result, data struct {
		TicketField TicketField `json:"ticket_field"`
	}

	data.TicketField = field

	body, err := z.Put(ctx, fmt.Sprintf("/ticket_fields/%d.json", ticketID), data)

	if err != nil {
		return TicketField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TicketField{}, err
	}

	return result.TicketField, err
}

// DeleteTicketField deletes the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_fields#delete-ticket-field
func (z *Client) DeleteTicketField(ctx context.Context, ticketID int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/ticket_fields/%d.json", ticketID))

	if err != nil {
		return err
	}

	return nil
}
