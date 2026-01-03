package model

// SpamBlockSetting 垃圾邮件插件设置
type SpamBlockSetting struct {
	ID        int     `xorm:"id pk autoincr comment('主键')" json:"-"`
	UserID    int     `xorm:"user_id int index('idx_uid') comment('用户id') unique('idx_uid')" json:"-"`
	ApiUrl    string  `xorm:"api_url varchar(255) index('idx_url') index comment('api url')" json:"api_url"`
	Timeout   int     `xorm:"timeout int not null default 10 comment('超时时间')" json:"timeout"`
	Threshold float64 `xorm:"threshold float64 not null default 0.2 comment('阈值')" json:"threshold"`
}

// TableName 表名
func (u *SpamBlockSetting) TableName() string {
	return "plugin_spam_block_setting"
}
