package main

import (
	"fmt"
	"github.com/mehanizm/airtable"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type AirTableApi struct {
	Client *airtable.Client
}

func (a *AirTableApi) InitializeClient() {
	apiKey := os.Getenv("AIRTABLE_API_KEY")

	client := airtable.NewClient(apiKey)
	a.Client = client
}

func (a *AirTableApi) GetKanji() []Kanji {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Kanji")

	sortQuery := struct {
		FieldName string
		Direction string
	}{"meaning", "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields("kanji", "meaning", "kun_reading", "on_reading").
		MaxRecords(3).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var kanji []Kanji

	for _, record := range records.Records {
		character := record.Fields["kanji"].(string)
		escapedKanji := url.QueryEscape(character)
		baseURL := "https://raw.githubusercontent.com/jcsirot/kanji.gif/master/kanji/gif/150x150/"
		gifURL := baseURL + escapedKanji + ".gif"

		kanji = append(kanji, Kanji{
			Kanji:      character,
			Meaning:    record.Fields["meaning"].(string),
			KunReading: commaSeparatedList(record.Fields["kun_reading"].(string)),
			OnReading:  commaSeparatedList(record.Fields["on_reading"].(string)),
			GifUrl:     gifURL,
		})

	}

	return kanji
}

func (a *AirTableApi) GetVocabulary() []Vocabulary {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Vocabulary")

	sortQuery := struct {
		FieldName string
		Direction string
	}{"english", "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields("kanji", "kana", "english", "ex1_ja_furigana", "ex1_en").
		MaxRecords(5).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var vocabulary []Vocabulary

	for _, record := range records.Records {

		vocabulary = append(vocabulary, Vocabulary{
			Kanji:     record.Fields["kanji"].(string),
			Kana:      record.Fields["kana"].(string),
			Meaning:   record.Fields["english"].(string),
			ExampleJp: generateHtmlWithFurigana(record.Fields["ex1_ja_furigana"].(string)),
			ExampleEn: record.Fields["ex1_en"].(string),
		})
	}

	return vocabulary
}

func (a *AirTableApi) GetListening() string {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Listening")

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		ReturnFields("url").
		MaxRecords(1).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var record = records.Records[0]

	return record.Fields["url"].(string)
}

func (a *AirTableApi) GetBook() string {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Books")

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		ReturnFields("url").
		MaxRecords(1).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var record = records.Records[0]

	return record.Fields["url"].(string)
}

func (a *AirTableApi) GetGrammar() []Grammar {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Grammar")

	sortQuery := struct {
		FieldName string
		Direction string
	}{"meaning", "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields("grammar", "meaning", "ex1_ja_furigana", "ex1_en", "ex2_ja_furigana", "ex2_en").
		MaxRecords(2).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var grammar []Grammar

	for _, record := range records.Records {

		grammar = append(grammar, Grammar{
			Grammar:     record.Fields["grammar"].(string),
			Explanation: record.Fields["meaning"].(string),
			Example1Jp:  generateHtmlWithFurigana(record.Fields["ex1_ja_furigana"].(string)),
			Example1En:  record.Fields["ex1_en"].(string),
			Example2Jp:  generateHtmlWithFurigana(record.Fields["ex2_ja_furigana"].(string)),
			Example2En:  record.Fields["ex2_en"].(string),
		})
	}

	return grammar
}

func removeHTMLTags(s string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, " ")
}

func commaSeparatedList(s string) string {
	text := removeHTMLTags(s)
	items := strings.Fields(text)
	return strings.Join(items, ", ")
}

func updateRecords(table airtable.Table) {
	toUpdateRecords := &airtable.Records{
		Records: []*airtable.Record{
			{
				Fields: map[string]any{
					"Field1": "value1",
					"Field2": true,
				},
			},
			{
				Fields: map[string]any{
					"Field1": "value1",
					"Field2": true,
				},
			},
		},
	}
	_, err := table.UpdateRecords(toUpdateRecords)
	if err != nil {
		// Handle error
	}
}
