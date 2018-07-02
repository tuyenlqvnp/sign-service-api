package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/tuyenlqvnp/sign-service-api/response"
	"net/http"
	"github.com/tuyenlqvnp/sign-service-api/bean"
	"bytes"
	"io"
)

type SignatureApi struct {
}

func (self SignatureApi) Init(router *gin.Engine) *gin.RouterGroup {
	signatureApiGroup := router.Group("/signature-api")
	{
		signatureApiGroup.POST("/sign", func(context *gin.Context) {
			self.SignWithCertificate(context);
		})
		signatureApiGroup.POST("/validate", func(context *gin.Context) {
			self.SignWithCertificate(context);
		})
	}
	return signatureApiGroup;
}

func (self SignatureApi) SignWithCertificate(context *gin.Context) {
	result := response.Base{}
	sourceFile, _, err := context.Request.FormFile("certificate")
	xmlDataString := context.PostForm("xmlData")
	password := context.PostForm("password")

	if (xmlDataString == "") {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	if (password == "") {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	defer sourceFile.Close()
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
	} else {
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, sourceFile); err != nil {
			result.SetStatus(bean.UnexpectedError)
		}
		cipherData, err := signatureService.EncryptDataWithCertificate(&xmlDataString, buf.Bytes(), password);
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
	}
	context.JSON(http.StatusOK, result)
	return
}

func (self SignatureApi) Validate(context *gin.Context) {
	result := response.Base{}
	context.JSON(http.StatusOK, result)
	return
}
