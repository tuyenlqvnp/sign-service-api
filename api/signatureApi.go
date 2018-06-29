package api

import (
	"github.com/gin-gonic/gin"
	"log"
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
	xmlDataString := context.PostForm("xmlData")
	cipherData, err := signatureService.EncryptDataWithCertificate(&xmlDataString);
	if (err == nil) {
		result, err := signatureService.InsertSignatureToXmlData(&xmlDataString, cipherData)
		if (err == nil) {
			log.Println(result);
		}
	}
}
