package middleware

import (
	"context"

	"github.com/editor3-astera-editora/chat-rag/internal/repository"
)

func CheckSimilarity(ctx context.Context, repo *repository.ChatRepository, embedding []float32, threshold float64) ([]repository.RetrievedDoc, bool, error) {
	docs, err := repo.SearchSimilarEmbeddings(ctx, embedding, 3, nil)
	if err != nil {
		return nil, false, err
	}

	if len(docs) == 0 || docs[0].Similarity < threshold {
		return nil, false, nil
	}

	return docs, true, nil
}
