package utxo

import (
	"errors"
	"fmt"
	"github.com/acexy/golang-toolkit/util/json"
	"strings"
)

const (
	Bitcoin Chain = iota
	Litecoin
)

type Chain int

type DataPlatform interface {
	chain() Chain
	convertRawData() (*StandardUtxoData, error)
}

type StandardUtxoData struct {
	Address string
	Txs     []*StandardUtxoTxn
}

type StandardUtxoTxn struct {
	BlockHash   string
	BlockHeight int
	BlockIndex  int
	Hash        string
	Inputs      []*StandardUtxoInput
	Outputs     []*StandardUtxoOutput
}

type StandardUtxoInput struct {
	PrevHash    string
	OutputIndex int
	OutputValue int
	Address     string
	Height      int
}

type StandardUtxoOutput struct {
	Address string
	Value   int
	Script  string
	SentBy  string
}

type StandardUtxoAnalyzeResult struct {
	Value  int
	Hash   string
	Height int
	Used   bool
}

func Analyze(data DataPlatform) error {

	rawData, err := data.convertRawData()

	if err != nil {
		return err
	}

	if rawData == nil {
		return errors.New("原始数据解析失败")
	}

	targetAddress := strings.ToLower(rawData.Address)
	txs := rawData.Txs
	if len(txs) == 0 {
		return errors.New("无交易数据")
	}
	fmt.Println("分析地址", targetAddress)

	allHistory := make(map[string]*StandardUtxoAnalyzeResult, 0)
	results := make([]*StandardUtxoAnalyzeResult, 0)

	for i := len(txs) - 1; i >= 0; i-- {
		tx := txs[i]
		inputs := tx.Inputs
		fmt.Println("交易 hash", tx.Hash, "高度", tx.BlockHeight)
		for _, input := range inputs {

			if strings.ToLower(input.Address) == targetAddress {
				analyze := StandardUtxoAnalyzeResult{
					Value:  input.OutputValue,
					Hash:   input.PrevHash,
					Height: input.Height,
					Used:   true,
				}

				utxoKey := fmt.Sprintf("hash:%s", analyze.Hash)
				fmt.Println("	消耗一块utxo hash", analyze.Hash, "value", analyze.Value)
				v, ok := allHistory[utxoKey]
				if ok {
					v.Used = true
				} else {
					allHistory[utxoKey] = &analyze
					results = append(results, &analyze)
				}
			}
		}

		outputs := tx.Outputs
		for _, output := range outputs {
			if strings.ToLower(output.Address) == targetAddress {
				analyze := StandardUtxoAnalyzeResult{
					Value:  output.Value,
					Hash:   tx.Hash,
					Height: tx.BlockHeight,
					Used:   false,
				}
				utxoKey := fmt.Sprintf("hash:%s", analyze.Hash)
				fmt.Println("	发现一块utxo hash", analyze.Hash, "value", analyze.Value)
				allHistory[utxoKey] = &analyze
				results = append(results, &analyze)
			}
		}
	}
	balance := 0
	for _, result := range results {
		if !result.Used {
			balance += result.Value
		}
	}
	fmt.Println("分析结果", "总余额", balance, "详情")
	fmt.Println(json.ToJsonFormat(results))
	return nil
}
