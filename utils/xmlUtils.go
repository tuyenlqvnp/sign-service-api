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

func (self XmlUtils) ParseFromInterfaceToString(data interface{}, prefix string, intent string, tags ...string) (string, error) {
	dataStr, err := mxj.AnyXmlIndent(data, prefix, intent, tags[0], tags[1])
	if (err == nil) {
		return string(dataStr), err
	}
	return "", err
}
