package verbeux_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	verbeux "github.com/verbeux-ai/generative-sdk/go"
)

func TestSendMessage(t *testing.T) {
	client := verbeux.NewClient(
		verbeux.WithApiKey("82041401-71da-11ef-a26e-42004e494300"),
		verbeux.WithBaseUrl("https://generative-api-dev.verbeux.com.br"),
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
