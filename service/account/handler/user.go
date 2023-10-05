package handler

import (
	"context"
	"github.com/lilhammer111/hammer-cloud/common"
	"github.com/lilhammer111/hammer-cloud/config"
	"github.com/lilhammer111/hammer-cloud/db"
	"github.com/lilhammer111/hammer-cloud/service/account/proto"
	"github.com/lilhammer111/hammer-cloud/util"
)

type User struct{}

func (User) Signup(ctx context.Context, req *proto.ReqSignup, res *proto.RespSignup) error {
	username := req.Username
	pwd := req.Password

	if len(username) < 3 || len(pwd) < 5 {
		res.Code = common.StatusParamInvalid
		res.Message = "register parameter invalid"
		return nil
	}

	encryptPWD := util.Sha1([]byte(pwd + config.PasswordSalt))

	if ok := db.UserSignUp(username, encryptPWD); ok {
		res.Code = common.StatusOK
		res.Message = "register succeeded"
	} else {
		res.Code = common.StatusRegisterFailed
		res.Message = "register failed"
	}
	return nil
}
