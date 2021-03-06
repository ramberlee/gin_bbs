package authorization

import (
	"gin_bbs/app/controllers"
	userModel "gin_bbs/app/models/user"
	authorizationRequest "gin_bbs/app/requests/api/authorization"
	"gin_bbs/pkg/errno"

	"github.com/gin-gonic/gin"
)

var (
	types = []string{"weixin"}
)

// Store 第三方登录
// @Summary 第三方登录
// @Tags authorization
// @Accept  json
// @Produce  json
// @Param social_type path string true "social_type in [weixin]"
// @Param json body authorization.Authorization true "微信 access_token openid 和 code，要么传 access_token openid 要么只传 code"
// @Success 200 {object} controllers.Response "{"token": 1}"
// @Router /api/socials/authorizations/{social_type} [post]
func Store(c *gin.Context) {
	var req *authorizationRequest.Authorization
	if err := c.ShouldBind(&req); err != nil {
		controllers.SendErrorResponse(c, errno.New(errno.ParamsError, err))
		return
	}
	req.SocialType = c.Param("social_type")
	weixinUserInfo, eno := req.Run()
	if eno != nil {
		controllers.SendErrorResponse(c, eno)
		return
	}

	var (
		user *userModel.User
		err  error
	)
	switch req.SocialType {
	case "weixin":
		if weixinUserInfo.Unionid != "" {
			user, err = userModel.GetByWeixinUnionID(weixinUserInfo.Unionid)
		} else {
			user, err = userModel.GetByWeixinOpenID(weixinUserInfo.OpenID)
		}

		if err != nil || user == nil {
			// 没有用户，默认创建一个用户
			user = &userModel.User{
				Name:         weixinUserInfo.NickName,
				Avatar:       weixinUserInfo.HeadImgURL,
				WeixinOpenID: weixinUserInfo.OpenID,
			}
			if weixinUserInfo.Unionid != "" {
				user.WeixinUnionID = weixinUserInfo.Unionid
			}

			if err := user.Create(); err != nil {
				controllers.SendErrorResponse(c, errno.New(errno.DatabaseError, err))
				return
			}
		}
	}

	controllers.SendOKResponse(c, map[string]uint{"token": user.ID})
}
