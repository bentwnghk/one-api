package common_test

import (
	"mime/multipart"
	"testing"

	"one-api/common"
	"one-api/types"

	"github.com/stretchr/testify/assert"
)

func TestGrokImagineTierRatio(t *testing.T) {
	assert.Equal(t, float64(1), common.GrokImagineTierRatio("1024x1024", "low"))
	assert.Equal(t, float64(1), common.GrokImagineTierRatio("1024x1024", ""))
	assert.Equal(t, float64(1.5), common.GrokImagineTierRatio("1024x1024", "medium"))
	assert.Equal(t, float64(1.5), common.GrokImagineTierRatio("2048x2048", "low"))
	assert.Equal(t, float64(2), common.GrokImagineTierRatio("2048x2048", "medium"))
	assert.Equal(t, float64(1), common.GrokImagineTierRatio("", ""))
	assert.Equal(t, float64(1.5), common.GrokImagineTierRatio("invalid", "medium"))
}

func TestIsPerImageBillingModel(t *testing.T) {
	assert.True(t, common.IsPerImageBillingModel("grok-imagine-image-2.0"))
	assert.True(t, common.IsPerImageBillingModel("x-ai/grok-imagine-image-2.0"))
	assert.False(t, common.IsPerImageBillingModel("dall-e-3"))
	assert.False(t, common.IsPerImageBillingModel(""))
}

func TestCountTokenImageGrokImagine(t *testing.T) {
	// 1K Low, n=1 -> $0.04
	tokens, err := common.CountTokenImage(types.ImageRequest{
		Model: "grok-imagine-image-2.0", Size: "1024x1024", N: 1, Quality: "low",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1000, tokens)

	// 2K Medium, n=1 -> $0.08
	tokens, err = common.CountTokenImage(types.ImageRequest{
		Model: "grok-imagine-image-2.0", Size: "2048x2048", N: 1, Quality: "medium",
	})
	assert.NoError(t, err)
	assert.Equal(t, 2000, tokens)

	// 1K Medium, n=2 -> 2 * $0.06
	tokens, err = common.CountTokenImage(types.ImageRequest{
		Model: "grok-imagine-image-2.0", Size: "1024x1024", N: 2, Quality: "medium",
	})
	assert.NoError(t, err)
	assert.Equal(t, 3000, tokens)

	// edit: unspecified quality defaults to low -> 1K Low + 3 input images -> $0.04 + 3 * $0.01
	tokens, err = common.CountTokenImage(types.ImageEditRequest{
		Model:  "grok-imagine-image-2.0",
		Size:   "1024x1024",
		N:      1,
		Image:  &multipart.FileHeader{},
		Images: []*multipart.FileHeader{{}, {}},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1750, tokens)

	// edit: unspecified quality defaults to low -> 2K Low + 1 input image -> $0.06 + $0.01
	tokens, err = common.CountTokenImage(types.ImageEditRequest{
		Model: "grok-imagine-image-2.0", Size: "2048x2048", N: 1, Image: &multipart.FileHeader{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1750, tokens)

	// edit: explicit low quality -> 1K Low + 1 input image -> $0.04 + $0.01
	tokens, err = common.CountTokenImage(types.ImageEditRequest{
		Model: "grok-imagine-image-2.0", Size: "1024x1024", N: 1, Quality: "low", Image: &multipart.FileHeader{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1250, tokens)

	// edit: explicit medium quality, case-insensitive -> 1K Medium
	tokens, err = common.CountTokenImage(types.ImageEditRequest{
		Model: "grok-imagine-image-2.0", Size: "1024x1024", N: 1, Quality: "Medium",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1500, tokens)

	// vendor-prefixed model name gets the same billing
	tokens, err = common.CountTokenImage(types.ImageEditRequest{
		Model: "x-ai/grok-imagine-image-2.0", Size: "1024x1024", N: 1, Quality: "low", Image: &multipart.FileHeader{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1250, tokens)

	tokens, err = common.CountTokenImage(types.ImageRequest{
		Model: "x-ai/grok-imagine-image-2.0", Size: "2048x2048", N: 1, Quality: "medium",
	})
	assert.NoError(t, err)
	assert.Equal(t, 2000, tokens)

	// edit: non-grok model keeps legacy behavior
	tokens, err = common.CountTokenImage(types.ImageEditRequest{
		Model: "dall-e-2", Size: "1024x1024", N: 1, Image: &multipart.FileHeader{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1250, tokens)
}
