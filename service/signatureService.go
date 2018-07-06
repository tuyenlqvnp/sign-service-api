package service

import (
	"crypto/rsa"
	"strings"
	"github.com/tuyenlqvnp/sign-service-api/bean"
	"encoding/base64"
	"github.com/tdewolff/minify"
	"regexp"
	"github.com/tdewolff/minify/xml"
)

type SignatureService struct {
}

func (self SignatureService) XmlMinify(data *string) *string {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	s, _ := m.String("", *data)
	return &s
}

func (self SignatureService) SignDataWithCertificate(data *string, certificateData []byte, password string) (*bean.CipherData, error) {
	xmlData, err := xmlUtils.ParseFromStringToInterface(data);
	if (err == nil) {
		invoice := xmlData["Invoice"].(map[string]interface{})
		xmlDataStringResult, err := xmlUtils.ParseFromInterfaceToString(invoice, "", "", "Invoice", "")
		xmlDataStringResult = strings.Replace(xmlDataStringResult, "\n", "", -1)
		minifiedXml := self.XmlMinify(&xmlDataStringResult)
		temp := strings.Split(*data, "<Invoice>")[0]
		if (err == nil) {
			temp = temp + *minifiedXml
			data = self.XmlMinify(&temp)
			cipherData := bean.CipherData{}
			private, certificate, err := pkcsUtils.ExtractData(certificateData, password);
			if (err == nil) {
				// check certificate
				err := pkcsUtils.VerifyCertificate(certificate);
				if (err == nil) {
					shaType := strings.Split(certificate.SignatureAlgorithm.String(), "-")[0]
					hashData := shaUtils.Hash(data, &shaType)
					signature, err := pkcsUtils.SignData(private.(*rsa.PrivateKey), hashData, &shaType);
					if (err == nil) {
						cipherData.PrivateKey = private
						cipherData.Certificate = certificate
						cipherData.Signature = signature
						hashDataText := base64.URLEncoding.EncodeToString(hashData)
						cipherData.HashData = &hashDataText
						return &cipherData, nil
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self SignatureService) ValidateSignature(signature map[string]interface{}, data *string, certificateData []byte, password string) (error) {
	private, certificate, err := pkcsUtils.ExtractData(certificateData, password);
	if (err == nil) {
		err := pkcsUtils.VerifyCertificate(certificate);
		if (err == nil) {
			shaType := strings.Split(certificate.SignatureAlgorithm.String(), "-")[0]
			hashData := shaUtils.Hash(data, &shaType)
			err := pkcsUtils.ValidateSignedData(&private.(*rsa.PrivateKey).PublicKey, signature["SignatureValue"].(string), hashData, &shaType)
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (self SignatureService) RemoveSignatureFromXmlData(xmlDataString *string) (string, map[string]interface{}, error) {
	path := "Invoice.Signature"
	xmlData, signatureInfo, err := xmlUtils.RemoveElement(xmlDataString, path);
	if (err == nil) {
		invoice := xmlData["Invoice"].(map[string]interface{})
		xmlDataStringResult, err := xmlUtils.ParseFromInterfaceToString(invoice, "", "", "Invoice", "")
		xmlDataStringResult = strings.Replace(xmlDataStringResult, "\n", "", -1)
		minifiedXml := self.XmlMinify(&xmlDataStringResult)
		temp := strings.Split(*xmlDataString, "<Invoice>")[0]
		if (err == nil) {
			return temp + *minifiedXml, signatureInfo, err
		} else {
			return "", signatureInfo, err
		}
	} else
	{
		return "", signatureInfo, err
	}
}

func (self SignatureService) InsertSignatureToXmlData(xmlDataString *string, cipherData *bean.CipherData) (string, error) {
	xmlData, err := xmlUtils.ParseFromStringToInterface(xmlDataString);
	if (err == nil) {
		invoice := xmlData["Invoice"].(map[string]interface{})
		signature := make(map[string]interface{})
		signature["SignatureValue"] = *cipherData.Signature

		x509Data := make(map[string]interface{})
		x509Data["X509Certificate"] = base64.URLEncoding.EncodeToString(cipherData.Certificate.Raw)
		keyInfo := make(map[string]interface{})
		keyInfo["X509Data"] = x509Data
		signature["KeyInfo"] = keyInfo

		signedInfo := make(map[string]interface{})
		canonicalizationMethod := make(map[string]interface{})
		canonicalizationMethod["-Algorithm"] = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
		signedInfo["CanonicalizationMethod"] = canonicalizationMethod
		signatureMethod := make(map[string]interface{})
		signatureMethod["-Algorithm"] = "http://www.w3.org/2000/09/xmldsig#rsa-sha256"
		signedInfo["SignatureMethod"] = signatureMethod
		Reference := make(map[string]interface{})
		Reference["-URI"] = "SigningData"
		Transform := make(map[string]interface{})
		Transform["-Algorithm"] = "http://www.w3.org/2000/09/xmldsig#enveloped-signature"
		Transforms := make(map[string]interface{})
		Transforms["Transform"] = Transform
		Reference["Transforms"] = Transforms
		signedInfo["Reference"] = Reference
		signature["SignedInfo"] = signedInfo
		DigestMethod := make(map[string]interface{})
		DigestMethod["-Algorithm"] = "http://www.w3.org/2000/09/xmldsig#sha256"
		signedInfo["DigestMethod"] = DigestMethod
		signedInfo["DigestValue"] = *cipherData.HashData

		invoice["Signature"] = signature
		xmlDataStringResult, err := xmlUtils.ParseFromInterfaceToString(invoice, "", "", "Invoice", "")
		xmlDataStringResult = strings.Replace(xmlDataStringResult, "\n", "", -1)
		minifiedXml := self.XmlMinify(&xmlDataStringResult)
		temp := strings.Split(*xmlDataString, "<Invoice>")[0]
		if (err == nil) {
			return temp + *minifiedXml, err
		}
	}
	return "", err
}
