package zendesk

import (
	"context"
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateBrand(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "brands.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateBrand(ctx, Brand{})
	if err != nil {
		t.Fatalf("Failed to send request to create brand: %s", err)
	}
}

func TestCreateBrandCanceledContext(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "brands.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	canceled, cancelFunc := context.WithCancel(ctx)
	cancelFunc()

	_, err := c.CreateBrand(canceled, Brand{})
	if err == nil {
		t.Fatalf("did not get expected error")
	}
}

func TestGetBrand(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "brand.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	brand, err := c.GetBrand(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get brand: %s", err)
	}

	expectedID := int64(360002143133)
	if brand.ID != expectedID {
		t.Fatalf("Returned brand does not have the expected ID %d. Brand ID is %d", expectedID, brand.ID)
	}
}

func TestUpdateBrand(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "brands.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	updatedBrand, err := c.UpdateBrand(ctx, int64(1234), Brand{})
	if err != nil {
		t.Fatalf("Failed to send request to create brand: %s", err)
	}

	expectedID := int64(360002143133)
	if updatedBrand.ID != expectedID {
		t.Fatalf("Updated brand %v did not have expected id %d", updatedBrand, expectedID)
	}
}

func TestUpdateBrandCanceledContext(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "brands.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	canceled, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	_, err := c.UpdateBrand(canceled, int64(1234), Brand{})
	if err == nil {
		t.Fatalf("did not get expected error")
	}
}

func TestDeleteBrand(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteBrand(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete brand: %s", err)
	}
}
