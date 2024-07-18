package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"one-api/common/config"
	"one-api/common/logger"
	"one-api/common/notify"
	"one-api/common/utils"
	"one-api/model"
	"one-api/providers"
	providers_base "one-api/providers/base"
	"one-api/types"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func testChannel(channel *model.Channel, testModel string) (err error, openaiErr *types.OpenAIErrorWithStatusCode) {
	if channel.TestModel == "" {
		return errors.New("請填寫測速模型後再試"), nil
	}

	// 创建一个 http.Request
	req, err := http.NewRequest("POST", "/v1/chat/completions", nil)
	if err != nil {
		return err, nil
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	request := buildTestRequest()

	if testModel != "" {
		request.Model = testModel
	} else {
		request.Model = channel.TestModel
	}

	provider := providers.GetProvider(channel, c)
	if provider == nil {
		return errors.New("channel not implemented"), nil
	}

	newModelName, err := provider.ModelMappingHandler(request.Model)
	if err != nil {
		return err, nil
	}

	request.Model = newModelName

	chatProvider, ok := provider.(providers_base.ChatInterface)
	if !ok {
		return errors.New("channel not implemented"), nil
	}

	chatProvider.SetUsage(&types.Usage{})

	response, openAIErrorWithStatusCode := chatProvider.CreateChatCompletion(request)

	if openAIErrorWithStatusCode != nil {
		return errors.New(openAIErrorWithStatusCode.Message), openAIErrorWithStatusCode
	}

	// 转换为JSON字符串
	jsonBytes, _ := json.Marshal(response)
	logger.SysLog(fmt.Sprintf("測試渠道 %s : %s 返回內容為：%s", channel.Name, request.Model, string(jsonBytes)))

	return nil, nil
}

func buildTestRequest() *types.ChatCompletionRequest {
	testRequest := &types.ChatCompletionRequest{
		Messages: []types.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "You just need to output 'hi' next.",
			},
		},
		Model:     "",
		MaxTokens: 2,
		Stream:    false,
	}
	return testRequest
}

func TestChannel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	channel, err := model.GetChannelById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	testModel := c.Query("model")
	tik := time.Now()
	err, openaiErr := testChannel(channel, testModel)
	tok := time.Now()
	milliseconds := tok.Sub(tik).Milliseconds()
	consumedTime := float64(milliseconds) / 1000.0

	success := false
	msg := ""
	if openaiErr != nil {
		if ShouldDisableChannel(channel.Type, openaiErr) {
			msg = fmt.Sprintf("測速失敗，已被禁用，原因：%s", err.Error())
			DisableChannel(channel.Id, channel.Name, err.Error(), false)
		} else {
			msg = fmt.Sprintf("測速失敗，原因：%s", err.Error())
		}
	} else if err != nil {
		msg = fmt.Sprintf("測速失敗，原因：%s", err.Error())
	} else {
		success = true
		msg = "測速成功"
		go channel.UpdateResponseTime(milliseconds)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"message": msg,
		"time":    consumedTime,
	})
}

var testAllChannelsLock sync.Mutex
var testAllChannelsRunning bool = false

func testAllChannels(isNotify bool) error {
	testAllChannelsLock.Lock()
	if testAllChannelsRunning {
		testAllChannelsLock.Unlock()
		return errors.New("測試已在運行中")
	}
	testAllChannelsRunning = true
	testAllChannelsLock.Unlock()
	channels, err := model.GetAllChannels()
	if err != nil {
		return err
	}
	var disableThreshold = int64(config.ChannelDisableThreshold * 1000)
	if disableThreshold == 0 {
		disableThreshold = 10000000 // a impossible value
	}
	go func() {
		var sendMessage string
		for _, channel := range channels {
			time.Sleep(config.RequestInterval)

			isChannelEnabled := channel.Status == config.ChannelStatusEnabled
			sendMessage += fmt.Sprintf("**通道 %s - #%d - %s** : \n\n", utils.EscapeMarkdownText(channel.Name), channel.Id, channel.StatusToStr())
			tik := time.Now()
			err, openaiErr := testChannel(channel, "")
			tok := time.Now()
			milliseconds := tok.Sub(tik).Milliseconds()
			// 通道为禁用状态，并且还是请求错误 或者 响应时间超过阈值 直接跳过，也不需要更新响应时间。
			if !isChannelEnabled {
				if err != nil {
					sendMessage += fmt.Sprintf("- 測試報錯: %s \n\n- 無需改變狀態，跳過\n\n", utils.EscapeMarkdownText(err.Error()))
					continue
				}
				if milliseconds > disableThreshold {
					sendMessage += fmt.Sprintf("- 響應時間 %.2fs 超過閾值 %.2fs \n\n- 無需改變狀態，跳過\n\n", float64(milliseconds)/1000.0, float64(disableThreshold)/1000.0)
					continue
				}
				// 如果已被禁用，但是请求成功，需要判断是否需要恢复
				// 手动禁用的通道，不会自动恢复
				if shouldEnableChannel(err, openaiErr) {
					if channel.Status == config.ChannelStatusAutoDisabled {
						EnableChannel(channel.Id, channel.Name, false)
						sendMessage += "- 已被啟用 \n\n"
					} else {
						sendMessage += "- 手動禁用的通道，不會自動恢復 \n\n"
					}
				}
			} else {
				// 如果通道启用状态，但是返回了错误 或者 响应时间超过阈值，需要判断是否需要禁用
				if milliseconds > disableThreshold {
					errMsg := fmt.Sprintf("響應時間 %.2fs 超過閾值 %.2fs ", float64(milliseconds)/1000.0, float64(disableThreshold)/1000.0)
					sendMessage += fmt.Sprintf("- %s \n\n- 禁用\n\n", errMsg)
					DisableChannel(channel.Id, channel.Name, errMsg, false)
					continue
				}

				if ShouldDisableChannel(channel.Type, openaiErr) {
					sendMessage += fmt.Sprintf("- 已被禁用，原因：%s\n\n", utils.EscapeMarkdownText(err.Error()))
					DisableChannel(channel.Id, channel.Name, err.Error(), false)
					continue
				}

				if err != nil {
					sendMessage += fmt.Sprintf("- 測試報錯: %s \n\n", utils.EscapeMarkdownText(err.Error()))
					continue
				}
			}
			channel.UpdateResponseTime(milliseconds)
			sendMessage += fmt.Sprintf("- 測試完成，耗時 %.2fs\n\n", float64(milliseconds)/1000.0)
		}
		testAllChannelsLock.Lock()
		testAllChannelsRunning = false
		testAllChannelsLock.Unlock()
		if isNotify {
			notify.Send("通道測試完成", sendMessage)
		}
	}()
	return nil
}

func TestAllChannels(c *gin.Context) {
	err := testAllChannels(true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func AutomaticallyTestChannels(frequency int) {
	if frequency <= 0 {
		return
	}

	for {
		time.Sleep(time.Duration(frequency) * time.Minute)
		logger.SysLog("testing all channels")
		_ = testAllChannels(false)
		logger.SysLog("channel test finished")
	}
}
