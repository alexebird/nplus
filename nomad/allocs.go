package nomad

import (
	"log"

	"github.com/hashicorp/nomad/api"
)

func (client *NomadClient) Allocations() ([]*api.Allocation, error) {
	allocStubs, _, err := client.api.Allocations().List(&api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	allocs := make([]*api.Allocation, 0)
	c := make(chan *api.Allocation)
	allocCount := len(allocStubs)

	for _, allocStub := range allocStubs {
		go client.singleAllocation(c, allocStub.ID)

	}

	for i := 0; i < allocCount; i++ {
		allocs = append(allocs, <-c)
	}

	return allocs, nil
}

func (client *NomadClient) singleAllocation(c chan *api.Allocation, allocID string) {
	alloc, _, err := client.api.Allocations().Info(allocID, &api.QueryOptions{})
	if err != nil {
		log.Fatal(err)
	}

	c <- alloc
}
