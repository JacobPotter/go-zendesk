package zendesk

import (
	"bytes"
	"context"
	"crypto/sha1"
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWrite(t *testing.T) {
	file := testhelper.ReadFixture(t, filepath.Join(http.MethodPost, "upload.json"))
	h := sha1.New()
	h.Write(file)
	expectedSum := h.Sum(nil)
	r := bytes.NewReader(file)
	var attachmentSum []byte
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := sha1.New()
		_, err := io.Copy(h, r.Body)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
		attachmentSum = h.Sum(nil)
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(file)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	w := c.UploadAttachment(ctx, "foo", "bar")

	_, err := io.Copy(w, r)
	if err != nil {
		t.Fatal("Received an error from write")
	}

	out, err := w.Close()
	if err != nil {
		t.Fatalf("Received an error from close %v", err)
	}

	expected := "6bk3gql82em5nmf"
	if out.Token != expected {
		t.Fatalf("Received an unexpected token %s expected %s", out.Token, expected)
	}

	if !reflect.DeepEqual(expectedSum, attachmentSum) {
		t.Fatalf("Check sum of the written file does not match the expected checksum")
	}
}

func TestWriteCancelledContext(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "ticket.json", 201)
	defer mockAPI.Close()

	c := NewTestClient(mockAPI)

	canceled, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	w := c.UploadAttachment(canceled, "foo", "bar")

	file := []byte("body")
	r := bytes.NewBuffer(file)

	_, err := io.Copy(w, r)
	if err == nil {
		t.Fatalf("did not recieve expected error")
	}

	_, err = w.Close()
	if err == nil {
		t.Fatal("Did not receive error when closing writer")
	}
}

func TestDeleteUpload(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteUpload(ctx, "foobar")
	if err != nil {
		t.Fatalf("Failed to delete ticket field: %s", err)
	}
}

func TestDeleteUploadCanceledContext(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	canceled, cancelFunc := context.WithCancel(ctx)
	cancelFunc()

	err := c.DeleteUpload(canceled, "foobar")
	if err == nil {
		t.Fatal("did not get expected error")
	}
}

func TestGetAttachment(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "attachment.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	attachment, err := c.GetAttachment(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get attachment: %s", err)
	}

	expectedID := int64(498483)
	if attachment.ID != expectedID {
		t.Fatalf("Returned attachment does not have the expected ID %d. Attachment id is %d", expectedID, attachment.ID)
	}
}

func TestRedactCommentAttachment(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodPut, "redact_ticket_comment_attachment.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	err := c.RedactCommentAttachment(ctx, 123, 456, 789)

	if err != nil {
		t.Fatalf("Failed to redact ticket comment attachment: %s", err)
	}
}
