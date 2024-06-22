package infocard

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

func (i *Infocard) XmlToText() ([]string, error) {
	return XmlToText(i.content)
}

func (i *Infocard) GetContent() string {
	return i.content
}

func XmlToText(raw string) ([]string, error) {
	prepared := strings.ReplaceAll(string(raw), `<?xml version="1.0" encoding="UTF-16"?>`, "")
	decoder := xml.NewDecoder(bytes.NewBufferString(prepared))

	lines := make([]string, 0, 10)
	var sb strings.Builder
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		switch tok := token.(type) {
		case xml.EndElement:
			if tok.Name.Local == "PARA" || tok.Name.Local == "POP" {
				lines = append(lines, sb.String())
				sb.Reset()
			}
		case xml.CharData:
			sb.WriteString(string(tok))
		default:
			continue
		}
	}

	return lines, nil
}
