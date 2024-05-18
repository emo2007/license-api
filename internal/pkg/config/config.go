package config

type Config struct {
	Common CommonConfig
	Rest   RestConfig
	DB     DBConfig
	Eth    EthConfig
}

type CommonConfig struct {
	LogLevel     string
	LogLocal     bool
	LogFile      string
	LogAddSource bool

	JWTSecret []byte
}

type RestConfig struct {
	Address string
	TLS     bool
}

type DBConfig struct {
	// persistent database config
	Host      string
	EnableSSL bool
	Database  string
	User      string
	Secret    string

	// cache config
	CacheHost   string
	CacheUser   string
	CacheSecret string
}

type EthConfig struct {
	// todo
}
