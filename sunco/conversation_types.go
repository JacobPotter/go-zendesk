package sunco

import "time"

// Participants in the Conversation
type Participants struct {
	UserID             string `json:"userId"`
	SubscribeSDKClient bool   `json:"subscribeSDKClient"`
}

// ConversationResponse is standard response for endpoints in ConversationsAPI
type ConversationResponse struct {
	Conversation Conversation `json:"conversation"`
}

// ActiveSwitchboardIntegration is the current switchboard that is connected to the Conversation
type ActiveSwitchboardIntegration struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	IntegrationID   string `json:"integrationId"`
	IntegrationType string `json:"integrationType"`
}

// PendingSwitchboardIntegration is the current switchboard that is connected to the Conversation
type PendingSwitchboardIntegration struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	IntegrationID   string `json:"integrationId"`
	IntegrationType string `json:"integrationType"`
}

type ConversationType string

const (
	Personal ConversationType = "personal"
	SDKGroup ConversationType = "sdkGroup"
)

// Conversation represents a sunco conversation
type Conversation struct {
	ID                            string                        `json:"id,omitempty"`
	Type                          ConversationType              `json:"type"`
	Metadata                      any                           `json:"metadata,omitempty"`
	Participants                  []Participants                `json:"participants,omitempty"`
	ActiveSwitchboardIntegration  ActiveSwitchboardIntegration  `json:"activeSwitchboardIntegration,omitempty"`
	PendingSwitchboardIntegration PendingSwitchboardIntegration `json:"pendingSwitchboardIntegration,omitempty"`
	IsDefault                     bool                          `json:"isDefault,omitempty"`
	DisplayName                   string                        `json:"displayName,omitempty"`
	Description                   string                        `json:"description,omitempty"`
	IconURL                       string                        `json:"iconUrl,omitempty"`
	BusinessLastRead              time.Time                     `json:"businessLastRead,omitempty"`
	LastUpdatedAt                 time.Time                     `json:"lastUpdatedAt,omitempty"`
	CreatedAt                     time.Time                     `json:"createdAt,omitempty"`
}
