package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/Jinnrry/pmail/hooks/spam_block/tools"
	"github.com/ydzydzydz/pmail_spam_block/logger"
)

type ModelResponse struct {
	Predictions [][]float64 `json:"predictions"`
}

type ModelRequest struct {
	Instances []InstanceItem `json:"instances"`
}

type InstanceItem struct {
	Token []string `json:"token"`
}

// getEmailContent 获取邮件内容
func getEmailContent(email *parsemail.Email) (string, error) {
	content := tools.Trim(tools.TrimHtml(string(email.HTML)))
	if content == "" {
		content = tools.Trim(string(email.Text))
	}
	if content == "" {
		return "", fmt.Errorf("邮件内容为空")
	}
	return content, nil
}

// getModelResponse 获取模型响应
func getModelResponse(content string, url string, timeout time.Duration) (respData *ModelResponse, err error) {
	reqData := ModelRequest{
		Instances: []InstanceItem{
			{
				Token: []string{
					content,
				},
			},
		},
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("序列化请求数据失败")
		return nil, err
	}

	if timeout == 0 {
		timeout = 5000 * time.Millisecond
	}
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("发送请求失败")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("读取响应体失败")
		return nil, err
	}
	if err := json.Unmarshal(body, &respData); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("解析响应数据失败")
		return nil, err
	}
	return respData, nil
}

// getClasses 获取分类结果
func getClasses(respData *ModelResponse) ([]float64, error) {
	if len(respData.Predictions) == 0 {
		return nil, fmt.Errorf("响应数据格式错误")
	}
	classes := respData.Predictions[0]
	if len(classes) != 3 {
		return nil, fmt.Errorf("响应数据格式错误")
	}
	return classes, nil
}

// maxScore 获取最大分数
func maxScore(classes []float64) float64 {
	var maxScore float64
	for _, score := range classes {
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}

// maxClass 获取最大分数对应的分类
func maxClass(classes []float64) Class {
	var maxScore float64
	var maxClass int
	for i, score := range classes {
		if score > maxScore {
			maxScore = score
			maxClass = i
		}
	}
	return Class(maxClass)
}
