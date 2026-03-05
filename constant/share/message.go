package share

import "github.com/gin-gonic/gin"

func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func ResponseSuccess(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"success": message})
}

func RespondDate(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"data": data})
}
