package config

import "time"

type DefaultConfig struct {
	Apps       Apps       `mapstructure:"apps"`
	Server     Server     `mapstructure:"server"`
	Database   Datasource `mapstructure:"database"`
	Mail       Mail       `mapstructure:"mail"`
	Service    Service    `mapstructure:"service"`
	AuthConfig AuthConfig `mapstructure:"authConfig"`
}

type Apps struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"enviroment"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Datasource struct {
	Url               string        `mapstructure:"url"`
	Port              string        `mapstructure:"port"`
	DatabaseName      string        `mapstructure:"databaseName"`
	Username          string        `mapstructure:"username"`
	Password          string        `mapstructure:"password"`
	Schema            string        `mapstructure:"schema"`
	ConnectionTimeout time.Duration `mapstructure:"connectionTimeout"`
	MaxIdleConnection int           `mapstructure:"maxIdleConnection"`
	MaxOpenConnection int           `mapstructure:"maxOpenConnection"`
	DebugMode         bool          `mapstructure:"debugMode"`
}

type Service struct {
	Gcp GcpService `mapstructure:"gcp"`
}

type GcpService struct {
	CredentialPath string `mapstructure:"credentialPath"`
	ProjectID      string `mapstructure:"projectId"`
	BucketName     string `mapstructure:"bucketName"`
}

type Mail struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Sender     string `mapstructure:"sender"`
	SenderName string `mapstructure:"senderName"`
	Password   string `mapstructure:"password"`
	SmtpAuth   string `mapstructure:"smtpAuth"`
	StartTLS   string `mapstructure:"startTLS"`
}

type AuthConfig struct {
	SessionMin         time.Duration `mapstructure:"sessionMin"`
	SessionRememberMe  string        `mapstructure:"sessionRememberMe"`
	MaximumLoginFailed int64         `mapstructure:"maximumLoginFailed"`
}
