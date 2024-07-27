package zendesk

import "regexp"

type Operator string

const (
	Is                       Operator = "is"
	IsNot                    Operator = "is_not"
	Includes                 Operator = "includes"
	NotIncludes              Operator = "not_includes"
	IncludesString           Operator = "includes_string"
	NotIncludesString        Operator = "not_includes_string"
	GreaterThan              Operator = "greater_than"
	LessThan                 Operator = "less_than"
	GreaterThanEqual         Operator = "greater_than_equal"
	LessThanEqual            Operator = "less_than_equal"
	IsBusinessHours          Operator = "is_business_hours"
	GreaterThanBusinessHours Operator = "greater_than_business_hours"
	LessThanBusinessHours    Operator = "less_than_business_hours"
	WithinPreviousNDays      Operator = "within_previous_n_days"
	Present                  Operator = "present"
	NotPresent               Operator = "not_present"
	Changed                  Operator = "changed"
	ChangedTo                Operator = "changed_to"
	ChangedFrom              Operator = "changed_from"
	NotChanged               Operator = "not_changed"
	NotChangedTo             Operator = "not_changed_to"
	NotChangedFrom           Operator = "not_changed_from"
	Value                    Operator = "value"
	ValuePrevious            Operator = "value_previous"
	NotValue                 Operator = "not_value"
	NotValuePrevious         Operator = "not_value_previous"
	EmptyOperator            Operator = ""
)

type ResourceType[T any] interface {
	ValidateResourceType() error
	ToValue() T
}

type ResourceTypes[T any] []ResourceType[T]

func (r ResourceTypes[T]) Elements() []T {
	elements := make([]T, len(r))

	for i, r2 := range r {
		elements[i] = r2.ToValue()
	}
	return elements
}

type ValidateValue[T any] interface {
	Validate(resourceType ResourceType[T]) error
}

type Validator[F any, T any] interface {
	ValidateValue(key F, value string, operator Operator, resourceType ResourceType[T]) error
	ValidKeys() []string
}

type ValueValidator[T any] struct {
	ValidationRegex *regexp.Regexp
	ResourceTypes   ResourceTypes[T]
	ValidOperators  []Operator
}
