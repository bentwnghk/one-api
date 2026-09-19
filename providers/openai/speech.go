package openai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"one-api/common"
	"one-api/common/config"
	"one-api/common/requester"
	"one-api/types"
)

type speechAudioResponse struct {
	Data     string                `json:"data"`
	Audio    string                `json:"audio"`
	Content  string                `json:"content"`
	B64Json  string                `json:"b64_json"`
	MimeType string                `json:"mime_type"`
	Usage    *types.ResponsesUsage `json:"usage"`
}

func (r *speechAudioResponse) AudioData() string {
	for _, s := range []string{r.Data, r.Audio, r.Content, r.B64Json} {
		if s != "" {
			return s
		}
	}
	return ""
}

func (p *OpenAIProvider) CreateSpeech(request *types.SpeechAudioRequest) (*http.Response, *types.OpenAIErrorWithStatusCode) {
	req, errWithCode := p.GetRequestTextBody(config.RelayModeAudioSpeech, request.Model, request)
	if errWithCode != nil {
		return nil, errWithCode
	}
	defer req.Body.Close()

	var resp *http.Response
	resp, errWithCode = p.Requester.SendRequestRaw(req)
	if errWithCode != nil {
		return nil, errWithCode
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return p.handleSpeechJsonResponse(resp, request)
	}

	p.Usage.TotalTokens = p.Usage.PromptTokens

	return resp, nil
}

func (p *OpenAIProvider) handleSpeechJsonResponse(resp *http.Response, request *types.SpeechAudioRequest) (*http.Response, *types.OpenAIErrorWithStatusCode) {
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, common.ErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError)
	}

	speechResponse := &speechAudioResponse{}
	if json.Unmarshal(body, speechResponse) == nil {
		if data := speechResponse.AudioData(); data != "" {
			audioBytes, decodeErr := base64.StdEncoding.DecodeString(data)
			if decodeErr != nil {
				audioBytes, decodeErr = base64.RawStdEncoding.DecodeString(data)
			}
			if decodeErr == nil {
				if speechResponse.Usage != nil {
					p.Usage.PromptTokens = speechResponse.Usage.InputTokens
					p.Usage.CompletionTokens = speechResponse.Usage.OutputTokens
					p.Usage.TotalTokens = speechResponse.Usage.TotalTokens
				} else {
					p.Usage.TotalTokens = p.Usage.PromptTokens
				}

				return buildAudioResponse(audioBytes, speechContentType(speechResponse.MimeType, request.ResponseFormat)), nil
			}
		}
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))
	return nil, requester.HandleErrorResp(resp, p.Requester.ErrorHandler, p.Requester.IsOpenAI)
}

func speechContentType(mimeType string, responseFormat string) string {
	if mimeType != "" {
		return mimeType
	}
	switch responseFormat {
	case "mp3":
		return "audio/mpeg"
	case "opus":
		return "audio/opus"
	case "aac":
		return "audio/aac"
	case "flac":
		return "audio/flac"
	case "wav":
		return "audio/wav"
	case "pcm":
		return "audio/pcm; rate=24000"
	default:
		return "application/octet-stream"
	}
}

func buildAudioResponse(audioBytes []byte, contentType string) *http.Response {
	audioResp := &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(bytes.NewReader(audioBytes)),
		Header:        make(http.Header),
		ContentLength: int64(len(audioBytes)),
	}
	audioResp.Header.Set("Content-Type", contentType)
	audioResp.Header.Set("Content-Length", strconv.Itoa(len(audioBytes)))

	return audioResp
}
