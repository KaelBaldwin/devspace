package cloud

import (
	"context"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

// GrapqhlRequest does a new graphql request and stores the result in the response
func (p *Provider) GrapqhlRequest(request string, vars map[string]interface{}, response interface{}) error {
	graphQlClient := graphql.NewClient(p.Host + GraphqlEndpoint)
	req := graphql.NewRequest(request)

	// Set vars
	if vars != nil {
		for key, val := range vars {
			req.Var(key, val)
		}
	}

	// Get token
	token, err := p.GetToken()
	if err != nil {
		return errors.Wrap(err, "get token")
	}

	// Set token
	req.Header.Set("Authorization", "Bearer "+token)

	// Run the graphql request
	return graphQlClient.Run(context.Background(), req, response)
}
