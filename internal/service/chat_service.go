package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/editor3-astera-editora/chat-rag/internal/llm"
	"github.com/editor3-astera-editora/chat-rag/internal/memory"
	"github.com/editor3-astera-editora/chat-rag/internal/model"
	"github.com/editor3-astera-editora/chat-rag/internal/repository"
)

type ChatService struct {
	Repo *repository.ChatRepository
}

var userBooks = map[string]string{
	"user123": "Matemática Financeira",
}

func NewChatService(repo *repository.ChatRepository) *ChatService {
	return &ChatService{Repo: repo}
}

func (s *ChatService) ProcessMessage(ctx context.Context, userID, content string) (string, []repository.RetrievedDoc, error) {
	memory.AddMessage(userID, model.Message{Role: "user", Content: content})

	bookName := userBooks[userID]

	// Detecta se a pergunta menciona "capítulo X"
	var chapterFilter *int
	reCap := regexp.MustCompile(`(?i)cap[ií]tulo\s+(\d+)`)
	if matches := reCap.FindStringSubmatch(content); len(matches) == 2 {
		if chap := matches[1]; chap != "" {
			var ch int
			fmt.Sscanf(chap, "%d", &ch)
			chapterFilter = &ch
		}
	}

	queryEmbedding, err := llm.GenerateEmbedding(ctx, content)
	if err != nil {
		return "", nil, fmt.Errorf("erro ao gerar embedding: %v", err)
	}

	// Busca apenas documentos do capítulo, se houver filtro
	docs, err := s.Repo.SearchSimilarEmbeddings(ctx, queryEmbedding, 3, chapterFilter)
	if err != nil {
		return "", nil, fmt.Errorf("erro ao buscar embeddings: %v", err)
	}

	// calcula a similaridade média (ou pega o primeiro resultado)
	if len(docs) == 0 || docs[0].Similarity < 0.75 {
		resposta := "Esse conteúdo não corresponde ao livro que estamos estudando. Tente reformular a pergunta dentro dos tópicos do material."
		memory.AddMessage(userID, model.Message{Role: "assistant", Content: resposta})
		return resposta, nil, nil
	}

	// Constrói contexto com origem
	var contextBuilder strings.Builder

	for _, d := range docs {
		if chapterFilter != nil {
			contextBuilder.WriteString(fmt.Sprintf("Livro: %s | Capítulo: %d\n%s\n---\n", bookName, *chapterFilter, d.Document))
		} else {
			contextBuilder.WriteString(fmt.Sprintf("Livro: %s\n%s\n---\n", bookName, d.Document))
		}
	}

	formulaRepo := repository.NewFormulaRepository(s.Repo.DB)
	formulaService := NewFormulaService(formulaRepo)

	chapterForFormula := 0
	if chapterFilter != nil {
		chapterForFormula = *chapterFilter
	} else if len(docs) > 0 {
		chapterForFormula = docs[0].Chapter
	}

	formulas, err := formulaService.GetRelevantFormulas(ctx, bookName, chapterForFormula, content)
	if err == nil && len(formulas) > 0 {
		formulaText := formulaService.FormatFormulas(formulas)
		contextBuilder.WriteString(formulaText)
	}

	history := memory.GetMemory(userID)
	var conversation strings.Builder
	conversation.WriteString("Você é um professor que explica conceitos de livros didáticos.\n")
	conversation.WriteString("\nContexto relevante:\n" + contextBuilder.String())

	for _, msg := range history {
		conversation.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	conversation.WriteString("\nPergunta: " + content)

	response, err := llm.GetLLMResponse(ctx, conversation.String())
	if err != nil {
		return "", nil, fmt.Errorf("erro ao chamar LLM: %v", err)
	}

	memory.AddMessage(userID, model.Message{Role: "assistant", Content: response})

	for i := range docs {
		docs[i].BookName = bookName
		if chapterFilter != nil {
			docs[i].Chapter = *chapterFilter
		}
	}

	return response, docs, nil
}
