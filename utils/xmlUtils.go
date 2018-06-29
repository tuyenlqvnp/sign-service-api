package utils

import (
	"github.com/clbanning/mxj"
)

type XmlUtils struct {
}

func (self XmlUtils) ParseFromStringToInterface(xmlData *string) (map[string]interface{}, error) {
	byteData := []byte(*xmlData)
	data, err := mxj.NewMapXml(byteData, )
	return data, err
}

func (self XmlUtils) ParseFromInterfaceToString(data interface{}) (string, error) {
	dataStr, err := mxj.AnyXmlIndent(data, "", "  ")
	if (err == nil) {
		return string(dataStr), err
	}
	return "", err
}
