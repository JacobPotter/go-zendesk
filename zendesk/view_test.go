package zendesk

import (
	"encoding/json"
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestGetView(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "view.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	view, err := c.GetView(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get view: %s", err)
	}

	expectedID := int64(360002440594)
	if view.ID != expectedID {
		t.Fatalf("Returned view does not have the expected ID %d. View ID is %d", expectedID, view.ID)
	}
}

func TestCreateView(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodPost, "view.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	fileName := "view_body.json"

	filePathAbs, err := filepath.Abs(filepath.Join("../fixture/body", fileName))

	if err != nil {
		t.Fatalf("Unable to generate filepath, error: %s", err.Error())
	}

	byteValue := openJsonFile(t, filePathAbs, fileName)

	var viewResp struct {
		View View `json:"view"`
	}

	err = json.Unmarshal(byteValue, &viewResp)

	if err != nil {
		t.Fatalf("Unable to unmarshal json, error: %s", err.Error())
	}

	view, err := c.CreateView(ctx, viewResp.View)
	if err != nil {
		t.Fatalf("Failed to get view: %s", err)
	}

	expectedTitle := "Kelly's tickets"
	if view.Title != expectedTitle {
		t.Fatalf("Returned view does not have the expected value %s. Actual value is %s", expectedTitle, view.Title)
	}
}

func TestUpdateView(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodPut, "view.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	fileName := "view_body.json"

	filePathAbs, err := filepath.Abs(filepath.Join("../fixture/body", fileName))

	if err != nil {
		t.Fatalf("Unable to generate filepath, error: %s", err.Error())
	}

	byteValue := openJsonFile(t, filePathAbs, fileName)

	var viewResp struct {
		View View `json:"view"`
	}

	err = json.Unmarshal(byteValue, &viewResp)

	if err != nil {
		t.Fatalf("Unable to unmarshal json, error: %s", err.Error())
	}

	view, err := c.UpdateView(ctx, 360002440594, viewResp.View)
	if err != nil {
		t.Fatalf("Failed to get view: %s", err)
	}

	expectedTitle := "Kelly's tickets"
	if view.Title != expectedTitle {
		t.Fatalf("Returned view does not have the expected value %s. Actual value is %s", expectedTitle, view.Title)
	}
}

func TestGetViews(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "views.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	views, _, err := c.GetViews(ctx)
	if err != nil {
		t.Fatalf("Failed to get views: %s", err)
	}

	if len(views) != 2 {
		t.Fatalf("expected length of views is 2, but got %d", len(views))
	}
}

func TestGetCountTicketsInViewsTestGetViews(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "views_ticket_count.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()
	ids := []string{"25", "78"}
	viewsCount, err := c.GetCountTicketsInViews(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to get views tickets count: %s", err)
	}

	if len(viewsCount) != 2 {
		t.Fatalf("expected length of views ticket counts is 2, but got %d", len(viewsCount))
	}
}

func TestDeleteView(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteView(ctx, 437)
	if err != nil {
		t.Fatalf("Failed to delete macro field: %s", err)
	}

}
