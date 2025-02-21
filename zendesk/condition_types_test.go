package zendesk

import "testing"

func TestCondition_Validate(t *testing.T) {
	cases := []struct {
		testName     string
		condition    Condition
		resourceType ConditionResourceType
		shouldPass   bool
	}{
		{
			testName: "should validate condition object for resource type",
			condition: Condition{
				Field:    string(ConditionFieldStatus),
				Operator: string(Is),
				Value:    ParsedValue{Data: "open"},
			},
			resourceType: TriggerConditionResource,
			shouldPass:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			if err := c.condition.Validate(c.resourceType); err != nil && c.shouldPass {
				t.Fatalf("condition validation returned unexpected error: %s", err)

			}
		})
	}
}
