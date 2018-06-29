package api

import (
	"github.com/gin-gonic/gin"
)

type SignatureApi struct {
}

func (self SignatureApi) Init(router *gin.Engine) *gin.RouterGroup {
	signatureApiGroup := router.Group("/signature-api")
	{
		signatureApiGroup.GET("/certificate", func(context *gin.Context) {
			self.getCertificateInfo(context);
		})
	}
	return signatureApiGroup;
}

func (self SignatureApi) getCertificateInfo(context *gin.Context) {
	//pkcsUtils := utils.PKCSUtils{};
	//pkcsUtils.
}
