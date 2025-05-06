package zendesk

import (
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"testing"
)

func TestGetGroupMemberships(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "group_memberships.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	groupMemberships, _, err := c.GetGroupMemberships(ctx, &GroupMembershipListOptions{GroupID: 123})
	if err != nil {
		t.Fatalf("Failed to get group memberships: %s", err)
	}

	if len(groupMemberships) != 2 {
		t.Fatalf("expected length of group memberships is 2, but got %d", len(groupMemberships))
	}
}
