package zendesk

import (
	"fmt"
	"golang.org/x/exp/maps"
	"regexp"
	"slices"
	"strings"
)

// Condition zendesk condition, see [Zendesk Conditions Reference]
//
// [Zendesk Conditions Reference]: https://developer.zendesk.com/documentation/ticketing/reference-guides/conditions-Reference/
type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator,omitempty"`
	Value    ParsedValue `json:"value,omitempty"`
}

func (c Condition) Validate(resourceType ResourceType[ConditionResourceType]) error {
	if err := resourceType.ValidateResourceType(); err != nil {
		return err
	}

	if err := ValidConditionOperatorValues.ValidateValue(
		ConditionField(c.Field),
		c.Value.Data,
		Operator(c.Operator),
		resourceType,
	); err != nil {
		return err
	}
	return nil

}

var _ ValidateValue[ConditionResourceType] = &Condition{}

type Conditions struct {
	All []Condition `json:"all,omitempty"`
	Any []Condition `json:"any,omitempty"`
}

type ConditionField string

func (c ConditionField) String() string {
	return string(c)
}

// condition field types which are defined by system
// https://developer.zendesk.com/rest_api/docs/core/triggers#conditions-reference
const (
	// ConditionFieldGroupID is alias for group_id
	ConditionFieldGroupID ConditionField = "group_id"
	// ConditionFieldAssigneeID is alias for assignee_id
	ConditionFieldAssigneeID ConditionField = "assignee_id"
	// ConditionFieldRequesterID is alias for requester_id
	ConditionFieldRequesterID ConditionField = "requester_id"
	// ConditionFieldOrganizationID is alias for organization_id
	ConditionFieldOrganizationID ConditionField = "organization_id"
	// ConditionFieldCurrentTags is alias for current_tags
	ConditionFieldCurrentTags ConditionField = "current_tags"
	// ConditionFieldViaID is alias for via_id
	ConditionFieldViaID ConditionField = "via_id"
	// ConditionFieldRecipient is alias for recipient
	ConditionFieldRecipient ConditionField = "recipient"
	// ConditionFieldCustomField is alias for custom_fields_ prefix
	ConditionFieldCustomField ConditionField = "custom_fields_"
	// ConditionFieldCustomFieldAlt is alias for ticket_fields_ prefix
	ConditionFieldCustomFieldAlt ConditionField = "ticket_fields_"
	// ConditionFieldType is alias for type
	ConditionFieldType ConditionField = "type"
	// ConditionFieldStatus is alias for status
	ConditionFieldStatus ConditionField = "status"
	// ConditionFieldPriority is alias for priority
	ConditionFieldPriority ConditionField = "priority"
	// ConditionFieldDescriptionIncludesWord is alias for description_includes_word
	ConditionFieldDescriptionIncludesWord ConditionField = "description_includes_word"
	// ConditionFieldLocaleID is alias for locale_id
	ConditionFieldLocaleID ConditionField = "locale_id"
	// ConditionFieldSatisfactionScore is alias for satisfaction_score
	ConditionFieldSatisfactionScore ConditionField = "satisfaction_score"
	// ConditionFieldSubjectIncludesWord is alias for subject_includes_word
	ConditionFieldSubjectIncludesWord ConditionField = "subject_includes_word"
	// ConditionFieldCommentIncludesWord is alias for comment_includes_word
	ConditionFieldCommentIncludesWord ConditionField = "comment_includes_word"
	// ConditionFieldCurrentViaID is alias for current_via_id
	ConditionFieldCurrentViaID ConditionField = "current_via_id"
	// ConditionFieldUpdateType is alias for update_type
	ConditionFieldUpdateType ConditionField = "update_type"
	// ConditionFieldCommentIsPublic is alias for comment_is_public
	ConditionFieldCommentIsPublic ConditionField = "comment_is_public"
	// ConditionFieldTicketIsPublic is alias for ticket_is_public
	ConditionFieldTicketIsPublic ConditionField = "ticket_is_public"
	// ConditionFieldReopens is alias for reopens
	ConditionFieldReopens ConditionField = "reopens"
	// ConditionFieldReplies is alias for reopens
	ConditionFieldReplies ConditionField = "replies"
	// ConditionFieldAgentStations is alias for agent_stations
	ConditionFieldAgentStations ConditionField = "agent_stations"
	// ConditionFieldGroupStations is alias for group_stations
	ConditionFieldGroupStations ConditionField = "group_stations"
	// ConditionFieldInBusinessHours is alias for in_business_hours
	ConditionFieldInBusinessHours ConditionField = "in_business_hours"
	// ConditionFieldRequesterTwitterFollowersCount is alias for requester_twitter_followers_count
	ConditionFieldRequesterTwitterFollowersCount ConditionField = "requester_twitter_followers_count"
	// ConditionFieldRequesterTwitterStatusesCount is alias for requester_twitter_statuses_count
	ConditionFieldRequesterTwitterStatusesCount ConditionField = "requester_twitter_statuses_count"
	// ConditionFieldRequesterTwitterVerified is alias for requester_twitter_verified
	ConditionFieldRequesterTwitterVerified ConditionField = "requester_twitter_verified"
	// ConditionFieldExactCreatedAt is alias for exact_created_at
	ConditionFieldExactCreatedAt ConditionField = "exact_created_at"
	// ConditionFieldNew is alias for NEW
	ConditionFieldNew ConditionField = "NEW"
	// ConditionFieldOpen is alias for OPEN
	ConditionFieldOpen ConditionField = "OPEN"
	// ConditionFieldPending is alias for PENDING
	ConditionFieldPending ConditionField = "PENDING"
	// ConditionFieldHold is alias for HOLD
	ConditionFieldHold ConditionField = "HOLD"
	// ConditionFieldSolved is alias for SOLVED
	ConditionFieldSolved ConditionField = "SOLVED"
	// ConditionFieldClosed is alias for CLOSED
	ConditionFieldClosed ConditionField = "CLOSED"
	// ConditionFieldAssignedAt is alias for assigned_at
	ConditionFieldAssignedAt ConditionField = "assigned_at"
	// ConditionFieldUpdatedAt is alias for updated_at
	ConditionFieldUpdatedAt ConditionField = "updated_at"
	// ConditionFieldRequesterUpdatedAt is alias for requester_updated_at
	ConditionFieldRequesterUpdatedAt ConditionField = "requester_updated_at"
	// ConditionFieldAssigneeUpdatedAt is alias for assignee_updated_at
	ConditionFieldAssigneeUpdatedAt ConditionField = "assignee_updated_at"
	// ConditionFieldDueDate is alias for due_date
	ConditionFieldDueDate ConditionField = "due_date"
	// ConditionFieldUntilDueDate is alias for until_due_date
	ConditionFieldUntilDueDate ConditionField = "until_due_date"
	// ConditionFieldBrandId is alias for brand_id
	ConditionFieldBrandId ConditionField = "brand_id"
	// ConditionFieldTicketFormId is alias for ticket_form_id
	ConditionFieldTicketFormId ConditionField = "ticket_form_id"
	// ConditionFieldUserCustomKey is a prefix alias for user.custom_fields.{key} where key is replaced with a key value
	ConditionFieldUserCustomKey      ConditionField = "user.custom_fields."
	ConditionFieldRequesterCustomKey ConditionField = "requester.custom_fields."
	// ConditionFieldOrganizationCustomKey is a prefix alias for organization.custom_fields.{key} where key is replaced with a key value
	ConditionFieldOrganizationCustomKey ConditionField = "organization.custom_fields."
	// ConditionFieldIsBusinessHours is an alias for is_business_hours
	ConditionFieldIsBusinessHours ConditionField = "is_business_hours"
	// ConditionFieldRequesterRole is alias for requester_role
	ConditionFieldRequesterRole ConditionField = "requester_role"
	// ConditionFieldAttachment is alias for attachment
	ConditionFieldAttachment ConditionField = "attachment"
	// ConditionFieldCC is an alias for cc
	ConditionFieldCC ConditionField = "cc"
	// ConditionFieldCustomStatusId  is an alias for custom_status_id
	ConditionFieldCustomStatusId ConditionField = "custom_status_id"
	// ConditionFieldTicketTypeId is an alias for ticket_type_id
	ConditionFieldTicketTypeId ConditionField = "ticket_type_id"
	// ConditionFieldSlaNextBreachAt is an alias for sla_next_breach_at
	ConditionFieldSlaNextBreachAt ConditionField = "sla_next_breach_at"
	ConditionFieldRole            ConditionField = "role"
	ConditionFieldWithinSchedule  ConditionField = "within_schedule"
)

type ConditionFields []ConditionField

var TimeBasedConditions = ConditionFields{
	ConditionFieldNew,
	ConditionFieldOpen,
	ConditionFieldPending,
	ConditionFieldSolved,
	ConditionFieldClosed,
	ConditionFieldHold,
	ConditionFieldAssignedAt,
	ConditionFieldUpdatedAt,
	ConditionFieldAssigneeUpdatedAt,
	ConditionFieldRequesterUpdatedAt,
	ConditionFieldDueDate,
	ConditionFieldUntilDueDate,
	ConditionFieldCustomField,
	ConditionFieldSlaNextBreachAt,
}

type ConditionResourceType string

var _ ResourceType[ConditionResourceType] = ConditionResourceType("")

const (
	TriggerConditionResource    ConditionResourceType = "trigger"
	AutomationConditionResource ConditionResourceType = "automation"
	ViewConditionResource       ConditionResourceType = "view"
	SlaConditionResource        ConditionResourceType = "sla"
)

var sharedConditionTypes = ResourceTypes[ConditionResourceType]{
	TriggerConditionResource,
	AutomationConditionResource,
	ViewConditionResource,
	SlaConditionResource,
}

var triggerAutomationViewConditionTypes = ResourceTypes[ConditionResourceType]{
	TriggerConditionResource,
	AutomationConditionResource,
	ViewConditionResource,
}

var triggerAutomationConditionTypes = ResourceTypes[ConditionResourceType]{
	TriggerConditionResource,
	AutomationConditionResource,
}

var slaConditionTypes = ResourceTypes[ConditionResourceType]{SlaConditionResource}

var timeBasedViewAutomationConditionTypes = ResourceTypes[ConditionResourceType]{
	ViewConditionResource,
	AutomationConditionResource,
}
var triggerConditionTypes = ResourceTypes[ConditionResourceType]{TriggerConditionResource}

func (c ConditionResourceType) ValidateResourceType() error {
	if !slices.Contains(sharedConditionTypes.Elements(), c) {
		return fmt.Errorf("invalid condition resource type: %s", c)
	}
	return nil
}

func (c ConditionResourceType) ToValue() ConditionResourceType {
	return c
}

type ConditionValueValidator ValueValidator[ConditionResourceType]

type ConditionsValueValidator map[ConditionField]ConditionValueValidator

func (c ConditionsValueValidator) ValidateValue(key ConditionField, value string, operator Operator, resourceType ResourceType[ConditionResourceType]) error {

	isCustomField := strings.HasPrefix(string(key), string(ConditionFieldCustomField))
	isCustomFieldAlt := strings.HasPrefix(string(key), string(ConditionFieldCustomFieldAlt))

	var newKey = key

	if isCustomField {
		newKey = ConditionFieldCustomField
	}

	if isCustomFieldAlt {
		newKey = ConditionFieldCustomFieldAlt
	}

	isOrgField := strings.HasPrefix(string(key), string(ConditionFieldOrganizationCustomKey))

	if isOrgField {
		newKey = ConditionFieldOrganizationCustomKey
	}

	isUserField := strings.HasPrefix(string(key), string(ConditionFieldUserCustomKey))

	if isUserField {
		newKey = ConditionFieldUserCustomKey
	}

	isRequesterField := strings.HasPrefix(string(key), string(ConditionFieldRequesterCustomKey))
	if isRequesterField {
		newKey = ConditionFieldRequesterCustomKey
	}

	if v, ok := c[newKey]; ok {

		keys := c.ValidKeys()

		if !slices.Contains(keys, string(newKey)) {
			return fmt.Errorf("invalid condition field %s", newKey)
		}

		if !slices.Contains(v.ResourceTypes.Elements(), resourceType.ToValue()) {
			return fmt.Errorf("invalid resource type for condition key: %s", newKey)
		}

		if len(v.ValidOperators) > 0 && !slices.Contains(v.ValidOperators, operator) {
			return fmt.Errorf("invalid operator for condition key: %s", newKey)
		}

		var result []byte

		if isCustomField {
			after, _ := strings.CutPrefix(string(key), string(ConditionFieldCustomField))
			result = v.ValidationRegex.Find([]byte(after))
		} else if isCustomFieldAlt {
			after, _ := strings.CutPrefix(string(key), string(ConditionFieldCustomFieldAlt))
			result = v.ValidationRegex.Find([]byte(after))
		} else {
			result = v.ValidationRegex.Find([]byte(value))

		}
		if result == nil {
			return fmt.Errorf(
				"invalid condition value %s. does not match regex: %s",
				string(result),
				v.ValidationRegex.String(),
			)
		}

		return nil
	}

	return fmt.Errorf("invalid condition field %s", newKey)

}

func (c ConditionsValueValidator) ValidKeys() []string {
	keys := maps.Keys(c)
	stringSlice := make([]string, len(keys))

	for i, field := range keys {
		stringSlice[i] = field.String()
	}

	slices.Sort(stringSlice)

	return stringSlice
}

var _ Validator[ConditionField, ConditionResourceType] = ConditionsValueValidator{}

var ValidConditionOperatorValues = ConditionsValueValidator{
	ConditionFieldGroupID: {
		ValidationRegex: regexp.MustCompile(`(^$|^\d+)`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			Changed,
			NotChanged,
			Value,
			ValuePrevious,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldAssigneeID: {
		ValidationRegex: regexp.MustCompile(`(^$|current_user|requester_id|^\d+)`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			Changed,
			NotChanged,
			Value,
			ValuePrevious,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldRequesterID: {
		ValidationRegex: regexp.MustCompile(`(^$|current_user|requester_id|^\d+)`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			Changed,
			NotChanged,
			Value,
			ValuePrevious,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldOrganizationID: {
		ValidationRegex: regexp.MustCompile(`(^$|^\d+$)`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			Changed,
			NotChanged,
			Value,
			ValuePrevious,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldCurrentTags: {
		ValidationRegex: regexp.MustCompile(`^\S+(?:\s\S+)*$`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators:  []Operator{Includes, NotIncludes},
	},
	ConditionFieldViaID: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators:  []Operator{Is, IsNot, Includes, NotIncludes},
	},
	ConditionFieldRecipient: {
		ValidationRegex: regexp.MustCompile(`.*`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators:  []Operator{EmptyOperator},
	},
	ConditionFieldCustomField: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators:  []Operator{Is, IsNot, WithinPreviousNDays, NotPresent, Present},
	},
	ConditionFieldCustomFieldAlt: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   slaConditionTypes,
		ValidOperators:  []Operator{Is, IsNot, WithinPreviousNDays},
	},
	ConditionFieldType: {
		ValidationRegex: regexp.MustCompile(`(question|incident|problem|task)`),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators:  []Operator{Is, IsNot},
	},
	ConditionFieldStatus: {
		ValidationRegex: regexp.MustCompile(`(new|open|pending|hold|solved|closed|^$)`),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			LessThan,
			GreaterThan,
			Changed,
			NotChanged,
			Value,
			ValuePrevious,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldPriority: {
		ValidationRegex: regexp.MustCompile(`(^$|low|normal|high|urgent)`),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			LessThan,
			GreaterThan,
			Changed,
			NotChanged,
			Value,
			ValuePrevious,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldLocaleID: {
		ValidationRegex: regexp.MustCompile("^[A-Za-z]{2,4}([_-][A-Za-z]{4})?([_-]([A-Za-z]{2}|[0-9]{3}))?$"),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators:  []Operator{Is, IsNot},
	},
	ConditionFieldSatisfactionScore: {
		ValidationRegex: regexp.MustCompile(`(good_with_comment|good|bad_with_comment|bad|false|true|offered|unoffered)`),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators: []Operator{
			Is,
			LessThan,
			GreaterThan,
			Changed,
			ChangedTo,
			ChangedFrom,
			NotChanged,
			NotChangedFrom,
			NotChangedTo,
			Value,
			ValuePrevious,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldBrandId: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators: []Operator{
			Is,
			LessThan,
			GreaterThan,
			Changed,
			ChangedTo,
			ChangedFrom,
			NotChanged,
			NotChangedFrom,
			NotChangedTo,
		},
	},
	ConditionFieldTicketFormId: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			LessThan,
			GreaterThan,
			Changed,
			ChangedTo,
			ChangedFrom,
			NotChanged,
			NotChangedFrom,
			NotChangedTo,
		},
	},
	ConditionFieldUserCustomKey: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			Present,
			NotPresent,
			Includes,
			NotIncludes,
			IncludesString,
			NotIncludesString,
		},
	},
	ConditionFieldRequesterCustomKey: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			Present,
			NotPresent,
			Includes,
			NotIncludes,
			IncludesString,
			NotIncludesString,
		},
	},
	ConditionFieldOrganizationCustomKey: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
			Present,
			NotPresent,
			Includes,
			NotIncludes,
			IncludesString,
			NotIncludesString,
		},
	},
	ConditionFieldSubjectIncludesWord: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			Includes,
			NotIncludes,
			Is,
			IsNot,
		},
	}, ConditionFieldCommentIncludesWord: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			Includes,
			NotIncludes,
			Is,
			IsNot,
		},
	},
	ConditionFieldCurrentViaID: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
		},
	},
	ConditionFieldUpdateType: {
		ValidationRegex: regexp.MustCompile(`(Create|Change)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators:  []Operator{EmptyOperator, Is},
	},
	ConditionFieldCommentIsPublic: {
		ValidationRegex: regexp.MustCompile(`(true|false|not_relevant|requester_can_see_comment)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators:  []Operator{EmptyOperator, Is, IsNot},
	},
	ConditionFieldTicketIsPublic: {
		ValidationRegex: regexp.MustCompile(`(public|private)`),
		ResourceTypes:   triggerAutomationConditionTypes,
		ValidOperators:  []Operator{Is, IsNot, EmptyOperator},
	},
	ConditionFieldReopens: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			GreaterThan,
			LessThan,
			Is,
		},
	},
	ConditionFieldReplies: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			GreaterThan,
			LessThan,
			Is,
		},
	},
	ConditionFieldAgentStations: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			GreaterThan,
			LessThan,
			Is,
		},
	}, ConditionFieldGroupStations: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			GreaterThan,
			LessThan,
			Is,
		},
	},
	ConditionFieldInBusinessHours: {
		ValidationRegex: regexp.MustCompile(`(true|false)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators:  []Operator{EmptyOperator},
	},
	ConditionFieldRequesterTwitterFollowersCount: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			GreaterThan,
			LessThan,
			Is,
		},
	},
	ConditionFieldRequesterTwitterStatusesCount: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			GreaterThan,
			LessThan,
			Is,
		},
	},
	ConditionFieldRequesterTwitterVerified: {
		ValidationRegex: regexp.MustCompile(`^$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators:  []Operator{EmptyOperator},
	},
	ConditionFieldRequesterRole: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
		},
	},
	ConditionFieldAttachment: {
		ValidationRegex: regexp.MustCompile(`^$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators:  []Operator{EmptyOperator},
	},
	ConditionFieldIsBusinessHours: {
		ValidationRegex: regexp.MustCompile(`(true|false)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators:  []Operator{EmptyOperator},
	},
	ConditionFieldCC: {
		ValidationRegex: regexp.MustCompile(`^$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators:  []Operator{EmptyOperator},
	},
	ConditionFieldCustomStatusId: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   sharedConditionTypes,
		ValidOperators: []Operator{
			Includes,
			NotIncludes,
			Is,
			IsNot,
			Changed,
			Value,
			ValuePrevious,
			NotChanged,
			NotValue,
			NotValuePrevious,
		},
	},
	ConditionFieldTicketTypeId: {
		ValidationRegex: regexp.MustCompile(`^([1234])$`),
		ResourceTypes:   slaConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
		},
	},
	ConditionFieldExactCreatedAt: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   slaConditionTypes,
		ValidOperators: []Operator{
			LessThan,
			LessThanEqual,
			GreaterThan,
			GreaterThanEqual,
		},
	},
	ConditionFieldNew: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldOpen: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldPending: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldHold: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldSolved: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldClosed: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldAssignedAt: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldUpdatedAt: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldRequesterUpdatedAt: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldAssigneeUpdatedAt: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldDueDate: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldUntilDueDate: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsBusinessHours,
			LessThan,
			LessThanBusinessHours,
			GreaterThan,
			GreaterThanBusinessHours,
		},
	},
	ConditionFieldDescriptionIncludesWord: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   timeBasedViewAutomationConditionTypes,
		ValidOperators:  []Operator{EmptyOperator},
	},
	ConditionFieldSlaNextBreachAt: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerAutomationViewConditionTypes,
		ValidOperators: []Operator{
			GreaterThan,
			LessThan,
			Is,
			IsNot,
		},
	},
	ConditionFieldRole: {
		ValidationRegex: regexp.MustCompile(`([\s\S]*)`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
		},
	},
	ConditionFieldWithinSchedule: {
		ValidationRegex: regexp.MustCompile(`^\d+$`),
		ResourceTypes:   triggerConditionTypes,
		ValidOperators: []Operator{
			Is,
			IsNot,
		},
	},
}
