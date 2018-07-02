package service

import (
	"crypto/rsa"
	"strings"
	"github.com/tuyenlqvnp/sign-service-api/bean"
	"encoding/base64"
)

type SignatureService struct {
}

func (self SignatureService) EncryptDataWithCertificate(data *string, certificateData []byte, password string) (*bean.CipherData, error) {
	cipherData := bean.CipherData{}
	private, certificate, err := pkcsUtils.ExtractData(certificateData, password);
	if (err == nil) {
		// check certificate
		err := pkcsUtils.VerifyCertificate(certificate);
		if (err == nil) {
			shaType := strings.Split(certificate.SignatureAlgorithm.String(), "-")[0]
			hashData := shaUtils.Hash(data, &shaType)
			hashDataText := base64.URLEncoding.EncodeToString(hashData)
			cipherText, err := pkcsUtils.SignData(private.(*rsa.PrivateKey), hashData, &shaType);
			if (err == nil) {
				//log.Println("Ciphertext: " + *cipherText);
				cipherData.PrivateKey = private
				cipherData.Certificate = certificate
				cipherData.CipherText = cipherText
				cipherData.HashData = &hashDataText
				return &cipherData, nil
			}
		}
	}
	return nil, err;
}

func (self SignatureService) InsertSignatureToXmlData(xmlDataString *string, cipherData *bean.CipherData) (string, error) {
	xmlData := make(map[string]interface{});
	xmlData, err := xmlUtils.ParseFromStringToInterface(xmlDataString);
	if (err == nil) {
		invoice := xmlData["Invoice"].(map[string]interface{})
		signature := make(map[string]interface{})
		signature["SignatureValue"] = *cipherData.CipherText

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
		xmlDataStringResult, err := xmlUtils.ParseFromInterfaceToString(invoice, "", "	", "Invoice", "")
		temp := strings.Split(*xmlDataString, "<Invoice>")[0]
		if (err == nil) {
			return temp + "\n" + xmlDataStringResult, err
		}
	}
	return "", err
}
