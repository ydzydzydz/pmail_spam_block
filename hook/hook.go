package hook

import (
	"fmt"
	"strings"
	"time"

	"github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/Jinnrry/pmail/models"
	"github.com/Jinnrry/pmail/utils/context"
	"github.com/ydzydzydz/pmail_spam_block/db"
	"github.com/ydzydzydz/pmail_spam_block/logger"
	"github.com/ydzydzydz/pmail_spam_block/model"
	"github.com/ydzydzydz/pmail_spam_block/service"
)

const (
	PLUGIN_NAME = "pmail_spam_block" // 插件名称
)

// Class 模型结果分类
type Class int

const (
	CLASS_NORMAL Class = iota // 正常邮件
	CLASS_AD                  // 广告邮件
	CLASS_SPAM                // 诈骗邮件
)

// EmailStatus 状态
type EmailStatus int

const (
	STATUS_NOT_SENT EmailStatus = iota // 未发送
	STATUS_SENT                        // 已发送
	STATUS_FAILED                      // 发送失败
	STATUS_DELETED                     // 已删除
	STATUS_AD       EmailStatus = 5    // 广告邮件
)

// SpamBlockHook 插件钩子
type SpamBlockHook struct {
	domain         string
	settingService *service.SettingService
	userService    *service.UserService
}

// NewSpamBlockHook 创建垃圾邮件插件钩子
func NewSpamBlockHook(cfg *config.Config) *SpamBlockHook {
	dataSource, err := db.NewDataSource(cfg)
	if err != nil {
		logger.PluginLogger.Fatal().Err(err).Msg("创建数据库连接失败")
	}
	logger.PluginLogger.Info().Msg("数据库初始化成功")

	settingService := service.NewSettingService(dataSource.SettingDao())
	userService := service.NewUserService(dataSource.UserDao())
	return &SpamBlockHook{
		domain:         cfg.Domain,
		settingService: settingService,
		userService:    userService,
	}
}

var _ framework.EmailHook = (*SpamBlockHook)(nil)

// GetName 获取插件名称
func (h *SpamBlockHook) GetName(ctx *context.Context) string { return PLUGIN_NAME }

// ReceiveSaveAfter 接收保存后的钩子
func (h *SpamBlockHook) ReceiveSaveAfter(ctx *context.Context, email *parsemail.Email, ue []*models.UserEmail) {
}

// ReceiveParseBefore 接收解析前的钩子
func (h *SpamBlockHook) ReceiveParseBefore(ctx *context.Context, email *[]byte) {}

// SendAfter 发送后的钩子
func (h *SpamBlockHook) SendAfter(ctx *context.Context, email *parsemail.Email, err map[string]error) {
}

// SendBefore 发送前的钩子
func (h *SpamBlockHook) SendBefore(ctx *context.Context, email *parsemail.Email) {}

// ReceiveParseAfter 接收解析后的钩子
func (h *SpamBlockHook) ReceiveParseAfter(ctx *context.Context, email *parsemail.Email) {
	// ctx 拿不到 user_id，也拿不到 is_admin，只能从 email 中获取
	for _, user := range email.To {
		// 只对本域用户进行处理
		account, domain := user.GetDomainAccount()
		if domain != h.domain {
			continue
		}
		userID, err := h.userService.GetUserID(account)
		// 对获取用户ID失败的用户，直接标记为已删除
		// 会影响到其他用户的邮件处理
		if err != nil {
			logger.PluginLogger.Error().Err(err).Str("account", account).Msg("获取用户ID失败")
			email.Status = int(STATUS_DELETED)
			break
		}
		// 对获取用户ID成功的用户，进行垃圾邮件处理逻辑
		h.spamBlock(userID, email)
	}
}

// SettingsHtml 获取设置 HTML
func (h *SpamBlockHook) SettingsHtml(ctx *context.Context, url string, requestData string) string {
	switch {
	// 获取用户设置
	case strings.Contains(url, "getSetting"):
		return h.getSettingResponse(ctx.UserID)
	// 更新设置
	case strings.Contains(url, "updateSetting"):
		return h.updateSetting(ctx.UserID, requestData)
	// 测试模型
	case strings.Contains(url, "testModel"):
		return h.testModelResponse(requestData)
	default:
		return SettingHtml
	}
}

// getSetting 获取用户设置
func (h *SpamBlockHook) getSetting(userID int) *model.SpamBlockSetting {
	setting, err := h.settingService.GetSetting(userID)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("获取用户设置失败")
		return nil
	}
	return setting
}

// spamBlock 垃圾邮件处理
func (h *SpamBlockHook) spamBlock(userID int, email *parsemail.Email) {
	setting := h.getSetting(userID)
	if setting == nil {
		logger.PluginLogger.Error().Int("user_id", userID).Msg("获取用户设置失败")
		return
	}
	if setting.ApiUrl == "" {
		logger.PluginLogger.Warn().Int("user_id", userID).Msg("模型API接口地址为空")
		return
	}
	content, err := getEmailContent(email)
	if err != nil {
		logger.PluginLogger.Warn().Err(err).Str("subject", email.Subject).Msg("获取邮件内容失败")
		return
	}

	respData, err := getModelResponse(fmt.Sprintf("%s %s", email.Subject, content), setting.ApiUrl, time.Duration(setting.Timeout)*time.Millisecond)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("获取模型响应失败")
		return
	}

	classes, err := getClasses(respData)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("获取分类结果失败")
		return
	}

	maxScore := maxScore(classes)
	maxClass := maxClass(classes)
	switch maxClass {
	case CLASS_NORMAL:
		logger.PluginLogger.Debug().Int("user_id", userID).Str("subject", email.Subject).Msg("邮件为正常邮件")
	case CLASS_AD:
		logger.PluginLogger.Info().Int("user_id", userID).Str("subject", email.Subject).Msg("邮件为广告邮件")
	case CLASS_SPAM:
		logger.PluginLogger.Info().Int("user_id", userID).Str("subject", email.Subject).Msg("邮件为垃圾邮件")
	}

	if setting.Threshold == 0 {
		setting.Threshold = 0.2
	}

	// 如果分类结果为正常邮件，直接返回
	if maxClass == CLASS_NORMAL {
		return
	}

	// 如果得分大于阈值，根据分类结果设置状态
	// 如果分类结果为诈骗邮件，设置状态为已删除
	// 如果分类结果为广告邮件，设置状态为广告邮件
	if maxScore > setting.Threshold {
		if maxClass == CLASS_SPAM {
			email.Status = int(STATUS_DELETED)
		} else {
			email.Status = int(STATUS_AD)
		}
	}
}
