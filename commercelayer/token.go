package commercelayer

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
)

type CachedTokenSource struct {
	innerTokenSource oauth2.TokenSource
	fileLocation     string
}

func NewCachedTokenSource(tokenSource oauth2.TokenSource, fileLocation string) *CachedTokenSource {
	cacheFile, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer cacheFile.Close()

	byteValue, _ := ioutil.ReadAll(cacheFile)

	var token oauth2.Token
	if len(byteValue) != 0 {
		err = json.Unmarshal(byteValue, &token)
		if err != nil {
			log.Fatal(err)
		}
	}

	innerTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	return &CachedTokenSource{
		innerTokenSource: innerTokenSource,
		fileLocation:     fileLocation,
	}
}

func (c *CachedTokenSource) Token() (*oauth2.Token, error) {
	token, err := c.innerTokenSource.Token()
	if err != nil {
		return nil, err
	}

	cacheFile, err := os.OpenFile(c.fileLocation, os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer cacheFile.Close()

	byteValue, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	_, err = cacheFile.Write(byteValue)
	if err != nil {
		return nil, err
	}

	return token, nil
}
