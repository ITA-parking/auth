package env

// ALL ENV MUST BE PREFIXED WITH SERVICE FOLLOWED BY THE VARIABLE NAME IN SNAKE_CASE
type Options struct {
	HttpPort           int    `help:"Http port" default:"8080"`
	DbConnectionString string `help:"PostgreSQL connection string"`
	JwtSecret          string `help:"JWT secret"`
}
