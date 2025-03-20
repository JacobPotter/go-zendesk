package sunco

import "time"

type AuthorType string

const (
	BusinessAuthor AuthorType = "business"
	UserAuthor     AuthorType = "user"
)

// Author of sunco Message
type Author struct {
	Type           AuthorType `json:"type"`
	AvatarUrl      string     `json:"avatarUrl,omitempty"`
	UserId         string     `json:"userId,omitempty"`
	UserExternalId string     `json:"userExternalId,omitempty"`
	DisplayName    string     `json:"displayName"`
}

// Content of Message
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Message object to be sent in Conversation
type Message struct {
	Id       string    `json:"id,omitempty"`
	Received time.Time `json:"received,omitempty"`
	Author   Author    `json:"author"`
	Content  Content   `json:"content"`
	Source   struct {
		Type string `json:"type"`
	} `json:"source,omitempty"`
}

// MessageResponse is the standard response for MessagesAPI endpoints
type MessageResponse struct {
	Messages []Message `json:"messages"`
	Meta     *Meta     `json:"meta,omitempty"`
	Links    *Links    `json:"links,omitempty"`
}

type Meta struct {
	HasMore      bool   `json:"hasMore"`
	AfterCursor  string `json:"afterCursor"`
	BeforeCursor string `json:"beforeCursor"`
}

type Links struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}
