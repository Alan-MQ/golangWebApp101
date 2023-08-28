type Config struct {
	DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri string `mapstructure:"REDIS_URI"`
	Port     string `mapstructure:"PORT"`
}