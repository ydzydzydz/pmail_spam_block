package controller

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ydzydzydz/pmail_spam_block/classifier"
	"github.com/ydzydzydz/pmail_spam_block/controller/response"
	"github.com/ydzydzydz/pmail_spam_block/logger"
	"github.com/ydzydzydz/pmail_spam_block/model"
	"github.com/ydzydzydz/pmail_spam_block/service"
)

// TestModelRequest 测试模型请求
type TestModelRequest struct {
	Setting model.SpamBlockSettingModel `json:"setting"`
	Content string                      `json:"content"`
}

// SettingController 设置控制器
type SettingController struct {
	settingService *service.SettingService
}

// NewSettingController 创建设置控制器实例
func NewSettingController(settingService *service.SettingService) *SettingController {
	return &SettingController{settingService: settingService}
}

// GetSetting 获取Spam Block设置
func (c *SettingController) GetSetting(userID int) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("获取Spam Block设置")
	setting, err := c.settingService.GetSetting(userID)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("获取Spam Block设置失败")
		return response.ErrorResponse("获取Spam Block设置失败").Json()
	}
	return response.SuccessResponse("获取Spam Block设置成功", setting).Json()
}

// UpdateSetting 更新Spam Block设置
func (c *SettingController) UpdateSetting(userID int, requestData string) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("更新Spam Block设置")

	var setting model.SpamBlockSettingModel
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("反序列化设置请求失败")
		return response.ErrorResponse("反序列化设置请求失败").Json()
	}

	setting.UserID = userID
	setting.ApiUrl = strings.TrimSpace(setting.ApiUrl)
	if err := c.settingService.UpdateSetting(userID, &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("更新Spam Block设置失败")
		return response.ErrorResponse("更新Spam Block设置失败").Json()
	}

	return response.SuccessResponse("更新Spam Block设置成功", nil).Json()
}

// TestModel 测试Spam Block模型
func (c *SettingController) TestModel(requestData string) string {
	logger.PluginLogger.Info().Msg("测试Spam Block模型")

	var testModelRequest TestModelRequest
	if err := json.Unmarshal([]byte(requestData), &testModelRequest); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("反序列化插件设置失败")
		return response.ErrorResponse("反序列化插件设置失败").Json()
	}

	content := testModelRequest.Content
	respData, err := classifier.GetModelResponse(content, testModelRequest.Setting.ApiUrl, time.Duration(testModelRequest.Setting.Timeout)*time.Millisecond)
	if err != nil {
		return response.ErrorResponse("获取模型响应失败").Json()
	}
	classes, err := classifier.GetClasses(respData)
	if err != nil {
		return response.ErrorResponse("解析模型响应失败").Json()
	}

	maxScore := classifier.MaxScore(classes)
	maxClass := classifier.MaxClass(classes)

	if maxClass == classifier.ClassNormal {
		return response.SuccessResponse("测试Spam Block模型成功: 正常邮件", respData).Json()
	}

	if maxScore > testModelRequest.Setting.Threshold {
		if maxClass == classifier.ClassSpam {
			return response.SuccessResponse("测试Spam Block模型成功: 诈骗邮件", respData).Json()
		}
		return response.SuccessResponse("测试Spam Block模型成功: 广告邮件", respData).Json()
	}

	return response.SuccessResponse("测试Spam Block模型成功", respData).Json()
}
