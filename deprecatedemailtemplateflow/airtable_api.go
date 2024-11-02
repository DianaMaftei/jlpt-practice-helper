package deprecatedemailtemplateflow

import (
	"fmt"
	"github.com/mehanizm/airtable"
	"log"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
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
	var kanji []Kanji

	table := a.Client.GetTable("appuDdot90oIORePj", "Kanji")
	records, err := getSortedUnseenRecords(table, "record_id", 6, "kanji", "meaning", "kun_reading", "on_reading")
	if err != nil {
		fmt.Println(err)
		return kanji
	}

	for _, record := range records.Records {
		character := getStringFieldOrDefaultEmpty(record, "kanji")
		kunReading := getStringFieldOrDefaultEmpty(record, "kun_reading")
		onReading := getStringFieldOrDefaultEmpty(record, "on_reading")
		meaning := getStringFieldOrDefaultEmpty(record, "meaning")

		kanjiAliveData, err := getKanjiAliveData(character)
		if err != nil {
			fmt.Println(err)
			continue
		}

		currentKanji := Kanji{
			Kanji:       character,
			Meaning:     meaning,
			KunReading:  commaSeparatedList(kunReading),
			OnReading:   commaSeparatedList(onReading),
			GifUrl:      getKanjiStrokesGifUrl(character),
			KanjiDetail: *kanjiAliveData,
		}
		currentKanji.KanjiDetail.Hint = parseHint(currentKanji.KanjiDetail.Hint)
		currentKanji.KanjiDetail.Examples = shuffle(currentKanji.KanjiDetail.Examples)

		kanji = append(kanji, currentKanji)
	}

	markRecordsAsSeen(table, records)
	return kanji
}

func (a *AirTableApi) GetVocabulary() []Vocabulary {
	var vocabulary []Vocabulary

	table := a.Client.GetTable("appuDdot90oIORePj", "Vocabulary")
	records, err := getSortedUnseenRecords(table, "record_id", 19, "kanji", "kana", "english", "ex1_ja_furigana", "ex1_en")
	if err != nil {
		fmt.Println(err)
		return vocabulary
	}

	for _, record := range records.Records {
		kanji := getStringFieldOrDefaultEmpty(record, "kanji")
		kana := getStringFieldOrDefaultEmpty(record, "kana")
		english := getStringFieldOrDefaultEmpty(record, "english")
		example1Jp := getStringFieldOrDefaultEmpty(record, "ex1_ja_furigana")
		example1En := getStringFieldOrDefaultEmpty(record, "ex1_en")
		example2Jp := getStringFieldOrDefaultEmpty(record, "ex2_ja_furigana")
		example2En := getStringFieldOrDefaultEmpty(record, "ex2_en")

		vocabulary = append(vocabulary, Vocabulary{
			Kanji:      kanji,
			Kana:       kana,
			Meaning:    english,
			Example1Jp: generateHtmlWithFurigana(example1Jp),
			Example1En: example1En,
			Example2Jp: generateHtmlWithFurigana(example2Jp),
			Example2En: example2En,
		})
	}

	markRecordsAsSeen(table, records)
	return vocabulary
}

func (a *AirTableApi) GetGrammar() []Grammar {
	var grammar []Grammar

	table := a.Client.GetTable("appuDdot90oIORePj", "Grammar")
	records, err := getSortedUnseenRecords(table, "record_id", 2, "grammar", "meaning", "ex1_ja_furigana", "ex1_en", "ex2_ja_furigana", "ex2_en")
	if err != nil {
		fmt.Println(err)
		return grammar
	}

	for _, record := range records.Records {
		exampleJp1 := getStringFieldOrDefaultEmpty(record, "ex1_ja_furigana")
		exampleEn1 := getStringFieldOrDefaultEmpty(record, "ex1_en")
		exampleJp2 := getStringFieldOrDefaultEmpty(record, "ex2_ja_furigana")
		exampleEn2 := getStringFieldOrDefaultEmpty(record, "ex2_en")

		grammar = append(grammar, Grammar{
			Grammar:     record.Fields["grammar"].(string),
			Explanation: record.Fields["meaning"].(string),
			Example1Jp:  generateHtmlWithFurigana(exampleJp1),
			Example1En:  exampleEn1,
			Example2Jp:  generateHtmlWithFurigana(exampleJp2),
			Example2En:  exampleEn2,
		})
	}

	markRecordsAsSeen(table, records)
	return grammar
}

func (a *AirTableApi) GetListening() string {
	table := a.Client.GetTable("appuDdot90oIORePj", "Listening")
	records, err := getSortedUnseenRecords(table, "url", 1, "url")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if len(records.Records) == 0 {
		return ""
	}

	var record = records.Records[0]
	_, err = record.UpdateRecordPartial(map[string]any{"seen": true})
	if err != nil {
		log.Fatal(err)
	}

	return record.Fields["url"].(string)
}

func (a *AirTableApi) GetBook() Book {
	table := a.Client.GetTable("appuDdot90oIORePj", "Books")

	records, err := getSortedUnseenRecords(table, "title", 1, "url", "img")
	if err != nil {
		fmt.Println(err)
		return Book{}
	}

	if len(records.Records) == 0 {
		return Book{}
	}

	markRecordsAsSeen(table, records)

	return Book{
		Url: records.Records[0].Fields["url"].(string),
		Img: records.Records[0].Fields["img"].(string),
	}
}

func getSortedUnseenRecords(table *airtable.Table, sortField string, maxRecords int, fields ...string) (*airtable.Records, error) {
	sortQuery := struct {
		FieldName string
		Direction string
	}{sortField, "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields(fields...).
		MaxRecords(maxRecords).
		Do()

	return records, err
}

func markRecordsAsSeen(table *airtable.Table, records *airtable.Records) {
	const maxRecordsPerRequest = 10

	for i := 0; i < len(records.Records); i += maxRecordsPerRequest {
		end := i + maxRecordsPerRequest
		if end > len(records.Records) {
			end = len(records.Records)
		}

		r := &airtable.Records{}
		r.Records = []*airtable.Record{}

		for _, record := range records.Records[i:end] {
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

func getStringFieldOrDefaultEmpty(record *airtable.Record, fieldName string) string {
	field, ok := record.Fields[fieldName].(string)
	if !ok {
		field = ""
	}

	return field
}

func getKanjiStrokesGifUrl(character string) string {
	baseURL := "https://raw.githubusercontent.com/jcsirot/kanji.gif/master/kanji/gif/150x150/"
	escapedKanji := url.QueryEscape(character)

	return baseURL + escapedKanji + ".gif"
}

func parseHint(hint string) string {
	re := regexp.MustCompile(`\[(\d+)\]`)
	return re.ReplaceAllStringFunc(hint, func(match string) string {
		number := match[1 : len(match)-1]
		imgTag := fmt.Sprintf(`<img width="14" src="https://media.kanjialive.com/mnemonic_hints/%s.svg"/>`, number)
		return imgTag
	})
}

func shuffle(list []Example) []Example {
	rand.Seed(time.Now().UnixNano())
	for i := len(list) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}

	return list
}
