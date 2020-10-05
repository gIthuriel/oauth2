package externalaccount

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	externalaccount "golang.org/x/oauth2/google/internal/oauth"
	"net/http"
	"time"
)

// The configuration for fetching tokens with external credentials.
type Config struct {
	Audience                       string
	SubjectTokenType               string
	TokenURL                       string
	TokenInfoURL                   string
	ServiceAccountImpersonationURL string
	ClientSecret                   string
	ClientID                       string
	CredentialSource               CredentialSource
	QuotaProjectID                 string

	Scopes                         []string
}

// Returns an external account TokenSource. This is to be called by package google to construct a google.Credentials.
func (c *Config) TokenSource(ctx context.Context) oauth2.TokenSource {
	ts := tokenSource {
		ctx: ctx,
		conf: c,
	}
	return ts
}

type format struct {
	// When not provided "text" type is assumed.
	Type string `json:"type"`
	// Only required for JSON.
	// This would be "access_token" for azure.
	SubjectTokenFieldName string `json:"subject_token_field_name"`
}

type CredentialSource struct {
	File string `json:"file"`

	URL string `json:"url"`
	Headers map[string]string `json:"headers"`

	EnvironmentID string `json:"environment_id"`
	RegionURL string `json:"region_url"`
	RegionalCredVerificationURL string `json:"regional_cred_verification_url"`
	CredVerificationURL string `json:"cred_verification_url"`
	Format format `json:"format"`

}

func (cs CredentialSource) instance() baseCredentialSource {
	if cs.EnvironmentID == "awsX" {
		return nil
		//return awsCredentialSource{EnvironmentID:cs.EnvironmentID, RegionURL:cs.RegionURL, RegionalCredVerificationURL: cs.RegionalCredVerificationURL, CredVerificationURL:cs.CredVerificationURL}
	} else if cs.File != "" {
		return fileCredentialSource{File:cs.File}
	} else if cs.URL != "" {
		return nil
		//return urlCredentialSource{URL:cs.URL, Headers:cs.Headers}
	} else {
		return nil
	}
}

type baseCredentialSource interface {
	retrieveSubjectToken(c *Config) (string, error)
}


// The following struct and methods are worth reviewing here though they're not exposed.

// tokenSource is the source that handles 3PI credentials.
type tokenSource struct {
	ctx  context.Context
	conf *Config
}

// This method is implemented so that tokenSource conforms to oauth2.TokenSource.
func (ts tokenSource) Token() (*oauth2.Token, error) {
	conf := ts.conf

	subjectToken, err := conf.CredentialSource.instance().retrieveSubjectToken(conf)
	if err != nil {
		return &oauth2.Token{}, err
	}
	stsRequest := externalaccount.STSTokenExchangeRequest{
		GrantType: "urn:ietf:params:oauth:grant-type:token-exchange",
		Audience: conf.Audience,
		Scope: conf.Scopes,
		RequestedTokenType: "urn:ietf:params:oauth:token-type:access_token",
		SubjectToken: subjectToken,
		SubjectTokenType: conf.SubjectTokenType,
	}
	header := make(http.Header)
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	clientAuth := externalaccount.ClientAuthentication{
		AuthStyle: oauth2.AuthStyleInHeader,
		ClientID: conf.ClientID,
		ClientSecret: conf.ClientSecret,
	}
	stsResp, err := externalaccount.ExchangeToken(ts.ctx, conf.TokenURL, &stsRequest, clientAuth, header, nil)
	if err != nil {
		fmt.Errorf("oauth2/google: %s", err.Error())
	}

	accessToken := &oauth2.Token{
		AccessToken: stsResp.AccessToken,
		TokenType: stsResp.TokenType,
	}
	if stsResp.ExpiresIn != 0 {
		if err != nil {
			fmt.Errorf("google/oauth2: got invalid expiry from security token service")
		}
		accessToken.Expiry = time.Now().Add(time.Duration(stsResp.ExpiresIn)*time.Second) //TODO: Does this actually work properly!?!??
	}
	if stsResp.RefreshToken != "" {
		accessToken.RefreshToken = stsResp.RefreshToken
	}
	fmt.Printf("Token result: %+v\n\n",accessToken)
	return accessToken, nil
}

// NOTE: this method doesn't exist yet. It is being investigated to add this method to oauth2.TokenSource.
//func (ts tokenSource) TokenInfo() (*oauth2.TokenInfo, error)