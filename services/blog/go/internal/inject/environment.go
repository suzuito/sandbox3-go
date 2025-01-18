package inject

type Environment struct {
	Env                 string `split_words:"true"`
	SiteOrigin          string `required:"true" split_words:"true"`
	GoogleTagManagerID  string `envconfig:"GOOGLE_TAG_MANAGER_ID"`
	AdminToken          string `required:"true" split_words:"true"`
	DirPathHTMLTemplate string `required:"true" split_words:"true"`
	DirPathCSS          string `required:"true" split_words:"true"`
	LogLevel            string `split_words:"true"`
	LoggerType          string `split_words:"true"`
}
