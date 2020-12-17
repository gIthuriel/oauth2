// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package externalaccount

import (
	"testing"
)

var testFileConfig = Config{
	Audience:                       "32555940559.apps.googleusercontent.com",
	SubjectTokenType:               "urn:ietf:params:oauth:token-type:jwt",
	TokenURL:                       "http://localhost:8080/v1/token",
	TokenInfoURL:                   "http://localhost:8080/v1/tokeninfo",
	ServiceAccountImpersonationURL: "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/service-gcs-admin@$PROJECT_ID.iam.gserviceaccount.com:generateAccessToken",
	ClientSecret:                   "notsosecret",
	ClientID:                       "rbrgnognrhongo3bi4gb9ghg9g",
}

func TestRetrieveFileSubjectToken(t *testing.T) {
	var fileSourceTests = []struct {
		name string
		cs   CredentialSource
		want string
	}{
		{
			name: "UntypedFileSource",
			cs: CredentialSource{
				File: "./testdata/3pi_cred.txt",
			},
			want: "street123",
		},
		{
			name: "TextFileSource",
			cs: CredentialSource{
				File:   "./testdata/3pi_cred.txt",
				Format: format{Type: fileTypeText},
			},
			want: "street123",
		},
		{
			name: "JSONFileSource",
			cs: CredentialSource{
				File:   "./testdata/3pi_cred.json",
				Format: format{Type: fileTypeJSON, SubjectTokenFieldName: "SubjToken"},
			},
			want: "321road",
		},
	}

	for _, test := range fileSourceTests {
		tfc := testFileConfig
		tfc.CredentialSource = test.cs

		out, err := tfc.parse().retrieveSubjectToken(&tfc)
		if err != nil {
			t.Errorf("Method retrieveSubjectToken for type fileCredentialSource in test %v failed; %e", test.name, err)
		}
		if out != test.want {
			t.Errorf("Test %v for method retrieveSubjectToken for type fileCredentialSouce failed: expected %v but got %v", test.name, test.want, out)
		}
	}
}
