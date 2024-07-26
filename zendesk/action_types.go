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
	Value interface{} `json:"value"`
}

var _ ValidateValue[ActionResourceType] = &Action{}

func (a Action) Validate(resourceType ResourceType[ActionResourceType]) error {

	if err := resourceType.ValidateResourceType(); err != nil {
		return err
	}

	switch a.Value.(type) {
	case string:
		if err := ValidActionValuesMap.ValidateValue(
			ActionField(a.Field),
			a.Value.(string),
			"",
			resourceType,
		); err != nil {
			return err
		}
		return nil
	case []string:
		if len(a.Value.([]string)) == 0 {
			return fmt.Errorf("no empty for action value for field %s", a.Field)
		}
		if err := ValidActionValuesMap.ValidateValue(
			ActionField(a.Field),
			a.Value.([]string)[0],
			"",
			resourceType,
		); err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf("invalid value type %T for field %s", a.Value, a.Field)
	}
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
	ActionFieldCustomField ActionField = "custom_field_"
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

func (a ActionsValueValidator) ValidateValue(
	key ActionField,
	value string,
	_ Operator,
	resourceType ResourceType[ActionResourceType],
) error {

	if v, ok := a[key]; ok {
		keys := a.ValidKeys()

		if !slices.Contains(keys, key) && !strings.HasPrefix(
			string(key),
			string(ActionFieldCustomField),
		) {
			return fmt.Errorf("invalid action field %s", key)
		}

		if !slices.Contains(v.ResourceTypes.Elements(), resourceType.ToValue()) {
			return fmt.Errorf("invalid resource type for action key: %s", resourceType)
		}

		var result []byte

		if strings.HasPrefix(string(key), string(ActionFieldCustomField)) {
			after, _ := strings.CutPrefix(string(key), string(ActionFieldCustomField))
			result = v.ValidationRegex.Find([]byte(after))
		} else {
			result = v.ValidationRegex.Find([]byte(value))

		}
		if result == nil {
			return fmt.Errorf(
				"invalid action value %s. does not match regex: %s",
				string(result),
				v.ValidationRegex.String(),
			)
		}

		return nil
	}

	return fmt.Errorf("invalid action field %s", key)

}

func (a ActionsValueValidator) ValidKeys() []ActionField {
	return maps.Keys(a)
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
		ValidationRegex: regexp.MustCompile("(^$|^[0-9]*$)"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldAssigneeID: {
		ValidationRegex: regexp.MustCompile("(^$|current_user|^[0-9]*$)"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldSetTags: {
		ValidationRegex: regexp.MustCompile(`^\w+(?:\s\w+)*$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldCurrentTags: {
		ValidationRegex: regexp.MustCompile(`^\w+(?:\s\w+)*$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldRemoveTags: {
		ValidationRegex: regexp.MustCompile(`^\w+(?:\s\w+)*$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldCustomStatusId: {
		ValidationRegex: regexp.MustCompile("^[0-9]*$"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldTicketFormID: {
		ValidationRegex: regexp.MustCompile("^[0-9]*$"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldFollower: {
		ValidationRegex: regexp.MustCompile("(^$|current_user|^[0-9]*$)"),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldCustomField: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedActionTypes,
	},
	ActionFieldSatisfactionScore: {
		ValidationRegex: regexp.MustCompile("^offered$"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationUser: {
		ValidationRegex: regexp.MustCompile("(all_agents|requester_id|assignee_id|current_user|^[0-9]*$)"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationGroup: {
		ValidationRegex: regexp.MustCompile("(group_id|^[0-9]*$)"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationTarget: {
		ValidationRegex: regexp.MustCompile("^[0-9]*$"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldNotificationWebhook: {
		ValidationRegex: regexp.MustCompile("^[A-Za-z0-9]*$"),
		ResourceTypes:   triggerAutomationActionTypes,
	},
	ActionFieldCC: {
		ValidationRegex: regexp.MustCompile("(^$|current_user|^[0-9]*$)"),
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
		ValidationRegex: regexp.MustCompile("(channel:all|channel:web|channel:chat)"),
		ResourceTypes:   macroActionTypes,
	},
	ActionFieldCommentValueHTML: {
		ValidationRegex: regexp.MustCompile("(channel:all|channel:web|channel:chat)"),
		ResourceTypes:   macroActionTypes,
	},
	ActionFieldCommentModeIsPublic: {
		ValidationRegex: regexp.MustCompile("(true|false)"),
		ResourceTypes:   macroActionTypes,
	},
}
