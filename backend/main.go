package main

import (
	"context"
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

type Conn struct {
	sessionId int
	context   openai.Client
}

//go:embed index.html
var fs embed.FS
var conns map[int]Conn = map[int]Conn{}

func main() {
	r := gin.Default()
	apiRouter := r.Group("/api")
	apiRouter.POST("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	apiRouter.POST("/chat", func(c *gin.Context) {
		id := c.GetInt("id")
		role := c.GetString("role")
		prompt := c.GetString("prompt")
		if _, ok := conns[id]; !ok {
			conns[id] = Conn{}
			conns[id].initConn("")
		}
		response, err := conns[id].send(role, prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"response": "error when handling conversation"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": response})
	})
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(fs))))
	panic(r.Run(":8080"))
}

func (c Conn) initConn(token string) error {
	c.context = *openai.NewClient(token)
	return nil
}

func (c Conn) send(role string, prompt string) (string, error) {
	ctx := context.Background()
	req := openai.CompletionRequest{
		Model:            openai.GPT3TextDavinci001,
		Temperature:      0.4,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		BestOf:           1,
		Prompt:           prompt,
	}
	resp, err := c.context.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}
