package externalaccount

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCredentialSource struct {

}

var testBaseCredSource = CredentialSource {
	File: "internalTestingFile",
}

var testConfig = Config{
	Audience: "32555940559.apps.googleusercontent.com",
	SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
	TokenURL: "http://localhost:8080/v1/token",
	TokenInfoURL: "http://localhost:8080/v1/tokeninfo",
	ServiceAccountImpersonationURL: "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/service-gcs-admin@$PROJECT_ID.iam.gserviceaccount.com:generateAccessToken",
	ClientSecret: "notsosecret",
	ClientID: "rbrgnognrhongo3bi4gb9ghg9g",
	CredentialSource: testBaseCredSource,
}

func TestTokenSource(t *testing.T) {
	want := tokenSource{
		ctx: context.Background(),
		conf: &testConfig,
	}
	got := testConfig.TokenSource(context.Background())
	if want != got {
		t.Errorf("Unexpected TokenSource; expected %+v but got %+v", want, got)
	}
}


func TestToken_Func(t *testing.T) {

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
}