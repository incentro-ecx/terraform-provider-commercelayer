package commercelayer

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
)

type cachedReusableTokenSource struct {
	innerTokenSource oauth2.TokenSource
	fileLocation     string
}

func newCachedTokenSource(tokenSource oauth2.TokenSource, fileLocation string) *cachedReusableTokenSource {
	cacheFile, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer cacheFile.Close()

	byteValue, _ := ioutil.ReadAll(cacheFile)

	var token oauth2.Token
	if len(byteValue) != 0 {
		err = json.Unmarshal(byteValue, &token)
		if err != nil {
			panic(err)
		}
	}

	innerTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	return &cachedReusableTokenSource{
		innerTokenSource: innerTokenSource,
		fileLocation:     fileLocation,
	}
}

func (c *cachedReusableTokenSource) Token() (*oauth2.Token, error) {
	token, err := c.innerTokenSource.Token()
	if err != nil {
		return token, err
	}

	cacheFile, err := os.OpenFile(c.fileLocation, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer cacheFile.Close()

	byteValue, err := json.Marshal(token)
	if err != nil {
		return token, err
	}

	_, err = cacheFile.Write(byteValue)
	if err != nil {
		return token, err
	}

	return token, nil
}
