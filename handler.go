package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type myHandler struct {
	sessionUserMap map[string]string
}

func newHandler() *myHandler {
	sessionUserMap := map[string]string{}
	return &myHandler{sessionUserMap}
}

const loginHTML = `<html>
 <body>
  <form action="/login" method="POST">
   <input type="text" name="userid">
   <input type="submit" value="submit">
  </form>
 </body>
</html>`

func (h *myHandler) getLogin(ectx echo.Context) error {
	return ectx.HTML(http.StatusOK, loginHTML)
}

func (h *myHandler) postLogin(ectx echo.Context) error {
	userID := ectx.FormValue("userid")
	sessionID := uuid.New().String()
	h.sessionUserMap[sessionID] = userID
	sess, err := session.Get("session", ectx)
	if err != nil {
		return err
	}
	sess.Values["sessionid"] = sessionID
	sess.Save(ectx.Request(), ectx.Response())
	return ectx.Redirect(http.StatusFound, "/auth/foo")
}

func (h *myHandler) logout(ectx echo.Context) error {
	sess, err := session.Get("session", ectx)
	if err != nil {
		return err
	}
	sessionID := sess.Values["sessionid"].(string)
	delete(h.sessionUserMap, sessionID)
	return ectx.Redirect(http.StatusFound, "/login")
}

func (h *myHandler) foo(ectx echo.Context) error {
	userID, err := h.userIDFromSession(ectx)
	if err != nil {
		return err
	}
	return ectx.String(http.StatusOK, fmt.Sprintf("Hello %s from foo", userID))
}

func (h *myHandler) bar(ectx echo.Context) error {
	userID, err := h.userIDFromSession(ectx)
	if err != nil {
		return err
	}
	return ectx.String(http.StatusOK, fmt.Sprintf("Hello %s from bar", userID))
}

func (h *myHandler) loginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	handler := func(ectx echo.Context) error {
		log.Print("loginMiddleware is being invoked")
		_, err := h.userIDFromSession(ectx)
		if err != nil {
			return ectx.Redirect(http.StatusFound, "/login")
		}
		return next(ectx)
	}
	return handler
}

func (h *myHandler) userIDFromSession(ectx echo.Context) (string, error) {
	sess, err := session.Get("session", ectx)
	if err != nil {
		return "", err
	}
	sessionID, ok := sess.Values["sessionid"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get sessionid")
	}
	userID, ok := h.sessionUserMap[sessionID]
	if !ok {
		return "", fmt.Errorf("session does not exist")
	}
	return userID, nil
}
