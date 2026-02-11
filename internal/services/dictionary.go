package services

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Vocabulary struct {
	words    []string
	byLength map[int][]string
}

type YandexDictResponse struct {
	Definition []struct {
		Word          string `json:"text"`
		PartOfSpeech  string `json:"pos"`
		Transcription string `json:"ts"`

		Tr []struct {
			Word         string `json:"text"`
			PartOfSpeech string `json:"pos"`

			Synonymous []struct {
				Word string `json:"text"`
			} `json:"syn"`

			Meaning []struct {
				Word string `json:"text"`
			} `json:"mean"`

			//Example []struct {
			//	Word string `json:"text"`
			//} `json:"ex"`
		} `json:"tr"`
	} `json:"def"`
}

func ExecuteRandomWordCommand(apiKey string) (string, error) {
	voc, err := LoadVocabulary("common_words.txt")
	if err != nil {
		return "", err
	}

	w := GetRandomWordFromVocabulary(voc)
	data, err := FetchWordWithTranslation(apiKey, w)
	if err != nil {
		return "", err
	}

	return FormatWordForTelegram(data), nil
}

// FetchWordWithTranslation A func for getting an exact word from YandexTranslate API.
// It returns YandexDictResponse struct which next will be converted to Telegram message in another func
func FetchWordWithTranslation(apiKey, word string) (*YandexDictResponse, error) {
	w := url2.QueryEscape(word)
	url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=en-ru&text=%s", apiKey, w)

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

func LoadVocabulary(fileName string) (*Vocabulary, error) {
	projDir := os.Getenv("PROJECT_DIR")
	//if err != nil {
	//	return nil, errors.New("error getting working dir")
	//}

	filePath := filepath.Join(projDir, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("failed to open file with words" + err.Error())
	}

	words := make([]string, 10000)
	byLength := make(map[int][]string, 10000)

	vocabulary := Vocabulary{
		words:    words,
		byLength: byLength,
	}

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); {
		wordLen := len(scanner.Text())
		vocabulary.words[i] = scanner.Text()
		byLength[wordLen] = append(byLength[wordLen], words[i])
		i++
	}

	return &vocabulary, nil
}

func GetRandomWordFromVocabulary(voc *Vocabulary) string {
	return voc.words[rand.Intn(len(voc.words))]
}

func FormatWordForTelegram(data *YandexDictResponse) string {
	if len(data.Definition) == 0 {
		return ""
	}

	var builder strings.Builder

	word := data.Definition[0]
	builder.WriteString(fmt.Sprintf("📚 %s", word.Word))

	if word.PartOfSpeech != "" {
		builder.WriteString(fmt.Sprintf("\n🔤 %s", translatePartOfSpeech(word.PartOfSpeech)))
	}

	if word.Transcription != "" {
		builder.WriteString(fmt.Sprintf("\n👄 [%s]", word.Transcription))
	}

	builder.WriteString("\n\n")

	for i, tr := range word.Tr {
		builder.WriteString(fmt.Sprintf("🇷🇺 *%d. %s*", i+1, tr.Word))

		if tr.PartOfSpeech != "" && tr.PartOfSpeech != word.PartOfSpeech {
			builder.WriteString(fmt.Sprintf(" (%s)", translatePartOfSpeech(tr.PartOfSpeech)))
		}

		if len(tr.Meaning) > 0 {
			builder.WriteString("\n   📖 _")
			for j, mean := range tr.Meaning {
				if j > 0 {
					builder.WriteString("; ")
				}
				builder.WriteString(mean.Word)

				// limit set to three values
				if j == 2 {
					builder.WriteString("...")
					break
				}
			}
		}

		if len(tr.Synonymous) > 0 {
			builder.WriteString("\n   🔄 ")
			for j, syn := range tr.Synonymous {
				if j > 0 {
					builder.WriteString("; ")
				}
				builder.WriteString(syn.Word)

				// limit set to five values
				if j == 4 {
					builder.WriteString("...")
					break
				}
			}
		}

		//if len(tr.Example) > 0 {
		//
		//}

		builder.WriteString("\n\n")

		// limit to 3 translations
		if i >= 3 {
			builder.WriteString("... и другие переводы\n")
			break
		}
	}

	return builder.String()
}

func translatePartOfSpeech(pos string) string {
	switch pos {
	case "noun":
		return "существительное"
	case "verb":
		return "глагол"
	case "adjective":
		return "прилагательное"
	case "adverb":
		return "наречие"
	case "pronoun":
		return "местоимение"
	case "preposition":
		return "предлог"
	case "conjunction":
		return "союз"
	case "interjection":
		return "междометие"
	case "numeral":
		return "числительное"
	case "participle":
		return "причастие"
	case "gerund":
		return "герундий"
	default:
		return pos
	}
}
