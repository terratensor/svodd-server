package qaparser

import (
	"net/url"
)

type FeedType int

const (
	// FeedTypeUnknown represents a feed that could not have its
	// type determiend.
	FeedTypeUnknown FeedType = iota
	// FeedTypeQA представляет Вопрос — Ответ https://xn----8sba0bbi0cdm.xn--p1ai/qa/video
	FeedTypeQA
	// FeedTypeQAQuestion представляет Вопрос — Ответ Список вопросов https://xn----8sba0bbi0cdm.xn--p1ai/qa/question
	FeedTypeQAQuestion
)

func DetectFeedType(link *url.URL) FeedType {

	if link.Path == "/qa/video" {
		return FeedTypeQA
	}
	if link.Path == "/qa/question" {
		return FeedTypeQAQuestion
	}
	return FeedTypeUnknown
}
