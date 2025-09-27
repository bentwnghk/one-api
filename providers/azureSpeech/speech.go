package azureSpeech

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"one-api/common"
	"one-api/common/config"
	"one-api/types"
	"strings"
)

var outputFormatMap = map[string]string{
	"mp3":  "audio-16khz-128kbitrate-mono-mp3",
	"opus": "audio-16khz-128kbitrate-mono-opus",
	"aac":  "audio-24khz-160kbitrate-mono-mp3",
	"flac": "audio-48khz-192kbitrate-mono-mp3",
}

func CreateSSML(text string, name string, role string, speed float64) string {
	ssmlTemplate := `<speak version='1.0' xml:lang='en-US'>
        <voice xml:lang='en-US' %s name='%s'%s>
            %s
        </voice>
    </speak>`

	roleAttribute := ""
	if role != "" {
		roleAttribute = fmt.Sprintf(" role='%s'", role)
	}

	rateAttribute := ""
	// Validate speed range (Azure typically supports 0.5x to 2.0x or 3.0x)
	if speed > 0 && speed >= 0.5 && speed <= 3.0 {
		// Try different rate formats that Azure supports
		// Format 1: multiplier (e.g., "0.5", "1.0", "2.0")
		rateAttribute = fmt.Sprintf(" rate='%.2f'", speed)
		fmt.Printf("DEBUG: Azure Speech - Speed parameter received: %.2f, rate attribute: %s\n", speed, rateAttribute)
	} else if speed > 0 {
		// Clamp speed to valid range if it's outside bounds
		clampedSpeed := math.Max(0.5, math.Min(3.0, speed))
		rateAttribute = fmt.Sprintf(" rate='%.2f'", clampedSpeed)
		fmt.Printf("DEBUG: Azure Speech - Speed %.2f clamped to %.2f, rate attribute: %s\n", speed, clampedSpeed, rateAttribute)
	}

	generatedSSML := fmt.Sprintf(ssmlTemplate, roleAttribute, name, rateAttribute, text)
	fmt.Printf("DEBUG: Azure Speech - Generated SSML: %s\n", generatedSSML)

	return generatedSSML
}

func (p *AzureSpeechProvider) GetVoiceMap() map[string][]string {
	defaultVoiceMapping := map[string][]string{
		"alloy":   {"zh-CN-YunxiNeural"},
		"echo":    {"zh-CN-YunyangNeural"},
		"fable":   {"zh-CN-YunxiNeural", "Boy"},
		"onyx":    {"zh-CN-YunyeNeural"},
		"nova":    {"zh-CN-XiaochenNeural"},
		"shimmer": {"zh-CN-XiaohanNeural"},
	}

	if p.Channel.Plugin == nil {
		return defaultVoiceMapping
	}

	customVoiceMapping, ok := p.Channel.Plugin.Data()["voice"]
	if !ok {
		return defaultVoiceMapping
	}

	for key, value := range customVoiceMapping {
		if _, exists := defaultVoiceMapping[key]; !exists {
			continue
		}
		customVoiceValue, isString := value.(string)
		if !isString || customVoiceValue == "" {
			continue
		}
		customizeVoice := strings.Split(customVoiceValue, "|")
		defaultVoiceMapping[key] = customizeVoice
	}

	return defaultVoiceMapping
}

func (p *AzureSpeechProvider) getRequestBody(request *types.SpeechAudioRequest) *bytes.Buffer {
	var voice, role string
	voiceMap := p.GetVoiceMap()
	if voiceMap[request.Voice] != nil {
		voice = voiceMap[request.Voice][0]
		if len(voiceMap[request.Voice]) > 1 {
			role = voiceMap[request.Voice][1]
		}
	} else {
		voice = request.Voice
	}

	// Debug log to verify speed value from request
	fmt.Printf("DEBUG: Azure Speech - Request speed value: %.2f\n", request.Speed)

	ssml := CreateSSML(request.Input, voice, role, request.Speed)

	return bytes.NewBufferString(ssml)

}

func (p *AzureSpeechProvider) CreateSpeech(request *types.SpeechAudioRequest) (*http.Response, *types.OpenAIErrorWithStatusCode) {
	url, errWithCode := p.GetSupportedAPIUri(config.RelayModeAudioSpeech)
	if errWithCode != nil {
		return nil, errWithCode
	}
	fullRequestURL := p.GetFullRequestURL(url)
	headers := p.GetRequestHeaders()
	responseFormatr := outputFormatMap[request.ResponseFormat]
	if responseFormatr == "" {
		responseFormatr = outputFormatMap["mp3"]
	}
	headers["X-Microsoft-OutputFormat"] = responseFormatr

	requestBody := p.getRequestBody(request)

	req, err := p.Requester.NewRequest(http.MethodPost, fullRequestURL, p.Requester.WithBody(requestBody), p.Requester.WithHeader(headers))
	if err != nil {
		return nil, common.ErrorWrapper(err, "new_request_failed", http.StatusInternalServerError)
	}
	defer req.Body.Close()

	var resp *http.Response
	resp, errWithCode = p.Requester.SendRequestRaw(req)
	if errWithCode != nil {
		return nil, errWithCode
	}

	p.Usage.TotalTokens = p.Usage.PromptTokens

	return resp, nil
}
