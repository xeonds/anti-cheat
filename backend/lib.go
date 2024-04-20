package main

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserClaim struct {
	Permission int       `json:"permission"`
	Name       string    `json:"name"`
	Expire     time.Time `json:"expire"`
}

// handlers for gorm
func create[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var d T
		if err := c.ShouldBindJSON(&d); err != nil {
			c.AbortWithStatus(404)
			log.Println("[gorm]parse creation data failed: ", err)
		} else {
			if err := db.Create(&d).Error; err != nil {
				c.AbortWithStatus(404)
				log.Println("[gorm]create data failed: ", err)
			} else {
				c.JSON(200, d)
			}
		}
	}
}
func get[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		params := c.Request.URL.Query()
		var d T
		if err := db.Where(params).Find(&d).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, d)
		}
	}
}
func getAll[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var d []T
		if err := db.Find(&d).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, d)
		}
	}
}
func update[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var d T
		if err := db.Save(&d).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, d)
		}
	}
}
func delete[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var d T
		params := c.Request.URL.Query()
		if err := db.Where(params).Delete(&d).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, d)
		}
	}
}
func handleFind[T any](mode string, db *gorm.DB) func(c *gin.Context) {
	if mode == "single" {
		return func(c *gin.Context) {
			var query T
			c.ShouldBindJSON(&query)
			err := db.Where(&query).Find(&query).Error
			if err != nil {
				c.JSON(200, gin.H{
					"message": "很抱歉，没有找到你要查询的内容哦",
				})
			} else {
				c.JSON(200, query)
			}
		}
	} else {
		return func(c *gin.Context) {
			var query T
			var results []T
			c.ShouldBindJSON(&query)
			err := db.Where(&query).Find(&results).Error
			if err != nil {
				c.JSON(200, gin.H{
					"message": "很抱歉，没有找到你要查询的内容哦",
				})
			} else {
				c.JSON(200, results)
			}
		}
	}
}

// 为Gin router 添加CRUD
func APIBuilder(router gin.IRouter, handlers ...func(*gin.RouterGroup) *gin.RouterGroup) func(gin.IRouter, string) *gin.RouterGroup {
	return func(router gin.IRouter, path string) *gin.RouterGroup {
		group := router.Group(path)
		for _, handler := range handlers {
			group = handler(group)
		}
		return group
	}
}
func AddCRUD[T any](router gin.IRouter, path string, db *gorm.DB) *gin.RouterGroup {
	return APIBuilder(router, func(group *gin.RouterGroup) *gin.RouterGroup {
		group.GET("", getAll[T](db))
		group.GET("/:id", get[T](db))
		group.POST("", create[T](db))
		group.PUT("/:id", update[T](db))
		group.DELETE("/:id", delete[T](db))
		return group
	})(router, path)
}
func AddStatic(router *gin.Engine, staticFileDir []string) {
	for _, dir := range staticFileDir {
		router.NoRoute(gin.WrapH(http.FileServer(http.Dir(dir))))
	}
}
func AddStaticFS(router *gin.Engine, fs embed.FS) {
	router.NoRoute(gin.WrapH(http.FileServer(http.FS(fs))))
}
func AddFindAPI[T any](router gin.IRouter, path string, mode string, db *gorm.DB) *gin.RouterGroup {
	return APIBuilder(router, func(group *gin.RouterGroup) *gin.RouterGroup {
		group.POST("", handleFind[T](mode, db))
		return group
	})(router, path)
}

// Gin中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		log.Print(latency)
	}
}

func JWTMiddleware(setUserToken func(c *gin.Context, userClaim UserClaim)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
			return
		}
		token := c.GetHeader("Authorization")
		parsed, err := ParseToken(token)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		if parsed.Valid() != nil {
			c.AbortWithStatus(401)
			return
		}
		if setUserToken != nil {
			// 设置用户信息字段
			setUserToken(c, *parsed)
		}
		c.Next()
	}
}

// 验证用户权限等级
// 权限位于[permLo, permHi]之间则为合理
func AuthMiddleware(permLo, permHi int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
		}
		token := c.GetHeader("Authorization")
		parsed, err := ParseToken(token)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		if parsed.Valid() != nil {
			c.AbortWithStatus(401)
			return
		}
		// 当权限变更时，使原来的token失效即可
		if parsed.Permission < permLo || parsed.Permission > permHi {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}

// 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// finish this golang cors middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Next()
	}
}

// 工具函数
// 生成token
func GenerateToken(userClaim *UserClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 验证token
func (uc *UserClaim) Valid() error {
	// 过期时间为一周
	if time.Now().Unix() > uc.Expire.Unix() {
		return errors.New("token expired")
	}
	return nil
}

// 解析token
func ParseToken(token string) (*UserClaim, error) {
	userClaim := new(UserClaim)
	claim, err := jwt.ParseWithClaims(token, userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !claim.Valid {
		return nil, errors.New("token validation failed")
	}
	return userClaim, nil
}

// 配置管理
// 有错误则返回nil
func LoadConfig() (*Config, error) {
	if _, err := os.Stat("config.yaml"); err != nil {
		os.WriteFile("config.yaml", []byte(""), 0644)
		return nil, errors.New("config file not found")
	}
	if err := func() error {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		return viper.ReadInConfig()
	}(); err != nil {
		return nil, errors.New("config file read failed")
	}
	if err := viper.Unmarshal(&Config{}); err != nil {
		return nil, errors.New("config file parse failed")
	}

	return viper.Get("config").(*Config), nil
}

// 验证码服务，使用redis存储
// 生成验证码
// length 验证码长度，id 验证码在redis的key，db redis客户端指针
func GenerateCaptcha(ctx context.Context, length int, id string, db *redis.Client) string {
	var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	code := func(chars string) string {
		return fmt.Sprintf("%s%s", "CAPT", func() string {
			rand.Seed(time.Now().UnixNano())
			result := make([]byte, length)
			for i := range result {
				result[i] = chars[rand.Intn(len(chars))]
			}
			return string(result)
		}())
	}(chars)
	db.Set(ctx, id, code, 5*time.Minute)
	return code
}

// 验证验证码
func VerifyCaptcha(ctx context.Context, id string, captcha string, db *redis.Client) bool {
	value := db.Get(ctx, id)
	return captcha == value.String()
}

func MailSendCode(mail, code string, config Config) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Get <%s>", config.MailAdmin)
	e.To = []string{mail}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("你的验证码为：<h1>" + code + "</h1>")
	return e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", config.MailAdmin, config.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}
