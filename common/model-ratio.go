package common

import (
	"strconv"
	"strings"
)

var DalleSizeRatios = map[string]map[string]float64{
	"dall-e-2": {
		"256x256":   1,
		"512x512":   1.125,
		"1024x1024": 1.25,
	},
	"dall-e-3": {
		"1024x1024": 1,
		"1024x1792": 2,
		"1792x1024": 2,
	},
}

var DalleGenerationImageAmounts = map[string][2]int{
	"dall-e-2": {1, 10},
	"dall-e-3": {1, 1}, // OpenAI allows n=1 currently.
}

// 按请求参数（分辨率/质量/数量）计费的图片模型，忽略上游返回的 usage
var PerImageBillingModels = map[string]bool{
	"grok-imagine-image-2.0": true,
}

// grok-imagine-image-2.0 计费倍率（相对基准价 $0.04/张）：
// 1K Low: 1  2K Low / 1K Medium: 1.5 ($0.06)  2K Medium: 2 ($0.08)
func GrokImagineTierRatio(size, quality string) float64 {
	is2K, isMedium := isLargeSize(size), quality == "medium"

	switch {
	case is2K && isMedium:
		return 2
	case is2K || isMedium:
		return 1.5
	default:
		return 1
	}
}

func isLargeSize(size string) bool {
	dims := strings.Split(size, "x")
	if len(dims) != 2 {
		return false
	}
	width, errW := strconv.Atoi(dims[0])
	height, errH := strconv.Atoi(dims[1])
	if errW != nil || errH != nil {
		return false
	}
	return width > 1024 || height > 1024
}

// 图片编辑时每张输入图片的附加费用（相对单张输出图片价格的倍率）
// grok-imagine-image-2.0: $0.01 / $0.04 = 0.25
var ImageEditInputImageRatios = map[string]float64{
	"grok-imagine-image-2.0": 0.25,
}

func IsPerImageBillingModel(model string) bool {
	return PerImageBillingModels[model]
}
