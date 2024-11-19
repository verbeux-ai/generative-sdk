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

	// Create a buffer to hold the multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add form fields
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

		// Infer the MIME type based on file extension
		mimeType := "application/octet-stream" // Default type
		if file.FileName != "" {
			mimeType = mime.TypeByExtension(filepath.Ext(file.FileName))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
		}

		// Create a form file with the correct headers
		err := createFormFileWithContentType(writer, fieldName, file.FileName, mimeType, file.Reader)
		if err != nil {
			return nil, err
		}
	}

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
