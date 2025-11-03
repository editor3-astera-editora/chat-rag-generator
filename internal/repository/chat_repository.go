package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type RetrievedDoc struct {
	Document   string
	BookName   string
	Unit       int
	Chapter    int
	Similarity float64
}

type ChatRepository struct {
	DB *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{DB: db}
}

// SearchSimilarEmbeddings retorna documentos mais próximos vetorialmente
func (r *ChatRepository) SearchSimilarEmbeddings(ctx context.Context, embedding []float32, limit int, chapterFilter *int) ([]RetrievedDoc, error) {
	query := `
		SELECT
			document,
			COALESCE(cmetadata->>'book_name', '') AS book_name,
			(cmetadata->>'unit')::int AS unit,
			(cmetadata->>'chapter')::int AS chapter,
			1 - (embedding <=> $1) / 2 AS similarity
		FROM langchain_pg_embedding
		WHERE ($3::int IS NULL OR (cmetadata->>'chapter')::int = $3::int)
		ORDER BY embedding <=> $1
		LIMIT $2;
	`

	vec := pgvector.NewVector(embedding)

	var chValue interface{}
	if chapterFilter != nil {
		chValue = *chapterFilter
	} else {
		chValue = nil
	}

	rows, err := r.DB.Query(ctx, query, vec, limit, chValue)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar busca vetorial: %w", err)
	}
	defer rows.Close()

	var docs []RetrievedDoc
	for rows.Next() {
		var doc RetrievedDoc
		err := rows.Scan(&doc.Document, &doc.BookName, &doc.Unit, &doc.Chapter, &doc.Similarity)
		if err != nil {
			log.Println("erro ao ler linha: ", err)
			continue
		}
		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar resultados: %w", err)
	}

	if len(docs) > 0 {
		log.Printf("🔍 %d documentos encontrados. Maior similaridade: %.3f", len(docs), docs[0].Similarity)
		for i, d := range docs {
			log.Printf("   → [%d] %s | Unidade %d, Capítulo %d | Sim=%.3f",
				i+1, d.BookName, d.Unit, d.Chapter, d.Similarity)
		}
	} else {
		log.Println("🔍 Nenhum documento similar encontrado.")
	}

	return docs, nil
}
