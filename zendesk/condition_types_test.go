package zendesk

import (
	"regexp"
	"testing"

	"golang.org/x/exp/maps"
)

func TestRegexConditionFields(t *testing.T) {
	keys := maps.Keys(ConditionMap)

	for _, k := range keys {

		values := ConditionMap[k]

		for _, v := range values.ValuesRegex {
			_, err := regexp.Compile(v)

			if err != nil {
				t.Fatalf("Error converting regex: %s", err.Error())
			}
		}
	}
}
