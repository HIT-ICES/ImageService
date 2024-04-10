package route

import (
	"containerd/imageOperation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func pullImages(context *gin.Context) {
	var images imageOperation.PullImagesJSON
	if err := context.BindJSON(&images); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := imageOperation.PullImages(&images)
	context.JSON(http.StatusOK, response)
}

func listImages(context *gin.Context) {
	imageArray := imageOperation.ListImages()
	context.JSON(http.StatusOK, imageArray)
}

func deleteImages(context *gin.Context) {
	var images imageOperation.DeleteImagesJSON
	if err := context.BindJSON(&images); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := imageOperation.DeleteImages(&images)
	context.JSON(http.StatusOK, *response)
}

func RouterHandler(router *gin.Engine) {
	router.GET("/listImages", listImages)
	router.POST("/pullImages", pullImages)
	router.POST("/deleteImages", deleteImages)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
