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
	keys []string
	mu   sync.Mutex
}

func NewAbbreviationStorage(file *os.File) *AbbreviationStorage {
	list := make(map[string]string)
	data, _ := ioutil.ReadFile(file.Name())
	_ = json.Unmarshal(data, &list)
	keys := make([]string, 0)
	for k := range list {
		keys = append(keys, k)
	}
	return &AbbreviationStorage{
		file: file,
		list: list,
		keys: keys,
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
			as.keys = append(as.keys, shorturl)
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
	countOfAbbreviationBegin := 0
	countOfAbbreviationEnd := 5 * page
	if countOfAbbreviationEnd > len(as.list) {
		countOfAbbreviationEnd = len(as.list)
	}

	if countOfAbbreviationEnd%5 == 0 {
		countOfAbbreviationBegin = countOfAbbreviationEnd - 5
	} else {
		countOfAbbreviationBegin = countOfAbbreviationEnd - (countOfAbbreviationEnd % 5)
	}

	res := make(map[string]string)
	for i := countOfAbbreviationBegin; i < countOfAbbreviationEnd; i++ {
		key := as.keys[i]
		res[key] = as.list[key]
	}
	return res, nil
}

func (as *AbbreviationStorage) DeleteAllShortenUrl() error {
	as.mu.Lock()
	defer as.mu.Unlock()
	for _, key := range as.keys {
		delete(as.list, key)
	}
	as.keys = nil
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
	for i, key := range as.keys {
		if key == shorturl {
			as.keys = append(as.keys[:i], as.keys[i+1:]...)
			break
		}
	}
	data, _ := json.Marshal(as.list)
	_ = ioutil.WriteFile(as.file.Name(), data, 0644)
	return nil
}
