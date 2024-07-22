package zendesk

// Action is definition of what the macro does to the ticket
//
// ref: https://develop.zendesk.com/hc/en-us/articles/360056760874-Support-API-Actions-reference
type Action struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

// action field types which defined by system
// https://developer.zendesk.com/rest_api/docs/core/triggers#actions-reference
type ActionField string

func NewActionField(key string) ActionField {

	a := ActionField(key)

	return a
}

// action field types which defined by system
// https://developer.zendesk.com/rest_api/docs/core/triggers#actions-reference
const (
	// ActionFieldStatus status
	ActionFieldStatus ActionField = "status"
	// ActionFieldType type
	ActionFieldType ActionField = "type"
	// ActionFieldPriority priority
	ActionFieldPriority ActionField = "priority"
	// ActionFieldGroupID group_id
	ActionFieldGroupID ActionField = "group_id"
	// ActionFieldAssigneeID assignee_id
	ActionFieldAssigneeID ActionField = "assignee_id"
	// ActionFieldSetTags set_tags
	ActionFieldSetTags ActionField = "set_tags"
	// ActionFieldCurrentTags current_tags
	ActionFieldCurrentTags ActionField = "current_tags"
	// ActionFieldRemoveTags remove_tags
	ActionFieldRemoveTags ActionField = "remove_tags"
	// ActionFieldSatisfactionScore satisfaction_score
	ActionFieldSatisfactionScore ActionField = "satisfaction_score"
	// ActionFieldNotificationUser notification_user
	ActionFieldNotificationUser ActionField = "notification_user"
	// ActionFieldNotificationGroup notification_group
	ActionFieldNotificationGroup ActionField = "notification_group"
	// ActionFieldNotificationTarget notification_target
	ActionFieldNotificationTarget ActionField = "notification_target"
	// ActionFieldTweetRequester tweet_requester
	ActionFieldTweetRequester ActionField = "tweet_requester"
	// ActionFieldCC cc
	ActionFieldCC ActionField = "cc"
	// ActionFieldLocaleID locale_id
	ActionFieldLocaleID ActionField = "locale_id"
	// ActionFieldSubject subject
	ActionFieldSubject ActionField = "subject"
	// ActionFieldCommentValue comment_value
	ActionFieldCommentValue ActionField = "comment_value"
	// ActionFieldCommentValueHTML comment_value_html
	ActionFieldCommentValueHTML ActionField = "comment_value_html"
	// ActionFieldCommentModeIsPublic comment_mode_is_public
	ActionFieldCommentModeIsPublic ActionField = "comment_mode_is_public"
	// ActionFieldTicketFormID ticket_form_id
	ActionFieldTicketFormID ActionField = "ticket_form_id"
)

// Map of action fields to possible values
var ActionMap = map[ActionField][]string{
	ActionFieldStatus:      {"(new|open|pending|hold|solved|closed)"},
	ActionFieldType:        {"(question|incident|problem|task)"},
	ActionFieldPriority:    {"(low|normal|high|urgent)"},
	ActionFieldGroupID:     {"^[0-9]*$"},
	ActionFieldAssigneeID:  {"^[0-9]*$"},
	ActionFieldSetTags:     {"^[\\w]+(?:\\s[\\w]+)*$"},
	ActionFieldCurrentTags: {"^[\\w]+(?:\\s[\\w]+)*$"},
	ActionFieldRemoveTags:  {"^[\\w]+(?:\\s[\\w]+)*$"},
}
