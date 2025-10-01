package main

import (
	"log"
	"umbandario-go/database"
	"umbandario-go/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB("./audio.db")
	defer database.DB.Close()

	router := gin.Default()
	v1 := router.Group("/api/v1")

	audio := v1.Group("/audio")
	audio.GET("", handlers.ListAudiosHandler)
	audio.GET("/play/:filename", handlers.StreamAudioHandler)
	audio.POST("", handlers.UploadAudioHandler)
	audio.DELETE("/:audioID", handlers.DeleteAudioHandler)

	lines := v1.Group("/lines")
	lines.GET("", handlers.ListLinesHandler)
	lines.POST("", handlers.CreateLineHandler)

	if err := router.Run(); err != nil {
		log.Panic("Não foi possível iniciar o servidor")
	}
}

// func streamAudioHandler(c *gin.Context) {
// 	filename := c.Param("filename")

// 	safeFilename := filepath.Base(filename)

// 	filePath := filepath.Join("audios", safeFilename)

// 	_, err := os.Stat(filePath)
// 	if os.IsNotExist(err) {
// 		c.String(http.StatusNotFound, "Arquivo de áudio não encontrado.")
// 		return
// 	}
// 	http.ServeFile(c.Writer, c.Request, filePath)
// }
