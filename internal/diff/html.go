package diff

import (
	"bytes"
	"io"

	htmldiff "github.com/dominikus1993/html-diff"
)

func toString(ioClose io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(ioClose)
	return buf.String() // Does a complete copy of the bytes in the buffer.
}

type BlobDiffer struct {
	cfg *htmldiff.Config
}

func NewBlobDiffer() *BlobDiffer {
	cfg := &htmldiff.Config{
		Granularity:  6,
		InsertedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: palegreen; text-decoration: underline;"}},
		DeletedSpan:  []htmldiff.Attribute{{Key: "style", Val: "background-color: lightpink; text-decoration: line-through;"}},
		ReplacedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: lightskyblue; text-decoration: overline;"}},
		CleanTags:    []string{"documize"},
	}
	return &BlobDiffer{cfg: cfg}
}

func (diff *BlobDiffer) DiffBlobs(old, new io.ReadCloser) (*string, error) {
	defer old.Close()
	defer new.Close()
	oldHtml := toString(old)
	newHtml := toString(new)
	res, err := diff.cfg.HTMLdiff([]string{oldHtml, newHtml})
	if err != nil {
		return nil, err
	}

	return &res[0], nil
}
