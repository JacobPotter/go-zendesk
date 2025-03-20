package sunco

import (
	"time"
)

type User struct {
	ExternalId   string    `json:"externalId"`
	SignedUpAt   time.Time `json:"signedUpAt"`
	ToBeRetained bool      `json:"toBeRetained,omitempty"`
	Profile      struct {
		GivenName string `json:"givenName,omitempty"`
		Surname   string `json:"surname,omitempty"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatarUrl,omitempty"`
		Locale    string `json:"locale,omitempty"`
	} `json:"profile"`
	Metadata interface{} `json:"metadata,omitempty"`
}

type UserResponse struct {
	User User `json:"user"`
}
