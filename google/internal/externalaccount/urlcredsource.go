// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	URL     string
	Headers map[string]string
	Format  format
}

func (cs urlCredentialSource) subjectToken() (string, error) {
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

	switch cs.Format.Type {
	case "json":
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(tokenBytes, &jsonData)
		if err != nil {
			return "", fmt.Errorf("oauth2/google: failed to unmarshal subject token file: %v", err)
		}
		val, ok := jsonData[cs.Format.SubjectTokenFieldName]
		if !ok {
			return "", errors.New("oauth2/google: provided subject_token_field_name not found in credentials")
		}
		token, ok := val.(string)
		if !ok {
			return "", errors.New("oauth2/google: improperly formatted subject token")
		}
		return token, nil
	case "text":
		return string(tokenBytes), nil
	case "":
		return string(tokenBytes), nil
	default:
		return "", errors.New("oauth2/google: invalid credential_source file format type")
	}

}
