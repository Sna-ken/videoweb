package base

import (
	"strconv"

	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Base struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

type RespWithData struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}

func ErrHTTPResp(c *app.RequestContext, err error) {
	e := errno.Convert(err)
	c.JSON(consts.StatusOK, Base{
		Code: strconv.Itoa(int(e.Code())),
		Msg:  e.Message(),
	})
}

func SuccessHTTPResp(c *app.RequestContext) {
	c.JSON(consts.StatusOK, Base{
		Code: strconv.Itoa(int(errno.Success.Code())),
		Msg:  errno.Success.Message(),
	})
}

func ErrHTTPRespWithData(c *app.RequestContext, err error, data any) {
	e := errno.Convert(err)
	c.JSON(consts.StatusOK, RespWithData{
		Code: strconv.Itoa(int(e.Code())),
		Msg:  e.Message(),
		Data: data,
	})
}

func SuccessHTTPRespWithData(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, RespWithData{
		Code: strconv.Itoa(int(errno.Success.Code())),
		Msg:  errno.Success.Message(),
		Data: data,
	})
}
