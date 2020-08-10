package osbckp

import (
	"github.com/joivo/osbckp/util"
	"strings"
)

const (
	zipSuffix = ".zip"
	backupBaseDir = "/home/backup/"

	kollaConfigDir = "/etc/kolla"

	dockerBaseDir  = "/home/docker/"
	volumesBaseDir = dockerBaseDir + "volumes/"

	containersDir  = dockerBaseDir + "containers/"
	builderDir  = dockerBaseDir + "builder/"
	imageDir = dockerBaseDir + "image/"
	networkDir = dockerBaseDir + "network/"
	overlay2Dir = dockerBaseDir + "overlay2/"
	pluginsDir = dockerBaseDir + "plugins/"
	runtimesDir = dockerBaseDir + "runtimes/"
	trustDir = dockerBaseDir + "trust/"
	cinderDir = volumesBaseDir + "cinder/"
	fluentdDir = volumesBaseDir + "fluentd_data/"
	glanceDir = volumesBaseDir + "glance/"
	grafanaDir = volumesBaseDir + "grafana/"
	haProxyDir = volumesBaseDir + "haproxy_socket/"
	iscsiDir = volumesBaseDir + "iscsi_info/"
	keystoneFernetTokensDir = volumesBaseDir + "keystone_fernet_tokens/"
	libvirtdDir = volumesBaseDir + "libvirtd/"
	magnumDir = volumesBaseDir + "magnum/"
	mariaDBDir = volumesBaseDir + "mariadb/"
	metadaFile = volumesBaseDir + "metadata.db"
	neutroMetadataDir = volumesBaseDir + "neutron_metadata_socket/"
	novaComputeDir = volumesBaseDir + "nova_compute/"
	novaLibvirtDir = volumesBaseDir + "nova_libvirt_qemu/"
	oVSwitchDir = volumesBaseDir + "openvswitch_db/"
	prometheusDir = volumesBaseDir + "prometheus/"
	rabbitmqDir = volumesBaseDir + "rabbitmq/"
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
	// Zip kolla config
	err := createZip(kollaConfigDir)
	util.HandleErr(err)

	err = createZip(containersDir)
	util.HandleErr(err)

	err = createZip(builderDir)
	util.HandleErr(err)

	err = createZip(imageDir)
	util.HandleErr(err)

	err = createZip(networkDir)
	util.HandleErr(err)

	err = createZip(overlay2Dir)
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

	err = createZip(haProxyDir)
	util.HandleErr(err)

	err = createZip(iscsiDir)
	util.HandleErr(err)

	err = createZip(keystoneFernetTokensDir)
	util.HandleErr(err)

	err = createZip(libvirtdDir)
	util.HandleErr(err)

	err = createZip(magnumDir)
	util.HandleErr(err)

	err = createZip(mariaDBDir)
	util.HandleErr(err)

	err = createZip(metadaFile)
	util.HandleErr(err)

	err = createZip(neutroMetadataDir)
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