package infocard

import (
	"encoding/xml"
	"strings"
)

type RDL struct {
	XMLName xml.Name `xml:"RDL"`
	TEXT    []string `xml:"TEXT"`
}

func (i *Infocard) XmlToText() ([]string, error) {
	return XmlToText(i.Content)
}

func XmlToText(xml_stuff string) ([]string, error) {
	var structy RDL

	prepared := strings.ReplaceAll(string(xml_stuff), `<?xml version="1.0" encoding="UTF-16"?>`, "")

	err := xml.Unmarshal([]byte(prepared), &structy)

	if _, ok := err.(*xml.SyntaxError); ok {
		replaced := strings.ReplaceAll(prepared, "&", "")
		err = xml.Unmarshal([]byte(replaced), &structy)
	}

	// logus.Log.CheckError(err, "unable converting xml to text", typelog.String("xml_stuff", string(xml_stuff)))
	if err != nil {
		return nil, err
	}
	return structy.TEXT, nil
}
