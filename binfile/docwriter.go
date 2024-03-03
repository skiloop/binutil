package binfile

import (
	"io"
)

type DocWriter interface {
	Write(doc *Doc) (int, error)
}

type docWriter struct {
	w io.Writer
}

func (w *docWriter) Write(doc *Doc) (int, error) {
	if nil == doc {
		return 0, nil
	}
	return doc.writeDoc(w.w)
}

func NewDocWriter(w io.Writer) DocWriter {
	return &docWriter{w: w}
}
