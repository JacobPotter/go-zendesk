package client

import (
	"context"
	"fmt"
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestBaseClient_Get(t *testing.T) {
	ctx := context.Background()

	testRetrySeconds := 3

	type args struct {
		fileName    string
		statusCode  int
		clientRetry bool
		headers     map[string]string
	}

	testRetrySecondsStr := strconv.Itoa(testRetrySeconds)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should retry on 429",
			args: args{
				fileName:   "retry_error.json",
				statusCode: http.StatusTooManyRequests,
				headers: map[string]string{
					"Content-Type":    "application/json",
					"retry-after":     testRetrySecondsStr,
					"ratelimit-reset": testRetrySecondsStr,
				},
				clientRetry: true,
			},
			wantErr: false,
		},
		{
			name: "should not retry on 429",
			args: args{
				fileName:   "retry_error.json",
				statusCode: http.StatusTooManyRequests,
				headers: map[string]string{
					"Content-Type":    "application/json",
					"retry-after":     testRetrySecondsStr,
					"ratelimit-reset": testRetrySecondsStr,
				},
				clientRetry: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodGet, tt.args.fileName, tt.args.statusCode, tt.args.headers, true)
			c := NewTestClient(mockAPI, tt.args.clientRetry)
			defer mockAPI.Close()

			start := time.Now()
			_, err := c.Get(ctx, "/test")
			elapsed := time.Since(start)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.clientRetry {
				assert.True(t, elapsed.Seconds() >= float64(testRetrySeconds), fmt.Sprintf("Elapsed time %f", elapsed.Seconds()))
			} else {
				testError := Error{
					ErrorBody: testhelper.ReadFixture(t, filepath.Join(http.MethodGet, "retry_error.json")),
				}
				assert.IsType(t, testError, err)

				assert.Equal(t, testError.ErrorBody, err.(Error).ErrorBody)

				assert.Equal(t, testRetrySecondsStr, err.(Error).Resp.Header.Get("retry-after"))

				assert.Equal(t, http.StatusTooManyRequests, err.(Error).Resp.StatusCode)
			}
		})
	}
}
