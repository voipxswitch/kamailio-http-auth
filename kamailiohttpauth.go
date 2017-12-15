package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/romana/rlog"
	"goji.io"

	"github.com/voipxswitch/kamailio-http-auth/internal/http"
	"github.com/voipxswitch/kamailio-http-auth/internal/userdata"
)

const (
	defaultListenAddress = "localhost:8000"
)

var (
	confPath = "config.json"
)

func init() {
	rlog.SetOutput(os.Stdout)
}

func main() {
	rlog.Info("started kamailio-http-auth")
	rlog.Debug("debug enabled")

	configFilePath := flag.String("config", "", "path to config file")
	flag.Parse()
	if *configFilePath != "" {
		confPath = *configFilePath
	}
	rlog.Infof("loading config from file [%s]", confPath)

	c, err := loadConfigFile(confPath)
	if err != nil {
		rlog.Errorf("could load config [%s]", err.Error())
		os.Exit(1)
	}

	listenAddress := c.ListenAddress
	if listenAddress == "" {
		listenAddress = defaultListenAddress
	}
	rlog.Infof("set listen address [%s]", listenAddress)

	err = userdata.Setup(c.UserFile)
	if err != nil {
		rlog.Errorf("could not setup user data [%s]", err.Error())
		os.Exit(1)
	}

	// start http
	http.New(listenAddress, goji.NewMux())
}


// struct used to unmarshal config.json
type serviceConfig struct {
	ListenAddress string `json:"listen_address"`
	UserFile      string `json:"user_file"`
}

func loadConfigFile(configFile string) (serviceConfig, error) {
	s := serviceConfig{}
	file, err := os.Open(configFile)
	if err != nil {
		return s, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}
