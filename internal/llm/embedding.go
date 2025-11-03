package llm

import (
	"context"
	"log"
	"math"
	"os"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))

	resp, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbedding3Large,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
	})
	if err != nil {
		return nil, err
	}

	emb64 := resp.Data[0].Embedding
	emb32 := make([]float32, len(emb64))
	var norm float64

	for _, v := range emb64 {
		norm += v * v
	}
	norm = math.Sqrt(norm)

	for i, v := range emb64 {
		emb32[i] = float32(v / norm)
	}
	log.Printf("Dimensão do embedding gerado: %d", len(emb32))
	return emb32, nil
}
