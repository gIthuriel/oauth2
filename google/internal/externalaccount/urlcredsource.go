package externalaccount

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type urlCredentialSource struct {
	URL string
	Headers map[string]string
}

func (cs urlCredentialSource) retrieveSubjectToken(c *Config) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", cs.URL, strings.NewReader(""))

	for key, val := range cs.Headers {
		req.Header.Add(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("oauth2/google: invalid response when retrieving subject token: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	tokenBytes, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		fmt.Errorf("oauth2/google: invalid body in subject token URL query: %v", err)
		return "", err
	}

	var output string
	switch c.CredentialSource.Format.Type {
	case "json":
		jsonData := make(map[string]interface{})
		json.Unmarshal(tokenBytes, &jsonData)
		if val, ok := jsonData[c.CredentialSource.Format.SubjectTokenFieldName]; !ok {
			return "", errors.New("oauth2/google: provided subject_token_field_name not found in credentials")
		} else {
			if token, ok := val.(string); !ok {
				return "", errors.New("oauth2/google: improperly formatted subject token")
			} else {
				output = token
			}

		}
	case "text":
		output = string(tokenBytes)
	case "":
		output = string(tokenBytes)
	default:
		return "", errors.New("oauth2/google: invalid credential_source file format type")
	}


	return output, nil
}