package zendesk

import (
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"testing"
)

func TestGetLocales(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "locales.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	locales, err := c.GetLocales(ctx)
	if err != nil {
		t.Fatalf("Failed to get locales: %s", err)
	}

	if len(locales) != 3 {
		t.Fatalf("expected length of groups is 3, but got %d", len(locales))
	}
}
