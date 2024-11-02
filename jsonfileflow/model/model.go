package model

type Metadata struct {
	Seen       bool   // `json:"seen"`
	UpdateDate string // `json:"update_date"`
}

type Grammar struct {
	Grammar string // `json:"grammar"`
	Meaning string // `json:"meaning"`
	Metadata
}

type Kanji struct {
	Kanji      string // `json:"kanji"`
	OnReading  string // `json:"on_reading"`
	KunReading string // `json:"kun_reading"`
	Meaning    string // `json:"meaning"`
	Examples   string // `json:"examples"`
	Components string // `json:"components"`
	Koohii1    string // `json:"koohii1"`
	Koohii2    string // `json:"koohii2"`
	Metadata
}

type Vocabulary struct {
	Kanji   string // `json:"kanji"`
	Reading string // `json:"reading"`
	Meaning string // `json:"meaning"`
	Metadata
}

type ShortTextResponse struct {
	TextJapanese                       string                         `json:"text_japanese"`
	TextEnglish                        string                         `json:"text_english"`
	ReadingComprehensionQuizInJapanese []ReadingComprehensionQuizItem `json:"reading_comprehension_quiz_in_japanese"`
}

type ReadingComprehensionQuizItem struct {
	Question string                       `json:"question"`
	Options  []ReadingComprehensionOption `json:"options"`
}

type ReadingComprehensionOption struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type SongLyricsResponse struct {
	LyricsJapanese           string `json:"lyrics_japanese"`
	LyricsEnglishTranslation string `json:"lyrics_english_translation"`
	SongTitle                string `json:"song_title"`
}

type KanjiImageResponse struct {
	KanjiList []KanjiImageItem `json:"kanji_list"`
}

type KanjiImageItem struct {
	Kanji            string `json:"kanji"`
	ImageDescription string `json:"image_description"`
	EncodedImage     string `json:"encoded_image"`
}

type GrammarSentencesResponse struct {
	GrammarPoint string            `json:"grammar_point"`
	Sentences    []GrammarSentence `json:"sentences"`
}

type GrammarSentence struct {
	SentenceJapanese           string `json:"sentence_japanese"`
	SentenceEnglishTranslation string `json:"sentence_english_translation"`
}

type Song struct {
	Ids        []string           `json:"id"`
	Style      string             `json:"style"`
	SongLyrics SongLyricsResponse `json:"song_lyrics"`
}

type SongResponse struct {
	Id string `json:"id"`
}

type GrammarInfo struct {
	Grammar          Grammar                  `json:"grammar"`
	GrammarSentences GrammarSentencesResponse `json:"grammar_sentences"`
}

type File struct {
	Kanji             []Kanji            `json:"kanji"`
	Vocabulary        []Vocabulary       `json:"vocabulary"`
	GrammarInfo       []GrammarInfo      `json:"grammar_info"`
	ShortTextWithQuiz ShortTextResponse  `json:"short_text_with_quiz"`
	Song              Song               `json:"song"`
	KanjiImages       KanjiImageResponse `json:"kanji_images"`
}
