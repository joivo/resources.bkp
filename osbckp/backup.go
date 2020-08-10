package osbckp

import (
	"github.com/joivo/osbckp/util"
	"strings"
)

const (
	zipSuffix     = ".zip"
	backupBaseDir = "/home/backup/"

	kollaConfigDir = "/etc/kolla"

	dockerBaseDir  = "/home/docker/"
	volumesBaseDir = dockerBaseDir + "volumes/"

	containersDir           = dockerBaseDir + "containers/"
	builderDir              = dockerBaseDir + "builder/"
	imageDir                = dockerBaseDir + "image/"
	networkDir              = dockerBaseDir + "network/"
	pluginsDir              = dockerBaseDir + "plugins/"
	runtimesDir             = dockerBaseDir + "runtimes/"
	trustDir                = dockerBaseDir + "trust/"
	cinderDir               = volumesBaseDir + "cinder/"
	fluentdDir              = volumesBaseDir + "fluentd_data/"
	glanceDir               = volumesBaseDir + "glance/"
	grafanaDir              = volumesBaseDir + "grafana/"
	iscsiDir                = volumesBaseDir + "iscsi_info/"
	keystoneFernetTokensDir = volumesBaseDir + "keystone_fernet_tokens/"
	magnumDir               = volumesBaseDir + "magnum/"
	metadataFile            = volumesBaseDir + "metadata.db"
	novaComputeDir          = volumesBaseDir + "nova_compute/"
	novaLibvirtDir          = volumesBaseDir + "nova_libvirt_qemu/"
	oVSwitchDir             = volumesBaseDir + "openvswitch_db/"
	prometheusDir           = volumesBaseDir + "prometheus/"
	rabbitmqDir             = volumesBaseDir + "rabbitmq/"
)

func createZip(base string) (err error) {
	dest := backupBaseDir + strings.ReplaceAll(
		base,
		"/",
		"-",
	) + zipSuffix
	err = Zip(base, dest)
	return
}

func CreateBackup() {
	const backupPath = "/home/backup"
	util.CreatePathIfNotExist(backupPath)

	// Zip kolla config
	err := createZip(kollaConfigDir)
	util.HandleErr(err)

	// Zip docker data
	err = createZip(containersDir)
	util.HandleErr(err)

	err = createZip(builderDir)
	util.HandleErr(err)

	err = createZip(imageDir)
	util.HandleErr(err)

	err = createZip(networkDir)
	util.HandleErr(err)

	err = createZip(pluginsDir)
	util.HandleErr(err)

	err = createZip(runtimesDir)
	util.HandleErr(err)

	err = createZip(trustDir)
	util.HandleErr(err)

	err = createZip(cinderDir)
	util.HandleErr(err)

	err = createZip(fluentdDir)
	util.HandleErr(err)

	err = createZip(glanceDir)
	util.HandleErr(err)

	err = createZip(grafanaDir)
	util.HandleErr(err)

	err = createZip(iscsiDir)
	util.HandleErr(err)

	err = createZip(keystoneFernetTokensDir)
	util.HandleErr(err)

	err = createZip(magnumDir)
	util.HandleErr(err)

	err = createZip(metadataFile)
	util.HandleErr(err)

	err = createZip(novaComputeDir)
	util.HandleErr(err)

	err = createZip(novaLibvirtDir)
	util.HandleErr(err)

	err = createZip(oVSwitchDir)
	util.HandleErr(err)

	err = createZip(prometheusDir)
	util.HandleErr(err)

	err = createZip(rabbitmqDir)
	util.HandleErr(err)
}
