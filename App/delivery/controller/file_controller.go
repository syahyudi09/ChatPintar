package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syahyudi09/ChatPintar/App/usecase"
)

type FileController struct {
	router  *gin.Engine
	usecase usecase.FileUsecase
}

func NewFileController(r *gin.Engine, usecase usecase.FileUsecase) *FileController {
	controller := FileController{
		router:  r,
		usecase: usecase,
	}

	r.POST("/upload-file", controller.SharingController)

	return &controller
}

func (fc *FileController) SharingController(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	uploader := c.PostForm("uploader") // Nama pengguna yang mengunggah
	fileName := header.Filename        // Nama file yang diunggah

	// Gunakan use case untuk mengunggah file
	err = fc.usecase.UploadFile(file, fileName, uploader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
