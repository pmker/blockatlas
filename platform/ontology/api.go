package ontology

import (
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/blockatlas/pkg/numbers"
	"strings"
)

type Platform struct {
	client Client
}

func (p *Platform) Init() error {
	p.client = Client{blockatlas.InitClient(viper.GetString("ontology.api"))}
	return nil
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[coin.ONT]
}

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	return p.GetTokenTxsByAddress(address, string(AssetONT))
}

func (p *Platform) GetTokenTxsByAddress(address string, token string) (blockatlas.TxPage, error) {
	txPage, err := p.client.GetTxsOfAddress(address, AssetType(token))
	if err != nil {
		logger.Error(err, "Ontology: Failed to get transactions for address and token",
			logger.Params{
				"address": address,
				"token":   token,
			})
		return blockatlas.TxPage{}, err
	}
	var txs []blockatlas.Tx
	for _, srcTx := range txPage.Result.TxnList {
		tx, ok := Normalize(&srcTx, AssetType(token))
		if !ok {
			continue
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

func (p *Platform) CurrentBlockNumber() (int64, error) {
	block, err := p.client.CurrentBlockNumber()
	if err != nil {
		logger.Error("CurrentBlockNumber", logger.Params{"platform": p.Coin().Symbol, "details": err.Error()})
		return 0, err
	}
	var height int64
	if block.Error != 0 {
		err = errors.E("explorer error")
	}
	if len(block.Result) > 0 {
		height = (int64)(block.Result[0].Height)
	}
	return height, nil
}

func (p *Platform) GetBlockByNumber(num int64) (*blockatlas.Block, error) {
	response, err := p.client.GetBlockByNumber(num)
	if err != nil {
		logger.Error("GetBlockByNumber", logger.Params{"platform": p.Coin().Symbol, "details": err.Error()})
		return nil, err
	}
	var (
		block blockatlas.Block
		txs   []blockatlas.Tx
	)
	if response.Error == 0 {
		block.ID = response.Result.Hash
		block.Number = int64(response.Result.Height)
		for _, txn := range response.Result.TxnList {
			tx := new(blockatlas.Tx)
			tx.ID = txn.TxnHash
			tx.Block = uint64(txn.Height)
			if txn.ConfirmFlag == 1 {
				tx.Status = blockatlas.StatusCompleted
			}
			tx.Date = int64(txn.TxnTime)
			tx.Coin = coin.Ontology().ID
			txs = append(txs, *tx)
		}
		block.Txs = txs
	}

	return &block, nil
}

func Normalize(srcTx *Tx, assetName AssetType) (tx blockatlas.Tx, ok bool) {
	if len(srcTx.TransferList) < 1 {
		return tx, false
	}
	fee := numbers.DecimalExp(srcTx.Fee, 9)
	status := blockatlas.StatusCompleted
	if srcTx.ConfirmFlag != 1 {
		status = blockatlas.StatusFailed
	}
	tx = blockatlas.Tx{
		ID:     srcTx.TxnHash,
		Coin:   coin.ONT,
		Fee:    blockatlas.Amount(fee),
		Date:   srcTx.TxnTime,
		Block:  srcTx.Height,
		Status: status,
	}

	switch assetName {
	case AssetONT:
		return normalizeONT(tx, srcTx.TransferList)
	case AssetONG:
		return normalizeONG(tx, srcTx.TransferList)
	}
	return tx, false
}

func normalizeONT(tx blockatlas.Tx, transfers Transfers) (blockatlas.Tx, bool) {
	transfer := transfers.getTransfer()
	if transfer == nil {
		return tx, false
	}
	i := strings.IndexRune(transfer.Amount, '.')
	value := transfer.Amount[:i]
	tx.From = transfer.FromAddress
	tx.To = transfer.ToAddress
	tx.Type = blockatlas.TxTransfer
	tx.Meta = blockatlas.Transfer{
		Value:    blockatlas.Amount(value),
		Symbol:   coin.Coins[coin.ONT].Symbol,
		Decimals: coin.Coins[coin.ONT].Decimals,
	}
	return tx, true
}

func normalizeONG(tx blockatlas.Tx, transfers Transfers) (blockatlas.Tx, bool) {
	transfer := transfers.getTransfer()
	if transfer == nil {
		return tx, false
	}

	from := transfer.FromAddress
	to := transfer.ToAddress
	tx.From = from
	tx.To = to
	tx.Type = blockatlas.TxNativeTokenTransfer
	value := numbers.DecimalExp(transfer.Amount, 9)
	if transfers.isClaimReward() {
		tx.Meta = blockatlas.AnyAction{
			Coin:     coin.Ontology().ID,
			Name:     "Ontology Gas",
			Symbol:   "ONG",
			TokenID:  string(AssetONG),
			Decimals: 9,
			Value:    blockatlas.Amount(value),
			Title:    blockatlas.AnyActionClaimRewards,
			Key:      blockatlas.KeyStakeClaimRewards,
		}
		return tx, true
	}
	tx.Meta = blockatlas.NativeTokenTransfer{
		Name:     "Ontology Gas",
		Symbol:   "ONG",
		TokenID:  string(AssetONG),
		Decimals: 9,
		Value:    blockatlas.Amount(value),
		From:     from,
		To:       to,
	}
	return tx, true
}
