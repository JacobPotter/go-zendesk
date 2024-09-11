package zendesk

//nolint
//go:generate mockgen -source=api.go -destination=mock/client.go -package=mock -mock_names=API=Client -aux_files github.com/JacobPotter/go-zendesk/zendesk=app.go,github.com/JacobPotter/go-zendesk/zendesk=attachment.go API

// API an interface containing all the zendesk client methods
type API interface {
	AppAPI
	AttachmentAPI
	AutomationAPI
	BaseAPI
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

var _ API = (*Client)(nil)
