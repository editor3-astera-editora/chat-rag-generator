package llm

import (
	"context"
	"os"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func GetLLMResponse(ctx context.Context, content string) (string, error) {
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Você é um professor que explica conceitos de livros didáticos, de forma clara e objetiva."),
			openai.UserMessage(content),
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "Sem resposta do modelo.", nil
	}

	return resp.Choices[0].Message.Content, nil
}
