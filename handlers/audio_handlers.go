package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"umbandario-go/database"
	"umbandario-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func ListAudiosHandler(c *gin.Context) {
	audioFiles, err := database.GetAllAudioFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar o banco de dados"})
		return
	}
	c.JSON(http.StatusOK, audioFiles)
}

func UploadAudioHandler(c *gin.Context) {
	lineID := c.PostForm("line_id")
	if lineID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID da linha (line_id) é obrigatório."})
		return
	}

	_, err := database.GetLineByID(lineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Linha não encontrada."})
		return
	}

	fileName := c.PostForm("file_name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O nome do arquivo é obrigatório."})
		return
	}

	file, err := c.FormFile("audio_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum arquivo enviado."})
		return
	}

	originalFilename := filepath.Base(file.Filename)
	ext := filepath.Ext(originalFilename)
	baseName := strings.TrimSuffix(originalFilename, ext)
	safeBaseName := slug.Make(baseName)
	safeFilename := safeBaseName + strings.ToLower(ext)

	savePath := filepath.Join("./audios", safeFilename)

	if _, err := os.Stat(savePath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Arquivo com este nome já existe."})
		return
	}

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível salvar o arquivo."})
		return
	}

	originalFilename = strings.TrimSpace(originalFilename)
	originalFilename = strings.ToLower(originalFilename)

	newAudioFile := models.AudioFile{
		ID:       uuid.NewString(),
		Filename: strings.TrimSuffix(originalFilename, ext),
		Filetype: strings.TrimPrefix(strings.ToLower(ext), "."),
		Path:     savePath,
		LineID:   lineID,
	}

	var audioFileCreated models.AudioFile

	if audioFileCreated, err = database.CreateAudioFile(newAudioFile); err != nil {
		os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar no banco de dados."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ponto cantado enviado com sucesso!",
		"audio":   audioFileCreated,
	})
}

func DeleteAudioHandler(c *gin.Context) {
	audioID := c.Param("audioID")

	audio, err := database.GetAudioFileByID(audioID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar o arquivo."})
		return
	}

	_, err = database.DeleteAudioFile(audioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar do banco de dados."})
		return
	}

	if err := os.Remove(audio.Path); err != nil {
		log.Printf("AVISO: Falha ao remover arquivo físico %s: %v", audio.Path, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Arquivo deletado com sucesso!", "id": audioID})
}

func StreamAudioHandler(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("./audios", filepath.Base(filename))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "Arquivo de áudio não encontrado.")
		return
	}
	http.ServeFile(c.Writer, c.Request, filePath)
}
