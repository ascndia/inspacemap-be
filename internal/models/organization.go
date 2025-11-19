package models

type OrgShortInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type OrganizationFilter struct {
	Name     *string `json:"name,omitempty"`
	Domain   *string `json:"domain,omitempty"`
	Slug     *string `json:"slug,omitempty"`
	Website  *string `json:"website,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type OrganizationQuery struct {
	OrganizationFilter
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type OrganizationQueryCursor struct {
	OrganizationFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}