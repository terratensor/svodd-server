package qaquestion

import (
	"io"
)

// Parser Вопрос — Ответ Список вопросов парсер
type Parser struct{}

func (p *Parser) Parse(r io.Reader) (*Feed, error) {
	var feed Feed // TODO: implement type
	return &feed, nil
}
