package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func bin(c *gin.Context) {
	platform := c.Query("platform")
	if platform == "" {
		c.String(403, "No Platform")
		return
	}

	if (platform != "linux") && (platform != "windows") {
		c.String(403, "Platform not Supported")
		return
	}

	result, err := serviceGlobal.bin(platform)
	if err != nil {
		println(err.Error())
		c.String(502, "Internal Error")
		return
	}

	c.String(200, strings.Join(result, "\n"))
	return
}
