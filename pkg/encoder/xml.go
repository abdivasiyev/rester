package encoder

import (
	"encoding/xml"
	"io"
)

var XmlEncoder Encoder = &xmlEncoder{}

type xmlEncoder struct {
	encoder *xml.Encoder
}

func (e *xmlEncoder) New(w io.Writer) Encoder {
	return &xmlEncoder{
		encoder: xml.NewEncoder(w),
	}
}

func (e *xmlEncoder) Encode(src any) error {
	return e.encoder.Encode(src)
}
