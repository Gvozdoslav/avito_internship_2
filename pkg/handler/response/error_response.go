package response

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, httpCode int, msg string) {
	logrus.Errorf(msg)
	c.AbortWithStatusJSON(httpCode, Error{msg}) // blocks next handlers, writes resp in JSON
}
