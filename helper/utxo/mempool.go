package utxo

import (
	"github.com/acexy/golang-toolkit/http"
	"github.com/acexy/golang-toolkit/logger"
	"github.com/jinzhu/copier"
)

type MempoolUtxoData struct {
	Address string
	Txs     []*MempoolUtxoTxn
}

type MempoolStatus struct {
	BlockHeight int `json:"block_height"`
}

type MempoolUtxoTxn struct {
	BlockHeight int
	BlockIndex  int
	Hash        string               `json:"txid"`
	Inputs      []*MempoolUtxoInput  `json:"vin"`
	Outputs     []*MempoolUtxoOutput `json:"vout"`
	Status      MempoolStatus        `json:"status"`
}

type MempoolPreVout struct {
	Value   int    `json:"value"`
	Index   int    `json:"vout"`
	Address string `json:"scriptpubkey_address"`
}

type MempoolUtxoInput struct {
	PrevHash    string `json:"txid"`
	OutputIndex int
	OutputValue int
	Address     string
	PreVout     MempoolPreVout `json:"prevout"`
}

type MempoolUtxoOutput struct {
	Address string `json:"scriptpubkey_address"`
	Value   int    `json:"value"`
	Script  string `json:"scriptpubkey"`
}

type MempoolPlatformData struct {
	chainId Chain
	address string
	client  *http.RestyClient
}

func (m *MempoolPlatformData) chain() Chain {
	return m.chainId
}

func (m *MempoolPlatformData) convertRawData() (*StandardUtxoData, error) {
	utxoData := make([]*MempoolUtxoTxn, 0)
	resp, err := m.client.R().SetReturnStruct(&utxoData).Get("https://mempool.space/api/address/" + m.address + "/txs")

	if err != nil || resp.String() == "{\"error\": \"Limits reached.\"}" || len(utxoData) == 0 {
		logger.Logrus().WithError(err).Println("解析原始json数据失败", resp.String())
		return nil, err
	}
	txs := make([]*StandardUtxoTxn, 0)
	for _, tx := range utxoData {
		tx.BlockHeight = tx.Status.BlockHeight
		tx.BlockIndex = -1

		inputs := tx.Inputs
		for _, input := range inputs {
			input.OutputIndex = input.PreVout.Index
			input.OutputValue = input.PreVout.Value
			input.Address = input.PreVout.Address
		}
	}

	err = copier.Copy(&txs, utxoData)
	if err != nil {
		return nil, err
	}
	result := &StandardUtxoData{
		Address: m.address,
		Txs:     txs,
	}

	return result, nil
}

func NewMempoolPlatform(chain Chain, address string, httpProxy ...string) *MempoolPlatformData {
	return &MempoolPlatformData{
		chainId: chain,
		address: address,
		client:  http.NewRestyClient(httpProxy...).SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15"),
	}
}
