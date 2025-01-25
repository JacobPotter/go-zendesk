package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Scopes map[string]struct {
	Scopes []string `json:"scopes"`
}

// Configuration is a dictionary of custom configuration fields
type Configuration struct {
	AssignTicketsToAnyBrand      bool   `json:"assign_tickets_to_any_brand"`
	AssignTicketsToAnyGroup      bool   `json:"assign_tickets_to_any_group"`
	ChatAccess                   bool   `json:"chat_access"`
	CustomObjects                Scopes `json:"custom_objects"`
	EndUserListAccess            string `json:"end_user_list_access"`
	EndUserProfileAccess         string `json:"end_user_profile_access"`
	ExploreAccess                string `json:"explore_access"`
	ForumAccess                  string `json:"forum_access"`
	ForumAccessRestrictedContent bool   `json:"forum_access_restricted_content"`
	LightAgent                   bool   `json:"light_agent"`
	MacroAccess                  string `json:"macro_access"`
	ManageAutomations            bool   `json:"manage_automations"`
	ManageBusinessRules          bool   `json:"manage_business_rules"`
	ManageContextualWorkspaces   bool   `json:"manage_contextual_workspaces"`
	ManageDynamicContent         bool   `json:"manage_dynamic_content"`
	ManageExtensionsAndChannels  bool   `json:"manage_extensions_and_channels"`
	ManageFacebook               bool   `json:"manage_facebook"`
	ManageGroupMemberships       bool   `json:"manage_group_memberships"`
	ManageGroups                 bool   `json:"manage_groups"`
	ManageOrganizationFields     bool   `json:"manage_organization_fields"`
	ManageOrganizations          bool   `json:"manage_organizations"`
	ManageRoles                  string `json:"manage_roles"`
	ManageSkills                 bool   `json:"manage_skills"`
	ManageSlas                   bool   `json:"manage_slas"`
	ManageSuspendedTickets       bool   `json:"manage_suspended_tickets"`
	ManageTeamMembers            string `json:"manage_team_members"`
	ManageTicketFields           bool   `json:"manage_ticket_fields"`
	ManageTicketForms            bool   `json:"manage_ticket_forms"`
	ManageTriggers               bool   `json:"manage_triggers"`
	ManageUserFields             bool   `json:"manage_user_fields"`
	OrganizationEditing          bool   `json:"organization_editing"`
	OrganizationNotesEditing     bool   `json:"organization_notes_editing"`
	ReportAccess                 string `json:"report_access"`
	SideConversationCreate       bool   `json:"side_conversation_create"`
	TicketAccess                 string `json:"ticket_access"`
	TicketCommentAccess          string `json:"ticket_comment_access"`
	TicketDeletion               bool   `json:"ticket_deletion"`
	TicketRedaction              bool   `json:"ticket_redaction"`
	ViewDeletedTickets           bool   `json:"view_deleted_tickets"`
	TicketEditing                bool   `json:"ticket_editing"`
	TicketMerge                  bool   `json:"ticket_merge"`
	TicketTagEditing             bool   `json:"ticket_tag_editing"`
	TwitterSearchAccess          bool   `json:"twitter_search_access"`
	ViewAccess                   string `json:"view_access"`
	VoiceAccess                  bool   `json:"voice_access"`
	VoiceDashboardAccess         bool   `json:"voice_dashboard_access"`
}

// CustomRole is zendesk CustomRole JSON payload format
// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/
type CustomRole struct {
	Description     string        `json:"description,omitempty"`
	ID              int64         `json:"id,omitempty"`
	TeamMemberCount int64         `json:"team_member_count"`
	Name            string        `json:"name"`
	Configuration   Configuration `json:"configuration"`
	RoleType        int64         `json:"role_type"`
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	UpdatedAt       time.Time     `json:"updated_at,omitempty"`
}

// CustomRoleAPI an interface containing all CustomRole related methods
type CustomRoleAPI interface {
	GetCustomRoles(ctx context.Context) ([]CustomRole, error)
	GetCustomRole(ctx context.Context, id int64) (CustomRole, error)
	CreateCustomRole(ctx context.Context, customRole CustomRole) (CustomRole, error)
	UpdateCustomRole(ctx context.Context, updatedId int64, customRole CustomRole) (CustomRole, error)
	DeleteCustomRole(ctx context.Context, id int64) error
}

// GetCustomRoles fetch CustomRoles list
func (z *Client) GetCustomRoles(ctx context.Context) ([]CustomRole, error) {
	var data struct {
		CustomRoles []CustomRole `json:"custom_roles"`
		Page
	}

	u := "/custom_roles.json"

	body, err := z.Get(ctx, u)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.CustomRoles, nil
}

func (z *Client) GetCustomRole(ctx context.Context, id int64) (CustomRole, error) {
	var data struct {
		CustomRole CustomRole `json:"custom_role"`
	}

	u := fmt.Sprintf("/custom_roles/%d.json", id)
	body, err := z.Get(ctx, u)
	if err != nil {
		return CustomRole{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return CustomRole{}, err
	}
	return data.CustomRole, nil
}

func (z *Client) CreateCustomRole(ctx context.Context, customRole CustomRole) (CustomRole, error) {
	var data, result struct {
		CustomRole CustomRole `json:"custom_role"`
	}
	u := "/custom_roles.json"

	data.CustomRole = customRole

	body, err := z.Post(ctx, u, data)
	if err != nil {
		return CustomRole{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomRole{}, err
	}
	return result.CustomRole, nil
}

func (z *Client) UpdateCustomRole(ctx context.Context, updatedId int64, customRole CustomRole) (CustomRole, error) {
	var data, result struct {
		CustomRole CustomRole `json:"custom_role"`
	}
	u := fmt.Sprintf("/custom_roles/%d.json", updatedId)
	data.CustomRole = customRole
	body, err := z.Put(ctx, u, data)
	if err != nil {
		return CustomRole{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomRole{}, err
	}
	return result.CustomRole, nil

}

func (z *Client) DeleteCustomRole(ctx context.Context, id int64) error {
	u := fmt.Sprintf("/custom_roles/%d.json", id)
	err := z.Delete(ctx, u)

	if err != nil {
		return err
	}
	return nil

}
