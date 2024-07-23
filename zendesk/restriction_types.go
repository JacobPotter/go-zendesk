package zendesk

type Restriction struct {
	Type string  `json:"type,omitempty"`
	ID   int64   `json:"id,omitempty"`
	IDS  []int64 `json:"ids,omitempty"`
}
