package middleware

import (
	"strings"
)

func CheckIdentity(question string) (bool, string) {
	lower := strings.ToLower(question)
	triggers := []string{
		"quem é você", "o que você é", "você é um", "você é real", "o que faz",
	}

	for _, t := range triggers {
		if strings.Contains(lower, t) {
			return true, "Sou seu assistente de inteligência artificial, preparado para te ajudar a estrudar com base nos livros e materiais didáticos do curso."
		}
	}
	return false, ""
}
