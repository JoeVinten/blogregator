package config

func (c *Config) SetUser(username string) error {

	c.CurrentUsername = username

	return write(*c)

}
