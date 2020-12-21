package docs

import "github.com/90lantran/github-star/internal/model"

// swagger:route POST /get-stars get-stars idOGetStarsEndpoint
// Return number of github stars for a list of originazation/repository.
// responses:
//   200: getStartsGoodResponse
//   400: invalidResponseWrapper
//	 500: internalServerResponseWrapper

// swagger:parameters idOGetStarsEndpoint
type request struct {
	// Each element of the array must be in form organizatio/repository.
	// in:body
	// required:true
	// example: {"input":["tinygo-org/tinygo-site"]}
	Body model.Request
}

// The reponse shows valid and invalid organization/repository and number of stars for valid ones.
// Example: {"totalStars":19,"invalidRepos":["tingo-org/homebrew-tools","tiygo-org/tinyfont","tinygo-org/tinyfnt"],"validRepos":{"tinygo-org/tinyfont":19}}"
// swagger:response getStartsGoodResponse
type getStartsGoodResponseWrapper struct {
	// in:body
	Body model.Response
}

// The reponse is invalid request message in case the request is not in the right format.
// Example: {"error":"invalid request"}
// swagger:response invalidResponseWrapper
type invalidResponseWrapper struct {
	// in:body
	Body model.Response
}

// The reponse is internal server when server cannot connect to github API.
// Example: {"error":"cannot connect to github"}
// swagger:response internalServerResponseWrapper
type internalServerResponseWrapper struct {
	// in:body
	Body model.Response
}

// swagger:route GET /health health idOHealthEndpoint
// Indication of health of the API.
// responses:
//   200: healthResponseWrapper

// Reponse shows the server is up and ready for use.
// example: "{"message":"the server is up!"}"
// swagger:response healthResponseWrapper
type healthResponseWrapper struct {
	// in:body
	Body model.HealthResponse
}
