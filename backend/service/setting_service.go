package service

import (
	"errors"

	"github.com/ydzydzydz/pmail_spam_block/model"
	"xorm.io/xorm"
)

// SettingService 设置服务
type SettingService struct {
	db *xorm.Engine
}

const (
	DefaultApiUrl    = ""    // 默认api url, 为空时, 则不发送
	DefaultTimeout   = 50000 // 默认超时时间
	DefaultThreshold = 0.2   // 默认阈值
)

// NewSettingService 创建设置服务实例
func NewSettingService(db *xorm.Engine) *SettingService {
	return &SettingService{db: db}
}

// GetSetting 获取设置
// 如果不存在, 则创建默认设置
func (s *SettingService) GetSetting(userID int) (*model.SpamBlockSettingModel, error) {
	has, err := s.db.Where("user_id = ?", userID).Exist(&model.SpamBlockSettingModel{})
	if err != nil {
		return nil, err
	}
	if !has {
		if err := s.CreateDefaultSetting(userID); err != nil {
			return nil, err
		}
	}
	setting := new(model.SpamBlockSettingModel)
	has, err = s.db.Where("user_id = ?", userID).Get(setting)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("setting not found")
	}
	return setting, nil
}

// UpdateSetting 更新设置
// 如果不存在, 则创建默认设置
func (s *SettingService) UpdateSetting(userID int, setting *model.SpamBlockSettingModel) error {
	has, err := s.db.Where("user_id = ?", setting.UserID).Exist(&model.SpamBlockSettingModel{})
	if err != nil {
		return err
	}
	if !has {
		if err := s.CreateDefaultSetting(setting.UserID); err != nil {
			return err
		}
	}
	_, err = s.db.Where("user_id = ?", userID).AllCols().Update(setting)
	return err
}

// CreateDefaultSetting 创建默认设置
func (s *SettingService) CreateDefaultSetting(userID int) error {
	setting := &model.SpamBlockSettingModel{
		UserID:    userID,
		ApiUrl:    DefaultApiUrl,
		Timeout:   DefaultTimeout,
		Threshold: DefaultThreshold,
	}
	_, err := s.db.Insert(setting)
	return err
}
