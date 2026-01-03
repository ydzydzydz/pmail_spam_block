package hook

import (
	_ "embed"
	"encoding/json"
	"strings"
	"time"

	"github.com/ydzydzydz/pmail_spam_block/logger"
	"github.com/ydzydzydz/pmail_spam_block/model"
)

var (
	//go:embed  dist/index.html
	SettingHtml string
)

type TestModelRequest struct {
	Setting model.SpamBlockSetting `json:"setting"`
	Content string                 `json:"content"`
}

// Response 响应体
type Response struct {
	Code    int    `json:"code"`    // 状态码 0 成功 -1 失败
	Message string `json:"message"` // 提示信息
	Data    any    `json:"data"`    // 数据
}

// SuccessResponse 成功响应
func SuccessResponse(message string, data any) *Response {
	return &Response{
		Code:    0,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(message string) *Response {
	return &Response{
		Code:    -1,
		Message: message,
	}
}

// Json 序列化响应体
func (r *Response) Json() string {
	json, err := json.Marshal(r)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("marshal response failed")
		return ""
	}
	return string(json)
}

// getSetting 获取Spam Block设置
func (h *SpamBlockHook) getSettingResponse(userID int) string {
	setting, err := h.settingService.GetSetting(userID)
	if err != nil {
		return ErrorResponse("获取Spam Block设置失败").Json()
	}

	return SuccessResponse("获取Spam Block设置成功", setting).Json()
}

// updateSetting 更新Spam Block设置
func (h *SpamBlockHook) updateSetting(userID int, requestData string) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("更新Spam Block设置")

	var setting model.SpamBlockSetting
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("反序列化设置请求失败")
		return ErrorResponse("反序列化设置请求失败").Json()
	}

	setting.UserID = userID
	setting.ApiUrl = strings.TrimSpace(setting.ApiUrl)
	if err := h.settingService.UpdateSetting(userID, &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("更新Spam Block设置失败")
		return ErrorResponse("更新Spam Block设置失败").Json()
	}

	return SuccessResponse("更新Spam Block设置成功", nil).Json()
}

// testModelResponse 测试Spam Block模型
func (h *SpamBlockHook) testModelResponse(requestData string) string {
	logger.PluginLogger.Info().Msg("测试Spam Block模型")
	var testModelRequest TestModelRequest
	if err := json.Unmarshal([]byte(requestData), &testModelRequest); err != nil {
		return ErrorResponse("反序列化插件设置失败").Json()
	}

	content := testModelRequest.Content
	respData, err := getModelResponse(content, testModelRequest.Setting.ApiUrl, time.Duration(testModelRequest.Setting.Timeout)*time.Millisecond)
	if err != nil {
		return ErrorResponse("获取模型响应失败").Json()
	}
	classes, err := getClasses(respData)
	if err != nil {
		return ErrorResponse("解析模型响应失败").Json()
	}

	maxScore := maxScore(classes)
	maxClass := maxClass(classes)

	if maxClass == CLASS_NORMAL {
		return SuccessResponse("测试Spam Block模型成功: 正常邮件", respData).Json()
	}

	if maxScore > testModelRequest.Setting.Threshold {
		if maxClass == CLASS_SPAM {
			return SuccessResponse("测试Spam Block模型成功: 诈骗邮件", respData).Json()
		} else {
			return SuccessResponse("测试Spam Block模型成功: 广告邮件", respData).Json()
		}
	}

	return SuccessResponse("测试Spam Block模型成功", respData).Json()
}
