package services

// filter condition
type FilterCondition struct {
	Field    string         `json:"field,omitempty" form:"field"`
	Fields   []string       `json:"fields,omitempty"  form:"fields"` // For fulltext operator
	Operator FilterOperator `json:"operator" form:"operator"`
	Value    interface{}    `json:"value" form:"value"`
	Values   []interface{}  `json:"values,omitempty" form:"values"` // For between/in operators
}

// Query structure
type UrlQuery struct {
	Limit     int               `json:"limit" form:"limit"`
	Offset    int               `json:"offset" form:"offset"`
	Filters   []FilterCondition `json:"filters" form:"filters"`
	SortBy    string            `json:"sort_by" form:"sort_by"`
	SortOrder string            `json:"sort_order" form:"sort_order"`
	UID       string            `json:"uid" form:"uid" binding:"required"`
}

// supported filter operators
type FilterOperator string

type ExistsRequest struct {
	Email string `json:"email" binding:"required"`
}

type CreateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	UID       string  `json:"uid" binding:"required"`
	Email     string  `json:"email" binding:"required"`
}

type CreateRequest struct {
	URL string `json:"url" binding:"required,url"`
	UID string `json:"uid" binding:"required"`
}

type GetUserUrlRequest struct {
	UID    string `json:"uid" binding:"required"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ClickEvent struct {
	ShortCode string `json:"shortCode"`
	UserId    string `json:"userId"`
}

type GetUserRequest struct {
	UID string `json:"uid" binding:"required"`
}
