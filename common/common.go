package common

import (
	"fmt"
	"math"
	"one-api/common/config"
)

func LogQuota(quota int) string {
	if config.DisplayInCurrencyEnabled {
		if quota < 0 {
			return fmt.Sprintf("-＄%.6f 額度", math.Abs(float64(quota)/config.QuotaPerUnit))
		}
		return fmt.Sprintf("＄%.6f 額度", float64(quota)/config.QuotaPerUnit)
	} else {
		return fmt.Sprintf("%d 點額度", quota)
	}
}
