package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JacobPotter/go-zendesk/internal/client"
	"time"
)

// UserField is struct for user_field payload
type UserField struct {
	ID                     int64               `json:"id,omitempty"`
	URL                    string              `json:"url,omitempty"`
	Key                    string              `json:"key,omitempty"`
	Type                   string              `json:"type"`
	Title                  string              `json:"title"`
	RawTitle               string              `json:"raw_title,omitempty"`
	Description            string              `json:"description,omitempty"`
	RawDescription         string              `json:"raw_description,omitempty"`
	Position               int64               `json:"position,omitempty"`
	Active                 bool                `json:"active"`
	System                 bool                `json:"system,omitempty"`
	RegexpForValidation    string              `json:"regexp_for_validation,omitempty"`
	Tag                    string              `json:"tag,omitempty"`
	CustomFieldOptions     []CustomFieldOption `json:"custom_field_options,omitempty"`
	CreatedAt              time.Time           `json:"created_at,omitempty"`
	UpdatedAt              time.Time           `json:"updated_at,omitempty"`
	RelationshipTargetType string              `json:"relationship_target_type,omitempty"`
	RelationshipFilter     RelationshipFilter  `json:"relationship_filter,omitempty"`
}

type UserFieldListOptions struct {
	PageOptions
}

type UserFieldAPI interface {
	GetUserFields(ctx context.Context, opts *UserFieldListOptions) ([]UserField, Page, error)
	CreateUserField(ctx context.Context, userField UserField) (UserField, error)
	GetUserField(ctx context.Context, id int64) (UserField, error)
	UpdateUserField(ctx context.Context, id int64, userField UserField) (UserField, error)
	DeleteUserField(ctx context.Context, id int64) error
	GetUserFieldsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[UserField]
	GetUserFieldsOBP(ctx context.Context, opts *OBPOptions) ([]UserField, Page, error)
	GetUserFieldsCBP(ctx context.Context, opts *CBPOptions) ([]UserField, client.CursorPaginationMeta, error)
}

// GetUserFields fetch trigger list
//
// https://developer.zendesk.com/rest_api/docs/support/user_fields#list-user-fields
func (z *Client) GetUserFields(ctx context.Context, opts *UserFieldListOptions) ([]UserField, Page, error) {
	var data struct {
		UserFields []UserField `json:"user_fields"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &UserFieldListOptions{}
	}

	u, err := client.AddOptions("/user_fields.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.UserFields, data.Page, nil
}

// CreateUserField creates new user field
// ref: https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#create-user-field
func (z *Client) CreateUserField(ctx context.Context, userField UserField) (UserField, error) {
	var data, result struct {
		UserField UserField `json:"user_field"`
	}
	data.UserField = userField

	body, err := z.Post(ctx, "/user_fields.json", data)
	if err != nil {
		return UserField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return UserField{}, err
	}
	return result.UserField, nil
}

func (z *Client) GetUserField(ctx context.Context, id int64) (UserField, error) {
	var result struct {
		UserField UserField `json:"user_field"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/user_fields/%d.json", id))

	if err != nil {
		return UserField{}, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return UserField{}, err
	}

	return result.UserField, nil
}

func (z *Client) UpdateUserField(ctx context.Context, id int64, userField UserField) (UserField, error) {
	var result, data struct {
		UserField UserField `json:"user_field"`
	}

	data.UserField = userField

	body, err := z.Put(ctx, fmt.Sprintf("/user_fields/%d.json", id), data)

	if err != nil {
		return UserField{}, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return UserField{}, err
	}

	return result.UserField, nil
}

func (z *Client) DeleteUserField(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/user_fields/%d.json", id))

	if err != nil {
		return err
	}

	return nil
}
