package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/editor3-astera-editora/chat-rag/internal/repository"
)

type FormulaService struct {
	Repo *repository.FormulaRepository
}

func NewFormulaService(repo *repository.FormulaRepository) *FormulaService {
	return &FormulaService{Repo: repo}
}

func (s *FormulaService) GetRelevantFormulas(ctx context.Context, bookName string, chapter int, question string) ([]repository.Formula, error) {
	formulas, err := s.Repo.SearchFormula(ctx, bookName, chapter, question)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fórmulas: %w", err)
	}
	return formulas, nil
}

func (s *FormulaService) FormatFormulas(formulas []repository.Formula) string {
	if len(formulas) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\nAs seguintes fórmulas podem ser úteis para responder esta pergunta:\n")
	for i, f := range formulas {
		b.WriteString(fmt.Sprintf(
			"\n%d) Fórmula: %s\nDescrição: %s\nFonte: Capítulo %s\n---\n",
			i+1, f.Formula, f.Description, f.SourceChapter,
		))
	}
	return b.String()
}
