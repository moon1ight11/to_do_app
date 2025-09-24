package api

import (
	"github.com/gin-gonic/gin"
)

func Router() {
	c := gin.Default()
	c.Run(":8080")
}
