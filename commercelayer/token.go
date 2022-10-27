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
	cacheFile        *os.File
}

func NewCachedTokenSource(tokenSource oauth2.TokenSource, cacheFile *os.File) *CachedTokenSource {
	var err error
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
		cacheFile:        cacheFile,
	}
}

func (c *CachedTokenSource) Token() (*oauth2.Token, error) {
	token, err := c.innerTokenSource.Token()
	if err != nil {
		return nil, err
	}

	byteValue, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	_, err = c.cacheFile.Write(byteValue)
	if err != nil {
		return nil, err
	}

	return token, nil
}
