package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	client2 "github.com/JacobPotter/go-zendesk/client"
	"time"
)

// SLA Policy metric values
//
// ref: https://developer.zendesk.com/rest_api/docs/support/sla_policies#metrics
const (
	AgentWorkTimeMetric      = "agent_work_time"
	FirstReplyTimeMetric     = "first_reply_time"
	NextReplyTimeMetric      = "next_reply_time"
	PausableUpdateTimeMetric = "pausable_update_time"
	PeriodicUpdateTimeMetric = "periodic_update_time"
	RequesterWaitTimeMetric  = "requester_wait_time"
)

type SLAPolicyMetric struct {
	Priority      string `json:"priority"`
	Metric        string `json:"metric"`
	Target        int    `json:"target"`
	BusinessHours bool   `json:"business_hours"`
}

type FirstReplyTime struct {
	ActivateOnTicketCreatedForEndUser                      bool `json:"activate_on_ticket_created_for_end_user,omitempty"`
	ActivateOnAgentTicketCreatedForEndUserWithInternalNote bool `json:"activate_on_agent_ticket_created_for_end_user_with_internal_note,omitempty"`
	ActivateOnLightAgentOnEmailForwardTicketFromEndUser    bool `json:"activate_on_light_agent_on_email_forward_ticket_from_end_user,omitempty"`
	ActivateOnAgentCreatedTicketForSelf                    bool `json:"activate_on_agent_created_ticket_for_self,omitempty"`
	FulfillOnAgentInternalNote                             bool `json:"fulfill_on_agent_internal_note,omitempty"`
}

type NextReplyTime struct {
	FulfillOnNonRequestingAgentInternalNoteAfterActivation        bool `json:"fulfill_on_non_requesting_agent_internal_note_after_activation,omitempty"`
	ActivateOnEndUserAddedInternalNote                            bool `json:"activate_on_end_user_added_internal_note,omitempty"`
	ActivateOnAgentRequestedTicketWithPublicCommentOrInternalNote bool `json:"activate_on_agent_requested_ticket_with_public_comment_or_internal_note,omitempty"`
	ActivateOnLightAgentInternalNote                              bool `json:"activate_on_light_agent_internal_note,omitempty"`
}

type PeriodicUpdateTime struct {
	ActivateOnAgentInternalNote bool `json:"activate_on_agent_internal_note,omitempty"`
}

type MetricSettings struct {
	FirstReplyTime     FirstReplyTime     `json:"first_reply_time,omitempty"`
	NextReplyTime      NextReplyTime      `json:"next_reply_time,omitempty"`
	PeriodicUpdateTime PeriodicUpdateTime `json:"periodic_update_time,omitempty"`
}

// SLAPolicy is zendesk slaPolicy JSON payload format
//
// ref: https://developer.zendesk.com/rest_api/docs/core/slas/policies#json-format
type SLAPolicy struct {
	ID             int64             `json:"id,omitempty"`
	Title          string            `json:"title"`
	Description    string            `json:"description,omitempty"`
	Position       int64             `json:"position,omitempty"`
	Active         bool              `json:"active"`
	Filter         Conditions        `json:"filter"`
	PolicyMetrics  []SLAPolicyMetric `json:"policy_metrics,omitempty"`
	MetricSettings MetricSettings    `json:"metric_settings,omitempty"`
	CreatedAt      *time.Time        `json:"created_at,omitempty"`
	UpdatedAt      *time.Time        `json:"updated_at,omitempty"`
}

// SLAPolicyListOptions is options for GetSLAPolicies
//
// ref: https://developer.zendesk.com/rest_api/docs/support/slas/policies#list-slas/policies
type SLAPolicyListOptions struct {
	PageOptions
	Active    bool   `url:"active"`
	SortBy    string `url:"sort_by,omitempty"`
	SortOrder string `url:"sort_order,omitempty"`
}

// SLAPolicyAPI an interface containing all slaPolicy related methods
type SLAPolicyAPI interface {
	GetSLAPolicies(ctx context.Context, opts *SLAPolicyListOptions) ([]SLAPolicy, Page, error)
	CreateSLAPolicy(ctx context.Context, slaPolicy SLAPolicy) (SLAPolicy, error)
	GetSLAPolicy(ctx context.Context, id int64) (SLAPolicy, error)
	UpdateSLAPolicy(ctx context.Context, id int64, slaPolicy SLAPolicy) (SLAPolicy, error)
	DeleteSLAPolicy(ctx context.Context, id int64) error
	GetSLAPoliciesIterator(ctx context.Context, opts *PaginationOptions) *Iterator[SLAPolicy]
	GetSLAPoliciesOBP(ctx context.Context, opts *OBPOptions) ([]SLAPolicy, Page, error)
	GetSLAPoliciesCBP(ctx context.Context, opts *CBPOptions) ([]SLAPolicy, client2.CursorPaginationMeta, error)
}

// GetSLAPolicies fetch slaPolicy list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/slas/policies#getting-slas/policies
func (z *Client) GetSLAPolicies(ctx context.Context, opts *SLAPolicyListOptions) ([]SLAPolicy, Page, error) {
	var data struct {
		SLAPolicies []SLAPolicy `json:"sla_policies"`
		Page
	}

	if opts == nil {
		return []SLAPolicy{}, Page{}, &client2.OptionsError{Opts: opts}
	}

	u, err := client2.AddOptions("/slas/policies.json", opts)
	if err != nil {
		return []SLAPolicy{}, Page{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return []SLAPolicy{}, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []SLAPolicy{}, Page{}, err
	}

	return data.SLAPolicies, data.Page, nil
}

// CreateSLAPolicy creates new slaPolicy
//
// ref: https://developer.zendesk.com/rest_api/docs/support/slas/policies#create-slaPolicy
func (z *Client) CreateSLAPolicy(ctx context.Context, slaPolicy SLAPolicy) (SLAPolicy, error) {
	var data, result struct {
		SLAPolicy SLAPolicy `json:"sla_policy"`
	}

	data.SLAPolicy = slaPolicy

	body, err := z.Post(ctx, "/slas/policies.json", data)
	if err != nil {
		return SLAPolicy{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return SLAPolicy{}, err
	}

	return result.SLAPolicy, nil
}

// GetSLAPolicy returns the specified slaPolicy
//
// ref: https://developer.zendesk.com/rest_api/docs/support/slas/policies#getting-slas/policies
func (z *Client) GetSLAPolicy(ctx context.Context, id int64) (SLAPolicy, error) {
	var result struct {
		SLAPolicy SLAPolicy `json:"sla_policy"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/slas/policies/%d.json", id))
	if err != nil {
		return SLAPolicy{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return SLAPolicy{}, err
	}

	return result.SLAPolicy, nil
}

// UpdateSLAPolicy updates the specified slaPolicy and returns the updated one
//
// ref: https://developer.zendesk.com/rest_api/docs/support/slas/policies#update-slaPolicy
func (z *Client) UpdateSLAPolicy(ctx context.Context, id int64, slaPolicy SLAPolicy) (SLAPolicy, error) {
	var data, result struct {
		SLAPolicy SLAPolicy `json:"sla_policy"`
	}

	data.SLAPolicy = slaPolicy

	body, err := z.Put(ctx, fmt.Sprintf("/slas/policies/%d.json", id), data)
	if err != nil {
		return SLAPolicy{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return SLAPolicy{}, err
	}

	return result.SLAPolicy, nil
}

// DeleteSLAPolicy deletes the specified slaPolicy
//
// ref: https://developer.zendesk.com/rest_api/docs/support/slas/policies#delete-slaPolicy
func (z *Client) DeleteSLAPolicy(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/slas/policies/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
