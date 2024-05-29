package encoder

import (
	"io"
)

type ContentTyper interface {
	ContentType() string
}

type Encoder interface {
	New(w io.Writer) Encoder
	Encode(src any) error
}
