package inject

type Environment struct {
	Env                 string `split_words:"true"`
	Port                int    `split_words:"true"`
	SiteOrigin          string `required:"true" split_words:"true"`
	GoogleTagManagerID  string `envconfig:"GOOGLE_TAG_MANAGER_ID"`
	AdminToken          string `required:"true" split_words:"true"`
	DirPathHTMLTemplate string `required:"true" split_words:"true"`
	DirPathCSS          string `required:"true" split_words:"true"`
	LogLevel            string `split_words:"true"`
	LoggerType          string `split_words:"true"`
	DBHost              string `envconfig:"DB_HOST" required:"true"`
	DBPort              uint16 `envconfig:"DB_PORT" required:"true"`
	DBName              string `envconfig:"DB_NAME" required:"true"`
	DBPassword          string `envconfig:"DB_PASSWORD" required:"true"`
	DBUser              string `envconfig:"DB_USER" required:"true"`
}
