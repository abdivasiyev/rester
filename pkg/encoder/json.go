package encoder

import (
	"encoding/json"
	"io"
)

var JsonEncoder Encoder = &jsonEncoder{}

type jsonEncoder struct {
	encoder *json.Encoder
}

func (d *jsonEncoder) New(w io.Writer) Encoder {
	return &jsonEncoder{
		encoder: json.NewEncoder(w),
	}
}

func (d *jsonEncoder) Encode(src any) error {
	return d.encoder.Encode(src)
}
