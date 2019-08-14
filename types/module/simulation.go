package module

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// SimulationManager defines a simulation manager that provides the high level utility
// for managing and executing simulation functionalities for a group of modules
type SimulationManager struct {
	Modules       map[string]AppModule
	StoreDecoders sdk.StoreDecoderRegistry
	ParamChanges  []simulation.ParamChange
}

// NewSimulationManager creates a new SimulationManager object
func NewSimulationManager(moduleMap map[string]AppModule) *SimulationManager {
	return &SimulationManager{
		Modules:       moduleMap,
		StoreDecoders: make(sdk.StoreDecoderRegistry),
		ParamChanges:  []simulation.ParamChange{},
	}
}

// RegisterStoreDecoders registers each of the modules' store decoders into a map
func (sm *SimulationManager) RegisterStoreDecoders() {
	for _, module := range sm.Modules {
		module.RegisterStoreDecoder(sm.StoreDecoders)
	}
}

// GenerateGenesisStates generates a randomized GenesisState for each of the
// registered modules
func (sm *SimulationManager) GenerateGenesisStates(cdc *codec.Codec, r *rand.Rand, genesisState map[string]json.RawMessage) {
	for _, module := range sm.Modules {
		module.GenerateGenesisState(cdc, r, genesisState)
	}
}

// RandomizedSimParamChanges generates randomized contents for creating params change
// proposal transactions
func (sm *SimulationManager) RandomizedSimParamChanges(cdc *codec.Codec, seed int64) {
	r := rand.New(rand.NewSource(seed))

	for _, module := range sm.Modules {
		sm.ParamChanges = append(sm.ParamChanges, module.RandomizedParams(cdc, r)...)
	}
}