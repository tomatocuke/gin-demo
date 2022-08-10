package service

import (
	"app/internal/common"
	"app/internal/common/constant"
	"app/internal/common/dto"
	"app/internal/model"
	"app/pkg"
	"app/pkg/response"
	"app/pkg/token"
	"sync"

	"github.com/jinzhu/copier"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

var (
	signupMu sync.Mutex
)

type User struct{}

// 发送验证码，自己去设计
func (User) SendCode(req dto.UserSendCodeReq, resp *dto.UserSendCodeResp) error {

	return nil
}

func (u User) SignUp(req dto.UserSignUpReq, resp *dto.UserSignUpResp) error {
	signupMu.Lock()
	defer signupMu.Unlock()

	// 检测账号重复
	var count int64
	model.UserPtr.DB().Where("account = ?", req.Account).Count(&count)
	if count > 0 {
		return response.Msg("账号已被注册")
	}

	// TODO 检查验证码

	salt := pkg.Rand.String(10)
	password := pkg.Encry.HmacSha256(req.Password, salt)

	return common.DB.Transaction(func(tx *gorm.DB) error {
		user := model.User{
			Account:  req.Account,
			Password: password,
			Salt:     salt,
			State:    constant.UserStateOk,
		}
		err := tx.Create(&user).Error
		if err != nil {
			return err
		}

		// TODO 其它逻辑，相关表

		log.Info("新注册账号: " + req.Account)

		return nil
	})
}

func (User) SignIn(req dto.UserSignInReq, resp *dto.UserSignInResp) error {
	user := model.UserPtr.GetByAccount(req.Account)
	if !user.IsValid() {
		return response.Msg("账号不存在")
	}
	if user.Password != pkg.Encry.HmacSha256(req.Password, user.Salt) {
		return response.Msg("密码错误")
	}
	if user.State == constant.UserStateFrozen {
		return response.Msg("账号已被禁用")
	}

	*resp = token.GenerateTokens(int(user.Id))

	log.Info("用户登录：", user.Id)

	return nil
}

func (u User) PostInfo(uid int, req dto.UserPostInfoReq, resp *dto.UserGetInfoResp) error {
	err := common.DB.Model(model.UserPtr).Take(&resp, uid).Error
	if err != nil {
		return err
	}

	copier.Copy(&resp, &req)

	return common.DB.Model(model.UserPtr).Where("id = ?", resp.Id).Updates(&req).Error
}

func (u User) GetInfo(uid int, req dto.UserGetInfoReq, resp *dto.UserGetInfoResp) error {
	err := model.UserPtr.DB().Take(&resp, uid).Error
	return err
}
