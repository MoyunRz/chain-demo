package web

import (
	"strconv"

	"chain-demo/cloudserver/hq-echo/model"
)

func UserHandler(c *Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		panic(err)
	}

	var User model.User
	u := User.GetUserById(id)

	c.Set("tmpl", "web/user")
	c.Set("data", map[string]interface{}{
		"title": "User",
		"user":  u,
	})

	return nil
}
