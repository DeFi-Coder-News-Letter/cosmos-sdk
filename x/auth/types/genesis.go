package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	Params   Params               `json:"params" yaml:"params"`
	Accounts []ValidatableAccount `json:"accounts" yaml:"accounts"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, accounts []ValidatableAccount) GenesisState {
	return GenesisState{params, accounts}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), []ValidatableAccount{}) // TODO could just create the struct and not initialize the account field.

}

// ValidatableAccount is an interface that accounts are unmarshalled into when the genesis file is read for the purpose of validating.
// It exists to avoid having to add a Validate method to the Account interface.
type ValidatableAccount interface {
	exported.Account // TODO Is this even needed? Could almost create a generic Validatable interface
	Validate() error
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	// Validate params
	if data.Params.TxSigLimit == 0 {
		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
	}
	if data.Params.SigVerifyCostED25519 == 0 {
		return fmt.Errorf("invalid ED25519 signature verification cost: %d", data.Params.SigVerifyCostED25519)
	}
	if data.Params.SigVerifyCostSecp256k1 == 0 {
		return fmt.Errorf("invalid SECK256k1 signature verification cost: %d", data.Params.SigVerifyCostSecp256k1)
	}
	if data.Params.MaxMemoCharacters == 0 {
		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
	}
	if data.Params.TxSizeCostPerByte == 0 {
		return fmt.Errorf("invalid tx size cost per byte: %d", data.Params.TxSizeCostPerByte)
	}

	// Validate accounts
	addrMap := make(map[string]bool, len(data.Accounts))
	for _, acc := range data.Accounts {

		// check for duplicated accounts
		addrStr := acc.GetAddress().String() // TODO why string?
		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}
		addrMap[addrStr] = true

		// check account specific validation
		if err := acc.Validate(); err != nil {
			return fmt.Errorf("invalid account found in genesis state; address: %s, error: %s", addrStr, err.Error())
		}

	}

	return nil
}
