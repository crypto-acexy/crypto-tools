package fee

import "testing"

func TestFeeDetail(t *testing.T) {
	NewFeeDetailRequest(
		21000,
		267,
		236.910013199,
		20,
	).Calc()
}

func TestFeeEstimation(t *testing.T) {
	NewFeeEstimationRequest(
		32,
		40,
		0.01,
	).Calc()
}
