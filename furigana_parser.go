package main

import (
	"strings"
)

func generateHtmlWithFurigana(sentence string) string {
	var html strings.Builder
	html.WriteString("<span>")

	sentenceParts := strings.Fields(sentence)

	for _, sentencePart := range sentenceParts {
		indexOfOpenParen := strings.Index(sentencePart, "[")
		indexOfCloseParen := strings.Index(sentencePart, "]")

		if indexOfCloseParen == -1 {
			html.WriteString(sentencePart)
			continue
		}
		kanji := sentencePart[:indexOfOpenParen]
		furigana := sentencePart[indexOfOpenParen+1 : indexOfCloseParen]
		html.WriteString(addRuby(kanji, furigana))
		html.WriteString(sentencePart[indexOfCloseParen+1:])
	}

	html.WriteString("</span>")
	return html.String()
}

func addRuby(kanji string, furigana string) string {
	var html strings.Builder

	html.WriteString("<ruby><span>")
	html.WriteString(furigana)
	html.WriteString("</span><rt>")
	html.WriteString(kanji)
	html.WriteString("</rt>")
	html.WriteString("</ruby>")

	return html.String()
}
