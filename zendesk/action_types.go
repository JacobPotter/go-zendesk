package zendesk

import (
	"fmt"
	"golang.org/x/exp/maps"
	"regexp"
	"slices"
	"strings"
)

// Action is definition of what the resource does to the ticket. [Zendesk Actions Reference]
//
// [Zendesk Actions Reference]: https://developer.zendesk.com/documentation/ticketing/reference-guides/actions-Reference/
type Action struct {
	Field string      `json:"field"`
	Value ParsedValue `json:"value,omitempty"`
}

var _ ValidateValue[ActionResourceType] = &Action{}

func (a Action) Validate(resourceType ResourceType[ActionResourceType]) error {

	if err := resourceType.ValidateResourceType(); err != nil {
		return err
	}

	if err := ValidActionValuesMap.ValidateValue(
		ActionField(a.Field),
		a.Value,
		"",
		resourceType,
	); err != nil {
		return err
	}
	return nil
}

// ActionField action field types which defined by system, see [Zendesk Actions Reference]
//
// [Zendesk Actions Reference]: https://developer.zendesk.com/documentation/ticketing/reference-guides/actions-Reference/
type ActionField string

func (a ActionField) String() string {
	return string(a)
}

// action field types which defined by system see [Zendesk Actions Reference]
//
// [Zendesk Actions Reference]: https://developer.zendesk.com/documentation/ticketing/reference-guides/actions-Reference/

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
	// ActionFieldNotificationWebhook notification_webhook
	ActionFieldNotificationWebhook ActionField = "notification_webhook"
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
	// ActionFieldCustomStatusId custom_status_id
	ActionFieldCustomStatusId ActionField = "custom_status_id"
	// ActionFieldFollower follower
	ActionFieldFollower ActionField = "follower"
	// ActionFieldBrandId brand_id
	ActionFieldBrandId ActionField = "brand_id"
	// ActionFieldAddSkills add_skills
	ActionFieldAddSkills ActionField = "add_skills"
	// ActionFieldSetSkills set_skills
	ActionFieldSetSkills ActionField = "set_skills"
	// ActionFieldRemoveSkills remove_skills
	ActionFieldRemoveSkills ActionField = "remove_skills"
	// ActionFieldCustomField custom_field_ prefix
	ActionFieldCustomField          ActionField = "custom_fields_"
	ActionSideConversationTicket    ActionField = "side_conversation_ticket"
	ActionSideConversationSlack     ActionField = "side_conversation_slack"
	ActionSetSchedule               ActionField = "set_schedule"
	ActionNotificationZis           ActionField = "notification_zis"
	ActionNotificationMessagingCsat ActionField = "notification_messaging_csat"
	ActionReplyPublic               ActionField = "reply_public"
	ActionReplyInternal             ActionField = "reply_internal"
)

// ActionResourceType String type of resource the action belongs to. Valid
// options are TriggerActionResource, AutomationActionResource, or
// MacroActionResource
type ActionResourceType string

var _ ResourceType[ActionResourceType] = ActionResourceType("")

const (
	TriggerActionResource    ActionResourceType = "trigger"
	AutomationActionResource ActionResourceType = "automation"
	MacroActionResource      ActionResourceType = "macro"
)

var sharedActionTypes = ResourceTypes[ActionResourceType]{
	TriggerActionResource,
	AutomationActionResource,
	MacroActionResource,
}

var triggerAutomationActionTypes = ResourceTypes[ActionResourceType]{
	TriggerActionResource,
	AutomationActionResource,
}

var triggerActionTypes = ResourceTypes[ActionResourceType]{TriggerActionResource}

var macroActionTypes = ResourceTypes[ActionResourceType]{MacroActionResource}

func (a ActionResourceType) ValidateResourceType() error {
	if !slices.Contains(sharedActionTypes.Elements(), a) {
		return fmt.Errorf("invalid action resource type: %s", a)
	}
	return nil
}

func (a ActionResourceType) ToValue() ActionResourceType {
	return a
}

type ActionsValueValidator map[ActionField]ValueValidator[ActionResourceType]

var _ Validator[ActionField, ActionResourceType] = ActionsValueValidator{}

func (a ActionsValueValidator) ValidateValue(key ActionField, value ParsedValue, _ Operator, resourceType ResourceType[ActionResourceType]) error {

	isCustomField := strings.HasPrefix(
		string(key),
		string(ActionFieldCustomField),
	)

	var newKey = key

	if isCustomField {
		newKey = ActionFieldCustomField
	}

	if v, ok := a[newKey]; ok {
		keys := a.ValidKeys()

		if !slices.Contains(keys, string(newKey)) {
			return fmt.Errorf("invalid action field %s", newKey)
		}

		if !slices.Contains(v.ResourceTypes.Elements(), resourceType.ToValue()) {
			return fmt.Errorf("invalid resource type for action key: %s", resourceType)
		}

		var found bool

		if isCustomField {
			after, _ := strings.CutPrefix(string(key), string(ActionFieldCustomField))
			found = v.ValidationRegex.Match([]byte(after))
		} else {
			if len(value.ListData) == 0 {
				found = v.ValidationRegex.Match([]byte(value.Data))
			} else {
				if newKey == ActionFieldNotificationUser || newKey == ActionFieldNotificationGroup || newKey == ActionFieldNotificationWebhook {
					found = v.ValidationRegex.Match([]byte(value.ListData[0]))
				} else {
					for _, val := range value.ListData {
						found = v.ValidationRegex.Match([]byte(val))
						if !found {
							return fmt.Errorf(
								"invalid condition value in list: %s. does not match regex: %s",
								val,
								v.ValidationRegex.String(),
							)
						}
					}
				}

			}

		}
		if !found {
			return fmt.Errorf(
				"invalid action value %v. does not match regex: %s",
				value,
				v.ValidationRegex.String(),
			)
		}

		return nil
	}

	return fmt.Errorf("invalid action field %s", newKey)

}

func (a ActionsValueValidator) ValidKeys() []string {

	keys := maps.Keys(a)

	stringSlice := make([]string, len(keys))

	for i, key := range keys {
		stringSlice[i] = string(key)
	}

	slices.Sort(stringSlice)

	return stringSlice
}

// ValidActionValuesMap Map of action fields to possible values, based on valid values from [Actions Reference]
//
// [Actions Reference]: https://developer.zendesk.com/documentation/ticketing/reference-guides/actions-reference/
var ValidActionValuesMap = ActionsValueValidator{
	ActionFieldStatus: {
		ValidationRegex: regexp.MustCompile("(new|open|pending|hold|solved|closed)"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldType: {
		ValidationRegex: regexp.MustCompile("(question|incident|problem|task)"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldPriority: {
		ValidationRegex: regexp.MustCompile("(low|normal|high|urgent)"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldGroupID: {
		ValidationRegex: regexp.MustCompile(`(^$|current_groups|^\d+$)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldAssigneeID: {
		ValidationRegex: regexp.MustCompile(`(^$|current_user|^\d+$)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldSetTags: {
		ValidationRegex: regexp.MustCompile(`^\S+(?:\s\S+)*$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldCurrentTags: {
		ValidationRegex: regexp.MustCompile(`^\S+(?:\s\S+)*$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldRemoveTags: {
		ValidationRegex: regexp.MustCompile(`^\S+(?:\s\S+)*$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldCustomStatusId: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldTicketFormID: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldFollower: {
		ValidationRegex: regexp.MustCompile(`(^$|current_user|^\d+$)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldCustomField: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldSatisfactionScore: {
		ValidationRegex: regexp.MustCompile("(good_with_comment|good|bad_with_comment|bad|false|true|offered|unoffered)"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationUser: {
		ValidationRegex: regexp.MustCompile(`(all_agents|requester_id|assignee_id|current_user|requester_and_ccs|^\d+$)`),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationGroup: {
		ValidationRegex: regexp.MustCompile(`(group_id|^\d+$)`),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationTarget: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationWebhook: {
		ValidationRegex: regexp.MustCompile("^.*$"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldCC: {
		ValidationRegex: regexp.MustCompile(`(^$|current_user|^\d+$)`),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldLocaleID: {
		ValidationRegex: regexp.MustCompile("^[A-Za-z]{2,4}([_-][A-Za-z]{4})?([_-]([A-Za-z]{2}|[0-9]{3}))?$"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldTweetRequester: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldBrandId: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerActionTypes,
	},
	ActionFieldAddSkills: {
		ValidationRegex: regexp.MustCompile(`^\w+(?:,\w+)*$`),
		ResourceTypes:   triggerActionTypes,
	},
	ActionFieldSetSkills: {
		ValidationRegex: regexp.MustCompile(`^\w+(?:,\w+)*$`),
		ResourceTypes:   triggerActionTypes,
	},
	ActionFieldRemoveSkills: {
		ValidationRegex: regexp.MustCompile(`^\w+(?:,\w+)*$`),
		ResourceTypes:   triggerActionTypes,
	},
	ActionFieldSubject: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   macroActionTypes,
	},
	ActionFieldCommentValue: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   macroActionTypes,
	},
	ActionFieldCommentValueHTML: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   macroActionTypes,
	},
	ActionFieldCommentModeIsPublic: {
		ValidationRegex: regexp.MustCompile("(true|false)"),
		ResourceTypes:   macroActionTypes,
	},
	ActionSideConversationTicket: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionSideConversationSlack: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionSetSchedule: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionNotificationZis: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionNotificationMessagingCsat: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionReplyPublic: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerActionTypes,
	},
	ActionReplyInternal: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerActionTypes,
	},
}
