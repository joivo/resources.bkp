package main

import (
	"flag"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"os"
	"time"

	log "github.com/google/logger"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
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
	handleErr(err, l)
	l.Infof("%d bytes read\n", count)
	return data
}

func main() {
	const logPath = "osbckp.log"
	flag.Parse()
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()
	l := log.Init("OpenStackBackup", *verbose, false, lf)
	defer l.Close()

	data := loadConfFile(l)

	var conf CloudsYaml

	err = yaml.Unmarshal(data, &conf)
	handleErr(err, l)

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: conf.Clouds.OpenStack.Auth.AuthUrl,
		UserID:           conf.Clouds.OpenStack.Auth.UserID,
		Password:         conf.Clouds.OpenStack.Auth.Password,
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	handleErr(err, l)

	clientOpts := gophercloud.EndpointOpts{
		Region:       conf.Clouds.OpenStack.RegionName,
		Availability: gophercloud.AvailabilityAdmin,
	}

	identityV3, err := openstack.NewIdentityV3(provider, clientOpts)

	handleErr(err, l)

	allPages, err := projects.List(identityV3, nil).AllPages()
	handleErr(err, l)
	extractedProjects, err := projects.ExtractProjects(allPages)




	computeV2, err := openstack.NewComputeV2(provider, clientOpts)
	handleErr(err, l)
	allPages, err = servers.List(computeV2, nil).AllPages()
	handleErr(err, l)
	srvs, err := servers.ExtractServers(allPages)
	handleErr(err, l)
	l.Infof("%d images were found", len(srvs))

	for _, srv := range srvs {
		const layout = "2006-01-02 15:04:05"
		snapshotName := srv.Name + "_" + time.Now().Format(layout)
		createImgOpts := servers.CreateImageOpts{
			Name: snapshotName,
		}
		l.Infoln("Snapshot name", snapshotName)
		l.Infof("Sending request to build image for %s", srv.Name)
		servers.CreateImage(computeV2, srv.ID, createImgOpts)
	}
}

func handleErr(err error, l *log.Logger) {
	if err != nil {
		l.Fatalf("error: %v", err)
	}
}
