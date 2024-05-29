package handle

import (
	g "backend/internal/global"
	"backend/internal/model"
	jwt "backend/internal/utils"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"regexp"
)

type UserAuth struct{}

type LoginVO struct {

	// 点赞 Set: 用于记录用户点赞过的文章, 评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
}

func isValidEmail(email string) bool {
	// 正则表达式检查邮箱格式
	// 这里使用一个简单的正则表达式，实际应用中可能需要更复杂的规则
	// 可以根据需求自行调整正则表达式
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}

type RegisterRequest struct {
	Email           string `json:"email" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

func (*UserAuth) Register(c *gin.Context) {
	var registerRequest RegisterRequest
	db := GetDB(c)
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		log.Println("参数格式错误")
		log.Println(err.Error())
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "参数格式错误",
		})
		return
	}
	if !isValidEmail(registerRequest.Email) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "邮箱格式错误...",
		})
		return
	}
	if registerRequest.Username == "" || registerRequest.Password == "" || registerRequest.Email == "" {
		log.Println("账号或密码不能为空")
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "账号或密码不能为空",
		})
		return
	}
	if emailExist, _ := model.EmailExist(db, registerRequest.Email); emailExist {
		log.Println("邮箱已经存在")
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "邮箱已经存在",
		})
		return
	}
	if usernameExist, _ := model.UsernameExist(db, registerRequest.Username); usernameExist {
		log.Println("用户名已经存在")
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "用户名已经存在",
		})
		return
	}
	if registerRequest.Password != registerRequest.ConfirmPassword {
		log.Println("账号和密码必须一致")
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "两次密码须一致",
		})
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		panic("密码加密错误...")
	}
	err = model.CreateUser(db, registerRequest.Username, string(passwordHash), registerRequest.Email)
	if err != nil {
		log.Println(err.Error())
		panic("注册用户失败")
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

type LoginRequest struct {
	AccountName string `json:"accountName" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func (*UserAuth) Token(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Println("参数错误")
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "参数错误",
		})
		return
	}
	db := GetDB(c)
	var user *model.User
	var err error
	if isValidEmail(loginRequest.AccountName) {
		//通过 email 进行登录
		log.Println("邮箱登录")
		user, err = model.GetPasswordByEmail(db, loginRequest.AccountName)
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				c.JSON(200, gin.H{
					"code":    200,
					"message": "error",
					"data":    "电子邮箱不存在",
				})
				return
			}
		}
	} else {
		//通过用户名 进行登录
		log.Println("用户名登录")
		user, err = model.GetPasswordByUsername(db, loginRequest.AccountName)
		if err != nil {
			log.Println(err.Error())
			if errors.Is(gorm.ErrRecordNotFound, err) {
				c.JSON(200, gin.H{
					"code":    200,
					"message": "error",
					"data":    "用户名不存在",
				})
				return
			}
		}
	}
	//检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "用户名或密码错误",
		})
		return
	}
	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, user.ID)
	err = session.Save()
	if err != nil {
		log.Println("session 保存错误")
		return
	}
	//rdb := GetRDB(c)
	//offlineKey := g.OFFLINE_USER + strconv.Itoa(user.ID)
	//rdb.Del(rctx, offlineKey).Result()
	token, err := jwt.GenToken(user.ID, user.Username, user.Email)
	c.JSON(200, gin.H{
		"code":     200,
		"message":  "success",
		"token":    token,
		"username": user.Username,
	})
}
