package main

import (
	"encoding/json"
	"os"

	"github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/ydzydzydz/pmail_spam_block/hook"
	"github.com/ydzydzydz/pmail_spam_block/logger"
)

const (
	MAIN_CONFIG_FILE = "./config/config.json"
)

func MustReadConfig() *config.Config {
	content, err := os.ReadFile(MAIN_CONFIG_FILE)
	if err != nil {
		logger.PluginLogger.Panic().Err(err).Msg("主配置文件读取失败")
	}
	var cfg config.Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		logger.PluginLogger.Panic().Err(err).Msg("主配置文件解析失败")
	}
	logger.PluginLogger.Info().Msg("主配置文件读取成功")
	return &cfg
}

func main() {
	cfg := MustReadConfig()
	logger.PluginLogger.Info().Msg("配置文件读取成功")

	// 启动插件
	framework.CreatePlugin(
		hook.PLUGIN_NAME,
		hook.NewSpamBlockHook(cfg),
	).Run()
}
