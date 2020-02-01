package ontology

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"testing"
)

const (
	srcOntTransfer = `
{
	"TxnType": 209,
	"ConfirmFlag": 1,
	"Fee": "0.010000000",
	"BlockIndex": 2,
	"TransferList": [
		{
			"FromAddress": "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
			"Amount": "2.000000000",
			"ToAddress": "AQ9kzzHNLCcyrPwJuVMrSPgGzqmuQNVwMF",
			"AssetName": "ont"
		}
	],
	"TxnTime": 1556952450,
	"TxnHash": "4804e1be63ebe1715d6b4a039cc9d84b86cde74c8a8c8411578e6dcadc1e5405",
	"Height": 3411115
}`

	srcOngTransfer = `
{
	"TxnType": 209,
	"ConfirmFlag": 1,
	"Fee": "0.010000000",
	"BlockIndex": 2,
	"TransferList": [
		{
			"FromAddress": "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
			"Amount": "2.1455",
			"ToAddress": "ASwdosWb2wH8Y3HYCYJpUJWWD3joyvvYGN",
			"AssetName": "ong"
		}
	],
	"TxnTime": 1555341286,
	"TxnHash": "a483d1d854e47a20692f472d72ff45b9a2bfc542f84dceb3171a48f68ba322cb",
	"Height": 2863855
}`
	srcRewardTransfer = `
{
	"TxnType": 209,
	"ConfirmFlag": 1,
	"Fee": "0.010000000",
	"BlockIndex": 1,
	"TransferList": [
		{
			"FromAddress": "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
			"Amount": "10.000000000",
			"ToAddress": "AQ9kzzHNLCcyrPwJuVMrSPgGzqmuQNVwMF",
			"AssetName": "ong"
		},
		{
			"FromAddress": "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
			"Amount": "0.010000000",
			"ToAddress": "AFmseVrdL9f9oyCzZefL9tG6UbviEH9ugK",
			"AssetName": "ong"
		}
	],
	"TxnTime": 1556952520,
	"TxnHash": "eccbfd040925a22884d87e73f818f30ab42d06046460b86e9a042a1e9cba7561",
	"Height": 3411141
}`
	srcFeeTransfer = `
{
	"TxnType": 209,
	"ConfirmFlag": 1,
	"Fee": "0.010000000",
	"BlockIndex": 1,
	"TransferList": [
		{
			"FromAddress": "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
			"Amount": "0.010000000",
			"ToAddress": "AFmseVrdL9f9oyCzZefL9tG6UbviEH9ugK",
			"AssetName": "ong"
		}
	],
	"TxnTime": 1556952520,
	"TxnHash": "eccbfd040925a22884d87e73f818f30ab42d06046460b86e9a042a1e9cba7561",
	"Height": 3411141
}`
)

var (
	dstOntTransfer = blockatlas.Tx{
		ID:     "4804e1be63ebe1715d6b4a039cc9d84b86cde74c8a8c8411578e6dcadc1e5405",
		Coin:   coin.ONT,
		From:   "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
		To:     "AQ9kzzHNLCcyrPwJuVMrSPgGzqmuQNVwMF",
		Fee:    "10000000",
		Date:   1556952450,
		Type:   "transfer",
		Status: blockatlas.StatusCompleted,
		Block:  3411115,
		Meta: blockatlas.Transfer{
			Value:    "2",
			Symbol:   "ONT",
			Decimals: 0,
		},
	}
	dstOngTransfer = blockatlas.Tx{
		ID:     "a483d1d854e47a20692f472d72ff45b9a2bfc542f84dceb3171a48f68ba322cb",
		Coin:   coin.ONT,
		From:   "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
		To:     "ASwdosWb2wH8Y3HYCYJpUJWWD3joyvvYGN",
		Fee:    "10000000",
		Date:   1555341286,
		Type:   blockatlas.TxNativeTokenTransfer,
		Status: blockatlas.StatusCompleted,
		Block:  2863855,
		Meta: blockatlas.NativeTokenTransfer{
			Name:     "Ontology Gas",
			Symbol:   "ONG",
			TokenID:  "ong",
			Decimals: 9,
			Value:    "2145500000",
			From:     "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
			To:       "ASwdosWb2wH8Y3HYCYJpUJWWD3joyvvYGN",
		},
	}
	dstRewardTransfer = blockatlas.Tx{
		ID:     "eccbfd040925a22884d87e73f818f30ab42d06046460b86e9a042a1e9cba7561",
		Coin:   coin.ONT,
		From:   "AUyL4TZ1zFEcSKDJrjFnD7vsq5iFZMZqT7",
		To:     "AQ9kzzHNLCcyrPwJuVMrSPgGzqmuQNVwMF",
		Fee:    "10000000",
		Date:   1556952520,
		Type:   blockatlas.TxNativeTokenTransfer,
		Status: blockatlas.StatusCompleted,
		Block:  3411141,
		Meta: blockatlas.AnyAction{
			Coin:     coin.Ontology().ID,
			Name:     "Ontology Gas",
			Symbol:   "ONG",
			TokenID:  "ong",
			Decimals: 9,
			Value:    "10000000000",
			Title:    blockatlas.AnyActionClaimRewards,
			Key:      blockatlas.KeyStakeClaimRewards,
		},
	}
)

func TestNormalize(t *testing.T) {
	var tests = []struct {
		name        string
		Transaction string
		AssetName   AssetType
		Expected    blockatlas.Tx
		wantErr     bool
	}{
		{"normalize ont", srcOntTransfer, AssetONT, dstOntTransfer, false},
		{"normalize ong", srcOngTransfer, AssetONG, dstOngTransfer, false},
		{"normalize claim reward", srcRewardTransfer, AssetONG, dstRewardTransfer, false},
		{"normalize fee", srcFeeTransfer, AssetONG, blockatlas.Tx{}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sourceTx Tx
			err := json.Unmarshal([]byte(test.Transaction), &sourceTx)
			assert.Nil(t, err)
			tx, ok := Normalize(&sourceTx, test.AssetName)
			if test.wantErr {
				assert.False(t, ok)
				return
			}
			assert.True(t, ok)
			assert.Equal(t, test.Expected, tx)
		})
	}
}
