package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Get JSON data from API and returns its body as []bytes
func (c *BaseClient) Get(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL.String()+path, nil)
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {

		if c.Retry.ClientRetry {
			duration := getRetryWaitTime(resp)
			WaitForRetry(ctx, duration)
			return c.Get(ctx, path)
		} else {
			duration := getRetryWaitTime(resp)
			c.Retry.WaitTime = duration
			return nil, NewError(body, resp)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, NewError(body, resp)
	}
	return body, nil
}

// Post send data to API and returns response body as []bytes
func (c *BaseClient) Post(ctx context.Context, path string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL.String()+path, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		if c.Retry.ClientRetry {
			duration := getRetryWaitTime(resp)
			WaitForRetry(ctx, duration)
			return c.Post(ctx, path, data)
		} else {
			retry := getRetryWaitTime(resp)
			c.Retry.WaitTime = retry
			return nil, NewError(body, resp)
		}
	}

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
		return nil, NewError(body, resp)
	}

	return body, nil
}

// Put sends data to API and returns response body as []bytes
func (c *BaseClient) Put(ctx context.Context, path string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.BaseURL.String()+path, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		if c.Retry.ClientRetry {
			duration := getRetryWaitTime(resp)
			WaitForRetry(ctx, duration)
			return c.Put(ctx, path, data)
		} else {
			retry := getRetryWaitTime(resp)
			c.Retry.WaitTime = retry
			return nil, NewError(body, resp)
		}
	}

	// NOTE: some webhook mutation APIs return status No Content.
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
		return nil, NewError(body, resp)
	}

	return body, nil
}

// Patch sends data to API and returns response body as []bytes
func (c *BaseClient) Patch(ctx context.Context, path string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, c.BaseURL.String()+path, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		if c.Retry.ClientRetry {
			duration := getRetryWaitTime(resp)
			WaitForRetry(ctx, duration)
			return c.Patch(ctx, path, data)
		} else {
			retry := getRetryWaitTime(resp)
			c.Retry.WaitTime = retry
			return nil, NewError(body, resp)
		}
	}

	// NOTE: some webhook mutation APIs return status No Content.
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
		return nil, NewError(body, resp)
	}

	return body, nil
}

// Delete sends data to API and returns an error if unsuccessful
func (c *BaseClient) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL.String()+path, nil)
	if err != nil {
		return err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		duration := getRetryWaitTime(resp)
		WaitForRetry(ctx, duration)
		return c.Delete(ctx, path)
	}

	if resp.StatusCode != http.StatusNoContent {
		return NewError(body, resp)
	}

	return nil
}

// GetData is a generic helper function that retrieves and unmarshals JSON data from a specified URL.
// It takes four parameters:
// - a pointer to a BaseClient (z) which is used to execute the GET request,
// - a context (ctx) for managing the request's lifecycle,
// - a string (url) representing the endpoint from which data should be retrieved,
// - and an empty interface (data) where the retrieved data will be stored after being unmarshalled from JSON.
//
// The function starts by sending a GET request to the specified URL. If the request is successful,
// the returned body in the form of a byte slice is unmarshalled into the provided empty interface using the json.Unmarshal function.
//
// If an error occurs during either the GET request or the JSON unmarshalling, the function will return this error.
func GetData(z BaseAPI, ctx context.Context, url string, data any) error {
	body, err := z.Get(ctx, url)
	if err == nil {
		err = json.Unmarshal(body, data)
		if err != nil {
			return err
		}
	}
	return err
}
