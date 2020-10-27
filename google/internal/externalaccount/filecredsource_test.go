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
	//CredentialSource: fileSource,
}
type fsTest struct {
	cs CredentialSource
	want string
}
var testFSUntyped = fsTest {
	cs: CredentialSource{
		File: "../../testdata/externalaccount/3pi_cred.txt",
	},
	want: "street123",
}
var testFSTypeText = fsTest {
	cs: CredentialSource{
		File: "../../testdata/externalaccount/3pi_cred.txt",
		Format: format{Type: fileTypeText},
	},
	want: "street123",
}
var testFSTypeJSON = fsTest {
	cs: CredentialSource{
		File: "../../testdata/externalaccount/3pi_cred.json",
		Format: format{Type: fileTypeJSON, SubjectTokenFieldName: "SubjToken"},
	},
	want: "321road",
}
var fileSourceTests = []fsTest{testFSUntyped, testFSTypeText, testFSTypeJSON}


func TestRetrieveFileSubjectToken_Untyped(t *testing.T) {
	for _, test := range fileSourceTests {
		testFileConfig.CredentialSource = test.cs

		out, err := test.cs.instance().retrieveSubjectToken(&testFileConfig)
		if err != nil {
			t.Errorf("Method retrieveSubjectToken for type fileCredentialSource failed; %e", err)
		}
		if out != test.want {
			t.Errorf("Method retrieveSubjectToken for type fileCredentialSouce failed: expected %v but got %v", "street123", out)
		}
	}
}