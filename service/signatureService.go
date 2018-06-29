package service

import (
	"crypto/rsa"
	"log"
	"strings"
)

type SignatureService struct {
}

func (self SignatureService) EncryptDataWithCertificate(data string) (error) {
	private, certificate, err := pkcsUtils.ExtractDataFromFile("/Users/thaibao/Desktop/my256.p12", "654321");
	if (err == nil) {
		// check certificate
		err := pkcsUtils.VerifyCertificate(certificate);
		if (err == nil) {
			shaType := strings.Split(certificate.SignatureAlgorithm.String(), "-")[0]
			hashData := shaUtils.Hash(data, shaType)
			cipherText, err := pkcsUtils.EncryptData(private.(*rsa.PrivateKey), []byte(hashData));
			if (err == nil) {
				log.Println("Ciphertext: " + *cipherText);
			}
		}
	}
	return err;
}
