// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period       time.Duration `config:"period"`
	Archive      bool          `config:"archive"`
	DiscordUsers []DiscordUser `config:"users"`
}

type DiscordUser struct {
	Token    string  `config:"token"`
	Username string  `config:"username"`
	Password string  `config:"password"`
	Guilds   []Guild `config:"guilds"`
}

type Guild struct {
	ID       string    `config:"id"`
	Channels []string `config:"channels"`
}

var DefaultConfig = Config{
	Period:  5 * time.Second,
	Archive: false,
}
