package zendesk

import (
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"testing"
)

func TestGetDynamicContentItems(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "dynamic_content/items.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	items, page, err := c.GetDynamicContentItems(ctx)
	if err != nil {
		t.Fatalf("Failed to get dynamic content items: %s", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected length of dynamic content items is 2, but got %d", len(items))
	}

	if len(items[0].Variants) != 3 {
		t.Fatalf("expected length of items[0].Variants is 3, but got %d", len(items[0].Variants))
	}

	if page.HasPrev() || page.HasNext() {
		t.Fatalf("page fields are wrong: %v", page)
	}
}

func TestCreateDynamicContentItem(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "dynamic_content/items.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	item, err := c.CreateDynamicContentItem(ctx, DynamicContentItem{})
	if err != nil {
		t.Fatalf("Failed to get valid response: %s", err)
	}
	if item.ID == 0 {
		t.Fatal("Failed to create dynamic content item")
	}
}
