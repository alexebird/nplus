package nomad

import (
	"github.com/hashicorp/nomad/api"
)

type NomadClient struct {
	api *api.Client
}

func Client() (*NomadClient, error) {
	cfg := api.DefaultConfig()

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &NomadClient{api: client}, nil
}
