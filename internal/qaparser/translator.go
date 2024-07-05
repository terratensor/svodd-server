package qaparser

type Translator interface {
	Translate(feed interface{}) (*[]Entry, error)
}

// DefaultRSSTranslator converts an rss.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between rss.Feed -> Feed
// for each of the fields in Feed.
type DefaultQAVideoTranslator struct{}

func (t *DefaultQAVideoTranslator) Translate(feed interface{}) (*[]Entry, error) {
	// rss, found := feed.(*qavideo.Feed)

	result := &[]Entry{}

	return result, nil
}

type DefaultQAQuestionTranslator struct{}

func (t *DefaultQAQuestionTranslator) Translate(feed interface{}) (*[]Entry, error) {
	result := &[]Entry{}

	return result, nil
}
