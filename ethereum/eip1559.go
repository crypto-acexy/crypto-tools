package ethereum

import (
	"fmt"
	"github.com/shopspring/decimal"
)

var ethDecimal = decimal.NewFromInt(1000000000)

type feeRequest struct {

	// 实际消耗Gas
	gas decimal.Decimal

	// 基础网络每单位燃料价格
	netFeePerGas decimal.Decimal

	// 愿意支付的最大手续费
	maxFeePerGas decimal.Decimal

	// 支付给矿工的每单位燃料的消费
	priorityFeeParGas decimal.Decimal
}

func NewFeeRequest(gas int32, maxFeePerGas, netFeePerGas, priorityFeeParGas float64) *feeRequest {
	gasDecimal := decimal.NewFromInt32(gas)
	netFeePerGasDecimal := decimal.NewFromFloat(netFeePerGas)
	maxFeePerGasDecimal := decimal.NewFromFloat(maxFeePerGas)
	priorityFeeParGasDecimal := decimal.NewFromFloat(priorityFeeParGas)
	return &feeRequest{
		gas:               gasDecimal,
		netFeePerGas:      netFeePerGasDecimal,
		maxFeePerGas:      maxFeePerGasDecimal,
		priorityFeeParGas: priorityFeeParGasDecimal,
	}
}

func (f *feeRequest) FeeDetail() {
	allowMaxFee := f.maxFeePerGas.Mul(f.gas)
	// 基础手续费
	baseFee := f.netFeePerGas.Mul(f.gas)
	minerFee := decimal.NewFromFloat32(0.0)
	if baseFee.GreaterThan(allowMaxFee) {
		fmt.Printf("手续费不足, maxFeePerGas至少为 %s Gwei \n", f.netFeePerGas.String())
		return
	} else if baseFee.Equal(allowMaxFee) || f.priorityFeeParGas.LessThanOrEqual(decimal.Zero) {
		fmt.Printf("矿工收入小费: 0\n")
	} else {
		minerFee = f.priorityFeeParGas.Mul(f.gas)
		if minerFee.Add(baseFee).GreaterThan(allowMaxFee) {
			minerFee = allowMaxFee.Sub(baseFee)
		}
		fmt.Printf("矿工收入小费: %s Gwei %s Ether\n", minerFee.String(), minerFee.DivRound(ethDecimal, 18).String())
	}

	fmt.Printf("基础网络费用为: %s Gwei %s Ether\n", baseFee.String(), baseFee.DivRound(ethDecimal, 18).String())
	fmt.Printf("总支付费用为: %s Gwei %s Ether\n", baseFee.Add(minerFee).String(), baseFee.Add(minerFee).DivRound(ethDecimal, 18).String())
}
