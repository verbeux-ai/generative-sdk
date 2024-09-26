package verbeux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (s *Client) SendMessage(request SendMessageRequest) (*SendMessageResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}
	requestURL.Path = fmt.Sprintf("%s/%s", SessionRoute, request.ID)

	// Cria um buffer para o corpo da requisição
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Adiciona campos ao form data
	err = writer.WriteField("id", request.ID)
	if err != nil {
		return nil, err
	}
	if len(request.Message) > 0 {
		err = writer.WriteField("message", request.Message)
		if err != nil {
			return nil, err
		}
	}

	for i, file := range request.Files {
		fieldName := fmt.Sprintf("file_%d", i)
		if file.FieldName != "" {
			fieldName = file.FieldName
		}
		part, err := writer.CreateFormFile(fieldName, file.FileName)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file.Reader)
		if err != nil {
			return nil, err
		}
	}

	// Finaliza o writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("PUT", requestURL.String(), &body)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", writer.FormDataContentType())
	httpRequest.Header.Set("api-key", s.apiKey)

	res, err := s.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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
