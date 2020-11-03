// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package externalaccount

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ExchangeToken(ctx context.Context, endpoint string, request *STSTokenExchangeRequest, authentication ClientAuthentication, headers http.Header, options map[string]interface{}) (*STSTokenExchangeResponse, error) {

	client := oauth2.NewClient(ctx, nil)

	data := url.Values{}
	data.Set("audience", request.Audience)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	data.Set("requested_token_type", "urn:ietf:params:oauth:token-type:access_token")
	data.Set("subject_token_type", request.SubjectTokenType)
	data.Set("subject_token", request.SubjectToken)
	data.Set("scope", strings.Join(request.Scope, " "))
	opts, err := json.Marshal(options)
	if err == nil {
		data.Set("options", string(opts))
	} else {
		fmt.Errorf("oauth2/google: failed to marshal additional options: %v", err)
	}

	authentication.InjectAuthentication(data, headers)
	encodedData := data.Encode()
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(encodedData))
	if err != nil {
		fmt.Errorf("oauth2/google: failed to properly build http request: %v", err)
	}
	for key, _ := range headers {
		for _, val := range headers.Values(key) {
			req.Header.Add(key, val)
		}
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("oauth2/google: invalid response from Secure Token Server: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		fmt.Errorf("oauth2/google: invalid body in Secure Token Server response: %v", err)
	}
	var stsResp STSTokenExchangeResponse
	err = json.Unmarshal(body, &stsResp)
	if err != nil {
		fmt.Errorf("oauth2/google: failed to unmarshal response body from STS server: %v", err)
	}
	return &stsResp, nil
}

type STSTokenExchangeRequest struct {
	ActingParty struct {
		ActorToken     string
		ActorTokenType string
	}
	GrantType          string
	Resource           string
	Audience           string
	Scope              []string
	RequestedTokenType string
	SubjectToken       string
	SubjectTokenType   string
}

type STSTokenExchangeResponse struct {
	AccessToken     string `json:"access_token"`
	IssuedTokenType string `json:"issued_token_type"`
	TokenType       string `json:"token_type"`
	ExpiresIn       int    `json:"expires_in"`
	Scope           string `json:"scope"`
	RefreshToken    string `json:"refresh_token"`
}
