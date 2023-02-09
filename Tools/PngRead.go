package Tools

import (
	"bytes"
	"github.com/GenesisAN/illusionsCard/Base"
)

func get_png(file []byte) int {
	res1 := bytes.Index(file, Base.PngEndChunk)
	return res1
}
