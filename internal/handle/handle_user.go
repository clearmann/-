package handle

import (
	jwt "backend/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserInfo struct {
}
type tokenRequest struct {
	Token string `json:"token"`
}

func (*UserInfo) Info(c *gin.Context) {
	var token tokenRequest
	err := c.ShouldBindJSON(&token)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "格式错误",
		})
		return
	}
	fmt.Println(token.Token)
	claim, err := jwt.ParseToken(token.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"id":       claim.ID,
		"username": claim.Username,
		"email":    claim.Email,
		"message":  "success",
	})
}
