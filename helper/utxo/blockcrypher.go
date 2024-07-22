package utxo

import (
	"errors"
	"fmt"
	"github.com/acexy/golang-toolkit/http"
	"github.com/acexy/golang-toolkit/logger"
	"github.com/jinzhu/copier"
	"time"
)

var blockcypherUrlConfig = map[Chain]string{
	Bitcoin:  "https://api.blockcypher.com/v1/btc/main/addrs/%s/full?limit=50&unspentOnly=true&includeScript=false",
	Litecoin: "https://api.blockcypher.com/v1/ltc/main/addrs/%s/full?limit=50&unspentOnly=true&includeScript=false",
	Dogecoin: "https://api.blockcypher.com/v1/doge/main/addrs/%s/full?limit=50&unspentOnly=true&includeScript=false",
}

type BlockcypherUtxoData struct {
	Address string                `json:"address"`
	Txs     []*BlockcypherUtxoTxn `json:"txs"`
	HasMore bool                  `json:"hasMore"`
}

type BlockcypherUtxoTxn struct {
	BlockHash   string                   `json:"block_hash"`
	BlockHeight int                      `json:"block_height"`
	BlockIndex  int                      `json:"block_index"`
	Hash        string                   `json:"hash"`
	Inputs      []*BlockcypherUtxoInput  `json:"inputs"`
	Outputs     []*BlockcypherUtxoOutput `json:"outputs"`
	NextInputs  string                   `json:"next_inputs"`
	NextOutputs string                   `json:"next_outputs"`
}

type BlockcypherUtxoInput struct {
	PrevHash    string   `json:"prev_hash"`
	OutputIndex int      `json:"output_index"`
	OutputValue int      `json:"output_value"`
	Addresses   []string `json:"addresses"`
	Address     string
	Height      int `json:"age"`
}

type BlockcypherUtxoOutput struct {
	Addresses []string `json:"addresses"`
	Address   string
	Value     int    `json:"value"`
	Script    string `json:"script"`
	SentBy    string `json:"sent_by"`
}

type BlockcrypherPlatformData struct {
	chainId Chain
	address string
	client  *http.RestyClient
}

func (b *BlockcrypherPlatformData) chain() Chain {
	return b.chainId
}

func (b *BlockcrypherPlatformData) convertRawData() (*StandardUtxoData, error) {
	var utxoData BlockcypherUtxoData

	uri := blockcypherUrlConfig[b.chain()]
	if uri == "" {
		return nil, errors.New("not support this chain now")
	}

	resp, err := b.client.R().SetReturnStruct(&utxoData).Get(fmt.Sprintf(uri, b.address))
	if err != nil || resp.String() == "{\"error\": \"Limits reached.\"}" {
		logger.Logrus().WithError(err).Println("解析原始json数据失败", resp.String())
		return nil, err
	}

	tx, err := b.loadTx(&utxoData)
	if err != nil {
		return nil, err
	}

	txSize := len(tx.Txs)
	minHeight := tx.Txs[txSize-1].BlockHeight
	hashMore := utxoData.HasMore
	for hashMore {
		hashMore = false
		time.Sleep(2 * time.Second)
		var moreUtxoData BlockcypherUtxoData
		resp, err = b.client.R().SetReturnStruct(&moreUtxoData).Get(fmt.Sprintf(uri+"&before=%d", b.address, minHeight))
		if err != nil || resp.String() == "{\"error\": \"Limits reached.\"}" {
			logger.Logrus().WithError(err).Println("解析原始json数据失败", resp.String())
			return nil, err
		}
		hashMore = moreUtxoData.HasMore
		moreTx, err := b.loadTx(&moreUtxoData)
		if err != nil {
			return nil, err
		}
		if moreTx == nil {
			return nil, errors.New("获取更多数据异常")
		}
		if len(moreTx.Txs) > 0 {
			tx.Txs = append(tx.Txs, moreTx.Txs...)
			txSize = len(moreTx.Txs)
			minHeight = moreTx.Txs[txSize-1].BlockHeight
		}
	}

	return tx, nil

}

func NewBlockcrypherPlatform(chain Chain, address string, httpProxy ...string) *BlockcrypherPlatformData {
	return &BlockcrypherPlatformData{
		chainId: chain,
		address: address,
		client:  http.NewRestyClient(httpProxy...).SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15"),
	}
}

func (b *BlockcrypherPlatformData) loadTx(utxoData *BlockcypherUtxoData) (*StandardUtxoData, error) {
	txs := utxoData.Txs
	for i := len(txs) - 1; i >= 0; i-- {
		tx := txs[i]
		//if tx.NextInputs != "" {
		//	nextInput := tx.NextInputs
		//	for {
		//		var moreTx BlockcypherUtxoTxn
		//		resp, err := b.client.R().SetReturnStruct(&moreTx).Get(nextInput)
		//		if err != nil || resp.String() == "{\"error\": \"Limits reached.\"}" {
		//			logger.Logrus().WithError(err).Error("查询更多utxo数据异常", nextInput)
		//			return nil, err
		//		} else {
		//			logger.Logrus().Debugln("补充查询output", nextInput)
		//		}
		//		if len(moreTx.Inputs) > 0 {
		//			tx.Inputs = append(tx.Inputs, moreTx.Inputs...)
		//		} else {
		//			break
		//		}
		//		if moreTx.NextInputs == "" {
		//			break
		//		}
		//		nextInput = moreTx.NextInputs
		//		time.Sleep(time.Second * 10 * time.Duration(rand.Intn(2)+1))
		//	}
		//}

		//if tx.NextOutputs != "" {
		//	nextOutput := tx.NextOutputs
		//	for {
		//		var moreTx BlockcypherUtxoTxn
		//		resp, err := b.client.R().SetReturnStruct(&moreTx).Get(nextOutput)
		//		if err != nil || resp.String() == "{\"error\": \"Limits reached.\"}" {
		//			logger.Logrus().WithError(err).Error("查询更多utxo数据异常", resp.String(), nextOutput)
		//			if err == nil {
		//				err = errors.New("请求触发限制")
		//			}
		//			return nil, err
		//		} else {
		//			logger.Logrus().Debugln("补充查询output", nextOutput)
		//		}
		//
		//		if len(moreTx.Outputs) > 0 {
		//			hitTargetAddress := false
		//			for _, output := range moreTx.Outputs {
		//				if strings.ToLower(output.Addresses[0]) == strings.ToLower(utxoData.Address) {
		//					hitTargetAddress = true
		//					break
		//				}
		//			}
		//			tx.Outputs = append(tx.Outputs, moreTx.Outputs...)
		//			if hitTargetAddress {
		//				logger.Logrus().Debugln("已找到目标UTXO地址忽略更多数据")
		//				break
		//			}
		//		} else {
		//			break
		//		}
		//
		//		if moreTx.NextOutputs == "" {
		//			break
		//		}
		//		nextOutput = moreTx.NextOutputs
		//		time.Sleep(time.Second * 10)
		//	}
		//}

		inputs := tx.Inputs
		for _, input := range inputs {
			input.Address = input.Addresses[0]
		}
		outputs := tx.Outputs
		for _, output := range outputs {
			output.Address = output.Addresses[0]
		}
	}
	var utxo StandardUtxoData
	err := copier.Copy(&utxo, &utxoData)
	if err != nil {
		logger.Logrus().WithError(err).Println("对象复制失败")
		return nil, err
	}
	return &utxo, nil
}
