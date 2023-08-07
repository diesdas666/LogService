package app

type DatabaseConfig struct {
	Password string
	Host     string
	Port     int
	Name     string
	User     string
}

type CacheConfig struct {
	Type  string
	Redis RedisConfig
}

type RedisConfig struct {
	Addr string
	DB   int
}

type Config struct {
	Deployment  string
	Credentials CredentialsConfig
	Server      ServerConfig
	Database    DatabaseConfig
	Cache       CacheConfig
}

type CredentialsConfig struct {
	Key    string
	Secret string
}

type ServerConfig struct {
	Port int
	Addr string
}
