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

const (
	SessionRoute = "session"
)
