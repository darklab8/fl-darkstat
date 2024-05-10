package linker

import (
	"fmt"

	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/front/chunk"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

func TemplateChunks[T any](build *builder.Builder, items []T, chunk_template func(chunk.Chunk[T]), base_url utils_types.FilePath) []chunk.Chunk[T] {
	var chunks []chunk.Chunk[T]
	for i := 0; i < len(items)/100 || i == 0; i++ {
		var CurUrl string = fmt.Sprintf("%s_%d", base_url.ToString(), i)
		data := chunk.Chunk[T]{
			Items:      items[i*100 : (i+1)*100],
			CurrentUrl: CurUrl,
		}
		chunks = append(chunks, data)
	}

	for _, chunk := range chunks {
		chunk_template(chunk)
	}

	return chunks
}
