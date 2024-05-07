package utxo

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Bitcoin Chain = iota
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
	Value     int
	Hash      string
	HashIndex int
	Height    int
	Used      bool
}

func Analyze(data DataPlatform) ([]*StandardUtxoAnalyzeResult, error) {

	rawData, err := data.convertRawData()

	if err != nil {
		return nil, err
	}

	if rawData == nil {
		return nil, errors.New("原始数据解析失败")
	}

	targetAddress := strings.ToLower(rawData.Address)
	txs := rawData.Txs
	if len(txs) == 0 {
		return nil, errors.New("无交易数据")
	}
	fmt.Println("分析地址", targetAddress)

	allHistory := make(map[string]*StandardUtxoAnalyzeResult, 0)
	results := make([]*StandardUtxoAnalyzeResult, 0)

	for i := len(txs) - 1; i >= 0; i-- {
		tx := txs[i]
		inputs := tx.Inputs
		fmt.Println("交易", tx.Hash, tx.BlockIndex, tx.BlockHeight)
		for _, input := range inputs {
			if strings.ToLower(input.Address) == targetAddress {
				analyze := StandardUtxoAnalyzeResult{
					Value:     input.OutputValue,
					Hash:      input.PrevHash,
					HashIndex: input.OutputIndex,
					Height:    input.Height,
					Used:      true,
				}

				utxoKey := fmt.Sprintf("hash:%s hashIndex:%d height:%d value:%d", analyze.Hash, analyze.HashIndex, analyze.Height, analyze.Value)
				fmt.Println("	消耗一块utxo", utxoKey)
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
					Value:     output.Value,
					Hash:      tx.Hash,
					HashIndex: tx.BlockIndex,
					Height:    tx.BlockHeight,
					Used:      false,
				}
				utxoKey := fmt.Sprintf("hash:%s hashIndex:%d height:%d value:%d", analyze.Hash, analyze.HashIndex, analyze.Height, analyze.Value)
				fmt.Println("	发现一块utxo", utxoKey)
				allHistory[utxoKey] = &analyze
				results = append(results, &analyze)

			}
		}
	}
	return results, nil
}
