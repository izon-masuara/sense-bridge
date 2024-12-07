package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server_lidm/apis"
	"server_lidm/db"
	"server_lidm/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadAudio(c *gin.Context) {
	code := c.Param("code")
	sec := c.Param("sec")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	defer file.Close()

	filePath := filepath.Join("tmp", fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename))
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}
	defer out.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read the file"})
		return
	}

	_, err = out.ReadFrom(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	speechToText, err := apis.SpeechToText(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resultDictionary := helpers.Dictionary(speechToText)

	os.Remove(out.Name())

	collection := db.GetCollection(db.Client, "sensebridge", code)
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"id": code}
	update := bson.M{
		"$push": bson.M{
			"sections." + sec + ".audio_text": speechToText,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     speechToText,
		"video_queue": resultDictionary,
	})
}
