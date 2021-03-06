package notification

import (
	notificationModel "gin_bbs/app/models/notification"
	userModel "gin_bbs/app/models/user"
	"gin_bbs/pkg/ginutils/pagination"
	"gin_bbs/app/controllers"

	"github.com/gin-gonic/gin"
	"gin_bbs/pkg/ginutils/flash"
)

// Index 消息列表
func Index(c *gin.Context, currentUser *userModel.User) {
	renderFunc, err := pagination.CreatePage(c, 20, "notifications", notificationModel.AllCount,
		func(offset, limit, _, _ int) (interface{}, error) {
			return notificationModel.List(userModel.TableName, currentUser.ID, offset, limit)
		})

	if err != nil {
		controllers.Render(c, "notifications/index", gin.H{"error": "错误: " + err.Error()})
		return
	}

	currentUser.NotificationCount = 0
	if err = currentUser.Update(); err != nil {
		flash.NewDangerFlash(c, "NotificationCount 更新失败！")
	}
	if err = notificationModel.Read(userModel.TableName, currentUser.ID); err != nil {
		flash.NewDangerFlash(c, "notificationModel 更新失败！")
	}

	controllers.Render(c, "notifications/index", renderFunc(gin.H{}))
}
