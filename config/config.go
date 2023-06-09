package config

type Config struct {
	PE struct {
		Console  string
		Token    string
		CACert   string
		HostCert string
		PrivKey  string
	}
}
