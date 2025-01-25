package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/client"
	"net/http"
)

// API an interface containing all the zendesk client methods
type API interface {
	AppAPI
	AttachmentAPI
	AutomationAPI
	client.BaseAPI
	BrandAPI
	CustomRoleAPI
	DynamicContentAPI
	GroupAPI
	GroupMembershipAPI
	LocaleAPI
	MacroAPI
	OrganizationAPI
	OrganizationFieldAPI
	OrganizationMembershipAPI
	ScheduleAPI
	SearchAPI
	SLAPolicyAPI
	TagAPI
	TargetAPI
	TicketAuditAPI
	TicketAPI
	TicketCommentAPI
	TicketFieldAPI
	TicketFormAPI
	TriggerAPI
	TriggerCategoryAPI
	UserAPI
	UserFieldAPI
	ViewAPI
	WebhookAPI
	CustomObjectAPI
}

type Client struct {
	*client.BaseClient
}

func NewClient(httpClient *http.Client) (*Client, error) {
	zdClient, err := client.NewBaseClient(httpClient, false)
	if err != nil {
		return nil, err
	}
	return &Client{BaseClient: zdClient}, nil
}

var _ API = (*Client)(nil)
