package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewShortenAndOriginal(t *testing.T) {
	url := "https://yandex.ru"
	requrl := "http://localhost:8080?url=" + url
	resp, err := http.Post(requrl, "application/json", nil)
	if err != nil {
		t.Error("Bad Request", err)
		return
	}
	defer resp.Body.Close()

	shortenres := make(map[string]interface{})
	_ = json.NewDecoder(resp.Body).Decode(&shortenres)

	res := shortenres["Result"].(string)
	requrl = "http://localhost:8080/" + res + "/original"
	resp, err = http.Get(requrl)
	if err != nil {
		t.Error("Bad Request")
		return
	}
	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(&shortenres)

	res = shortenres["Result"].(string)
	// res += "ru" //Для проверки ошибки
	if url == res {
		t.Logf("Url: %s is equal from short url: %s", url, res)
	} else {
		t.Errorf("Url: %s is not equal from short url: %s", url, res)
	}
}

func TestNewShortenAndDeleteShorten(t *testing.T) {
	url := "https://habr.com"
	requrl := "http://localhost:8080?url=" + url
	resp, err := http.Post(requrl, "application/json", nil)
	if err != nil {
		t.Error("Bad Request", err)
		return
	}
	defer resp.Body.Close()

	shortenres := make(map[string]interface{})
	_ = json.NewDecoder(resp.Body).Decode(&shortenres)

	res := shortenres["Result"].(string)
	expectedRes := "Short url " + res + " was delete"
	requrl = "http://localhost:8080/" + res + "/delete"
	req, err := http.NewRequest(http.MethodDelete, requrl, nil)
	if err != nil {
		t.Error("Bad Request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Error("Bad Request")
		return
	}
	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(&shortenres)
	res = shortenres["Message"].(string)

	if expectedRes == res {
		t.Log(res)
	} else {
		t.Errorf("Result messages is not equal")
	}
}

func TestDeleteAllAndGetList(t *testing.T) {
	requrl := "http://localhost:8080/delete"
	req, err := http.NewRequest(http.MethodDelete, requrl, nil)
	if err != nil {
		t.Error("Bad Request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Bad Request")
		return
	}
	defer resp.Body.Close()

	shortenres := make(map[string]interface{})
	_ = json.NewDecoder(resp.Body).Decode(&shortenres)
	res := shortenres["Message"].(string)

	expectedRes := "Every abbreviations was delete"
	if res == expectedRes {
		t.Log(res)
	} else {
		t.Error("Some problems with delete")
		return
	}

	url := "https://habr.com"
	requrl = "http://localhost:8080?url=" + url
	resp, err = http.Post(requrl, "application/json", nil)
	if err != nil {
		t.Error("Bad Request", err)
		return
	}
	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(&shortenres)
	shortres := shortenres["Result"].(string)

	requrl = "http://localhost:8080/list/1"
	resp, err = http.Get(requrl)
	if err != nil {
		t.Error("Bad Request", err)
		return
	}
	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(&shortenres)
	listres := shortenres["Result"].(map[string]interface{})
	if val, ok := listres[shortres]; ok && val == url {
		t.Log("Abbreviation List is true")
	} else {
		t.Error("Problem with abbreviation list")
	}

}

func TestNewShortenAndRedirect(t *testing.T) {
	urls := []string{"https://yandex.ru", "https://www.google.com/", "https://google.com", "https://mail.ru", "twitter.com"}
	results := []bool{true, true, false, true, false} //Если изменить 3-й элемент в срезе на true, то в цикле сработает else
	for i, url := range urls {
		if forTestingNewShortenAndRedirect(url, t) == results[i] {
			t.Log("Result:", results[i], "for url:", url)
		} else {
			t.Error("Result:", false, "when expected", results[i], "for url:", url)
		}
	}
}

func forTestingNewShortenAndRedirect(url string, t *testing.T) bool {
	requrl := "http://localhost:8080?url=" + url
	resp, err := http.Post(requrl, "application/json", nil)
	if err != nil {
		t.Log("Bad Request", err)
		return false
	}
	defer resp.Body.Close()

	shortenres := make(map[string]interface{})
	_ = json.NewDecoder(resp.Body).Decode(&shortenres)

	if resp.StatusCode != http.StatusOK {
		err := shortenres["Error"].(string)
		t.Log(err)
		return false
	}

	res := shortenres["Result"].(string)
	requrl = "http://localhost:8080/" + res
	resp, err = http.Get(requrl)
	if err != nil {
		t.Log("Bad Request")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusOK {
		shortenres := make(map[string]interface{})
		_ = json.NewDecoder(resp.Body).Decode(&shortenres)
		err := shortenres["Error"].(string)
		t.Log(err)
		return false
	}

	if resp.Request.URL.String() == url {
		return true
	} else {
		return false
	}
}
