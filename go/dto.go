package verbeux

import (
	"io"
	"time"
)

type SessionCreateRequest struct {
	SessionHistory
	SessionAssistantID
	ClientDataBody
}

type SessionAssistantID struct {
	AssistantId int `json:"assistant_id"`
}

type SessionHistory struct {
	History []HistoryRecord `json:"history"`
}

type SessionUpdateRequest struct {
	SessionID                 string          `json:"-"`
	CurrentConversation       []HistoryRecord `json:"current_conversation,omitempty"`
	CurrentConversationAppend []HistoryRecord `json:"current_conversation_append,omitempty"`
	RestrictedByContext       *bool           `json:"restricted_by_context,omitempty"`
	ClientDataBody
}

type SessionCreateResponse struct {
	ID                  string              `json:"id"`
	AssistantID         uint                `json:"assistant_id"`
	CurrentConversation []Message           `json:"current_conversation"`
	Description         string              `json:"description"`
	CompanyID           uint                `json:"company_id"`
	RestrictedByContext bool                `json:"restricted_by_context"`
	GenerativeTriggers  []GenerativeTrigger `json:"integration_calls"`
	ClientDataBody

	// Message is returned when has an error
	Message interface{} `json:"message"`
}

type ClientDataBody struct {
	ClientData map[string]string `json:"client_data,omitempty"`
}

type OneShotRequest struct {
	SendMessageBody
	SessionHistory
	SessionAssistantID
	SeedSession string `json:"seed_session,omitempty" form:"seed_session" validate:"omitempty,uuid"`
}

type Channel string

type SendMessageRequest struct {
	ID string `json:"id"`
	SendMessageBody
}

type SendMessageBody struct {
	Message               string           `json:"message"`
	Channel               Channel          `json:"channel"`
	Files                 []FileAttachment `json:"files"`
	FilesURL              []string         `json:"files_url" validate:"omitempty,max=3"`
	Debug                 bool             `json:"debug"`
	IgnoreTriggerResponse bool             `json:"ignore_trigger_response"`
	ForceTriggerCall      bool             `json:"force_trigger_call"`
	Copilot               bool             `json:"copilot"`
	ClientDataBody
}

type FileAttachment struct {
	FieldName string
	FileName  string
	MimeType  string
	Reader    io.Reader
}

type SendMessageResponse struct {
	ID             string                       `json:"id"`
	Response       []SendMessageResponseContent `json:"response"`
	IsAnythingElse bool                         `json:"isAnythingElse"`

	// Message is returned when has an error
	Message interface{} `json:"message"`
}

type SendMessageResponseContent struct {
	// TODO: use interface on data
	Data any          `json:"data"`
	Type ResponseType `json:"type"`
}

type ResponseType string

type HistoryRecord struct {
	Content []any  `json:"content"`
	Role    string `json:"role"`
}

type HistoryRecordContentText struct {
	Text string `json:"text"`
}

type ChatMessageType string

type PresetModelType string

type GenerativeActionType string

type Message struct {
	Role ChatMessageType `json:"role"`
	// TODO: change from undefined interface to defined interface
	Content []interface{} `json:"content"`
}

type GenerativeTrigger struct {
	ID        uint      `json:"id" query:"id"`
	CreatedAt time.Time `json:"createdAt" query:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" query:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt" query:"deletedAt"`

	AssistantID       int64              `json:"assistantID" validate:"required" query:"assistantID"`
	IsActive          *bool              `json:"isActive" query:"isActive" gorm:"default:true"`
	GenerativeActions []GenerativeAction `json:"generativeActions" gorm:"constraint:OnDelete:CASCADE;"`
	PresetModel       PresetModelType    `json:"presetModel" query:"presetModel"`

	Args *FunctionCalling `json:"args" gorm:"-"`
}

type GenerativeAction struct {
	ID        uint      `json:"id" query:"id"`
	CreatedAt time.Time `json:"createdAt" query:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" query:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt" query:"deletedAt"`

	Type   GenerativeActionType `json:"type" validate:"required,oneof=WEBHOOK" query:"type"`
	Action interface{}          `json:"action" gorm:"-"`

	GenerativeTriggerID int64              `json:"generativeTriggerID" validate:"required" query:"generativeTriggerID"`
	GenerativeTrigger   *GenerativeTrigger `json:"generativeTrigger" gorm:"-"`
}

type FunctionCalling struct {
	Name        string                     `json:"name" validate:"required"`
	Description string                     `json:"description" validate:"required"`
	Parameters  *FunctionCallingParameters `json:"parameters,omitempty" validate:"omitempty"`
}

type FunctionCallingParameters struct {
	Type       string                                        `json:"type" validate:"required"`
	Properties map[string]FunctionCallingParametersPropriety `json:"properties"`
	Required   []string                                      `json:"required,omitempty"`
}

type FunctionCallingParametersPropriety struct {
	Type        string   `json:"type" validate:"required"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}
