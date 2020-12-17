package docs

import "github.com/90lantran/github-star/internal/model"

// swagger:route POST /foobar foobar-tag idOfFoobarEndpoint
// Foobar does some amazing stuff.
// responses:
//   200: foobarResponse

// This text will appear as description of your response body.
// swagger:response foobarResponse
type foobarResponseWrapper struct {
	// in:body
	Body model.Response
}

// swagger:parameters idOfFoobarEndpoint
type foobarParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body model.Request
}
