package main

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed index.html
var fs embed.FS

type LLM struct {
	SessionId int
	Client    *openai.Client
	Context   context.Context
}

// 环境变量
var JwtSecret = os.Getenv("JWT_SECRET")
var MailPassword = os.Getenv("MAIL_PASSWORD")
var MailAdmin = os.Getenv("MAIL_ADMIN")

// 应用配置数据
type Config struct {
	MailAdmin    string
	MailPassword string
	JwtSecret    string
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type ChatRecord struct {
	gorm.Model
}

type ScreenShots struct {
	gorm.Model
}

type PhoneRecord struct {
	gorm.Model
}

func main() {
	r := gin.Default()
	conns := map[int]LLM{}
	db, _ := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	db.AutoMigrate(&ChatRecord{}, &ScreenShots{}, &PhoneRecord{})

	apiRouter := r.Group("/api/v1")
	apiRouter.POST("/check", handleCheck)
	apiRouter.POST("/chat", handleChat(conns))
	modelRouter := r.Group("/model/v1")
	AddCRUD[ChatRecord](modelRouter, "/chat_record", db)
	AddStaticFS(r, fs)
	panic(r.Run(":8080"))
}

func NewLLM(sessionId int) LLM {
	return LLM{SessionId: sessionId, Client: openai.NewClient("YOUR_API_KEY"), Context: context.Background()}
}

func (c LLM) send(role string, prompt string) (string, error) {
	req := openai.CompletionRequest{
		Model:            openai.GPT3TextDavinci001,
		Temperature:      0.4,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		BestOf:           1,
		Prompt:           "[Role: " + role + "] " + prompt,
	}
	resp, err := c.Client.CreateCompletion(c.Context, req)
	if err != nil {
		return "", errors.New("openai.CreateCompletion err: " + err.Error())
	}
	return resp.Choices[0].Text, nil
}

// handlers

func handleCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func handleChat(conns map[int]LLM) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		role := c.GetString("role")
		prompt := c.GetString("prompt")
		if _, ok := conns[id]; !ok {
			conns[id] = NewLLM(id)
		}
		response, err := conns[id].send(role, prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"response": "error when handling conversation"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": response})
	}
}
