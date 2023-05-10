package goutils

import (
	"github.com/shopspring/decimal"
)

var oneHundredDecimal decimal.Decimal = decimal.NewFromInt(100)

// 分转元
func Fen2Yuan(fen int64) float64 {
	y, _ := decimal.NewFromInt(fen).Div(oneHundredDecimal).Truncate(2).Float64()
	return y
}

// 元转分
func Yuan2Fen(yuan float64) int64 {

	f, _ := decimal.NewFromFloat(yuan).Mul(oneHundredDecimal).Truncate(0).Float64()
	return int64(f)

}

// 计算手续费 接收参数 金额 和 费率 返回 float 64
func ChargeFee(totalFee int, feeRate string) int {
	// 100 * 0.04
	fTotalFee := float64(totalFee) // 浮点数
	v2 := decimal.NewFromFloat(fTotalFee)

	f, _ := decimal.NewFromFloat(String2Float64(feeRate)).Mul(v2).Truncate(0).Float64()
	return int(f)
}
