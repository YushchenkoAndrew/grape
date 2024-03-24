package config

// type EnvType struct {
// 	Host     string `mapstructure:"HOST"`
// 	Port     string `mapstructure:"PORT"`
// 	BasePath string `mapstructure:"BASE_PATH"`

// 	// DataBase
// 	DBType string `mapstructure:"DB_TYPE"`
// 	DBName string `mapstructure:"DB_NAME"`
// 	DBHost string `mapstructure:"DB_HOST"`
// 	DBPort string `mapstructure:"DB_PORT"`
// 	DBUser string `mapstructure:"DB_USER"`
// 	DBPass string `mapstructure:"DB_PASS"`

// 	// Redis
// 	RedisHost string `mapstructure:"REDIS_HOST"`
// 	RedisPort string `mapstructure:"REDIS_PORT"`
// 	RedisPass string `mapstructure:"REDIS_PASS"`

// 	// K3s
// 	K3sPath string `mapstructure:"K3S_PATH"`

// 	// JWT
// 	AccessSecret  string `mapstructure:"ACCESS_SECRET"`
// 	RefreshSecret string `mapstructure:"REFRESH_SECRET"`

// 	// Root User Login + Pass & Pepper
// 	ID     string `mapstructure:"grape_ID"`
// 	URL    string `mapstructure:"grape_URL"`
// 	User   string `mapstructure:"grape_USER"`
// 	Pass   string `mapstructure:"grape_PASS"`
// 	Pepper string `mapstructure:"grape_PEPPER"`

// 	// Pagination setting
// 	LiveTime int64 `mapstructure:"LIVE_TIME"`
// 	Items    int   `mapstructure:"ITEMS"`
// 	Limit    int   `mapstructure:"LIMIT"`

// 	// Rate Info
// 	RateLimit int `mapstructure:"RATE_LIMIT"`
// 	RateTime  int `mapstructure:"RATE_TIME"`

// 	BotUrl    string `mapstructure:"BOT_URL"`
// 	BotKey    string `mapstructure:"BOT_KEY"`
// 	BotPepper string `mapstructure:"BOT_Pepper"`

// 	// Migration Settings
// 	ForceMigrate  bool   `mapstructure:"FORCE_MIGRATE"`
// 	MigrationPath string `mapstructure:"MIGRATION_PATH"`
// }

// // FIXME: I should fix this one day
// var ENV EnvType

// type envConfig struct {
// 	path string
// 	name string
// }

// type ConfigT interface {
// 	Init()
// }

// func NewEnvConfig(path, name string) func() ConfigT {
// 	return func() ConfigT {
// 		return &envConfig{path, name}
// 	}
// }

// func (c *envConfig) Init() {
// 	viper.AddConfigPath(c.path)
// 	if c.name == "" {
// 		viper.SetConfigFile(".env")
// 	} else {
// 		viper.SetConfigName(c.name)
// 		viper.SetConfigType("env")
// 	}

// 	viper.AutomaticEnv()
// 	if err := viper.ReadInConfig(); err != nil {
// 		panic(fmt.Errorf("Failed on reading .env file %v", err))
// 	}

// 	if err := viper.Unmarshal(&ENV); err != nil {
// 		panic("Failed on reading .env file")
// 	}
// }
