package zendesk

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

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

type ParsedValue struct {
	Data string
}

var _ json.Unmarshaler = &ParsedValue{}
var _ json.Marshaler = &ParsedValue{}

func (v *ParsedValue) UnmarshalJSON(bytes []byte) error {

	var valueRaw any

	if err := json.Unmarshal(bytes, &valueRaw); err != nil {
		return err
	}

	var value string

	switch newValue := valueRaw.(type) {
	case string:
		value = newValue
	case int:
		value = strconv.Itoa(newValue)
	case int32:
	case int64:
		value = strconv.Itoa(int(newValue))
	case float64:
		value = strconv.FormatFloat(newValue, 'f', -1, 64)
	case float32:
		value = strconv.FormatFloat(float64(newValue), 'f', -1, 32)
	case bool:
		value = strconv.FormatBool(newValue)
	case time.Time:
		value = newValue.Format(time.RFC3339)
	case nil:
		value = ""
	default:
		return fmt.Errorf("invalid value type: %T", newValue)
	}

	v.Data = value

	return nil

}

func (v *ParsedValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Data)
}
