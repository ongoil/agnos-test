package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogError(
	c *gin.Context,
	message string,
	err error,
	fields logrus.Fields,
) {
	logrus.WithFields(fields).WithError(err).Error(message)
}
