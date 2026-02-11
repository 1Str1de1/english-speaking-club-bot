package services

func ExecuteTranslateText(apiKey, text string) (string, error) {
	data, err := FetchWordWithTranslation(apiKey, text)
	if err != nil {
		return "", err
	}
	return FormatWordForTelegram(data), nil
}
