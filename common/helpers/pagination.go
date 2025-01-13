package helpers

type PaginationData struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	// The total number of entities that can be retrieved
	Total int `json:"total"`
}

type ResponseEntities struct {
	Data       interface{}     `json:"data,omitempty"`
	Message    string          `json:"message"`
	Pagination *PaginationData `json:"pagination,omitempty"`
}

type OrderedOffsetEntitiesQuery struct {
	OrderBy    string         `json:"order_by"`
	Descending bool           `json:"descending"`
	Term       string         `json:"term,omitempty"`
	Pagination PaginationData `json:"pagination,omitempty"`
}
