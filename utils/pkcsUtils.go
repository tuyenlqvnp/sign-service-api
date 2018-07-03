package utils

import (
	"log"
	"golang.org/x/crypto/pkcs12"
	"crypto/x509"
	"io/ioutil"
	"errors"
	"crypto/rsa"
	//"math/big"
	//"encoding/base64"
	"crypto/rand"
	"crypto"
	"encoding/base64"
	"strings"
)

type PKCSUtils struct {
}

func (self PKCSUtils) ExtractData(data []byte, password string) (interface{}, *x509.Certificate, error) {
	privateKey, certificate, err := pkcs12.Decode(data, password)
	return privateKey, certificate, err;
}

func (self PKCSUtils) ExtractDataFromFile(filePath string, password string) (interface{}, *x509.Certificate, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return self.ExtractData(data, password);
}

func (self PKCSUtils) VerifyCertificate(cert *x509.Certificate) error {
	_, err := cert.Verify(x509.VerifyOptions{})
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case x509.CertificateInvalidError:
		switch e.Reason {
		case x509.Expired:
			return errors.New("certificate has expired or is not yet valid")
		default:
			return err
		}
	case x509.UnknownAuthorityError:
		// Apple cert isn't in the cert pool
		// ignoring this error
		return nil
	default:
		return err
	}
}

func (self PKCSUtils) SignData(priv *rsa.PrivateKey, data []byte, shaType *string) (*string, error) {
	cryptoHash := crypto.SHA256
	if (strings.ToUpper(*shaType) == "SHA1") {
		cryptoHash = crypto.SHA1
	} else if (strings.ToUpper(*shaType) == "SHA512") {
		cryptoHash = crypto.SHA512
	}
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, priv, cryptoHash, data)
	if (err != nil) {
		return nil, err
	}
	signature := base64.URLEncoding.EncodeToString(signatureBytes)
	return &signature, nil
}

func (self PKCSUtils) ValidateSignedData(pub *rsa.PublicKey, signature string, data []byte, shaType *string) (error) {
	cryptoHash := crypto.SHA256
	if (strings.ToUpper(*shaType) == "SHA1") {
		cryptoHash = crypto.SHA1
	} else if (strings.ToUpper(*shaType) == "SHA512") {
		cryptoHash = crypto.SHA512
	}
	signatureBytes, err := base64.URLEncoding.DecodeString(signature)
	if (err != nil) {
		return err
	}
	err = rsa.VerifyPKCS1v15(pub, cryptoHash, data, signatureBytes)
	if (err != nil) {
		return err
	}
	return nil
}

/*func (self PKCSUtils) EncryptData(priv *rsa.PrivateKey, data []byte) (*string, error) {
	k := (priv.N.BitLen() + 7) / 8
	tLen := len(data)
	// rfc2313, section 8:
	// The length of the data D shall not be more than k-11 octets
	if tLen > k-11 {
		err := errors.New("input size too large")
		return nil, err
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], data)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(priv.N) > 0 {
		err := errors.New("encryption error")
		return nil, err
	}
	var m *big.Int
	var ir *big.Int
	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}
	enc := m.Bytes()
	cipherText := base64.URLEncoding.EncodeToString(enc)
	return &cipherText, nil
}*/
