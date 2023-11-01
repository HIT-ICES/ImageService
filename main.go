package main

import (
	"containerd/route"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	route.RouterHandler(router)
}
