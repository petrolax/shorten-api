package AbbreviationStorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/teris-io/shortid"
)

type AbbreviationStorage struct {
	file *os.File
	list map[string]string
	mu   sync.Mutex
}

func NewAbbreviationStorage(file *os.File) *AbbreviationStorage {
	var mutex sync.Mutex
	list := make(map[string]string)
	data, _ := ioutil.ReadFile(file.Name())
	if len(data) > 1 {
		_ = json.Unmarshal(data, &list)
	}
	return &AbbreviationStorage{
		file: file,
		list: list,
		mu:   mutex,
	}
}

func (as *AbbreviationStorage) CreateNewShortenUrl(url string) (string, error) {
	as.mu.Lock()
	defer as.mu.Unlock()
	shorturl := ""
	sid, _ := shortid.New(1, shortid.DefaultABC, 2444)
	for {
		shorturl = sid.MustGenerate()
		if _, ok := as.list[shorturl]; !ok {
			as.list[shorturl] = url
			break
		}
	}

	data, _ := json.Marshal(as.list)
	_ = ioutil.WriteFile(as.file.Name(), data, 0644)
	return shorturl, nil
}

func (as *AbbreviationStorage) GetLengthenUrl(shorturl string) (string, error) {
	as.mu.Lock()
	defer as.mu.Unlock()
	url, ok := as.list[shorturl]
	if !ok {
		return "", errors.New(fmt.Sprintf("Abbreviation %s doesn't exist", shorturl))
	}
	return url, nil
}

func (as *AbbreviationStorage) GetListOfAbbreviations(page int) (map[string]string, error) {
	as.mu.Lock()
	defer as.mu.Unlock()
	countOfAbbreviationEnd := 5 * page
	if countOfAbbreviationEnd > len(as.list) {
		countOfAbbreviationEnd = len(as.list)
	}
	countOfAbbreviationBegin := countOfAbbreviationEnd - 5

	i := 0
	res := make(map[string]string)
	for shorturl, url := range as.list {
		if i < countOfAbbreviationBegin {
			i++
			continue
		}
		if i > countOfAbbreviationEnd {
			break
		}
		res[shorturl] = url
		i++
	}
	return res, nil
}

func (as *AbbreviationStorage) DeleteAllShortenUrl() error {
	as.mu.Lock()
	defer as.mu.Unlock()
	for key := range as.list {
		delete(as.list, key)
	}
	data, _ := json.Marshal(as.list)
	_ = ioutil.WriteFile(as.file.Name(), data, 0644)
	return nil
}

func (as *AbbreviationStorage) DeleteShortenUrl(shorturl string) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	if _, ok := as.list[shorturl]; !ok {
		return errors.New(fmt.Sprintf("Abbreviation %s doesn't exist", shorturl))
	}
	delete(as.list, shorturl)
	data, _ := json.Marshal(as.list)
	_ = ioutil.WriteFile(as.file.Name(), data, 0644)
	return nil
}
