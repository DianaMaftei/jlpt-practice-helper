package deprecatedemailtemplateflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getKanjiAliveData(kanji string) (*KanjiDetail, error) {
	rapidApiKey := os.Getenv("RAPID_API_KEY")

	url := fmt.Sprintf("https://kanjialive-api.p.rapidapi.com/api/public/kanji/%s", kanji)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", rapidApiKey)
	req.Header.Add("X-RapidAPI-Host", "kanjialive-api.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var KanjiDetail KanjiDetail
	err = json.Unmarshal(body, &KanjiDetail)
	if err != nil {
		return nil, err
	}

	return &KanjiDetail, nil
}
