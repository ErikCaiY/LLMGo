package main

import (
	"context"
	"github.com/tmc/langchaingo/llms/ollama"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnswerRequest struct {
	Question string `json:"question" bing:"required"`
	Model    string `json:"model"`
}

type AnswerResponse struct {
	Answer string `json:"answer"`
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("./templates/*")
	r.GET("/chat", getHTML)
	r.POST("/api/answer", getAnswerFromModel)

	r.Run(":8080")
}

func getHTML(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}

func getAnswerFromModel(ctx *gin.Context) {
	var request AnswerRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	llm, err := ollama.New(ollama.WithModel(request.Model))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := llm.Call(context.Background(), request.Question)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Answer": response})
}
