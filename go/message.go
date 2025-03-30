package verbeux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"strings"
)

func (s *Client) SendMessage(request SendMessageRequest) (*SendMessageResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	requestURL.Path = fmt.Sprintf("%s/%s", SessionRoute, request.ID)

	writer, body, err := s.buildSendMessageBody(request.SendMessageBody)
	if err != nil {
		return nil, err
	}

	if err := writer.WriteField("id", request.ID); err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("PUT", requestURL.String(), body)
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

func (s *Client) buildSendMessageBody(request SendMessageBody) (*multipart.Writer, *bytes.Buffer, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if len(request.Message) > 0 {
		err := writer.WriteField("message", request.Message)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(request.Channel) > 0 {
		err := writer.WriteField("channel", string(request.Channel))
		if err != nil {
			return nil, nil, err
		}
	}

	if len(request.ClientData) > 0 {
		clientDataJSON, err := json.Marshal(request.ClientData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to serialize clientData: %w", err)
		}
		if err := writer.WriteField("clientData", string(clientDataJSON)); err != nil {
			return nil, nil, err
		}
	}

	if request.ForceTriggerCall != nil {
		err := writer.WriteField("force_trigger_call", fmt.Sprintf("%t", *request.ForceTriggerCall))
		if err != nil {
			return nil, nil, err
		}
	}
	if request.Debug != nil {
		err := writer.WriteField("debug", fmt.Sprintf("%t", *request.Debug))
		if err != nil {
			return nil, nil, err
		}
	}
	if request.IgnoreTriggerResponse != nil {
		err := writer.WriteField("ignore_trigger_response", fmt.Sprintf("%t", *request.IgnoreTriggerResponse))
		if err != nil {
			return nil, nil, err
		}
	}
	if len(request.FilesURL) > 0 {
		fileURLJSON, err := json.Marshal(request.FilesURL)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to serialize filesURL: %w", err)
		}
		if err := writer.WriteField("files_url", string(fileURLJSON)); err != nil {
			return nil, nil, err
		}
	}

	for i, file := range request.Files {
		fieldName := fmt.Sprintf("file_%d", i)
		if file.FieldName != "" {
			fieldName = file.FieldName
		}

		// Infer the MIME type based on file extension
		mimeType := file.MimeType // Default type
		if file.FileName != "" && mimeType == "" {
			mimeType = mime.TypeByExtension(filepath.Ext(file.FileName))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
		}

		// Create a form file with the correct headers
		err := createFormFileWithContentType(writer, fieldName, file.FileName, mimeType, file.Reader)
		if err != nil {
			return nil, nil, err
		}
	}

	return writer, &body, nil
}

func (s *Client) OneShot(request OneShotRequest) (*SendMessageResponse, error) {
	requestURL, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	requestURL.Path = fmt.Sprintf("%s/one-shot", SessionRoute)
	writer, body, err := s.buildSendMessageBody(request.SendMessageBody)
	if err != nil {
		return nil, err
	}

	if len(request.SeedSession) > 0 {
		if err := writer.WriteField("seed_session", request.SeedSession); err != nil {
			return nil, err
		}
	}
	if err := writer.WriteField("assistant_id", fmt.Sprintf("%d", request.AssistantId)); err != nil {
		return nil, err
	}

	if len(request.History) > 0 {
		historyDataJson, err := json.Marshal(request.History)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize history: %w", err)
		}
		if err := writer.WriteField("history", string(historyDataJson)); err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("POST", requestURL.String(), body)
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

// Custom function to create a form file with the correct Content-Type
func createFormFileWithContentType(w *multipart.Writer, fieldname, filename, contentType string, reader io.Reader) error {
	// Create the form-data header with the proper Content-Disposition and Content-Type
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)

	// Create the part with the custom headers
	part, err := w.CreatePart(h)
	if err != nil {
		return err
	}

	// Copy the file data into the part
	_, err = io.Copy(part, reader)
	return err
}

// Helper function to escape quotes in field names and file names
func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
