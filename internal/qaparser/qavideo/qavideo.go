package qavideo

import (
	"io"

	"github.com/terratensor/svodd-server/internal/entities/answer"
)

// Parser Вопрос — Ответ видео выпуски парсер
type Parser struct{}

func (p *Parser) Parse(r io.Reader) (*[]answer.Entry, error) {
	var entries []answer.Entry // TODO: implement type
	return &entries, nil
}