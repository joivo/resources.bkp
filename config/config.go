package config

import (
	"os"
	"sync"

	"github.com/joivo/osbckp/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type OpenStackConfig struct {
	Mutex  *sync.Mutex
	Clouds struct {
		OpenStack struct {
			Auth struct {
				AuthUrl        string `yaml:"auth_url"`
				Username       string `yaml:"username"`
				UserID         string `yaml:"userid"`
				Password       string `yaml:"password"`
				ProjectID      string `yaml:"project_id"`
				ProjectName    string `yaml:"project_name"`
				UserDomainName string `yaml:"user_domain_name"`
			} `yaml:"auth"`
			RegionName         string `yaml:"region_name"`
			Interface          string `yaml:"interface"`
			IdentityAPIVersion int    `yaml:"identity_api_version"`
		} `yaml:"openstack"`
	} `yaml:"clouds"`
}

var openStackConfig OpenStackConfig

func LoadConfig() {
	data := loadConfFile()
	err := yaml.Unmarshal(data, &openStackConfig)
	util.HandleErr(err)
	openStackConfig.Mutex = new(sync.Mutex)
}

func GetOpenStackConfig() *OpenStackConfig {
	openStackConfig.Mutex.Lock()
	defer openStackConfig.Mutex.Unlock()
	return &openStackConfig
}

func loadConfFile() []byte {
	log.Info("Loading clouds file")
	const cloudsFile = "clouds.yaml"
	file, err := os.Open(cloudsFile)

	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Infoln("Could not obtain stat: ", err.Error())
	}
	log.Infof("Loaded %d bytes from %s\n", fi.Size(), cloudsFile)

	data := make([]byte, fi.Size())

	count, err := file.Read(data)
	util.HandleErr(err)

	log.Infof("%d bytes read\n", count)
	return data
}
