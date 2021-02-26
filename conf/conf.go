package conf

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

var Config *Conf

type Server struct {
	Listen             string `json:"listen"`
	Port               int    `json:"port"`
	JwtKey             string `json:"jwt_key"`
	RecaptchaSecretKey string `json:"recaptchaKey"`
}

type DBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
	Charset  string `json:"charset"`
}

type MailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

type Restrict struct {
	RestrictionCount int64 `json:"restrictionCount"`
	RestrictionTime  int64 `json:"restrictionTime"`
}

type Redis struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Conf struct {
	Database DBConfig   `json:"database"`
	Server   Server     `json:"server"`
	Mail     MailConfig `json:"mail"`
	Redis    Redis      `json:"redis"`
	Restrict Restrict   `json:"restrict"`
}

func init() {
	var confDir string
	flag.StringVar(&confDir, "c", "./conf.json", "location of configuration file")
	flag.Parse()
	Config = &Conf{
		Database: DBConfig{},
		Server:   Server{},
	}
	file, err := os.Open(confDir)
	if err != nil {
		log.Fatal(err)
	}
	confBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(confBytes, Config)
	if err != nil {
		log.Panicln("error while unmarshalling config json", err)
	}
	log.Println("configuration loaded successfully")
}
