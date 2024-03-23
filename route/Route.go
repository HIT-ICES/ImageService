package route

import (
	"containerd/imageOperation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func pullImage(context *gin.Context) {
	imageName := context.Query("imageName")
	response := imageOperation.PullImage(imageName)
	context.String(200, response)
}

func listImages(context *gin.Context) {
	imageArray := imageOperation.ListImages()
	context.JSON(200, imageArray)
}

func deleteImage(context *gin.Context) {
	var images []imageOperation.DeleteImages
	if err := context.BindJSON(&images); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := imageOperation.DeleteImage(&images)
	context.JSON(http.StatusOK, *response)
}

func RouterHandler(router *gin.Engine) {
	router.GET("/listImages", listImages)
	router.GET("/pullImage", pullImage)
	router.POST("/deleteImage", deleteImage)

	router.Run(":8080")
}
