package api

type Config struct {
	TlsCaLocation string
}

type TokenFunc func() (string, error)

type GenerateClusterTokenRequest struct {
	ClusterName string `json:"cluster_name" binding:"required"`
}

type GenerateClusterTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
}

type IntrospectClusterTokenRequest struct{}

type IntrospectClusterTokenResponse struct{}

type AttachRolesToUserRequest struct {
	Roles []Role `json:"roles" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type AttachRolesToUserResponse struct {
	Message string `json:"message"`
}

type SyncRolesRequest struct {
	Users []User `json:"users" binding:"required"`
}

type SyncRolesResponse struct {
	Message string `json:"message"`
}

type User struct {
	Email string `json:"email" binding:"required"`
	Roles []Role `json:"roles"`
}

type Role struct {
	Resource    string `json:"resource" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

type Cluster struct {
	Name string `json:"name"`
}

type Account struct {
	ID       string    `json:"id"`
	Alias    string    `json:"alias"`
	Clusters []Cluster `json:"clusters"`
}
