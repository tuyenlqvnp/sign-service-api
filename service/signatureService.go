package service

import (
	"crypto/rsa"
	"log"
	"strings"
)

type SignatureService struct {
}

func (self SignatureService) EncryptDataWithCertificate(data *string) (string, error) {
	private, certificate, err := pkcsUtils.ExtractDataFromFile("/Users/thaibao/Desktop/my256.p12", "654321");
	if (err == nil) {
		// check certificate
		err := pkcsUtils.VerifyCertificate(certificate);
		if (err == nil) {
			shaType := strings.Split(certificate.SignatureAlgorithm.String(), "-")[0]
			hashData := shaUtils.Hash(data, &shaType)
			cipherText, err := pkcsUtils.EncryptData(private.(*rsa.PrivateKey), []byte(hashData));
			if (err == nil) {
				log.Println("Ciphertext: " + *cipherText);
				return *cipherText, nil
			}
		}
	}
	return "", err;
}

func (self SignatureService) InsertSignatureToXmlData(xmlDataString *string, signature *string) (string, error) {
	xmlData := make(map[string]interface{});
	xmlData, err := xmlUtils.ParseFromStringToInterface(xmlDataString);
	if (err == nil) {
		//invoice := xmlData["Invoice"].(map[string]interface{})
		signatureValue := make(map[string]interface{})
		signatureValue["SignatureValue"] = *signature
		xmlData["Signature"] = signatureValue
		xmlDataStringResult, err := xmlUtils.ParseFromInterfaceToString(xmlData)
		if (err == nil) {
			return xmlDataStringResult, err
		}
	}
	return "", err
}
