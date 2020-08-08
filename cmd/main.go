package main

import (
	"flag"
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

func CreateServersSnapshots(l *log.Logger, conf CloudsYaml, provider *gophercloud.ProviderClient) {

	computeOpts := gophercloud.EndpointOpts{
		Region:       conf.Clouds.OpenStack.RegionName,
		Availability: gophercloud.AvailabilityAdmin,
	}

	computeV2, err := openstack.NewComputeV2(provider, computeOpts)
	handleErr(err, l)

	allPages, err := servers.List(computeV2, servers.ListOpts{
		AllTenants:   true,
	}).AllPages()
	handleErr(err, l)

	srvs, err := servers.ExtractServers(allPages)
	handleErr(err, l)

	l.Infof("%d images were found", len(srvs))

	var srvc = make(chan servers.Server)

	go func() {
		for srv := range srvc {
			const layout = "2006-01-02"
			snapshotName := srv.Name + "_" + time.Now().Format(layout)
			createImgOpts := servers.CreateImageOpts{
				Name: snapshotName,
			}
			l.Infoln("Snapshot name", snapshotName)
			l.Infof("Sending request to build image for %s", srv.Name)
			r := servers.CreateImage(computeV2, srv.ID, createImgOpts)
			l.Infoln(r.Result.PrettyPrintJSON())
		}
	}()

	for _, srv := range srvs {
		timeout := (time.Duration(2) * time.Second).Seconds()
		err := gophercloud.WaitFor(int(timeout), func() (bool, error) {
			srvc <- srv
			return true, nil
		})
		handleErr(err, l)
	}

	time.Sleep(120 * time.Second)
}

func CreateClientProvider(conf CloudsYaml, l *log.Logger) (*gophercloud.ProviderClient, error) {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: conf.Clouds.OpenStack.Auth.AuthUrl,
		UserID:           conf.Clouds.OpenStack.Auth.UserID,
		Password:         conf.Clouds.OpenStack.Auth.Password,
		Scope: &gophercloud.AuthScope{
			ProjectID: conf.Clouds.OpenStack.Auth.ProjectID,
		},
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	return provider, err
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
	l.Infoln("**** Starting Backup Service ****")

	data := loadConfFile(l)

	var conf CloudsYaml

	err = yaml.Unmarshal(data, &conf)
	handleErr(err, l)

	provider, err := CreateClientProvider(conf, l)
	handleErr(err, l)
	CreateServersSnapshots(l, conf, provider)
}

func handleErr(err error, l *log.Logger) {
	if err != nil {
		l.Fatalf("error: %v", err)
	}
}
