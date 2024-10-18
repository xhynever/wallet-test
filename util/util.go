package util

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var GlobalConfig *Config

func init() {
	rand.Seed(time.Now().Unix())
}

type Config struct {
	PgUrl             string `mapstructure:"POSTGRES_URL"`
	PgUsername        string `mapstructure:"POSTGRES_USER"`
	PgPassword        string `mapstructure:"POSTGRES_PASSWORD"`
	PgHost            string `mapstructure:"POSTGRES_HOST"`
	PgName            string `mapstructure:"POSTGRES_DB"`
	PgPool            int    `mapstructure:"POSTGRES_POOL"`
	PgPort            string `mapstructure:"POSTGRES_PORT"`
	DbDriver          string `mapstructure:"DB_DRIVER"`
	ServersAddress    string `mapstructure:"ADDRESS"`
	AuthorizationPort int    `mapstructure:"AUTHORIZATION_PORT"`
}
type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	Pool     string
}

func InitConfig() *Config {
	config, err := LoadConfig(".")
	if err != nil {
		logrus.Fatal("cannot load config: ", err)
	}
	GlobalConfig = &config
	return GlobalConfig
}
func InitDB() (*sqlx.DB, error) {
	dbConfig := DbConfig{
		Host:     GlobalConfig.PgHost,
		Username: GlobalConfig.PgUsername,
		Password: GlobalConfig.PgPassword,
		DbName:   GlobalConfig.PgName,
		Port:     GlobalConfig.PgPort,
	}
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.DbName)
	// 指定db，可替换
	fmt.Println("db：", dbUrl)
	fmt.Println("db：", GlobalConfig.PgUrl)
	conn, err := sqlx.Open(GlobalConfig.DbDriver, GlobalConfig.PgUrl)
	if err != nil {
		logrus.Fatal("cannot connect to db", err)
		return nil, err
	}
	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
  CHECK (balance >= 0)
);
CREATE TABLE IF NOT EXISTS "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint,
  "amount" bigint ,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE IF NOT EXISTS "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint,
  "to_account_id" bigint,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	return conn, nil
}

func AuthPort() int {
	config, err := LoadConfig(".")
	if err != nil {
		logrus.Fatal("cannot load config: ", err)
	}
	return config.AuthorizationPort
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EUR, USD, CAD:
		return true
	}
	return false
}
