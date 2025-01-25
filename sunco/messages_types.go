package sunco

import "time"

type AuthorType string

const (
	BusinessAuthor AuthorType = "business"
	UserAuthor     AuthorType = "user"
)

// Author of sunco Message
type Author struct {
	Type      AuthorType `json:"type"`
	AvatarUrl string     `json:"avatarUrl,omitempty"`
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
}
