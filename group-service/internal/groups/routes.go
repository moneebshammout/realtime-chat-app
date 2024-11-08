package groups

import (
	"group-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Router(app *echo.Group) {
	logger.Infof("Adding Group Routes")
	groups := app.Group("/groups")

	groups.GET("", list)
	groups.GET("/:id", getGroup, middleware.ValidationMiddleware(&GroupGetSerizliser{}))
	groups.POST("", create, middleware.ValidationMiddleware(&GroupCreateSerizliser{}))
	groups.PUT("/:id/members/add", addMembers, middleware.ValidationMiddleware(&AddMembersSerizliser{}))
	groups.PUT("/:id/members/remove", removeMembers, middleware.ValidationMiddleware(&removeMembersSerizliser{}))
	groups.DELETE("/:id", deleteGroup, middleware.ValidationMiddleware(&IDParam{}))
}
