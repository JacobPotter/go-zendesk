package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"testing"
)

func TestGetTicketTags(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "tags.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	tags, err := c.GetTicketTags(ctx, int64(2))
	if err != nil {
		t.Fatalf("Failed to get ticket tags: %s", err)
	}

	expectedLength := 2
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
}

func TestGetOrganizationTags(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "tags.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	tags, err := c.GetOrganizationTags(ctx, int64(2))
	if err != nil {
		t.Fatalf("Failed to get organization tags: %s", err)
	}

	expectedLength := 2
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
}

func TestGetUserTags(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "tags.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	tags, err := c.GetUserTags(ctx, int64(2))
	if err != nil {
		t.Fatalf("Failed to get user tags: %s", err)
	}

	expectedLength := 2
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
}

func TestAddTicketTags(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "tags.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	tag := Tag("example")

	tags, err := c.AddTicketTags(ctx, 2, []Tag{tag})
	if err != nil {
		t.Fatalf("Failed to add ticket tags: %s", err)
	}

	expectedLength := 3
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
	if tags[2] != tag {
		t.Fatalf("Returned tags does not have the expexted tag %s. %s given", "important", tags[0])
	}
}

func TestAddOrganizationTags(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "tags.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	tag := Tag("example")

	tags, err := c.AddOrganizationTags(ctx, 2, []Tag{tag})
	if err != nil {
		t.Fatalf("Failed to add ticket tags: %s", err)
	}

	expectedLength := 3
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
	if tags[2] != tag {
		t.Fatalf("Returned tags does not have the expexted tag %s. %s given", "important", tags[0])
	}
}

func TestAddUserTags(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "tags.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	tag := Tag("example")

	tags, err := c.AddUserTags(ctx, 2, []Tag{tag})
	if err != nil {
		t.Fatalf("Failed to add ticket tags: %s", err)
	}

	expectedLength := 3
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
	if tags[2] != tag {
		t.Fatalf("Returned tags does not have the expexted tag %s. %s given", "important", tags[0])
	}
}
