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

func (self XmlUtils) RemoveElements(xmlData *string, paths []string) (map[string]interface{}, []map[string]interface{}, error) {
	byteData := []byte(*xmlData)
	data, err := mxj.NewMapXml(byteData)
	if (err != nil) {
		return data, nil, err
	}
	elements := make([]map[string]interface{}, len(paths))
	for _, path := range paths {
		e, _ := data.ValueForPath(path)
		elements = append(elements, e.(map[string]interface{}))
		err := data.Remove(path)
		if err != nil {
			return data, nil, err
		}
	}
	return data, elements, nil
}

func (self XmlUtils) RemoveElement(xmlData *string, path string) (map[string]interface{}, map[string]interface{}, error) {
	byteData := []byte(*xmlData)
	//byteData, err := mxj.BeautifyXml(byteData, "", "")
	//if (err != nil) {
	//	return nil, nil, err
	//}
	data, err := mxj.NewMapXml(byteData)
	if (err != nil) {
		return data, nil, err
	}
	e, err := data.ValueForPath(path)
	if (err != nil) {
		return data, nil, err
	}
	err = data.Remove(path)
	if (err != nil) {
		return data, nil, err
	}
	return data, e.(map[string]interface{}), nil
}
