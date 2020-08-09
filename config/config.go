package config

import (
	"os"

	"github.com/joivo/osbckp/util"

	"github.com/nuveo/log"
	"gopkg.in/yaml.v3"
)

type OpenStackConfig struct {
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
}

func GetOpenStackConfig() *OpenStackConfig {
	return &openStackConfig
}

func loadConfFile() []byte {
	log.Println("Loading clouds file")
	const cloudsFile = "clouds.yaml"
	file, err := os.Open(cloudsFile)

	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Errorln("Could not obtain stat: ", err.Error())
	}
	log.Printf("Loaded %d bytes from %s\n", fi.Size(), cloudsFile)

	data := make([]byte, fi.Size())

	count, err := file.Read(data)
	util.HandleErr(err)

	log.Printf("%d bytes read\n", count)
	return data
}
