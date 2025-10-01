package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"umbandario-go/database"

	"github.com/gin-gonic/gin"
)

func CreateLineHandler(c *gin.Context) {
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome da linha é obrigatório."})
		return
	}

	lineName := strings.TrimSpace(requestBody.Name)
	if lineName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome da linha está vazio."})
		return
	}

	lineName = strings.ToLower(lineName)

	_, err := database.GetLineByName(lineName)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Uma linha com este nome já existe."})
		return
	}

	newLine, err := database.CreateLine(lineName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, newLine)
}

func ListLinesHandler(c *gin.Context) {
	lines, err := database.GetAllLines()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao buscar as linhas."})
		return
	}
	c.JSON(http.StatusOK, lines)
}
