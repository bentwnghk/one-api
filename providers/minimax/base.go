package minimax

import (
	"encoding/json"
	"fmt"
	"net/http"
	"one-api/common/requester"
	"one-api/model"
	"one-api/providers/base"
	"one-api/providers/openai"
	"one-api/types"
	"strings"
)

type MiniMaxProviderFactory struct{}

// 创建 MiniMaxProvider
func (f MiniMaxProviderFactory) Create(channel *model.Channel) base.ProviderInterface {
	return &MiniMaxProvider{
		OpenAIProvider: openai.OpenAIProvider{
			BaseProvider: base.BaseProvider{
				Config:    getConfig(),
				Channel:   channel,
				Requester: requester.NewHTTPRequester(*channel.Proxy, requestErrorHandle),
			},
		},
	}
}

type MiniMaxProvider struct {
	openai.OpenAIProvider
}

// GetFullRequestURL overrides the base provider's method to add GroupId for MiniMax
func (p *MiniMaxProvider) GetFullRequestURL(requestURL string, modelName string) string {
	baseURL := p.GetBaseURL()
	if modelName != "" {
		requestURL = strings.Replace(requestURL, "{model}", modelName, 1)
	}

	keyParts := strings.Split(p.Channel.Key, ":")
	groupID := ""
	if len(keyParts) > 1 { // Assuming key is "group_id:api_key"
		groupID = keyParts[0]
	} else if len(keyParts) == 1 && !strings.HasPrefix(p.Channel.Key, "sk-") { // Or just group_id if no colon
		// This condition might need adjustment based on how users store single group_id vs api_key
		// For now, assume if it's not like an sk- key, it could be a group_id if only one part.
		// A more robust way would be to have separate fields or a clearer convention.
		// groupID = keyParts[0] // Commenting out for now to prefer "group_id:api_key"
	}


	if groupID == "" && p.Channel.Other != "" { // Fallback to Other field if configured
		groupID = p.Channel.Other
	}


	if groupID != "" {
		return fmt.Sprintf("%s%s?GroupId=%s", baseURL, requestURL, groupID)
	}
	return fmt.Sprintf("%s%s", baseURL, requestURL)
}

// GetRequestHeaders overrides the base provider's method to use the correct API key part for MiniMax
func (p *MiniMaxProvider) GetRequestHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	// Default to OpenAI's way of getting headers, then adjust for MiniMax
	// p.OpenAIProvider.GetRequestHeaders() // This would call the embedded one, let's be more direct

	apiKey := p.Channel.Key
	keyParts := strings.Split(p.Channel.Key, ":")
	if len(keyParts) > 1 { // Assuming key is "group_id:api_key"
		apiKey = keyParts[1]
	}
	// If only one part, assume it's the api_key if it starts with sk-, otherwise it might be ambiguous
	// or it's just the api_key without a group_id (though API needs group_id in URL)

	headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	return headers
}

func getConfig() base.ProviderConfig {
	return base.ProviderConfig{
		BaseURL:         "https://api.minimax.io",
		ChatCompletions: "/v1/chat/completions",
		AudioSpeech:     "/v1/t2a_v2",
		// Embeddings:      "/v1/embeddings",
		// ModelList:       "/v1/models",
	}
}

// 请求错误处理
func requestErrorHandle(resp *http.Response) *types.OpenAIError {
	minimaxError := &MiniMaxBaseResp{}
	err := json.NewDecoder(resp.Body).Decode(minimaxError)
	if err != nil {
		return nil
	}

	return errorHandle(&minimaxError.BaseResp)
}

// 错误处理
func errorHandle(minimaxError *BaseResp) *types.OpenAIError {
	if minimaxError.StatusCode == 0 {
		return nil
	}
	return &types.OpenAIError{
		Message: minimaxError.StatusMsg,
		Type:    "minimax_error",
		Code:    minimaxError.StatusCode,
	}
}
