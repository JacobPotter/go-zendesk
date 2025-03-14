// Code generated by Script. DO NOT EDIT.
// Source: script/codegen/main.go
//
// Generated by this command:
//
//	go run script/codegen/main.go

package zendesk

import (
	"context"
	"github.com/JacobPotter/go-zendesk/client"
)

func (z *Client) GetAllTicketAuditsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[TicketAudit] {
	return &Iterator[TicketAudit]{
		CommonOptions: opts.CommonOptions,
		pageSize:      opts.PageSize,
		hasMore:       true,
		isCBP:         opts.IsCBP,
		pageAfter:     "",
		pageIndex:     1,
		ctx:           ctx,
		obpFunc:       z.GetAllTicketAuditsOBP,
		cbpFunc:       z.GetAllTicketAuditsCBP,
	}
}

func (z *Client) GetAllTicketAuditsOBP(ctx context.Context, opts *OBPOptions) ([]TicketAudit, Page, error) {
	var data struct {
		TicketAudits []TicketAudit `json:"audits"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &OBPOptions{}
	}

	u, err := client.AddOptions("/ticket_audits.json", tmp)

	if err != nil {
		return nil, Page{}, err
	}

	err = client.GetData(z, ctx, u, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.TicketAudits, data.Page, nil
}

func (z *Client) GetAllTicketAuditsCBP(ctx context.Context, opts *CBPOptions) ([]TicketAudit, client.CursorPaginationMeta, error) {
	var data struct {
		TicketAudits []TicketAudit               `json:"audits"`
		Meta         client.CursorPaginationMeta `json:"meta"`
	}

	tmp := opts
	if tmp == nil {
		tmp = &CBPOptions{}
	}

	u, err := client.AddOptions("/ticket_audits.json", tmp)

	if err != nil {
		return nil, data.Meta, err
	}

	err = client.GetData(z, ctx, u, &data)
	if err != nil {
		return nil, data.Meta, err
	}
	return data.TicketAudits, data.Meta, nil
}
