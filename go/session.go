package verbeux

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (s *Client) CreateSession(ctx context.Context, request SessionCreateRequest) (*SessionCreateResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	requestURL.Path = SessionRoute
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", requestURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("api-key", s.apiKey)

	res, err := s.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	returnedBody := SessionCreateResponse{}
	if res.Body != nil {
		if err := json.NewDecoder(res.Body).Decode(&returnedBody); err != nil {
			return nil, err
		}
	}

	if res.StatusCode > 399 {
		return nil, fmt.Errorf("%w: %v", ErrCreateSession, returnedBody)
	}

	return &returnedBody, nil
}

func (s *Client) UpdateSession(ctx context.Context, request SessionUpdateRequest) (*SessionCreateResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	requestURL.Path = fmt.Sprintf("%s/%s", SessionRoute, request.SessionID)
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, "PATCH", requestURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("api-key", s.apiKey)

	res, err := s.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	returnedBody := SessionCreateResponse{}
	if res.Body != nil {
		if err := json.NewDecoder(res.Body).Decode(&returnedBody); err != nil {
			return nil, err
		}
	}

	if res.StatusCode > 399 {
		return nil, fmt.Errorf("%w: %v", ErrCreateSession, returnedBody)
	}

	return &returnedBody, nil
}
