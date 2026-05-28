package verbeux_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/require"
	verbeux "github.com/verbeux-ai/generative-sdk/go"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	os.Exit(m.Run())

}
func TestSendMessage(t *testing.T) {
	client := verbeux.NewClient(
		verbeux.WithApiKey(os.Getenv("API_KEY")),
		verbeux.WithBaseUrl(os.Getenv("BASE_URL")),
	)

	resCreateSession, err := client.CreateSession(context.Background(), verbeux.SessionCreateRequest{
		SessionAgentID: verbeux.SessionAgentID{
			AgentId: 865,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resCreateSession)

	result, err := client.SendMessage(context.Background(), verbeux.SendMessageRequest{
		ID: resCreateSession.ID,
		SendMessageBody: verbeux.SendMessageBody{
			Message: "Ola",
			ClientDataBody: verbeux.ClientDataBody{
				ClientData: map[string]string{
					"phone": "phone testing",
				},
			},
			Channel: verbeux.EvolutionChannel,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

func TestSendMessageAudio(t *testing.T) {
	client := verbeux.NewClient(
		verbeux.WithApiKey(os.Getenv("API_KEY")),
		verbeux.WithBaseUrl(os.Getenv("BASE_URL")),
	)

	resCreateSession, err := client.CreateSession(context.Background(), verbeux.SessionCreateRequest{
		SessionAgentID: verbeux.SessionAgentID{
			AgentId: 865,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, resCreateSession)

	fileOpen, err := os.Open("test.mp3")
	require.NoError(t, err)

	file := verbeux.FileAttachment{
		FileName: "test.mp3",
		Reader:   fileOpen,
	}

	result, err := client.SendMessage(context.Background(), verbeux.SendMessageRequest{
		ID: resCreateSession.ID,
		SendMessageBody: verbeux.SendMessageBody{
			ClientDataBody: verbeux.ClientDataBody{
				ClientData: map[string]string{
					"phone": "phone testing",
				},
			},
			Files: []verbeux.FileAttachment{file},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

func TestSendMessageImage(t *testing.T) {
	client := verbeux.NewClient(
		verbeux.WithApiKey(os.Getenv("API_KEY")),
		verbeux.WithBaseUrl(os.Getenv("BASE_URL")),
	)

	resCreateSession, err := client.CreateSession(context.Background(), verbeux.SessionCreateRequest{
		SessionAgentID: verbeux.SessionAgentID{
			AgentId: 865,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, resCreateSession)

	fileOpen, err := os.Open("test.jpeg")
	require.NoError(t, err)

	file := verbeux.FileAttachment{
		FileName: "test.jpeg",
		Reader:   fileOpen,
	}

	result, err := client.SendMessage(context.Background(), verbeux.SendMessageRequest{
		ID: resCreateSession.ID,
		SendMessageBody: verbeux.SendMessageBody{
			Message: "oq tem na imagem?",
			ClientDataBody: verbeux.ClientDataBody{
				ClientData: map[string]string{
					"phone": "phone testing",
				},
			},

			Files: []verbeux.FileAttachment{file},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

func TestSendMessageOneShot(t *testing.T) {
	client := verbeux.NewClient(
		verbeux.WithApiKey(os.Getenv("API_KEY")),
		verbeux.WithBaseUrl(os.Getenv("BASE_URL")),
	)

	fileOpen, err := os.Open("test.jpeg")
	require.NoError(t, err)

	file := verbeux.FileAttachment{
		FileName: "test.jpeg",
		Reader:   fileOpen,
	}

	result, err := client.OneShot(context.Background(), verbeux.OneShotRequest{
		SessionAgentID: verbeux.SessionAgentID{
			AgentId: 865,
		},
		SendMessageBody: verbeux.SendMessageBody{
			Message: "oq tem na imagem?",
			ClientDataBody: verbeux.ClientDataBody{
				ClientData: map[string]string{
					"phone": "phone testing",
				},
			},

			Files: []verbeux.FileAttachment{file},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

func TestDeleteMessagesNoContent(t *testing.T) {
	client := newDeleteMessagesTestClient(t, func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/session/session-123/message", r.URL.Path)
		require.Equal(t, "test-api-key", r.Header.Get("api-key"))

		var request verbeux.DeleteMessagesRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, []string{"msg-1"}, request.IDs)

		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Header:     make(http.Header),
		}, nil
	})

	result, err := client.DeleteMessages(context.Background(), verbeux.DeleteMessagesRequest{
		SessionID: "session-123",
		IDs:       []string{"msg-1"},
	})

	require.NoError(t, err)
	require.Empty(t, result)
}

func TestDeleteMessagesErrorResponse(t *testing.T) {
	client := newDeleteMessagesTestClient(t, func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/session/session-123/message", r.URL.Path)

		return jsonResponse(http.StatusNotFound, map[string]any{
			"message": "message not found",
		})
	})

	result, err := client.DeleteMessages(context.Background(), verbeux.DeleteMessagesRequest{
		SessionID: "session-123",
		IDs:       []string{"msg-1"},
	})

	require.Nil(t, result)
	require.ErrorIs(t, err, verbeux.ErrDeleteMessage)
	require.ErrorContains(t, err, "message not found")
}

func TestDeleteMessagesSuccessResponse(t *testing.T) {
	client := newDeleteMessagesTestClient(t, func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/session/session-123/message", r.URL.Path)

		return jsonResponse(http.StatusOK, []map[string]any{
			{
				"role":    "ai",
				"content": []any{"deleted"},
			},
		})
	})

	result, err := client.DeleteMessages(context.Background(), verbeux.DeleteMessagesRequest{
		SessionID: "session-123",
		IDs:       []string{"msg-1"},
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, verbeux.ChatMessageTypeAI, result[0].Role)
	require.Equal(t, []interface{}{"deleted"}, result[0].Content)
}

func newDeleteMessagesTestClient(t *testing.T, roundTrip roundTripFunc) *verbeux.Client {
	t.Helper()

	return verbeux.NewClient(
		verbeux.WithApiKey("test-api-key"),
		verbeux.WithBaseUrl("http://delete-messages.test"),
		verbeux.WithHttpClient(&http.Client{Transport: roundTrip}),
	)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(statusCode int, body any) (*http.Response, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(payload)),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}, nil
}
