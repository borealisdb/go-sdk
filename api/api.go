package api

type API interface {
	GenerateClusterToken(request GenerateClusterTokenRequest) (GenerateClusterTokenResponse, error)
}
