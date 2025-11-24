package callback

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	"gorm.io/gorm"
)

func getInt64FromParam(c *gin.Context, key string) int64 {
	_value := c.Param(key)
	value, err := strconv.ParseInt(_value, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad param, err: " + err.Error(),
		})
	}
	return value
}

func getUserIdFromToken(c *gin.Context) int64 {
	userId := c.GetInt64("userId")
	return userId
}

func getModeFromToken(c *gin.Context) string {
	return c.GetString("mode")
}

func handleDBError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
	} else {
		elog.Error("DB error", l.E(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "DB error" + err.Error(),
		})
	}
}
