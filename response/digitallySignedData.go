package response

import "log"

type DigitallySignedData struct {
	SignedData string `json:"signed_data"`
}

func (self DigitallySignedData) Init(signedData string) (*DigitallySignedData) {
	log.Println(signedData)
	self.SignedData = signedData
	return &self
}
