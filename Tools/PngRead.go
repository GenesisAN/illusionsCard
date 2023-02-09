package Tools

import (
	"bytes"
	"illusionCard/Base"
)

func get_png(file []byte) int {
	res1 := bytes.Index(file, Base.PngEndChunk)
	return res1
}
