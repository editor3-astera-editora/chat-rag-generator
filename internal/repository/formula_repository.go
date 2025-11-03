package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Formula struct {
	Formula       string
	Description   string
	Concepts      map[string]any
	SourceChapter string
}

type FormulaRepository struct {
	DB *pgxpool.Pool
}

func NewFormulaRepository(db *pgxpool.Pool) *FormulaRepository {
	return &FormulaRepository{DB: db}
}

func (r *FormulaRepository) SearchFormula(ctx context.Context, bookName string, chapter int, query string) ([]Formula, error) {
	sql := `
		SELECT formula, description, concepts, source_chapter
		FROM formulas_map
		WHERE book_name = $1
			AND chapter = $2
			AND (formula ILIKE '%' || $3 || '%' OR description ILIKE '%' || $3 || '%')
		LIMIT 3;
	`

	rows, err := r.DB.Query(ctx, sql, bookName, chapter, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fórmulas: %w", err)
	}
	defer rows.Close()

	var formulas []Formula
	for rows.Next() {
		var f Formula
		err := rows.Scan(&f.Formula, &f.Description, &f.Concepts, &f.SourceChapter)
		if err == nil {
			formulas = append(formulas, f)
		}
	}
	return formulas, nil
}
