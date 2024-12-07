package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"server_lidm/apis"
	"server_lidm/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadImage(c *gin.Context) {
	code := c.Param("code")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	defer file.Close()

	uploadPath := "./tmp"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		err := os.Mkdir(uploadPath, os.ModePerm)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("create upload directory err: %s", err.Error()))
			return
		}
	}

	filePath := filepath.Join(uploadPath, fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename))
	out, err := os.Create(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("create file err: %s", err.Error()))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("save file err: %s", err.Error()))
		return
	}

	res := apis.ImageAnlyze(filePath)

	collection := db.GetCollection(db.Client, "sensebridge", code)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"id": code}

	var result bson.M
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("find document err: %s", err.Error()))
		return
	}

	sections, ok := result["sections"].(primitive.M)
	if !ok {
		c.String(http.StatusInternalServerError, "sections tidak ditemukan atau formatnya salah")
		return
	}

	var lastSectionKey string
	var maxSectionID int

	for key := range sections {
		sectionID, err := strconv.Atoi(key)
		if err != nil {
			c.String(http.StatusInternalServerError, "format section key tidak valid")
			return
		}
		if sectionID > maxSectionID {
			maxSectionID = sectionID
			lastSectionKey = key
		}
	}

	newSectionID := strconv.Itoa(maxSectionID + 1)
	newSection := map[string]interface{}{
		"audio_text": []string{},
		"image_text": "",
	}

	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("sections.%s.image_text", lastSectionKey): res,
			fmt.Sprintf("sections.%s", newSectionID):              newSection,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("update document err: %s", err.Error()))
		return
	}

	os.Remove(filePath)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
