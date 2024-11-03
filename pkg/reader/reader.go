package reader

import (
	"fmt"
	"net/url"
	"os"
)

func ReadFile(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		fmt.Printf("Ошибка при парсинге URI: %v\n", err)
		os.Exit(1)
	}

	filePath := u.Path

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reader.Read(uri: %s): %v\n", uri, err)
	}

	return string(data), nil
}
