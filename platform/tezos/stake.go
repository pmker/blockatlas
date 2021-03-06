package tezos

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
	services "github.com/trustwallet/blockatlas/services/assets"
	"math"
	"strconv"
)

func (p *Platform) GetDelegations(address string) (blockatlas.DelegationsPage, error) {
	account, err := p.stakeClient.GetAccount(address)
	if err != nil {
		return nil, err
	}
	if !account.IsDelegated {
		return make(blockatlas.DelegationsPage, 0), nil
	}
	validators, err := services.GetValidatorsMap(p)
	if err != nil {
		return nil, err
	}
	return NormalizeDelegation(account, validators)
}

func NormalizeDelegation(account Account, validators blockatlas.ValidatorMap) (blockatlas.DelegationsPage, error) {
	validator, ok := validators[account.Delegate]
	if !ok {
		return nil, errors.E("Validator not found",
			errors.Params{"Address": account.Address, "Delegate": account.Delegate, "Balance": account.Balance})
	}
	balance := removeDecimals(account.Balance)
	return blockatlas.DelegationsPage{
		{
			Delegator: validator,
			Value:     balance,
			Status:    blockatlas.DelegationStatusActive,
		},
	}, nil
}

func (p *Platform) GetValidators() (blockatlas.ValidatorPage, error) {
	results := make(blockatlas.ValidatorPage, 0)
	validators, err := p.rpcClient.GetValidators()

	if err != nil {
		return results, err
	}

	for _, v := range validators {
		results = append(results, normalizeValidator(v))
	}

	return results, nil
}

func (p *Platform) GetDetails() blockatlas.StakingDetails {
	return getDetails()
}

func (p *Platform) UndelegatedBalance(address string) (string, error) {
	account, err := p.stakeClient.GetAccount(address)
	if err != nil {
		return "0", err
	}
	return removeDecimals(account.Balance), nil
}

func getDetails() blockatlas.StakingDetails {
	return blockatlas.StakingDetails{
		Reward:        blockatlas.StakingReward{Annual: Annual},
		MinimumAmount: "0",
		LockTime:      0,
		Type:          blockatlas.DelegationTypeDelegate,
	}
}

func normalizeValidator(v Validator) (validator blockatlas.Validator) {
	// How to calculate Tezos APR? I have no idea. Tezos team does not know either. let's assume it's around 7% - no way to calculate in decentralized manner
	// Delegation rewards distributed by the validators manually, it's up to them to do it.

	return blockatlas.Validator{
		Status:  true,
		ID:      v.Address,
		Details: getDetails(),
	}
}

func removeDecimals(volume float64) string {
	decimals := coin.Coins[coin.XTZ].Decimals
	d := math.Pow10(int(decimals))
	v := volume * d
	return strconv.Itoa(int(v))
}
