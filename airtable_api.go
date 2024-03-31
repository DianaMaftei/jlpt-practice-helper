package main

import (
	"fmt"
	"github.com/mehanizm/airtable"
	"log"
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
	}{"record_id", "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields("kanji", "meaning", "kun_reading", "on_reading").
		MaxRecords(5).
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

		kunReading, ok := record.Fields["kun_reading"].(string)
		if !ok {
			kunReading = ""
		}

		onReading, ok := record.Fields["on_reading"].(string)
		if !ok {
			onReading = ""
		}

		kanji = append(kanji, Kanji{
			Kanji:      character,
			Meaning:    record.Fields["meaning"].(string),
			KunReading: commaSeparatedList(kunReading),
			OnReading:  commaSeparatedList(onReading),
			GifUrl:     gifURL,
		})

	}

	updateRecords(table, records)

	return kanji
}

func (a *AirTableApi) GetVocabulary() []Vocabulary {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Vocabulary")

	sortQuery := struct {
		FieldName string
		Direction string
	}{"record_id", "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields("kanji", "kana", "english", "ex1_ja_furigana", "ex1_en").
		MaxRecords(30).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var vocabulary []Vocabulary

	for _, record := range records.Records {

		exampleJp, ok := record.Fields["ex1_ja_furigana"].(string)
		if !ok {
			exampleJp = ""
		}

		exampleEn, ok := record.Fields["ex1_en"].(string)
		if !ok {
			exampleEn = ""
		}

		vocabulary = append(vocabulary, Vocabulary{
			Kanji:     record.Fields["kanji"].(string),
			Kana:      record.Fields["kana"].(string),
			Meaning:   record.Fields["english"].(string),
			ExampleJp: generateHtmlWithFurigana(exampleJp),
			ExampleEn: exampleEn,
		})
	}

	updateRecords(table, records)

	return vocabulary
}

func (a *AirTableApi) GetGrammar() []Grammar {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Grammar")

	sortQuery := struct {
		FieldName string
		Direction string
	}{"record_id", "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields("grammar", "meaning", "ex1_ja_furigana", "ex1_en", "ex2_ja_furigana", "ex2_en", "bunpro").
		MaxRecords(2).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var grammar []Grammar

	for _, record := range records.Records {

		exampleJp1, ok := record.Fields["ex1_ja_furigana"].(string)
		if !ok {
			exampleJp1 = ""
		}

		exampleEn1, ok := record.Fields["ex1_en"].(string)
		if !ok {
			exampleEn1 = ""
		}

		exampleJp2, ok := record.Fields["ex2_ja_furigana"].(string)
		if !ok {
			exampleJp2 = ""
		}

		exampleEn2, ok := record.Fields["ex2_en"].(string)
		if !ok {
			exampleEn2 = ""
		}

		bunpro, ok := record.Fields["bunpro"].(string)
		if !ok {
			bunpro = ""
		}

		grammar = append(grammar, Grammar{
			Grammar:     record.Fields["grammar"].(string),
			Explanation: record.Fields["meaning"].(string),
			Example1Jp:  generateHtmlWithFurigana(exampleJp1),
			Example1En:  exampleEn1,
			Example2Jp:  generateHtmlWithFurigana(exampleJp2),
			Example2En:  exampleEn2,
			Bunpro:      bunpro,
		})
	}

	updateRecords(table, records)

	return grammar
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
	_, err = record.UpdateRecordPartial(map[string]any{"seen": true})
	if err != nil {
		log.Fatal(err)
	}

	return record.Fields["url"].(string)
}

func (a *AirTableApi) GetBook() Book {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Books")

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		ReturnFields("url", "img").
		MaxRecords(1).
		Do()
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	var record = records.Records[0]
	_, err = record.UpdateRecordPartial(map[string]any{"seen": true})
	if err != nil {
		log.Fatal(err)
	}

	return Book{
		Url: record.Fields["url"].(string),
		Img: record.Fields["img"].(string),
	}
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

func updateRecords(table *airtable.Table, records *airtable.Records) {

	r := &airtable.Records{}
	r.Records = []*airtable.Record{}

	for _, record := range records.Records {
		newRecord := airtable.Record{
			ID: record.ID,
			Fields: map[string]interface{}{
				"seen": true,
			},
		}

		r.Records = append(r.Records, &newRecord)
	}

	_, err := table.UpdateRecordsPartial(r)
	if err != nil {
		log.Fatal(err)
	}
}
