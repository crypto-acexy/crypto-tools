package ethereum

import "testing"

func TestFeeDetail(t *testing.T) {
	NewFeeDetailRequest(
		21000,
		40,
		12.921534465,
		0.01,
	).Calc()
}

func TestFeeEstimation(t *testing.T) {
	NewFeeEstimationRequest(
		32,
		40,
		0.01,
	).Calc()
}
