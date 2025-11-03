package middleware

import (
	"context"
	"strings"

	"github.com/editor3-astera-editora/chat-rag/internal/repository"
)

func CheckFormula(ctx context.Context, repo *repository.FormulaRepository, bookName string, chapter int, question string) ([]repository.Formula, bool, error) {
	// Palavras indicativas de fórmula
	triggers := []string{"fórmula", "equação", "calcular", "como calcular", "como achar", "como encontrar", "calcule"}

	lower := strings.ToLower(question)
	relevant := false
	for _, t := range triggers {
		if strings.Contains(lower, t) {
			relevant = true
			break
		}
	}

	if !relevant {
		return nil, false, nil
	}

	formulas, err := repo.SearchFormula(ctx, bookName, chapter, question)
	if err != nil {
		return nil, false, err
	}

	return formulas, true, nil
}
