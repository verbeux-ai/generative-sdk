package verbeux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (s *Client) CreateSession(request SessionCreateRequest) (*SessionCreateResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	requestURL.Path = SessionRoute
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("POST", requestURL.String(), bytes.NewBuffer(body))
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
