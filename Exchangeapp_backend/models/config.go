package models

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
	Redis struct {
		Addr     string
		DB       int
		Password string
	}
	TgBot struct {
		Token  string
		Hour   int
		Min    int
		ChatID int64
	}
}
