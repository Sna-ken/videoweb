package base

import (
	"github.com/Sna-ken/videoweb/kitex_gen/model"
	"github.com/Sna-ken/videoweb/pkg/errno"
)

func BaseRPCResp(err error) *model.BaseResp {
	if err == nil {
		return &model.BaseResp{
			Code: errno.SuccessCode,
			Msg:  errno.Success.Message(),
		}
	}

	bizErr := errno.Convert(err)
	return &model.BaseResp{
		Code: bizErr.Code(),
		Msg:  bizErr.Message(),
	}
}

func SuccessRPCResp() *model.BaseResp {
	return BaseRPCResp(nil)
}
