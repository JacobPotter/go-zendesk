package zendesk

// CustomFieldOption is struct for value of `custom_field_options`
type CustomFieldOption struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Position int64  `json:"position,omitempty"`
	RawName  string `json:"raw_name,omitempty"`
	URL      string `json:"url,omitempty"`
	Value    string `json:"value"`
}

type RelationshipFilterObject struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// RelationshipFilter is struct for value of `relationship_filter`
type RelationshipFilter struct {
	All []RelationshipFilterObject `json:"all"`
	Any []RelationshipFilterObject `json:"any"`
}
