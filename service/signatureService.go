package service

import (
	"crypto/rsa"
	"log"
	"strings"
	"crypto/x509"
)

type SignatureService struct {
}

func (self SignatureService) EncryptDataWithCertificate(data *string) (string, *x509.Certificate, error) {
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
				return *cipherText, certificate, nil
			}
		}
	}
	return "", nil, err;
}

func (self SignatureService) InsertSignatureToXmlData(xmlDataString *string, signatureStr *string, certificate *x509.Certificate) (string, error) {
	xmlData := make(map[string]interface{});
	xmlData, err := xmlUtils.ParseFromStringToInterface(xmlDataString);
	if (err == nil) {
		signature := make(map[string]interface{})
		signature["SignatureValue"] = *signatureStr

		x509Data := make(map[string]interface{})
		x509Data["X509Certificate"] = "adafasfsfasf"
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
		signedInfo["DigestValue"] = "dCuHcDtDhzByOLEhCqgwTkhELBE"

		xmlData["Signature"] = signature
		xmlDataStringResult, err := xmlUtils.ParseFromInterfaceToString(xmlData)
		if (err == nil) {
			return xmlDataStringResult, err
		}
	}
	return "", err
}
