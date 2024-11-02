package dbs

import (
	"encoding/csv"
	"jlpt-practice-helper/jsonfileflow/model"
	"os"
	"strconv"
)

var grammarFilePath = "static/data/grammar.csv"

func LoadGrammarFromFile() ([]model.Grammar, error) {
	var grammars []model.Grammar

	f, err := os.Open(grammarFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records[1:] { // Skip header
		var seen bool
		if record[2] == "" {
			seen = false
		} else {
			seen, err = strconv.ParseBool(record[2])
			if err != nil {
				return nil, err
			}
		}

		if err != nil {
			return nil, err
		}
		grammars = append(grammars, model.Grammar{
			Grammar: record[0],
			Meaning: record[1],
			Metadata: model.Metadata{
				Seen:       seen,
				UpdateDate: record[3],
			},
		})
	}

	return grammars, nil
}

func SaveGrammarToFile(grammars []model.Grammar) error {
	file, err := os.Create(grammarFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Grammar", "Meaning", "Seen", "Update_Date"})

	for _, g := range grammars {
		err := writer.Write([]string{g.Grammar, g.Meaning, strconv.FormatBool(g.Metadata.Seen), g.Metadata.UpdateDate})
		if err != nil {
			return err
		}
	}

	return nil
}
