package response

type DigitallySignedData struct {
	SignedData string `json:"signed_data"`
}

func (self DigitallySignedData) Init(signedData string) (*DigitallySignedData) {
	self.SignedData = signedData
	return &self
}
