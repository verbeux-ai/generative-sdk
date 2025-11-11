package verbeux_test

import (
	"context"
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
		SessionAssistantID: verbeux.SessionAssistantID{
			AssistantId: 865,
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
		SessionAssistantID: verbeux.SessionAssistantID{
			AssistantId: 865,
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
		SessionAssistantID: verbeux.SessionAssistantID{
			AssistantId: 865,
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
		SessionAssistantID: verbeux.SessionAssistantID{
			AssistantId: 865,
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
