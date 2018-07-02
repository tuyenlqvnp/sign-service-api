package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/tuyenlqvnp/sign-service-api/response"
	"net/http"
	"github.com/tuyenlqvnp/sign-service-api/bean"
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
	result := response.Base{}
	cipherData, err := signatureService.EncryptDataWithCertificate(&xmlDataString);
	if (err == nil) {
		signedData, err := signatureService.InsertSignatureToXmlData(&xmlDataString, cipherData)
		if (err == nil) {
			log.Println(result);
			digitallySignedData := response.DigitallySignedData{}.Init(signedData)
			result.Data = digitallySignedData
			result.SetStatus(bean.Success)
		} else {
			result.SetStatus(bean.UnexpectedError)
		}
	} else {
		result.SetStatus(bean.UnexpectedError)
	}
	context.JSON(http.StatusOK, result)
	return
}
