package zendesk

import (
	"testing"
)

func TestAction_ValidateAction(t *testing.T) {
	cases := []struct {
		testName     string
		action       Action
		resourceType ActionResourceType
		shouldPass   bool
	}{
		{
			testName: "should validate action object for resource type",
			action: Action{
				Field: string(ActionFieldStatus),
				Value: "open",
			},
			resourceType: TriggerActionResource,
			shouldPass:   true,
		}, {
			testName: "should not validate action object for resource type",
			action: Action{
				Field: string(ActionFieldStatus),
				Value: "blah",
			},
			resourceType: TriggerActionResource,
			shouldPass:   false,
		},
		{
			testName: "should not validate action object for invalid resource type",
			action: Action{
				Field: string(ActionFieldNotificationWebhook),
				Value: "ABCDEFG12345",
			},
			resourceType: MacroActionResource,
			shouldPass:   false,
		},
		{
			testName: "should not validate action object for invalid field id",
			action: Action{
				Field: "some_key",
				Value: "some_value",
			},
			resourceType: MacroActionResource,
			shouldPass:   false,
		},
		{
			testName: "should validate action object with array of values",
			action: Action{
				Field: string(ActionFieldNotificationWebhook),
				Value: []string{"ABCDEFG12345", "webhook body string"},
			},
			resourceType: TriggerActionResource,
			shouldPass:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := c.action.ValidateAction(c.resourceType); err != nil && c.shouldPass {
				t.Fatalf("action validation returned unexpected error: %s", err)
			}

		})
	}

}

func TestActionValueValidator_ValidateActionValue(t *testing.T) {
	cases := []struct {
		testName     string
		key          ActionField
		value        string
		resourceType ActionResourceType
		shouldPass   bool
	}{
		{
			testName:     "should pass with valid status value",
			key:          ActionFieldStatus,
			value:        "open",
			resourceType: TriggerActionResource,
			shouldPass:   true,
		},
		{
			testName:     "should fail with invalid status value",
			key:          ActionFieldStatus,
			value:        "blah",
			resourceType: TriggerActionResource,
			shouldPass:   false,
		},
		{
			testName:     "should fail with invalid field",
			key:          "some_key",
			value:        "blah",
			resourceType: TriggerActionResource,
			shouldPass:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := ValidActionValuesMap.ValidateActionValue(c.key, c.value, c.resourceType); err != nil && c.shouldPass {
				t.Fatalf("Validation failed for key: %s, value: %s.\n Error: %s", c.key, c.value, err.Error())
			}
		})
	}
}

func TestActionValueValidator_ValidateFieldId(t *testing.T) {
	cases := []struct {
		testName   string
		key        ActionField
		shouldPass bool
	}{
		{
			testName:   "should pass with valid id value",
			key:        ActionFieldStatus,
			shouldPass: true,
		},
		{
			testName:   "should fail with invalid id value",
			key:        "some_bad_value",
			shouldPass: false,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := ValidActionValuesMap.ValidateFieldId(c.key); err != nil && c.shouldPass {
				t.Fatalf("Validation failed for key: %s, Error: %s", c.key, err.Error())
			}
		})
	}
}
