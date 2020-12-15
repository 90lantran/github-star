package constants

const (
	// APIEndpoint is the ALB of API deployment (empty for local)
	APIEndpoint = ""
	// APIHeathEndpoint is the health check url
	APIHeathEndpoint = APIEndpoint + "/health"
	// APIGetStarsEndpoint is the get stars for list of organization/repository url
	APIGetStarsEndpoint = APIEndpoint + "/get-stars"
)
