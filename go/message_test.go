package verbeux_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	verbeux "github.com/verbeux-ai/generative-sdk/go"
)

func TestSendMessage(t *testing.T) {
	client := verbeux.NewClient(
		verbeux.WithApiKey(""),
		verbeux.WithBaseUrl(""),
	)

	resCreateSession, err := client.CreateSession(verbeux.SessionCreateRequest{
		AssistantId: 66,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resCreateSession)

	result, err := client.SendMessage(verbeux.SendMessageRequest{
		ID:      resCreateSession.ID,
		Message: "Ola",
		ClientData: map[string]string{
			"phone": "phone testing",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
}
