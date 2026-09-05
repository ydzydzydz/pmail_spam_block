package service

import (
	"testing"

	"github.com/Jinnrry/pmail/models"
)

// TestUserService_GetUserID_Exists 用户存在时，GetUserID 返回用户 ID
func TestUserService_GetUserID_Exists(t *testing.T) {
	eng := newTestEngine(t)
	svc := NewUserService(eng)

	user := &models.User{
		Account:  "alice",
		Name:     "Alice",
		Password: "secret",
	}
	if _, err := eng.Insert(user); err != nil {
		t.Fatalf("插入预置用户失败: %v", err)
	}

	id, err := svc.GetUserID("alice")
	if err != nil {
		t.Fatalf("GetUserID 返回错误: %v", err)
	}
	if id != user.ID {
		t.Errorf("GetUserID = %d, 期望 %d", id, user.ID)
	}
}

// TestUserService_GetUserID_NotFound 用户不存在时，GetUserID 返回错误
func TestUserService_GetUserID_NotFound(t *testing.T) {
	eng := newTestEngine(t)
	svc := NewUserService(eng)

	id, err := svc.GetUserID("ghost")
	if err == nil {
		t.Fatalf("用户不存在时期望返回错误，实际返回 ID = %d", id)
	}
	if err.Error() != "user not found" {
		t.Errorf("错误信息 = %q, 期望 %q", err.Error(), "user not found")
	}
}
