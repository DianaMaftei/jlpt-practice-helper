package main

type Kanji struct {
	Kanji       string
	Meaning     string
	OnReading   string
	KunReading  string
	GifUrl      string
	KanjiDetail KanjiDetail
}

type Vocabulary struct {
	Kanji      string
	Kana       string
	Meaning    string
	Example1Jp string
	Example1En string
	Example2Jp string
	Example2En string
}

type Grammar struct {
	Grammar     string
	Explanation string
	Example1Jp  string
	Example1En  string
	Example2Jp  string
	Example2En  string
}

type Book struct {
	Url string
	Img string
}

type KanjiDetail struct {
	Examples []Example `json:"examples"`
	Hint     string    `json:"mn_hint"`
}

type Example struct {
	Japanese string `json:"japanese"`
	Meaning  struct {
		English string `json:"english"`
	} `json:"meaning"`
}
