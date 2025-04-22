package connection

import "log"

type Settings struct {
	DriverName string
	Ip         string
	Port       int
	DbName     string
	Username   string
	Password   string
	Ssl        bool
	SslMode    string
}

func NewSettings(options ...Option) Settings {
	settings := Settings{
		DriverName: "",
		Ip:         "localhost",
		Port:       0,
		DbName:     "",
		Username:   "",
		Password:   "",
		Ssl:        false,
		SslMode:    "",
	}

	for _, option := range options {
		option(&settings)
	}

	if settings.DriverName == "" {
		log.Panic("driver not set")
	}

	return settings
}

type Option func(*Settings)

func WithDriver(name string) Option {
	return func(settings *Settings) {
		settings.DriverName = name
	}
}

func WithIp(ip string) Option {
	return func(settings *Settings) {
		settings.Ip = ip
	}
}

func WithPort(port int) Option {
	return func(settings *Settings) {
		settings.Port = port
	}
}

func WithDb(name string) Option {
	return func(settings *Settings) {
		settings.DbName = name
	}
}

func WithUser(user string) Option {
	return func(settings *Settings) {
		settings.Username = user
	}
}

func WithPass(pass string) Option {
	return func(settings *Settings) {
		settings.Password = pass
	}
}

func WithSslMode(mode string) Option {
	return func(settings *Settings) {
		settings.SslMode = mode
	}
}
