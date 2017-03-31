package device 

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// AppConf is the type to decompose the configuration
type AppConf struct {
	MongoDbReplicaList []string `yaml:"MongoDbReplicaList"`
	KafkaBrokerList    []string `yaml:"KafkaBrokerList"`
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
