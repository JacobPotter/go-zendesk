package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMacros(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "macros.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	macros, _, err := c.GetMacros(ctx, &MacroListOptions{
		PageOptions: PageOptions{
			Page:    1,
			PerPage: 10,
		},
		SortBy:    "id",
		SortOrder: "asc",
	})
	if err != nil {
		t.Fatalf("Failed to get macros: %s", err)
	}

	expectedLength := 2
	if len(macros) != expectedLength {
		t.Fatalf("Returned macros does not have the expected length %d. Macros length is %d", expectedLength, len(macros))
	}
}

func TestGetMacro(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "macro.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	macro, err := c.GetMacro(ctx, 2)
	if err != nil {
		t.Fatalf("Failed to get macro: %s", err)
	}

	expectedID := int64(360111062754)
	if macro.ID != expectedID {
		t.Fatalf("Returned macro does not have the expected ID %d. Macro id is %d", expectedID, macro.ID)
	}

}

func TestCreateMacro(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "macro.json", http.StatusCreated)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	macro, err := c.CreateMacro(ctx, Macro{
		Title: "nyanyanyanya",
		// Comment: MacroComment{
		// 	Body: "(●ↀ ω ↀ )",
		// },
	})
	if err != nil {
		t.Fatalf("Failed to create macro: %s", err)
	}

	expectedID := int64(4)
	if macro.ID != expectedID {
		t.Fatalf("Returned macro does not have the expected ID %d. Macro id is %d", expectedID, macro.ID)
	}
}

func TestUpdateMacro(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "macro.json", http.StatusOK)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	macro, err := c.UpdateMacro(ctx, 2, Macro{})
	if err != nil {
		t.Fatalf("Failed to update macro: %s", err)
	}

	expectedID := int64(2)
	if macro.ID != expectedID {
		t.Fatalf("Returned macro does not have the expected ID %d. Macro id is %d", expectedID, macro.ID)
	}
}

func TestUpdateMacroFailure(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "macro.json", http.StatusInternalServerError)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.UpdateMacro(ctx, 2, Macro{})
	if err == nil {
		t.Fatal("BaseClient did not return error when api failed")
	}
}

func TestDeleteMacro(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteMacro(ctx, 437)
	if err != nil {
		t.Fatalf("Failed to delete macro field: %s", err)
	}
}
