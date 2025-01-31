package azure

import (
	"one-api/common/requester"
	"one-api/model"
	"one-api/providers/base"
	"one-api/providers/openai"
)

type AzureProviderFactory struct{}

// 创建 AzureProvider
func (f AzureProviderFactory) Create(channel *model.Channel) base.ProviderInterface {
	config := getAzureConfig()
	provider := &AzureProvider{
		OpenAIProvider: openai.OpenAIProvider{
			BaseProvider: base.BaseProvider{
				Config:    config,
				Channel:   channel,
				Requester: requester.NewHTTPRequester(*channel.Proxy, openai.RequestErrorHandle),
			},
			IsAzure:              true,
			BalanceAction:        false,
			SupportStreamOptions: true, // Default to true
		},
	}

    // Check the model name and set SupportStreamOptions accordingly
    if channel.modelName == "o1-mini" || channel.modelName == "o1-preview" {
        provider.SupportStreamOptions = false
    }

    return provider
}

func getAzureConfig() base.ProviderConfig {
	return base.ProviderConfig{
		BaseURL:             "",
		Completions:         "/completions",
		ChatCompletions:     "/chat/completions",
		Embeddings:          "/embeddings",
		AudioSpeech:         "/audio/speech",
		AudioTranscriptions: "/audio/transcriptions",
		AudioTranslations:   "/audio/translations",
		ImagesGenerations:   "/images/generations",
		ChatRealtime:        "/realtime",
	}
}

type AzureProvider struct {
	openai.OpenAIProvider
}
