package service

import (
	"testing"

	"github.com/Jinnrry/pmail/models"
	"github.com/ydzydzydz/pmail_spam_block/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

// newTestEngine 创建基于内存 SQLite 的 xorm 引擎并同步表结构，测试结束自动关闭
func newTestEngine(t *testing.T) *xorm.Engine {
	t.Helper()
	eng, err := xorm.NewEngine("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("创建内存数据库失败: %v", err)
	}
	// 内存库每个连接独立，强制单连接以保证表结构与数据可见
	eng.SetMaxOpenConns(1)
	eng.SetMaxIdleConns(1)
	if err := eng.Sync2(new(model.SpamBlockSettingModel), new(models.User)); err != nil {
		t.Fatalf("同步表结构失败: %v", err)
	}
	t.Cleanup(func() { _ = eng.Close() })
	return eng
}

// TestSettingService_GetSetting_CreatesDefaultWhenNotExist 设置不存在时，GetSetting 自动创建默认设置
func TestSettingService_GetSetting_CreatesDefaultWhenNotExist(t *testing.T) {
	eng := newTestEngine(t)
	svc := NewSettingService(eng)

	got, err := svc.GetSetting(1001)
	if err != nil {
		t.Fatalf("GetSetting 返回错误: %v", err)
	}
	if got.UserID != 1001 {
		t.Errorf("UserID = %d, 期望 1001", got.UserID)
	}
	if got.ApiUrl != DefaultApiUrl {
		t.Errorf("ApiUrl = %q, 期望 %q", got.ApiUrl, DefaultApiUrl)
	}
	if got.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %d, 期望 %d", got.Timeout, DefaultTimeout)
	}
	if got.Threshold != DefaultThreshold {
		t.Errorf("Threshold = %v, 期望 %v", got.Threshold, DefaultThreshold)
	}

	// 默认设置已落库，且只有一条记录
	cnt, err := eng.Where("user_id = ?", 1001).Count(new(model.SpamBlockSettingModel))
	if err != nil {
		t.Fatalf("统计设置记录失败: %v", err)
	}
	if cnt != 1 {
		t.Errorf("设置记录数 = %d, 期望 1", cnt)
	}

	// 再次获取应返回同一条记录，不会重复创建
	got2, err := svc.GetSetting(1001)
	if err != nil {
		t.Fatalf("第二次 GetSetting 返回错误: %v", err)
	}
	if got2.ID != got.ID {
		t.Errorf("第二次获取的 ID = %d, 期望与首次一致 %d", got2.ID, got.ID)
	}
}

// TestSettingService_GetSetting_ReturnsExisting 设置已存在时，GetSetting 直接返回已有记录
func TestSettingService_GetSetting_ReturnsExisting(t *testing.T) {
	eng := newTestEngine(t)
	svc := NewSettingService(eng)

	seed := &model.SpamBlockSettingModel{
		UserID:    2002,
		ApiUrl:    "http://model.example/api",
		Timeout:   1234,
		Threshold: 0.8,
	}
	if _, err := eng.Insert(seed); err != nil {
		t.Fatalf("插入预置设置失败: %v", err)
	}

	got, err := svc.GetSetting(2002)
	if err != nil {
		t.Fatalf("GetSetting 返回错误: %v", err)
	}
	if got.ID != seed.ID {
		t.Errorf("ID = %d, 期望 %d", got.ID, seed.ID)
	}
	if got.ApiUrl != seed.ApiUrl {
		t.Errorf("ApiUrl = %q, 期望 %q", got.ApiUrl, seed.ApiUrl)
	}
	if got.Timeout != seed.Timeout {
		t.Errorf("Timeout = %d, 期望 %d", got.Timeout, seed.Timeout)
	}
	if got.Threshold != seed.Threshold {
		t.Errorf("Threshold = %v, 期望 %v", got.Threshold, seed.Threshold)
	}
}

// TestSettingService_UpdateSetting_Existing 设置已存在时，UpdateSetting 更新该记录且不新增
func TestSettingService_UpdateSetting_Existing(t *testing.T) {
	eng := newTestEngine(t)
	svc := NewSettingService(eng)

	if _, err := svc.GetSetting(3003); err != nil {
		t.Fatalf("预置默认设置失败: %v", err)
	}

	updated := &model.SpamBlockSettingModel{
		UserID:    3003,
		ApiUrl:    "http://new.example/api",
		Timeout:   999,
		Threshold: 0.5,
	}
	if err := svc.UpdateSetting(3003, updated); err != nil {
		t.Fatalf("UpdateSetting 返回错误: %v", err)
	}

	got := new(model.SpamBlockSettingModel)
	has, err := eng.Where("user_id = ?", 3003).Get(got)
	if err != nil {
		t.Fatalf("查询更新后的设置失败: %v", err)
	}
	if !has {
		t.Fatal("更新后设置记录不存在")
	}
	if got.ApiUrl != updated.ApiUrl {
		t.Errorf("ApiUrl = %q, 期望 %q", got.ApiUrl, updated.ApiUrl)
	}
	if got.Timeout != updated.Timeout {
		t.Errorf("Timeout = %d, 期望 %d", got.Timeout, updated.Timeout)
	}
	if got.Threshold != updated.Threshold {
		t.Errorf("Threshold = %v, 期望 %v", got.Threshold, updated.Threshold)
	}

	// 仍然只有一条记录
	cnt, err := eng.Where("user_id = ?", 3003).Count(new(model.SpamBlockSettingModel))
	if err != nil {
		t.Fatalf("统计设置记录失败: %v", err)
	}
	if cnt != 1 {
		t.Errorf("设置记录数 = %d, 期望 1", cnt)
	}
}

// TestSettingService_UpdateSetting_NotExist_CreatesAndUpdates 设置不存在时，UpdateSetting 先建默认再更新
func TestSettingService_UpdateSetting_NotExist_CreatesAndUpdates(t *testing.T) {
	eng := newTestEngine(t)
	svc := NewSettingService(eng)

	setting := &model.SpamBlockSettingModel{
		UserID:    4004,
		ApiUrl:    "http://x.example/api",
		Timeout:   1,
		Threshold: 0.1,
	}
	if err := svc.UpdateSetting(4004, setting); err != nil {
		t.Fatalf("UpdateSetting 返回错误: %v", err)
	}

	got := new(model.SpamBlockSettingModel)
	has, err := eng.Where("user_id = ?", 4004).Get(got)
	if err != nil {
		t.Fatalf("查询设置失败: %v", err)
	}
	if !has {
		t.Fatal("设置记录未创建")
	}
	if got.ApiUrl != setting.ApiUrl {
		t.Errorf("ApiUrl = %q, 期望 %q", got.ApiUrl, setting.ApiUrl)
	}
	if got.Timeout != setting.Timeout {
		t.Errorf("Timeout = %d, 期望 %d", got.Timeout, setting.Timeout)
	}
	if got.Threshold != setting.Threshold {
		t.Errorf("Threshold = %v, 期望 %v", got.Threshold, setting.Threshold)
	}
}

// TestSettingService_CreateDefaultSetting_DuplicateUserFails 同一用户重复创建默认设置应失败（user_id 唯一索引）
func TestSettingService_CreateDefaultSetting_DuplicateUserFails(t *testing.T) {
	eng := newTestEngine(t)
	svc := NewSettingService(eng)

	if err := svc.CreateDefaultSetting(5005); err != nil {
		t.Fatalf("首次创建默认设置失败: %v", err)
	}
	if err := svc.CreateDefaultSetting(5005); err == nil {
		t.Error("重复创建默认设置期望返回错误，实际为 nil")
	}
}
