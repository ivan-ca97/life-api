package handler

type createShareRequest struct {
	GranteeEmail string `json:"grantee_email"`
	ResourceType string `json:"resource_type"`
	CanWrite     bool   `json:"can_write"`
}
