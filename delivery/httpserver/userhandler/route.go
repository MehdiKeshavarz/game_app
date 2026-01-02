package userhandler

import "github.com/labstack/echo/v4"
 
func (h Handler) SetUserRoutes(e *echo.Echo) {
	groups := e.Group("/users")

	groups.POST("/register", h.userRegister)
	groups.POST("/login", h.userLogin)
	groups.GET("/profile", h.userProfile)

}
