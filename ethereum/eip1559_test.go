package ethereum

import "testing"

func TestFeeDetail(t *testing.T) {
	request := NewFeeRequest(
		21000,
		40,
		12.921534465,
		0.01)
	request.FeeDetail()
}
