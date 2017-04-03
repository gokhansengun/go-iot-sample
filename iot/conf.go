package iot

import (
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

// PostgresDetails keeps the database user, name, password
type PostgresDetails struct {
	Host   string `yaml:"Host"`
	Port	int `yaml:"Port"`
	User   string `yaml:"User"`
	Password string `yaml:"Password"`
	DbName string `yaml:"DbName"`
}

// AppConf is the type to decompose the configuration
type AppConf struct {
	MongoDbReplicaList []string        `yaml:"MongoDbReplicaList"`
	KafkaBrokerList    []string        `yaml:"KafkaBrokerList"`
	PgDetails          PostgresDetails `yaml:"PostgresDetails"`
	configRead         bool
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
	return "mongodb://" + appConf.MongoDbReplicaList[0]
}

// GetKafkaConnStr returns the KafkaBrokerList in the format expected
func (appConf AppConf) GetKafkaConnStr() []string {
	return appConf.KafkaBrokerList
}

// GetPostgresConnStr returns the PostgresConn str in the format expected
func (appConf AppConf) GetPostgresConnStr() string {
	pgDetails := appConf.PgDetails

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgDetails.Host, pgDetails.Port, pgDetails.User, pgDetails.Password, pgDetails.DbName)

	// fmt.Printf("The Postgres conn str is %s\n", connStr)

	return connStr
}
