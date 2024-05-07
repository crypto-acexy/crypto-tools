package utxo

import (
	"github.com/acexy/golang-toolkit/http"
	"github.com/acexy/golang-toolkit/logger"
)

type MempoolUtxoData struct {
	Address string
	Txs     []*MempoolUtxoTxn
}

type MempoolStatus struct {
	BlockHeight uint64 `json:"block_height"`
}

type MempoolUtxoTxn struct {
	BlockHeight int
	BlockIndex  int
	Hash        string               `json:"txid"`
	Inputs      []*MempoolUtxoInput  `json:"vin"`
	Outputs     []*MempoolUtxoOutput `json:"vout"`
	Status      MempoolStatus        `json:"status"`
}

type MempoolUtxoInput struct {
	PrevHash    string   `json:"txid"`
	OutputIndex int      `json:"output_index"`
	OutputValue int      `json:"output_value"`
	Addresses   []string `json:"addresses"`
	Age         int      `json:"age"`
}

type MempoolUtxoOutput struct {
	Addresses []string `json:"addresses"`
	Value     int      `json:"value"`
	Script    string   `json:"script"`
	SentBy    string   `json:"sent_by"`
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
	var utxoData BlockcypherUtxoData
	resp, err := m.client.R().SetReturnStruct(&utxoData).Get("https://mempool.space/api/address/" + m.address + "/tx")
	if err != nil || resp.String() == "{\"error\": \"Limits reached.\"}" {
		logger.Logrus().WithError(err).Println("解析原始json数据失败", resp.String())
		return nil, err
	}
	return nil, nil
}

func NewMempoolPlatform(chain Chain, address string, httpProxy ...string) *BlockcrypherPlatformData {
	return &BlockcrypherPlatformData{
		chainId: chain,
		address: address,
		client:  http.NewRestyClient(httpProxy...).SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15"),
	}
}
