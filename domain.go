package main

type Kanji struct {
	Kanji      string
	Meaning    string
	OnReading  string
	KunReading string
	GifUrl     string
}

type Vocabulary struct {
	Kanji     string
	Kana      string
	Meaning   string
	ExampleJp string
	ExampleEn string
}

type Grammar struct {
	Grammar     string
	Explanation string
	Example1Jp  string
	Example1En  string
	Example2Jp  string
	Example2En  string
	Bunpro      string
}

type Book struct {
	Url string
	Img string
}
