package verbeux

const defaultURL = "https://generative-api.verbeux.com.br"

// ChatMessageType types
const (
	ChatMessageTypeAI     ChatMessageType = "ai"
	ChatMessageTypeHuman  ChatMessageType = "human"
	ChatMessageTypeSystem ChatMessageType = "system"
	ChatMessageTypeTool   ChatMessageType = "tool"
)

// GenerativeActionType types
const (
	GenActionWebhook GenerativeActionType = "WEBHOOK"
)

// ResponseType types
const (
	ReferenceResponseType ResponseType = "reference"
	TextResponseType      ResponseType = "text"
	TriggerResponseType   ResponseType = "trigger"
)

// Routes to request in genai
const (
	SessionRoute = "session"
)