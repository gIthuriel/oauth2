package externalaccount

import "testing"

var testFileConfig = Config{
	Audience: "32555940559.apps.googleusercontent.com",
	SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
	TokenURL: "http://localhost:8080/v1/token",
	TokenInfoURL: "http://localhost:8080/v1/tokeninfo",
	ServiceAccountImpersonationURL: "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/service-gcs-admin@$PROJECT_ID.iam.gserviceaccount.com:generateAccessToken",
	ClientSecret: "notsosecret",
	ClientID: "rbrgnognrhongo3bi4gb9ghg9g",
	CredentialSource: fileSource,
}
var fileSource = CredentialSource{
	File: "../../testdata/externalaccount/file_credentials.json",
}

func TestRetrieveFileSubjectToken(t *testing.T) {

	out, err := fileSource.instance().retrieveSubjectToken(&testFileConfig)
	if err != nil {
		t.Errorf("Method retrieveSubjectToken for type fileCredentialSource failed; %e", err)
	}


}

