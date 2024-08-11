package file_tranport

import "github.com/gin-gonic/gin"

type GetProgressInfoRequest struct {
	Info FileInfo `json:"info" form:"info"`
}

func (s *FileTransportServer) GetProgressInfo(ctx *gin.Context) {
	var req GetProgressInfoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"bind error": err.Error()})
		return
	}

	// get data handler
	dh, err := s.adviser.GetServerDataHandler(&req.Info)
	if err != nil {
		ctx.JSON(500, gin.H{"get data handler error": err.Error()})
		return
	}

	// get progress info
	info := dh.GetProgressBar()
	ctx.JSON(200, gin.H{"info": info})
}

func (s *FileTransportServer) UploadFile(ctx *gin.Context) {

}
