package zendesk

type ConditionField string

// condition field types which are defined by system
// https://developer.zendesk.com/rest_api/docs/core/triggers#conditions-reference
const (
	// ConditionFieldGroupID group_id
	ConditionFieldGroupID ConditionField = "group_id"
	// ConditionFieldAssigneeID assignee_id
	ConditionFieldAssigneeID ConditionField = "assignee_id"
	// ConditionFieldRequesterID requester_id
	ConditionFieldRequesterID ConditionField = "requester_id"
	// ConditionFieldOrganizationID organization_id
	ConditionFieldOrganizationID ConditionField = "organization_id"
	// ConditionFieldCurrentTags current_tags
	ConditionFieldCurrentTags ConditionField = "current_tags"
	// ConditionFieldViaID via_id
	ConditionFieldViaID ConditionField = "via_id"
	// ConditionFieldRecipient recipient
	ConditionFieldRecipient ConditionField = "recipient"
	// ConditionFieldType type
	ConditionFieldType ConditionField = "type"
	// ConditionFieldStatus status
	ConditionFieldStatus ConditionField = "status"
	// ConditionFieldPriority priority
	ConditionFieldPriority ConditionField = "priority"
	// ConditionFieldDescriptionIncludesWord description_includes_word
	ConditionFieldDescriptionIncludesWord ConditionField = "description_includes_word"
	// ConditionFieldLocaleID locale_id
	ConditionFieldLocaleID ConditionField = "locale_id"
	// ConditionFieldSatisfactionScore satisfaction_score
	ConditionFieldSatisfactionScore ConditionField = "satisfaction_score"
	// ConditionFieldSubjectIncludesWord subject_includes_word
	ConditionFieldSubjectIncludesWord ConditionField = "subject_includes_word"
	// ConditionFieldCommentIncludesWord comment_includes_word
	ConditionFieldCommentIncludesWord ConditionField = "comment_includes_word"
	// ConditionFieldCurrentViaID current_via_id
	ConditionFieldCurrentViaID ConditionField = "current_via_id"
	// ConditionFieldUpdateType update_type
	ConditionFieldUpdateType ConditionField = "update_type"
	// ConditionFieldCommentIsPublic comment_is_public
	ConditionFieldCommentIsPublic ConditionField = "comment_is_public"
	// ConditionFieldTicketIsPublic ticket_is_public
	ConditionFieldTicketIsPublic ConditionField = "ticket_is_public"
	// ConditionFieldReopens reopens
	ConditionFieldReopens ConditionField = "reopens"
	// ConditionFieldReplies
	ConditionFieldReplies ConditionField = ""
	// ConditionFieldAgentStations agent_stations
	ConditionFieldAgentStations ConditionField = "agent_stations"
	// ConditionFieldGroupStations group_stations
	ConditionFieldGroupStations ConditionField = "group_stations"
	// ConditionFieldInBusinessHours in_business_hours
	ConditionFieldInBusinessHours ConditionField = "in_business_hours"
	// ConditionFieldRequesterTwitterFollowersCount requester_twitter_followers_count
	ConditionFieldRequesterTwitterFollowersCount ConditionField = "requester_twitter_followers_count"
	// ConditionFieldRequesterTwitterStatusesCount requester_twitter_statuses_count
	ConditionFieldRequesterTwitterStatusesCount ConditionField = "requester_twitter_statuses_count"
	// ConditionFieldRequesterTwitterVerified requester_twitter_verified
	ConditionFieldRequesterTwitterVerified ConditionField = "requester_twitter_verified"
	// ConditionFieldTicketTypeID ticket_type_id
	ConditionFieldTicketTypeID ConditionField = "ticket_type_id"
	// ConditionFieldExactCreatedAt exact_created_at
	ConditionFieldExactCreatedAt ConditionField = "exact_created_at"
	// ConditionFieldNew NEW
	ConditionFieldNew ConditionField = "NEW"
	// ConditionFieldOpen OPEN
	ConditionFieldOpen ConditionField = "OPEN"
	// ConditionFieldPending PENDING
	ConditionFieldPending ConditionField = "PENDING"
	// ConditionFieldSolved SOLVED
	ConditionFieldSolved ConditionField = "SOLVED"
	// ConditionFieldClosed CLOSED
	ConditionFieldClosed ConditionField = "CLOSED"
	// ConditionFieldAssignedAt assigned_at
	ConditionFieldAssignedAt ConditionField = "assigned_at"
	// ConditionFieldUpdatedAt updated_at
	ConditionFieldUpdatedAt ConditionField = "updated_at"
	// ConditionFieldRequesterUpdatedAt requester_updated_at
	ConditionFieldRequesterUpdatedAt ConditionField = "requester_updated_at"
	// ConditionFieldAssigneeUpdatedAt
	ConditionFieldAssigneeUpdatedAt ConditionField = "assignee_updated_at"
	// ConditionFieldDueDate due_date
	ConditionFieldDueDate ConditionField = "due_date"
	// ConditionFieldUntilDueDate until_due_date
	ConditionFieldUntilDueDate ConditionField = "until_due_date"
)

type ConditionOperator string

const (
	GreaterThan           ConditionOperator = "greater_than"
	Includes              ConditionOperator = "includes"
	Is                    ConditionOperator = "is"
	IsNot                 ConditionOperator = "is_not"
	LessThan              ConditionOperator = "less_than"
	LessThanBusinessHours ConditionOperator = "less_than_business_hours"
	NotIncludes           ConditionOperator = "not_includes"
	NotPresent            ConditionOperator = "not_present"
	Present               ConditionOperator = "present"
	WithinPreviousNDays   ConditionOperator = "within_previous_n_days"
)

var ConditionMap = map[ConditionField]struct {
	Operator    ConditionOperator
	ValuesRegex []string
}{}

// Condition zendesk automation condition
//
// ref: https://developer.zendesk.com/rest_api/docs/core/automations#conditions-reference
type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}
type Conditions struct {
	All []Condition `json:"all,omitempty"`
	Any []Condition `json:"any,omitempty"`
}
