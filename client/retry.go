package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// WaitForRetry will wait for set duration before returning, takes ctx
// context.Context as first arg so wait can be cancelled if context is cancelled
func WaitForRetry(ctx context.Context, retryWaitDuration time.Duration) {

	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled")
		return
	case <-time.After(retryWaitDuration):
		return
	}

}

func getRetryWaitTime(resp *http.Response) time.Duration {
	retryWaitTime, err := strconv.ParseInt(resp.Header.Get("retry-after"), 10, 64)

	if err != nil {
		fmt.Printf("Error getting retry header, trying secondary header")
		retryWaitTime, err = strconv.ParseInt(resp.Header.Get("ratelimit-reset"), 10, 64)
		if err != nil {
			fmt.Printf("Error getting retry header, setting retry after 60 sec by default")
			retryWaitTime = 60
		}
	}

	return time.Duration(retryWaitTime) * time.Second
}
