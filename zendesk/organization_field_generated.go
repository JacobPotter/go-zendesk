// Code generated by Script. DO NOT EDIT.
// Source: script/codegen/main.go
//
// Generated by this command:
//
//	go run script/codegen/main.go

package zendesk

import (
	"context"
	"github.com/JacobPotter/go-zendesk/internal/client"
)

func (z *Client) GetOrganizationFieldsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[OrganizationField] {
	return &Iterator[OrganizationField]{
		CommonOptions: opts.CommonOptions,
		pageSize:      opts.PageSize,
		hasMore:       true,
		isCBP:         opts.IsCBP,
		pageAfter:     "",
		pageIndex:     1,
		ctx:           ctx,
		obpFunc:       z.GetOrganizationFieldsOBP,
		cbpFunc:       z.GetOrganizationFieldsCBP,
	}
}

func (z *Client) GetOrganizationFieldsOBP(ctx context.Context, opts *OBPOptions) ([]OrganizationField, Page, error) {
	var data struct {
		OrganizationFields []OrganizationField `json:"organization_fields"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &OBPOptions{}
	}

	u, err := client.AddOptions("/organization_fields.json", tmp)

	if err != nil {
		return nil, Page{}, err
	}

	err = client.GetData(z, ctx, u, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.OrganizationFields, data.Page, nil
}

func (z *Client) GetOrganizationFieldsCBP(ctx context.Context, opts *CBPOptions) ([]OrganizationField, client.CursorPaginationMeta, error) {
	var data struct {
		OrganizationFields []OrganizationField         `json:"organization_fields"`
		Meta               client.CursorPaginationMeta `json:"meta"`
	}

	tmp := opts
	if tmp == nil {
		tmp = &CBPOptions{}
	}

	u, err := client.AddOptions("/organization_fields.json", tmp)

	if err != nil {
		return nil, data.Meta, err
	}

	err = client.GetData(z, ctx, u, &data)
	if err != nil {
		return nil, data.Meta, err
	}
	return data.OrganizationFields, data.Meta, nil
}
