package domain

type BotConfig struct {
	TelebotToken  string
	DatabaseToken string
	DatabaseURL   string
}

func (c *BotConfig) Validate() error {
	if c.TelebotToken == "" {
		return ErrMissingTeletoken
	}
	if c.DatabaseToken == "" {
		return ErrMissingDBToken
	}
	if c.DatabaseURL == "" {
		return ErrMissingDBURL
	}
	return nil
}
