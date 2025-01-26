package zendesk

import (
	"context"
	"github.com/JacobPotter/go-zendesk/client"
)

func (z *Client) GetDynamicContentItemsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[DynamicContentItem] {
	return &Iterator[DynamicContentItem]{
		CommonOptions: opts.CommonOptions,
		pageSize:      opts.PageSize,
		hasMore:       true,
		isCBP:         opts.IsCBP,
		pageAfter:     "",
		pageIndex:     1,
		ctx:           ctx,
		obpFunc:       z.GetDynamicContentItemsOBP,
		cbpFunc:       z.GetDynamicContentItemsCBP,
	}
}

func (z *Client) GetDynamicContentItemsOBP(ctx context.Context, opts *OBPOptions) ([]DynamicContentItem, Page, error) {
	var data struct {
		DynamicContentItems []DynamicContentItem `json:"items"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &OBPOptions{}
	}

	u, err := client.AddOptions("/dynamic_content/items.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	err = client.GetData(z, ctx, u, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.DynamicContentItems, data.Page, nil
}

func (z *Client) GetDynamicContentItemsCBP(ctx context.Context, opts *CBPOptions) ([]DynamicContentItem, client.CursorPaginationMeta, error) {
	var data struct {
		DynamicContentItems []DynamicContentItem        `json:"items"`
		Meta                client.CursorPaginationMeta `json:"meta"`
	}

	tmp := opts
	if tmp == nil {
		tmp = &CBPOptions{}
	}

	u, err := client.AddOptions("/dynamic_content/items.json", tmp)
	if err != nil {
		return nil, data.Meta, err
	}

	err = client.GetData(z, ctx, u, &data)
	if err != nil {
		return nil, data.Meta, err
	}
	return data.DynamicContentItems, data.Meta, nil
}
