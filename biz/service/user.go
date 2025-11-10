package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/cache"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/biz/service/taskqueue"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/crypt"
	"judgeMore/pkg/errno"
	"judgeMore/pkg/utils"
	"strings"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{
		ctx: ctx,
		c:   c,
	}
}

func (svc *UserService) Login(req *model.User) (*model.User, error) {
	exist, err := mysql.IsUserExist(svc.ctx, req)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceUserNotExistCode, "user not exist")
	}
	// 密码检验
	userInfo, err := mysql.GetUserInfoByRoleId(svc.ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	// 激活检验
	if userInfo.Status == 0 {
		return nil, errno.NewErrNo(errno.ServiceUserDeathCode, "user not active ")
	}
	if !crypt.VerifyPassword(req.Password, userInfo.Password) {
		return nil, errno.Errorf(errno.ServiceUserPasswordError, "password not match")
	}
	return userInfo, nil
}

func (svc *UserService) Register(req *model.User) (string, error) {
	exist, err := mysql.IsUserExist(svc.ctx, req)
	if err != nil {
		return "", err
	}
	if exist {
		return "", errno.NewErrNo(errno.ServiceUserExistCode, "user already exist")
	}
	req.Password, err = crypt.PasswordHash(req.Password)
	if err != nil {
		return "", fmt.Errorf("hash password failed: %w", err)
	}
	//验证邮箱
	err = svc.SendEmail(req)
	if err != nil {
		return "", err
	}
	// 创建账户
	req.Role = "student"
	req.Status = 0
	uid, err := mysql.CreateUser(svc.ctx, req)
	if err != nil {
		return "", err
	}
	return uid, nil
}
func (svc *UserService) QueryUserInfo(u *model.User) (UserInfo *model.User, err error) {
	//存在性检验
	exist, err := mysql.IsUserExist(svc.ctx, u)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceUserNotExistCode, "user not exist")
	}
	userInfo, err := mysql.GetUserInfoByRoleId(svc.ctx, u.Uid)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (svc *UserService) UpdateUserInfo(u *model.User) (UserInfo *model.User, err error) {
	// 由于uid读取自token 所以不做存在性检验
	id := GetUserIDFromContext(svc.c)
	u.Uid = id
	// 检验学院和专业
	exist, err := IsMajorExist(svc.ctx, u.Major)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceMajorNotExistCode, "major not exist")
	}
	exist, err = IsCollegeExist(svc.ctx, u.College)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceCollegeExistCode, "college not exist")
	}
	// 检验年级
	if !utils.IsGradeValid(u.Grade) {
		return nil, errno.NewErrNo(errno.ServiceGradeNotExistCode, "invalid grade")
	}
	userInfo, err := svc.UpdateUser(svc.ctx, u)
	if err != nil {
		return nil, err
	}
	taskqueue.AddStuToAdminRelation(svc.ctx, constants.UserKey, u)
	return userInfo, nil
}

func (svc *UserService) UpdateUserPassword(u *model.User, c string) error {
	exist, err := mysql.IsUserExist(svc.ctx, u)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceUserNotExistCode, "user not exist")
	}
	// 验证code
	email := fmt.Sprintf("%s%s", u.Uid, constants.EmailSuffix)
	key := fmt.Sprintf("Email:%s", email)
	exist = cache.IsKeyExist(svc.ctx, key)
	if !exist {
		return errors.New("code expired")
	}
	code, err := cache.GetCodeCache(svc.ctx, key)
	if err != nil {
		return err
	}
	if code != c {
		return errors.New("code not match")
	}
	//
	u.Password, err = crypt.PasswordHash(u.Password)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)

	}
	return mysql.UpdateUserPassword(svc.ctx, u)
}

func (svc *UserService) Email(userInfo *model.User) error {
	key := fmt.Sprintf("Email:%s", userInfo.Email)
	exist := cache.IsKeyExist(svc.ctx, key)
	if exist {
		return errno.NewErrNo(errno.ServiceEmailWaitCode, "wait tow minute")
	}
	return svc.SendEmail(userInfo)
}

func (svc *UserService) SendEmail(user *model.User) error {
	// 首先进行验证 学号即Uid 与fzu邮箱强绑定
	Correct := strings.HasSuffix(user.Email, constants.EmailSuffix) && len(user.Email) == constants.EmailLength
	if !Correct {
		return errno.NewErrNo(errno.ServiceEmailIncorrectCode, "Uid do not match email")
	}
	key := fmt.Sprintf("Email:%s", user.Email)
	code, err := cache.PutCodeToCache(svc.ctx, key)
	if err != nil {
		return err
	}
	err = utils.MailSendCode(user.Email, code)
	if err != nil {
		return err
	}
	return nil
}
func (svc *UserService) Logout() error {
	token_id := GetTokenIdFromContext(svc.c)
	err := cache.PutTokenIdToCache(svc.ctx, token_id)
	if err != nil {
		return err
	}
	return nil
}

func (svc *UserService) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// 由于这里需要对需更新的内容做选择 在svc处处理
	var updateParams []string
	if user.Major != "" {
		updateParams = append(updateParams, user.Major)
	}
	if user.College != "" {
		updateParams = append(updateParams, user.College)
	}
	if user.Grade != "" {
		updateParams = append(updateParams, user.Grade)
	}
	// 如果有需要更新的字段才执行
	if len(updateParams) > 0 {
		return mysql.UpdateInfoByRoleId(ctx, user.Uid, updateParams...)
	}

	return nil, nil

}

func (svc *UserService) VerifyEmail(data *model.EmailAuth) (err error) {
	// 判断存不存在
	key := fmt.Sprintf("Email:%s", data.Email)
	exist := cache.IsKeyExist(svc.ctx, key)
	if !exist {
		return errno.NewErrNo(errno.ServiceCodeExpired, "code expired")
	}
	emailParts := strings.Split(data.Email, "@")
	localPart := emailParts[0]
	data.Uid = localPart
	code, err := cache.GetCodeCache(svc.ctx, key)
	if err != nil {
		return err
	}
	if code != data.Code {
		return errno.NewErrNo(errno.ServiceCodeNotMatched, "code not match")
	}
	// 更新user表的信息
	err = mysql.ActivateUser(svc.ctx, data.Uid)
	if err != nil {
		return err
	}
	err = cache.DeleteCodeCache(svc.ctx, key)
	if err != nil {
		return err
	}
	return nil
}
