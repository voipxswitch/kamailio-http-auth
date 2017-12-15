package userdata

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/romana/rlog"
)

var (
	userFile string
	users    map[string]userData
)

type userData struct {
	Password string `json:"password"`
}

func Setup(cfgUserFile string) error {
	var err error
	userFile = cfgUserFile
	users = make(map[string]userData)
	err = reloadUserData()
	if err != nil {
		return err
	}
	rlog.Infof("set user data file [%s]", userFile)
	return nil
}

func loadUserDataFromFile() (map[string]userData, error) {
	type userDirectoryData struct {
		Users map[string]userData `json:"users"`
	}
	s := userDirectoryData{}
	file, err := os.Open(userFile)
	if err != nil {
		return s.Users, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	if err != nil {
		return s.Users, err
	}

	return s.Users, nil
}

func reloadUserData() error {
	rlog.Debugf("reloading user data from file [%s]", userFile)
	newUserData, err := loadUserDataFromFile()
	if err != nil {
		return err
	}
	users = newUserData
	return nil
}

func LoadPassword(username string, realm string) (string, bool) {
	key := fmt.Sprintf("%s@%s", username, realm)
	rlog.Debugf("loading password for user [%s]", key)
	err := reloadUserData()
	if err != nil {
		rlog.Errorf("could not reload list [%s]", err.Error())
	}
	user, ok := users[key]
	if !ok {
		rlog.Infof("user [%s] not found", key)
		return "", false
	}
	return user.Password, true
}
