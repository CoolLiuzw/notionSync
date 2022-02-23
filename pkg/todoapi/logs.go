package todoapi

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SaveStatus(resp *http.Response) error {
	bodyBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	currentTime := time.Now()
	fileName := "requestlog" + currentTime.String() + ".txt"

	file, fileErr := os.Create(filepath.Join("logs", fileName))
	if fileErr != nil {
		return fileErr
	}

	_, writeErr := file.Write(bodyBytes[:len(bodyBytes)])
	if writeErr != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return closeErr
		}

		return writeErr
	}

	closeErr := file.Close()
	if closeErr != nil {
		return closeErr
	}

	return nil
}
