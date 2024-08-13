package verbeux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

func (s *Client) CreateSession(request SessionCreateRequest) (*SessionResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	requestURL.Path = SessionRoute
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Post(requestURL.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	returnedBody := SessionResponse{}
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
