package middleware

import (
	"context"
	"one-api/common/logger"
	"one-api/common/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := utils.GetTimeString() + utils.GetRandomString(8)
		c.Set(logger.RequestIdKey, id)
		c.Set("requestStartTime", time.Now())
		ctx := context.WithValue(c.Request.Context(), logger.RequestIdKey, id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(logger.RequestIdKey, id)
		c.Next()
	}
}
