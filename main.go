package main

import (
	"flag"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"os"

	log "github.com/google/logger"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"gopkg.in/yaml.v3"
)

type CloudsYaml struct {
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

const logPath = "osbckp.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func loadConfFile(l *log.Logger) []byte {
	l.Info("Loading clouds file")
	const cloudsFile = "clouds.yaml"
	file, err := os.Open(cloudsFile)

	if err != nil {
		l.Fatal(err.Error())
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		l.Infoln("Could not obtain stat: ", err.Error())
	}
	l.Infof("Loaded %d bytes from %s\n", fi.Size(), cloudsFile)

	data := make([]byte, fi.Size())

	count, err := file.Read(data)
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Infof("%d bytes read\n data:\n %s \n", count, string(data))
	return data
}

func main() {
	flag.Parse()
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()
	l := log.Init("OpenStackBackup", *verbose, true, lf)
	defer l.Close()

	data := loadConfFile(l)

	var conf CloudsYaml

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		l.Fatalf("error: %v", err)
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: conf.Clouds.OpenStack.Auth.AuthUrl,
		UserID:           conf.Clouds.OpenStack.Auth.UserID,
		Password:         conf.Clouds.OpenStack.Auth.Password,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		l.Fatalf("error: %v ", err)
	}

	l.Infoln("token ", provider.TokenID)

	client, err := openstack.NewComputeV2(
		provider,
		gophercloud.EndpointOpts{
			Region:       conf.Clouds.OpenStack.RegionName,
			Availability: gophercloud.AvailabilityAdmin,
		})
	if err != nil {
		l.Fatalf("error: %v", err)
	}

	server, err := servers.Get(client, "a414b947-4e9a-4cff-9294-c89dffc617b1").Extract()
	if err != nil {
		l.Fatalf("error: %v", err)
	}
	l.Infoln(server.Status)
}
