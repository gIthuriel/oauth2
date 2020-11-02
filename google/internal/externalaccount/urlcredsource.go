package externalaccount

type urlCredentialSource struct {
	URL string
	Headers map[string]interface{}
}

func (cs urlCredentialSource) retrieveSubjectToken(c *Config) string {

	return ""
}