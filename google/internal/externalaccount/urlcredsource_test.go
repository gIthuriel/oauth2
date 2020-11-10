package externalaccount

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var myURLToken = "testTokenValue"


/* I started trying to implement table-driven testing here but I don't think it's right.  The server needs to have different performance depending on which test is being run.
*/

/* type urlsTest struct {
	name string
	cs CredentialSource
	want string
}

var testURLsUntyped = urlsTest{
	name: "UntypedURLSource",
	cs: CredentialSource{},
	want: "testTokenValue",
}
var testURLsTypeText = urlsTest{
	name: "TextURLSource",
	cs: CredentialSource{},
	want: "testTokenValue",
}
var testURLsTypeJSON = urlsTest{
	name: "JSON_URLSource",
	cs: CredentialSource{
		Format: format{Type: fileTypeJSON, SubjectTokenFieldName: "SubjToken"},
	},
	want: "testTokenValue",
}

var urlSourceTests = []urlsTest{testURLsUntyped, testURLsTypeText, testURLsTypeJSON} */

func TestRetrieveURLSubjectToken_Text(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Unexpected request method, %v is found", r.Method)
		}
		w.Write([]byte("testTokenValue"))
	}))
	cs := CredentialSource{
		URL: ts.URL,
		Format: format{Type: fileTypeText},
	}
	tfc := testFileConfig
	tfc.CredentialSource = cs

	out, err := cs.instance().retrieveSubjectToken(&tfc)
	if err != nil {
		t.Fatalf("Failed to retrieve URL subject token: %v", err)
	}
	if out != myURLToken {
		t.Errorf("Recieved wrong subject token from URL; got %v but want %v", out, myURLToken)
	}

}

// Checking that retrieveSubjectToken properly defaults to type text
func TestRetrieveURLSubjectToken_Untyped(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Unexpected request method, %v is found", r.Method)
		}
		w.Write([]byte("testTokenValue"))
	}))
	cs := CredentialSource{
		URL: ts.URL,
	}
	tfc := testFileConfig
	tfc.CredentialSource = cs

	out, err := cs.instance().retrieveSubjectToken(&tfc)
	if err != nil {
		t.Fatalf("Failed to retrieve USRL subject token: %v", err)
	}
	if out != myURLToken {
		t.Errorf("Recieved wrong subject token from URL; got %v but want %v", out, myURLToken)
	}

}

func TestRetrieveURLSubjectToken_JSON(t *testing.T) {
	type tokenResponse struct {
		TestToken string `json:"SubjToken"`
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Unexpected request method, %v is found", r.Method)
		}
		resp := tokenResponse{TestToken: "testTokenValue"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Failed to marshal values: %v", err)
		}
		w.Write(jsonResp)
	}))
	cs := CredentialSource{
		URL: ts.URL,
		Format: format{Type: fileTypeJSON, SubjectTokenFieldName: "SubjToken"},
	}
	tfc := testFileConfig
	tfc.CredentialSource = cs

	out, err := cs.instance().retrieveSubjectToken(&tfc)


	if err != nil {
		t.Fatalf("Failed to retrieve URL subject token: %v", err)
	}
	if out != myURLToken {
		t.Errorf("Recieved wrong subject token from URL; got %v but want %v", out, myURLToken)
	}

}