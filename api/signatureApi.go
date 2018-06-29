package api

import (
	"github.com/gin-gonic/gin"
)

type SignatureApi struct {
}

func (self SignatureApi) Init(router *gin.Engine) *gin.RouterGroup {
	signatureApiGroup := router.Group("/signature-api")
	{
		signatureApiGroup.POST("/certificate", func(context *gin.Context) {
			self.SignWithCertificate(context);
		})
	}
	return signatureApiGroup;
}

func (self SignatureApi) SignWithCertificate(context *gin.Context) {
	signatureService.EncryptDataWithCertificate("xml data in here");
}
