package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

func HandlePOSTHttpRequest(ctx context.Context, url string, payload any, method string) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString("TAZAPAY_AUTH_TOKEN"),
	}

	jsonBody, ok1 := json.Marshal(payload)
	if ok1 != nil {
		return nil, fmt.Errorf("error creating request body: %w", ok1)
	}

	req, ok2 := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if ok2 != nil {
		return nil, fmt.Errorf("error creating request: %w", ok2)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, ok3 := http.DefaultClient.Do(req)
	if ok3 != nil {
		return nil, fmt.Errorf("error making request: %w", ok3)
	}
	defer resp.Body.Close()

	if resp.StatusCode < constants.HTTPStatusOKMin || resp.StatusCode >= constants.HTTPStatusOKMax {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: %v, body: <failed to read body: %w>",
				constants.ErrNonSuccessStatus, resp.Status, readErr)
		}

		return nil, fmt.Errorf("%w: %v, body: %s", constants.ErrNonSuccessStatus, resp.Status, string(bodyBytes))
	}

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body: %w", readErr)
	}

	var result map[string]any
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}

func HandleGETHttpRequest(ctx context.Context, url, method string) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString("TAZAPAY_AUTH_TOKEN"),
	}

	req, ok1 := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if ok1 != nil {
		return nil, fmt.Errorf("error creating request: %w", ok1)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, ok2 := http.DefaultClient.Do(req)
	if ok2 != nil {
		return nil, fmt.Errorf("error making request: %w", ok2)
	}
	defer resp.Body.Close()

	if resp.StatusCode < constants.HTTPStatusOKMin ||
		resp.StatusCode >= constants.HTTPStatusOKMax {

		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: %v, body: <failed to read body: %w>",
				constants.ErrNonSuccessStatus, resp.Status, readErr)
		}

		return nil, fmt.Errorf("%w: %v, body: %s", constants.ErrNonSuccessStatus,
			resp.Status, string(bodyBytes))
	}

	bodyBytes, ok3 := io.ReadAll(resp.Body)
	if ok3 != nil {
		return nil, fmt.Errorf("error reading response body: %w", ok3)
	}

	var result map[string]any
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}
