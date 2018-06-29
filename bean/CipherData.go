package bean

import "crypto/x509"

type CipherData struct {
	HashData    *string
	PrivateKey  interface{}
	Certificate *x509.Certificate
	CipherText  *string
}
