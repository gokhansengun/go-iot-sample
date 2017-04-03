package iot

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type tPostgres struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DbName   string `yaml:"DbName"`
}

type tMongo struct {
	Replicas []struct {
		Host string `yaml:"Host"`
		Port int    `yaml:"Port"`
	} `yaml:"Replicas"`
}

type tKafka struct {
	Brokers []struct {
		Host string `yaml:"Host"`
		Port int    `yaml:"Port"`
	} `yaml:"Brokers"`
}

type tRedis struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
}

// AppConf is the type to decompose the configuration
type AppConf struct {
	Mongo      tMongo    `yaml:"MongoDb"`
	Kafka      tKafka    `yaml:"Kafka"`
	Postgres   tPostgres `yaml:"Postgres"`
	Redis      tRedis    `yaml:"Redis"`
	configRead bool
}

// NewConfig creates a new configuration instance
func NewConfig(confPath string) *AppConf {
	appConf := new(AppConf)

	buf, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatalf("error in reading the configuration file, msg: %v", err)
	}

	err = yaml.Unmarshal(buf, &appConf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return appConf
}

// GetMongoDbConnStr returns the MongoDbReplicaList in the format
// expected by the mgo style connection string
func (appConf AppConf) GetMongoDbConnStr() string {
	// TODO: gseng - only supporting single server now
	return "mongodb://" + appConf.Mongo.Replicas[0].Host
}

// GetKafkaConnStr returns the KafkaBrokerList in the format expected
func (appConf AppConf) GetKafkaConnStr() []string {
	list := make([]string, len(appConf.Kafka.Brokers))

	for i := 0; i < len(list); i++ {
		list[i] = fmt.Sprintf("%s:%d", appConf.Kafka.Brokers[0].Host, appConf.Kafka.Brokers[0].Port)
	}

	return list
}

// GetPostgresConnStr returns the PostgresConn str in the format expected
func (appConf AppConf) GetPostgresConnStr() string {
	pgConf := appConf.Postgres

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgConf.Host, pgConf.Port, pgConf.Username, pgConf.Password, pgConf.DbName)

	// fmt.Printf("The Postgres conn str is %s\n", connStr)

	return connStr
}

// GetRedisAddr returns the RedisConn str in the format expected
func (appConf AppConf) GetRedisAddr() string {
	redisConf := appConf.Redis

	addr := fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)

	return addr
}
