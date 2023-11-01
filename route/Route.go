package route

import (
	"containerd/imageOperation"
	"github.com/gin-gonic/gin"
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
	imageName := context.Query("imageName")
	response := imageOperation.DeleteImage(imageName)
	context.String(200, response)
}

func RouterHandler(router *gin.Engine) {
	router.GET("/listImages", listImages)
	router.GET("/pullImage", pullImage)
	router.GET("/deleteImage", deleteImage)

	router.Run(":8080")
}
