package saleRecordFee

import "github.com/pangpanglabs/goutils/number"

func GetToFixedPrice(price float64, BaseTrimCode string) float64 {
	if BaseTrimCode == "" || BaseTrimCode == "A" {
		// 原价
		return number.ToFixed(price, nil)
	}

	var setting *number.Setting
	switch BaseTrimCode {
	case "C":
		// 按元向上取整
		setting = &number.Setting{
			RoundStrategy: "ceil",
		}
	case "O":
		// 按角向下取整
		setting = &number.Setting{
			RoundStrategy: "floor",
			RoundDigit:    1,
		}
	case "P":
		// 按角四舍五入
		setting = &number.Setting{
			RoundStrategy: "round",
			RoundDigit:    1,
		}
	case "Q":
		// 按角向上取整
		setting = &number.Setting{
			RoundStrategy: "ceil",
			RoundDigit:    1,
		}
	case "R":
		// 按元四舍五入
		setting = &number.Setting{
			RoundStrategy: "round",
		}
	case "T":
		// 按元向下取整
		setting = &number.Setting{
			RoundStrategy: "floor",
		}
	case "feeAmount":
		// 按分向上取整
		setting = &number.Setting{
			RoundStrategy: "ceil",
			RoundDigit:    2,
		}
	}

	return number.ToFixed(price, setting)
}
