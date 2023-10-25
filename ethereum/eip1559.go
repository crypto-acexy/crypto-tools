package ethereum

import (
	"fmt"
	"github.com/shopspring/decimal"
)

var ethDecimal = decimal.NewFromInt(1000000000)

type feeDetailRequest struct {
	// 实际消耗Gas
	gas decimal.Decimal

	// 基础网络每单位燃料价格
	netFeePerGas decimal.Decimal

	// 愿意支付的最大手续费
	maxFeePerGas decimal.Decimal

	// 支付给矿工的每单位燃料的消费
	priorityFeeParGas decimal.Decimal
}

type feeEstimationRequest struct {
	netMinFeePerGas   decimal.Decimal
	netMaxFeePerGas   decimal.Decimal
	priorityFeeParGas decimal.Decimal
}

func NewFeeEstimationRequest(netMinFeePerGas, netMaxFeePerGas, priorityFeeParGas float64) *feeEstimationRequest {
	return &feeEstimationRequest{
		netMinFeePerGas:   decimal.NewFromFloat(netMinFeePerGas),
		netMaxFeePerGas:   decimal.NewFromFloat(netMaxFeePerGas),
		priorityFeeParGas: decimal.NewFromFloat(priorityFeeParGas),
	}
}

func NewFeeDetailRequest(gas int32, maxFeePerGas, netFeePerGas, priorityFeeParGas float64) *feeDetailRequest {
	return &feeDetailRequest{
		gas:               decimal.NewFromInt32(gas),
		netFeePerGas:      decimal.NewFromFloat(netFeePerGas),
		maxFeePerGas:      decimal.NewFromFloat(maxFeePerGas),
		priorityFeeParGas: decimal.NewFromFloat(priorityFeeParGas),
	}
}

func (f *feeDetailRequest) Calc() {
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

func (f *feeEstimationRequest) Calc() {
	fmt.Printf("maxFeePerGas 的范围应当在 %s Gwei ~ %s Gwei 之间", f.netMinFeePerGas.Add(f.priorityFeeParGas).String(), f.netMaxFeePerGas.Add(f.priorityFeeParGas))
}
