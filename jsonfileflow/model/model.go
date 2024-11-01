package model

type Metadata struct {
	Seen       bool
	UpdateDate string
}

type Grammar struct {
	Grammar string
	Meaning string
	Metadata
}

type Kanji struct {
	Kanji      string
	OnReading  string
	KunReading string
	Meaning    string
	Examples   string
	Components string
	Koohii1    string
	Koohii2    string
	Metadata
}

type Vocabulary struct {
	Kanji   string
	Reading string
	Meaning string
	Metadata
}

type ShortTextResponse struct {
	TextJapanese                       string                       `json:"text_japanese"`
	TextEnglish                        string                       `json:"text_english"`
	ReadingComprehensionQuizInJapanese ReadingComprehensionQuizItem `json:"reading_comprehension_quiz_in_japanese"`
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

type SongsMetadata struct {
	Metadata []SongResponse `json:"metadata"`
	Style    string         `json:"style"`
}

type SongResponse struct {
	AudioURL string `json:"audio_url"`
	ImageURL string `json:"image_url"`
}

type File struct {
	Kanji             []Kanji                  `json:"kanji"`
	Vocabulary        []Vocabulary             `json:"vocabulary"`
	Grammar           []Grammar                `json:"grammar"`
	ShortTextWithQuiz ShortTextResponse        `json:"short_text_with_quiz"`
	SongLyrics        SongLyricsResponse       `json:"song_lyrics"`
	SongsMetadata     SongsMetadata            `json:"song_metadata"`
	KanjiImages       KanjiImageResponse       `json:"kanji_images"`
	GrammarSentences  GrammarSentencesResponse `json:"grammar_sentences"`
}
