package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type (
	// View is struct for group membership payload
	// https://developer.zendesk.com/api-reference/ticketing/business-rules/views/
	View struct {
		URL         string      `json:"url,omitempty"`
		ID          int64       `json:"id,omitempty"`
		Title       string      `json:"title,omitempty"`
		Active      bool        `json:"active,omitempty"`
		UpdatedAt   string      `json:"updated_at,omitempty"`
		CreatedAt   string      `json:"created_at,omitempty"`
		Default     bool        `json:"default,omitempty"`
		Position    int64       `json:"position,omitempty"`
		Description string      `json:"description,omitempty"`
		Execution   interface{} `json:"execution,omitempty"`
		Conditions  Conditions  `json:"conditions,omitempty"`
		Restriction interface{} `json:"restriction,omitempty"`
		RawTitle    string      `json:"raw_title,omitempty"`
		All         []Condition `json:"all,omitempty"`
		Any         []Condition `json:"any,omitempty"`
		Output      ViewOutput  `json:"output,omitempty"`
	}

	ViewOutput struct {
		Columns    []string `json:"columns,omitempty"`
		GroupBy    string   `json:"group_by,omitempty"`
		GroupOrder string   `json:"group_order,omitempty"`
		SortBy     string   `json:"sort_by,omitempty"`
		SortOrder  string   `json:"sort_order,omitempty"`
	}

	ViewExecution struct {
		GroupBy      string                `json:"group_by,omitempty"`
		GroupOrder   string                `json:"group_order,omitempty"`
		SortBy       string                `json:"sort_by,omitempty"`
		SortOrder    string                `json:"sort_order,omitempty"`
		Group        ViewExecColumn        `json:"group,omitempty"`
		Sort         ViewExecField         `json:"sort,omitempty"`
		Columns      []ViewExecColumn      `json:"columns,omitempty"`
		Fields       []ViewExecField       `json:"fields,omitempty"`
		CustomFields []CustomFieldViewExec `json:"custom_fields,omitempty"`
	}

	ViewExecColumn struct {
		ID         string  `json:"id"`
		Title      *string `json:"title,omitempty"`
		Filterable *bool   `json:"filterable,omitempty"`
		Sortable   *bool   `json:"sortable,omitempty"`
		Type       string  `json:"type,omitempty"`
		URL        *string `json:"url,omitempty"`
		Order      string  `json:"order,omitempty"`
	}

	CustomFieldViewExec struct {
		ID         int64  `json:"id"`
		Title      string `json:"title"`
		Type       string `json:"type"`
		URL        string `json:"url"`
		Filterable bool   `json:"filterable"`
		Sortable   bool   `json:"sortable"`
	}

	ViewExecField struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		Filterable bool   `json:"filterable"`
		Sortable   bool   `json:"sortable"`
		Order      string `json:"order,omitempty"`
	}

	ViewCount struct {
		ViewID int64  `json:"view_id"`
		URL    string `json:"url"`
		Value  int64  `json:"value"`
		Pretty string `json:"pretty"`
		Fresh  bool   `json:"fresh"`
	}

	ViewColumn string

	// ViewAPI encapsulates methods on view
	ViewAPI interface {
		GetView(context.Context, int64) (View, error)
		GetViews(context.Context) ([]View, Page, error)
		CreateView(ctx context.Context, newView View) (View, error)
		UpdateView(ctx context.Context, updatedId int64, udpatedView View) (View, error)
		DeleteView(context.Context, int64) error
		GetTicketsFromView(context.Context, int64, *TicketListOptions) ([]Ticket, Page, error)
		GetCountTicketsInViews(ctx context.Context, ids []string) ([]ViewCount, error)
		GetTicketsFromViewIterator(ctx context.Context, opts *PaginationOptions) *Iterator[Ticket]
		GetTicketsFromViewOBP(ctx context.Context, opts *OBPOptions) ([]Ticket, Page, error)
		GetTicketsFromViewCBP(ctx context.Context, opts *CBPOptions) ([]Ticket, CursorPaginationMeta, error)
		GetViewsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[View]
		GetViewsOBP(ctx context.Context, opts *OBPOptions) ([]View, Page, error)
		GetViewsCBP(ctx context.Context, opts *CBPOptions) ([]View, CursorPaginationMeta, error)
	}
)

const (
	AssignedColumn          ViewColumn = "assigned"
	AssigneeColumn          ViewColumn = "assignee"
	DueDateColumn           ViewColumn = "due_date"
	GroupColumn             ViewColumn = "group"
	IDColumn                ViewColumn = "nice_id"
	UpdatedColumn           ViewColumn = "updated"
	UpdatedAssigneeColumn   ViewColumn = "updated_assignee"
	UpdatedRequesterColumn  ViewColumn = "updated_requester"
	UpdatedByTypeColumn     ViewColumn = "updated_by_type"
	OrganizationColumn      ViewColumn = "organization"
	PriorityColumn          ViewColumn = "priority"
	CreatedColumn           ViewColumn = "created"
	RequesterColumn         ViewColumn = "requester"
	LocaleIdColumn          ViewColumn = "locale_id"
	SatisfactionScoreColumn ViewColumn = "satisfaction_score"
	SolvedColumn            ViewColumn = "solved"
	StatusColumn            ViewColumn = "status"
	DescriptionColumn       ViewColumn = "description"
	SubmitterColumn         ViewColumn = "submitter"
	TicketFormColumn        ViewColumn = "ticket_form"
	TypeColumn              ViewColumn = "type"
	BrandColumn             ViewColumn = "brand"
	TicketStatusColumn      ViewColumn = "custom_status_id"
)

var ValidViewColumns = []ViewColumn{
	AssignedColumn,
	AssigneeColumn,
	DueDateColumn,
	GroupColumn,
	IDColumn,
	UpdatedColumn,
	UpdatedAssigneeColumn,
	UpdatedRequesterColumn,
	UpdatedByTypeColumn,
	OrganizationColumn,
	PriorityColumn,
	CreatedColumn,
	RequesterColumn,
	LocaleIdColumn,
	SatisfactionScoreColumn,
	SolvedColumn,
	StatusColumn,
	DescriptionColumn,
	SubmitterColumn,
	TicketFormColumn,
	TypeColumn,
	BrandColumn,
	TicketStatusColumn,
}

// GetViews gets all views
// ref: https://developer.zendesk.com/api-reference/ticketing/business-rules/views/#list-views
func (z *Client) GetViews(ctx context.Context) ([]View, Page, error) {
	var result struct {
		Views []View `json:"views"`
		Page
	}

	body, err := z.get(ctx, "/views.json")

	if err != nil {
		return []View{}, Page{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return []View{}, Page{}, err
	}

	return result.Views, result.Page, nil
}

// GetView gets a given view
// ref: https://developer.zendesk.com/api-reference/ticketing/business-rules/views/#show-view
func (z *Client) GetView(ctx context.Context, viewID int64) (View, error) {
	var result struct {
		View View `json:"view"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/views/%d.json", viewID))

	if err != nil {
		return View{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return View{}, err
	}

	return result.View, nil
}

func (z *Client) CreateView(ctx context.Context, newView View) (View, error) {
	var data, result struct {
		View View `json:"view"`
	}

	data.View = newView

	body, err := z.post(ctx, "/views.json", newView)

	if err != nil {
		return View{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return View{}, err
	}

	return result.View, nil

}
func (z *Client) UpdateView(ctx context.Context, updatedViewId int64, updatedView View) (View, error) {
	var data, result struct {
		View View `json:"view"`
	}

	data.View = updatedView
	body, err := z.put(ctx, fmt.Sprintf("/views/%d.json", updatedViewId), data)

	if err != nil {
		return View{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return View{}, err
	}

	return result.View, nil

}

func (z *Client) DeleteView(ctx context.Context, viewId int64) error {
	err := z.delete(ctx, fmt.Sprintf("/views/%d.json", viewId))

	if err != nil {
		return err
	}

	return nil
}

// GetTicketsFromView gets the tickets of the specified view
// ref: https://developer.zendesk.com/api-reference/ticketing/business-rules/views/#list-tickets-from-a-view
func (z *Client) GetTicketsFromView(ctx context.Context, viewID int64, opts *TicketListOptions,
) ([]Ticket, Page, error) {
	var result struct {
		Tickets []Ticket `json:"tickets"`
		Page
	}
	tmp := opts
	if tmp == nil {
		tmp = &TicketListOptions{}
	}

	path := fmt.Sprintf("/views/%d/tickets.json", viewID)
	url, err := addOptions(path, tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, url)

	if err != nil {
		return []Ticket{}, Page{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return []Ticket{}, Page{}, err
	}

	return result.Tickets, result.Page, nil
}

// GetCountTicketsInViews count tickets in views using views ids
// ref https://developer.zendesk.com/api-reference/ticketing/business-rules/views/#count-tickets-in-views
func (z *Client) GetCountTicketsInViews(ctx context.Context, ids []string) ([]ViewCount, error) {
	var result struct {
		ViewCounts []ViewCount `json:"view_counts"`
	}
	idsURLParameter := strings.Join(ids, ",")
	body, err := z.get(ctx, fmt.Sprintf("/views/count_many?ids=%s", idsURLParameter))

	if err != nil {
		return []ViewCount{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return []ViewCount{}, err
	}
	return result.ViewCounts, nil
}
