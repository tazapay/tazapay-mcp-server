//nolint:sloglint // slog attributes can be used
package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

func HandlePOSTHttpRequest(ctx context.Context, logger *slog.Logger, url string,
	payload any, method string,
) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString(constants.StrTAZAPAYAuthToken),
	}

	var reqBody io.Reader
	if payload == nil {
		reqBody = http.NoBody

		logger.InfoContext(ctx, "Sending POST request with empty body")
	} else {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			logger.ErrorContext(ctx, constants.StrFailedToCreateHTTPRequest, slog.Any(constants.Error, err))
			return nil, fmt.Errorf(constants.StrErrorCreatingRequest, err)
		}

		logger.InfoContext(ctx, "Sending POST request",
			slog.Any("payload", payload),
		)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrFailedToCreateHTTPRequest, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorCreatingRequest, err)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrHTTPRequestFailed, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorMakingRequest, err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logger.ErrorContext(ctx, constants.StrFailedToReadResponseBody, slog.Any(constants.Error, readErr))
		return nil, fmt.Errorf(constants.StrErrorReadingResponseBody, readErr)
	}

	if resp.StatusCode < constants.HTTPStatusOKMin || resp.StatusCode >= constants.HTTPStatusOKMax {
		logger.ErrorContext(ctx, constants.StrNonSuccessHTTPResponse,
			slog.Int(constants.StrStatusCode, resp.StatusCode),
		)

		return nil, fmt.Errorf(constants.StrWrappedErrorWithBody,
			constants.ErrNonSuccessStatus, resp.Status, string(bodyBytes))
	}

	var result map[string]any
	if ok := json.Unmarshal(bodyBytes, &result); ok != nil {
		logger.ErrorContext(ctx, constants.StrFailedToDecodeResponseJSON, slog.Any(constants.Error, ok))
		return nil, fmt.Errorf(constants.StrErrorDecodingResponse, ok)
	}

	logger.InfoContext(ctx, "POST request successful")

	return result, nil
}

func HandleGETHttpRequest(ctx context.Context, logger *slog.Logger,
	url, method string,
) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString(constants.StrTAZAPAYAuthToken),
	}

	logger.InfoContext(ctx, "Sending GET request", 
	slog.Any("headers", headers))

	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrFailedToCreateHTTPRequest, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorCreatingRequest, err)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrHTTPRequestFailed, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorMakingRequest, err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logger.ErrorContext(ctx, constants.StrFailedToReadResponseBody, slog.Any(constants.Error, readErr))
		return nil, fmt.Errorf(constants.StrErrorReadingResponseBody, readErr)
	}

	if resp.StatusCode < constants.HTTPStatusOKMin || resp.StatusCode >= constants.HTTPStatusOKMax {
		logger.ErrorContext(ctx, constants.StrNonSuccessHTTPResponse,
			slog.Int(constants.StrStatusCode, resp.StatusCode),
		)

		return nil, fmt.Errorf(constants.StrWrappedErrorWithBody,
			constants.ErrNonSuccessStatus, resp.Status, string(bodyBytes))
	}

	var result map[string]any
	if ok := json.Unmarshal(bodyBytes, &result); ok != nil {
		logger.ErrorContext(ctx, constants.StrFailedToDecodeResponseJSON, slog.Any(constants.Error, ok))
		return nil, fmt.Errorf(constants.StrErrorDecodingResponse, ok)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrFailedToCreateHTTPRequest, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorCreatingRequest, err)
	}

	logger.InfoContext(ctx, "GET request successful",
		slog.String("result", string(resultJSON)),
	)

	return result, nil
}

func HandlePUTHttpRequest(ctx context.Context, logger *slog.Logger,
	url string, payload any, method string,
) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString(constants.StrTAZAPAYAuthToken),
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrFailedToCreateHTTPRequest, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorCreatingRequest, err)
	}

	logger.InfoContext(ctx, "Sending PUT request",
		slog.Any("payload", payload),
	)

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.ErrorContext(ctx, constants.StrFailedToCreateHTTPRequest, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorCreatingRequest, err)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrHTTPRequestFailed, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorMakingRequest, err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logger.ErrorContext(ctx, constants.StrFailedToReadResponseBody, slog.Any(constants.Error, readErr))
		return nil, fmt.Errorf(constants.StrErrorReadingResponseBody, readErr)
	}

	if resp.StatusCode < constants.HTTPStatusOKMin ||
		resp.StatusCode >= constants.HTTPStatusOKMax {

		logger.ErrorContext(ctx, constants.StrNonSuccessHTTPResponse,
			slog.Int(constants.StrStatusCode, resp.StatusCode),
		)

		return nil, fmt.Errorf(constants.StrWrappedErrorWithBody, constants.ErrNonSuccessStatus,
			resp.Status, string(bodyBytes))
	}

	var result map[string]any
	if ok := json.Unmarshal(bodyBytes, &result); ok != nil {
		logger.ErrorContext(ctx, constants.StrFailedToDecodeResponseJSON,
			slog.Any(constants.Error, ok))

		return nil, fmt.Errorf(constants.StrErrorDecodingResponse, ok)
	}

	logger.InfoContext(ctx, "PUT request successful")

	return result, nil
}

func HandleDELETEHttpRequest(ctx context.Context, logger *slog.Logger,
	url, method string,
) (map[string]any, error) {
	headers := map[string]string{
		constants.HeaderAccept:        constants.AcceptJSON,
		constants.HeaderAuthorization: constants.AuthSchemeBasic + viper.GetString(constants.StrTAZAPAYAuthToken),
	}

	logger.InfoContext(ctx, "Sending DELETE request")

	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrFailedToCreateHTTPRequest, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorCreatingRequest, err)
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, constants.StrHTTPRequestFailed, slog.Any(constants.Error, err))
		return nil, fmt.Errorf(constants.StrErrorMakingRequest, err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logger.ErrorContext(ctx, constants.StrFailedToReadResponseBody, slog.Any(constants.Error, readErr))
		return nil, fmt.Errorf(constants.StrErrorReadingResponseBody, readErr)
	}

	if resp.StatusCode < constants.HTTPStatusOKMin || resp.StatusCode >= constants.HTTPStatusOKMax {
		logger.ErrorContext(ctx, constants.StrNonSuccessHTTPResponse,
			slog.Int(constants.StrStatusCode, resp.StatusCode),
		)

		return nil, fmt.Errorf(constants.StrWrappedErrorWithBody, constants.ErrNonSuccessStatus,
			resp.Status, string(bodyBytes))
	}

	var result map[string]any
	if ok := json.Unmarshal(bodyBytes, &result); ok != nil {
		logger.ErrorContext(ctx, constants.StrFailedToDecodeResponseJSON, slog.Any(constants.Error, ok))
		return nil, fmt.Errorf(constants.StrErrorDecodingResponse, ok)
	}

	logger.InfoContext(ctx, "DELETE request successful")

	return result, nil
}


// AuthHeaderHTTPContextFunc is a function that adds the authorization header to the context from the incoming requests.
func AuthHeaderHTTPContextFunc(ctx context.Context, r *http.Request) context.Context {
    authHeader := r.Header.Get(constants.HeaderAuthorization)
    var basicToken string
    if after, ok :=strings.CutPrefix(authHeader, "Bearer Basic "); ok  {
        basicToken = after
    } else if after, ok :=strings.CutPrefix(authHeader, "Basic "); ok  {
        basicToken = after
    }
    viper.Set(constants.StrTAZAPAYAuthToken, basicToken)
    return r.Context()
}

