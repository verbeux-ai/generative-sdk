package verbeux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (s *Client) SendMessage(request SendMessageRequest) (*SendMessageResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	requestURL.Path = fmt.Sprintf("%s/%s", SessionRoute, request.ID)
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("PUT", requestURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("api-key", s.apiKey)

	res, err := s.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	returnedBody := SendMessageResponse{}
	if res.Body != nil {
		if err := json.NewDecoder(res.Body).Decode(&returnedBody); err != nil {
			return nil, err
		}
	}

	if res.StatusCode > 399 {
		return nil, fmt.Errorf("%w: %v", ErrSendMessage, returnedBody)
	}

	return &returnedBody, nil
}
