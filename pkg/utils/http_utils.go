package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

func HandlePOSTHttpRequest(ctx context.Context, logger *slog.Logger, url string,
	payload any, method string,
) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString("TAZAPAY_AUTH_TOKEN"),
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal request payload", slog.Any("error", err))
		return nil, fmt.Errorf("error creating request body: %w", err)
	}

	logger.Info("Sending POST request", slog.Any("payload", payload))

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error("Failed to create HTTP request", slog.Any("error", err))
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("HTTP request failed", slog.Any("error", err))
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logger.Error("Failed to read response body", slog.Any("error", readErr))
		return nil, fmt.Errorf("error reading response body: %w", readErr)
	}

	if resp.StatusCode < constants.HTTPStatusOKMin || resp.StatusCode >= constants.HTTPStatusOKMax {
		logger.Error("Non-success HTTP response",
			slog.Int("status_code", resp.StatusCode),
			slog.String("body", string(bodyBytes)),
		)

		return nil, fmt.Errorf("%w: %v, body: %s", constants.ErrNonSuccessStatus, resp.Status, string(bodyBytes))
	}

	var result map[string]any
	if ok := json.Unmarshal(bodyBytes, &result); ok != nil {
		logger.Error("Failed to decode response JSON", slog.Any("error", ok))
		return nil, fmt.Errorf("error decoding response: %w", ok)
	}

	logger.Info("POST request successful")

	return result, nil
}

func HandleGETHttpRequest(ctx context.Context, logger *slog.Logger, url, method string) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString("TAZAPAY_AUTH_TOKEN"),
	}

	logger.Info("Sending GET request")

	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		logger.Error("Failed to create HTTP request", slog.Any("error", err))
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("HTTP request failed", slog.Any("error", err))
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logger.Error("Failed to read response body", slog.Any("error", readErr))
		return nil, fmt.Errorf("error reading response body: %w", readErr)
	}

	if resp.StatusCode < constants.HTTPStatusOKMin || resp.StatusCode >= constants.HTTPStatusOKMax {
		logger.Error("Non-success HTTP response",
			slog.Int("status_code", resp.StatusCode),
			slog.String("body", string(bodyBytes)),
		)

		return nil, fmt.Errorf("%w: %v, body: %s", constants.ErrNonSuccessStatus, resp.Status, string(bodyBytes))
	}

	var result map[string]any
	if ok := json.Unmarshal(bodyBytes, &result); ok != nil {
		logger.Error("Failed to decode response JSON", slog.Any("error", ok))
		return nil, fmt.Errorf("error decoding response: %w", ok)
	}

	logger.Info("GET request successful")

	return result, nil
}
