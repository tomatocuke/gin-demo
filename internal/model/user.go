package model

var (
	UserPtr *User
)

type User struct {
	BaseID

	// 主要字段放在前边
	Nickname string `json:"nickname" gorm:"type:varchar(20); not null; default:''; comment:昵称;"`
	Account  string `json:"account" gorm:"uniqueIndex:uk_account; type:varchar(32); not null; default:''; comment:账号;"`

	// 数字尽量小，尽量注意内存对齐
	State  uint8 `json:"state" gorm:"type:tinyint unsigned; not null; default:3; comment:状态 1正常2冻结3未验证;"`
	Gender uint8 `json:"gender" gorm:"type:tinyint; not null; default:3; comment:性别 1男2女3未知;"`
	Age    uint8 `json:"age" gorm:"type:tinyint; not null; default:0; comment:年龄;"`

	Avatar   string `json:"avatar" gorm:"type:varchar(100); not null; default:''; comment:头像;"`
	Password string `json:"-" gorm:"column:passwd; type:varchar(100); not null; default:''; comment:密码;"`
	Salt     string `json:"-" gorm:"type:varchar(10); not null; default:''; comment:密码盐;"`

	BaseTimes
}

func (*User) GetByAccount(account string) (row User) {
	DB().Where("account = ?", account).Take(&row)
	return
}
