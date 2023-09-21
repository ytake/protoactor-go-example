package source

import "github.com/asynkron/protoactor-go/persistence"

type InMemory struct {
	providerState persistence.ProviderState
}

func NewInMemory(snapshotInterval int) *InMemory {
	return &InMemory{
		providerState: persistence.NewInMemoryProvider(snapshotInterval),
	}
}

func (p *InMemory) GetState() persistence.ProviderState {
	return p.providerState
}
