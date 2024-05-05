package schemas

import (
	"cronJob/internal/global"
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Validatable interface {
	BindValidParam(c *gin.Context) error
}

type FormPage struct {
	PageNo   int `json:"pageNo" form:"pageNo" comment:"pageNo页数" example:"1" default:"1"`           //页数
	PageSize int `json:"pageSize" form:"pageSize" comment:"pageSize每页条数" example:"20" default:"10"` //每页条数
}

// DefaultGetValidParams 获取Get、Post、Uri请求参数
func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	var (
		errUri   error
		errQuery error
		errForm  error
	)

	errUri = c.ShouldBindUri(params)
	errQuery = c.ShouldBindQuery(params)
	errForm = c.ShouldBind(params)

	if errUri != nil && errQuery != nil && errForm != nil {
		return errForm
	}

	//获取验证器
	valid, err := GetValidator(c)
	if err != nil {
		return err
	}
	//获取翻译器
	trans, err := GetTranslation(c)
	if err != nil {
		return err
	}
	err = valid.Struct(params)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}

		return errors.New(strings.Join(sliceErrs, ","))
	}

	return nil
}

func GetValidator(c *gin.Context) (*validator.Validate, error) {
	val, ok := c.Get(global.ValidatorKey)
	if !ok {
		return nil, errors.New("未设置验证器")
	}
	validator, ok := val.(*validator.Validate)
	if !ok {
		return nil, errors.New("获取验证器失败")
	}
	return validator, nil
}

func GetTranslation(c *gin.Context) (ut.Translator, error) {
	trans, ok := c.Get(global.TranslatorKey)
	if !ok {
		return nil, errors.New("未设置翻译器")
	}
	translator, ok := trans.(ut.Translator)
	if !ok {
		return nil, errors.New("获取翻译器失败")
	}
	return translator, nil
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  "",
		Data: data,
	})
}

func ResponseError(ctx *gin.Context, code int, err error) {
	ctx.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  err.Error(),
		Data: "",
	})
}

const (
	TaskListParamInvalid = 1001 + iota
	TaskCreateError
	TaskUpdateError
	TaskNotExist // 服务ID错误
	TaskAddParamInvalid
	TaskSearchParamInvalid
	TaskSearchListInvalid
	TaskAlreadyExist
	TaskParamInvalid
	TaskLogParamInvalid
	TaskLogListInvalid
	TaskLogDeleteInvalid

	UserSearchParamInvalid
	UserSearchListInvalid
	UserEditParamInvalid
	UserNotExist
	UserPasswordError
	UserUpdateFailed
	UserCreateFailed
	UserDeleteFailed

	RoleCreateFailed
	RoleListFailed
	RoleNotExist

	LoginParamInvalid
)
