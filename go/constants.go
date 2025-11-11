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
	ButtonResponseType    ResponseType = "button"
	ListResponseType      ResponseType = "list"
	ImageResponseType     ResponseType = "image"
	ImageV2ResponseType   ResponseType = "imageV2"
	AudioResponseType     ResponseType = "audio"
)

// Routes to request in genai
const (
	SessionRoute = "session"
)

// Channels
const (
	WhatsappChannel  Channel = "WHATSAPP"
	ZAPIChannel      Channel = "Z_API"
	EvolutionChannel Channel = "EVOLUTION"
	MyzapChannel     Channel = "MY_ZAP"
	RitaChannel      Channel = "RITA"
	APIChannel       Channel = "API"
)
