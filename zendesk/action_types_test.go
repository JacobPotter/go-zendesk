package zendesk

import (
	"regexp"
	"testing"

	"golang.org/x/exp/maps"
)

func TestRegexActionFields(t *testing.T) {

	keys := maps.Keys(ActionMap)

	for _, k := range keys {

		values := ActionMap[k]

		for _, v := range values {
			_, err := regexp.Compile(v)

			if err != nil {
				t.Fatalf("Error converting regex: %s", err.Error())
			}
		}
	}

}
