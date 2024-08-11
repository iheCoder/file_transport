package request_helper

import (
	"file_tranport"
	"github.com/gin-gonic/gin"
)

type RequestHelper struct {
	ginEngine *gin.Engine
}

func NewRequestHelper(port string, server *file_tranport.FileTransportServer) *RequestHelper {
	// create gin engine
	r := gin.Default()

	// register router
	r.GET("/progress_info", server.GetProgressInfo)
	r.POST("/upload", server.UploadFile)

	// run server
	go func() {
		r.Run(":" + port)
	}()

	return &RequestHelper{
		ginEngine: r,
	}
}
