package questionanswer

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/lib/httpclient"
	"golang.org/x/net/html"
)

type Entry struct {
	Url         *url.URL
	Title       string
	Video       *url.URL
	Datetime    *time.Time
	Content     []QuestionAnswer
	Fragments   []Fragment
	Comments    []Comment
	Html        string
	moderator   []string
	responsible []string
}

type Fragment struct {
	QuestionAnswer string
	Chunk          int
}

type QuestionAnswer struct {
	Question []string
	Answer   []string
}

type Comment struct {
	Username   string
	Text       string
	AvatarFile *url.URL
	Role       string
	Datetime   *time.Time
	DataID     string
	Type       int
	Position   int
}

const TypeComment = 3

var defaultModerator = []string{"Ведущий", "Ведущая:", "Дмитрий Таран:", "Сергей Будков:", "ВедущийЯ:", "Айнис Казимирович Петкус"}
var defaultResponsible = []string{"Валерий Викторович Пякин:", "Валерий Викторович", "Пякин Валерий Викторович", "В.В. Пякин:", "Валерий"}

func NewEntry(url *url.URL, cfg *config.Config) *Entry {
	moderator := cfg.Questionanswer.Moderator
	if condition := len(moderator) == 0; condition {
		moderator = defaultModerator
	}
	responsible := cfg.Questionanswer.Responsible
	if condition := len(responsible) == 0; condition {
		responsible = defaultResponsible
	}
	return &Entry{
		Url:         url,
		moderator:   moderator,
		responsible: responsible,
	}
}

func (entry *Entry) FetchData(client *httpclient.HttpClient) error {

	responseBytes, err := client.Get(entry.Url)
	if err != nil {
		return err
	}
	err = entry.Parse(responseBytes)
	if err != nil {
		return err
	}

	entry.splitAnswers()

	return nil
}

func (e *Entry) Parse(resBytes []byte) error {

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resBytes))
	if err != nil {
		return err
	}
	el := doc.Find(".block").First()

	e.Title = el.Find("h1").First().Text()

	e.Video, err = url.Parse(el.Find(".embed-responsive iframe").AttrOr("src", ""))
	if err != nil {
		log.Printf("failed to parse video url: %v", err)
	}

	dateStr := el.Find(".datetime").First().Text()
	dateStr = strings.TrimSpace(dateStr)
	datetime, err := time.Parse("15:04 02.01.2006", dateStr)
	if err != nil {
		log.Printf("failed to parse datetime: %v", err)
	} else {
		e.Datetime = &datetime
	}

	els := doc.Find("#answer-content").First()
	// Populate the `Html` field of the `Entry` struct with the outer HTML of the
	// "#answer-content" element.
	//
	// This field will be used later to create a new QuestionAnswer struct.
	//
	// Parameters:
	//   - els: The goquery Selection representing the "#answer-content" element.
	//
	// Return type: None.
	e.Html = populateOuterHtml(els)
	e.SplitIntoChunks(els)

	els = doc.Find(".comment-list").First()
	els.Find(".comment-item").Each(func(i int, s *goquery.Selection) {
		// сохраняем текст комментария вместе с html кодом
		text, err := goquery.OuterHtml(s.Find(".comment-text"))
		if err != nil {
			log.Printf("failed to get comment html: %v", err)
			text = s.Find(".comment-text").Text()
		}

		e.Comments = append(e.Comments, Comment{
			Username:   strings.TrimSpace(s.Find(".username").Text()),
			Text:       strings.TrimSpace(text),
			AvatarFile: parseAvatarFile(s.Find(".ava-80").AttrOr("src", "")),
			Role:       strings.TrimSpace(s.Find(".role").Text()),
			Datetime:   parseDatetime(s.Find(".datetime").Text()),
			DataID:     s.Find(".comment-text").AttrOr("data-id", ""),
			Type:       TypeComment,
			Position:   i + 1,
		})
	})
	return nil
}

func parseDatetime(datetime string) *time.Time {
	t, err := time.Parse("15:04 02.01.2006", strings.TrimSpace(datetime))
	if err != nil {
		log.Printf("failed to parse datetime: %v", err)
	}
	return &t
}

func parseAvatarFile(avatarFile string) *url.URL {
	u, err := url.Parse(strings.TrimSpace(avatarFile))
	if err != nil {
		log.Printf("failed to parse avatar url: %v", err)
	}
	return u
}

// SplitIntoChunks разбивает текст на вопросы и ответы.
// Он основан на поиске конкретных строк в тексте.
// Если нашелся текст "Ведущий:", то он начинает добавлять текст в массив вопросов.
// Если нашелся текст "Валерий Викторович Пякин:", то он начинает добавлять текст в массив ответов.
// Если нашелся текст "Ведущий:" или "Валерий Викторович Пякин:" в середине массива текста,
// то он создает новый QuestionAnswer, добавляет его в массив Content и начинает новый цикл.
func (e *Entry) SplitIntoChunks(els *goquery.Selection) {
	// isQuestion - это флаг, который говорит, что мы находимся в вопросе.
	isQuestion := false
	// isAnswer - это флаг, который говорит, что мы находимся в ответе.
	isAnswer := false

	// questionAnswer - это массив, который содержит вопросы и ответы.
	var questionAnswer []QuestionAnswer
	// question - это массив, который содержит текст вопроса.
	var question []string
	// answer - это массив, который содержит текст ответа.
	var answer []string

	// Мы идем по всем абзацам текста.
	els.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())

		// Мы ищем текст "Ведущий:".
		moderatorIndex, curModerator := checkStrIndex(text, e.moderator)
		// Мы ищем текст "Валерий Викторович Пякин:".
		responsibleIndex, curResponsible := checkStrIndex(text, e.responsible)

		// Если нашелся текст "Ведущий:", то мы начинаем новый вопрос.
		if moderatorIndex == 0 {
			// Если у нас уже есть вопрос и ответ, то мы создаем новый QuestionAnswer.
			if len(question) > 0 && len(answer) > 0 {
				questionAnswer = append(questionAnswer, QuestionAnswer{Question: question, Answer: answer})
				question = nil
				answer = nil
			}

			// Мы добавляем текст в массив вопроса.
			text = WrapPhrase(*curModerator, text)
			question = append(question, text)
			isQuestion = true
			isAnswer = false
		}

		// Если нашелся текст "Валерий Викторович Пякин:", то мы начинаем новый ответ.
		if responsibleIndex == 0 {

			text = WrapPhrase(*curResponsible, text)

			// Мы добавляем текст в массив ответа.
			answer = append(answer, text)
			isAnswer = true
			isQuestion = false
		}

		// Если у нас есть текст "Ведущий:", то мы добавляем его в массив вопроса.
		if moderatorIndex != 0 && isQuestion {
			question = append(question, text)
		}
		// Если у нас есть текст "Валерий Викторович Пякин:", то мы добавляем его в массив ответа.
		if responsibleIndex != 0 && isAnswer {
			answer = append(answer, text)
		}
	})

	// Если у нас есть вопрос и ответ, то мы создаем новый QuestionAnswer.
	if len(question) > 0 && len(answer) > 0 {
		questionAnswer = append(questionAnswer, QuestionAnswer{Question: question, Answer: answer})
		question = nil
		answer = nil
	}

	e.Content = questionAnswer
}

func (e *Entry) splitAnswers() {

	var result []Fragment
	for _, qa := range e.Content {

		isNewFragment := true
		startAnswer := false
		chunk := 1
		fragment := Fragment{QuestionAnswer: "", Chunk: chunk}

		for _, ans := range qa.Answer {

			if isNewFragment {
				for _, q := range qa.Question {
					if strings.TrimSpace(q) == "" {
						continue
					}
					fragment.QuestionAnswer += fmt.Sprintf("<p class=\"question\">%v</p>", q)
				}
				isNewFragment = false
				startAnswer = true
			}

			// Check if the answer starts with "Валерий Викторович Пякин:"
			// and if it doesn't already have a <strong> tag.
			// If so, add a <strong> tag to the answer.
			// If not, add the answer without any changes.
			if startAnswer {
				// Check if the answer starts with "Валерий Викторович Пякин:".
				responsibleIndex, _ := checkStrIndex(ans, e.responsible)
				if responsibleIndex != 0 && strings.Index(ans, "<strong>") != 0 {
					// Add a <strong> tag to the answer.
					// Условие если ans не пустой, во избежание висящих … без текста, когда были вырезаны изображения из исходника.
					if ans != "" {
						fragment.QuestionAnswer += fmt.Sprintf(
							"<p class=\"answer\"><strong>Валерий Викторович: … </strong>%v</p>",
							ans,
						)
						startAnswer = false
					}
				} else {
					// Add the answer without any changes.
					fragment.QuestionAnswer += fmt.Sprintf("<p class=\"answer\">%v</p>", ans)
					startAnswer = false
				}
			} else {
				// Add the answer without any changes.
				fragment.QuestionAnswer += fmt.Sprintf("<p class=\"answer\">%v</p>", ans)
			}

			if (utf8.RuneCountInString(fragment.QuestionAnswer)) > 2700 {
				result = append(result, fragment)
				chunk++
				fragment = Fragment{QuestionAnswer: "", Chunk: chunk}
				isNewFragment = true
			}
		}

		if utf8.RuneCountInString(fragment.QuestionAnswer) > 0 {
			result = append(result, fragment)
		}
	}

	e.Fragments = result
}

func checkStrIndex(text string, str []string) (int, *string) {
	// Iterate over the strings in the array
	for _, r := range str {
		// Check if the current string is found at the beginning of the text
		if strings.Index(text, r) == 0 {
			return 0, &r // Return the index and a pointer to the matched string
		}
	}
	return -1, nil // Return -1 if no match is found
}

// WrapPhrase wraps a specific phrase in a text with a strong tag.
func WrapPhrase(phrase, text string) string {
	index := strings.Index(text, phrase)
	if index == -1 {
		return text
	}

	prefix := text[:index]
	suffix := text[index+len(phrase):]
	wrapped := fmt.Sprintf("%v<strong>%v</strong>%v", prefix, phrase, suffix)

	return wrapped
}

// populateOuterHtml generates the outer HTML representation of the given goquery Selection.
// This func keeps track of whether the current node is within a `<table>`,
// and only removes "style" attributes from nodes outside the table.
//
// Parameters:
// - els: a pointer to the goquery Selection
// Returns:
// - a string containing the generated outer HTML
func populateOuterHtml(els *goquery.Selection) string {
	htmlStr, err := goquery.OuterHtml(els)
	if err != nil {
		log.Printf("failed to get html: %v", err)
	}

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		panic(err)
	}

	removeStyle(doc, false)

	var buf strings.Builder
	html.Render(&buf, doc)
	return buf.String()
}

// removeStyle removes the "style" attribute from HTML elements.
// It skips elements inside tables because they often have inline styles.
//
// Parameters:
// - n: the HTML node to remove the "style" attribute from.
// - inTable: a boolean indicating whether the current node is inside a table.
func removeStyle(n *html.Node, inTable bool) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "table":
			inTable = true
		case "/table":
			inTable = false
		default:
			if !inTable {
				for i := len(n.Attr) - 1; i >= 0; i-- {
					if n.Attr[i].Key == "style" {
						copy(n.Attr[i:], n.Attr[i+1:])
						n.Attr = n.Attr[:len(n.Attr)-1]
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeStyle(c, inTable)
	}
}
