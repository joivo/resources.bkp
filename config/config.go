package config

import (
	"os"
	"sync"
	"time"

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

const (
	DateLayout         = "2006-01-02"
	SnapshotSuffix     = "snapshot_"
	UsefulVolumeStatus = "available"
	UsefulServerStatus = "active"
	PoolingInterval    = 10 * time.Second
	FifteenDaysInMin   = 360
)

var (
	openStackConfig OpenStackConfig
	mu              = new(sync.Mutex)
)

const (
	cloudsFile = "clouds.yaml"
)

func LoadConfig() {
	cloudsConf := getBytesOfFile(cloudsFile)
	err := yaml.Unmarshal(cloudsConf, &openStackConfig)
	util.HandleFatal(err)
}

func GetOpenStackConfig() *OpenStackConfig {
	defer mu.Unlock()
	mu.Lock()
	return &openStackConfig
}

func getBytesOfFile(fileName string) []byte {
	log.Printf("Loading %s file\n", fileName)

	file, err := os.Open(fileName)
	util.HandleFatal(err)
	defer file.Close()

	fi, err := file.Stat()
	util.HandleErr(err)
	log.Printf("Loaded %d bytes from %s\n", fi.Size(), fileName)

	data := make([]byte, fi.Size())
	count, err := file.Read(data)
	util.HandleFatal(err)

	log.Printf("%d bytes read\n", count)
	return data
}
