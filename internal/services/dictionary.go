package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WordWithTranslation struct {
	EngWord       string
	RuWord        string
	EngDefinition string
	RuDefinition  string
}

type YandexDictResponse struct {
	Definition []struct {
		Word         string `json:"text"`
		PartOfSpeech string `json:"pos"`

		Tr []struct {
			Word         string `json:"text"`
			PartOfSpeech string `json:"pos"`

			Synonymous []struct {
				Word         string `json:"text"`
				PartOfSpeech string `json:"pos"`
			} `json:"syn"`

			Meaning []struct {
				Word         string `json:"text"`
				PartOfSpeech string `json:"pos"`
			} `json:"mean"`
		} `json:"tr"`
	} `json:"def"`
}

func GetWordWithTranslation(apiKey, word string) (*YandexDictResponse, error) {
	url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=en-ru&text=%s", apiKey, word)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.New("error while getting word: " + err.Error())
	}
	defer resp.Body.Close()

	var data YandexDictResponse
	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.New("error unmarshalling JSON: " + err.Error())
	}

	return &data, nil
}
