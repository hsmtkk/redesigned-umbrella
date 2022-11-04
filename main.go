package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	h := newHandler()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	loginRequiredGroup := e.Group("/auth/*")
	loginRequiredGroup.Use(h.loginMiddleware)

	// Routes
	e.GET("/login", h.getLogin)
	e.POST("/login", h.postLogin)
	e.GET("/logout", h.logout)

	loginRequiredGroup.GET("/auth/foo", h.foo)
	loginRequiredGroup.GET("/auth/bar", h.bar)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
