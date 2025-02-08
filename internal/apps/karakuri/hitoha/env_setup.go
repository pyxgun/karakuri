package hitoha

import (
	"encoding/json"
	"futaba"
	"karakuri_mod"
	"karakuripkgs"
	"os"
	"os/exec"
)

const (
	INIT_FLAG = "/tmp/karakuri_init"
)

func setupDirectory() {
	if _, stat := os.Stat(karakuripkgs.FUTABA_ROOT); stat != nil {
		if err := os.MkdirAll(karakuripkgs.FUTABA_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}
	if _, stat := os.Stat(karakuripkgs.HITOHA_ROOT); stat != nil {
		if err := os.MkdirAll(karakuripkgs.HITOHA_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}
	if _, stat := os.Stat(karakuripkgs.IMAGE_ROOT); stat != nil {
		if err := os.MkdirAll(karakuripkgs.IMAGE_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(karakuripkgs.HITOHA_CONTAINER_LIST); stat != nil {
		newContainerList()
	}
	if _, stat := os.Stat(karakuripkgs.HITOHA_IMAGE_LIST); stat != nil {
		newImageList()
	}
	if _, stat := os.Stat(karakuripkgs.HITOHA_NAMESPACE_LIST); stat != nil {
		createNamespaceList()
	}
	if _, stat := os.Stat(karakuripkgs.FUTABA_CONTAINER_LIST); stat != nil {
		futaba.NewContainerList()
	}

	// module
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_ROOT); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_LIST); stat != nil {
		karakuri_mod.NewModList()
	}
}

func setupNetworkInterface() {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		createNamespaceList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	// setup interfaces
	for _, ns := range namespace_list_data.Namespaces {
		// create interface
		createNetworkInterface(
			ns.Network.DevType,
			ns.Network.Name,
			ns.Network.Address,
			ns.Network.Subnet,
		)
		// allow traffic rule
		allowContainerTrafficRule(ns.Network.Name)
	}
}

func setupNat() {
	// masquerade
	cmd1 := exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", "10.157.0.0/16", "-j", "MASQUERADE")
	if err := cmd1.Run(); err != nil {
		panic(err)
	}
}

func setupIpForward() {
	cmd := exec.Command("/sbin/sysctl", "-w", "net.ipv4.ip_forward=1")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func setupBridgeTrafficRule() {
	// iptables -I FORWARD -m physdev --physdev-is-bridged -j ACCEPT
	cmd := exec.Command("iptables", "-I", "FORWARD", "-m", "physdev", "--physdev-is-bridged", "-j", "ACCEPT")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func setupSystemModTrafficRule() {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		createNamespaceList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	for _, entry := range namespace_list_data.Namespaces[2:] {
		allowContainerModTrafficRule(entry.Network.Name)
	}
}

// set all container status to "stopped"
func changeContainerStatusToStop() {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		panic(err)
	}

	var container_list_data ContainerList
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	for i, _ := range container_list_data.List {
		container_list_data.List[i].Status = "stopped"
	}

	data, _ := json.MarshalIndent(container_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_CONTAINER_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func checkInitStatus() bool {
	if _, err := os.ReadFile(INIT_FLAG); err != nil {
		return false
	}
	return true
}

// setup network
func setupNetworks() {
	// setup futaba's bridge interface
	setupNetworkInterface()
	// setup nat using iptables
	setupNat()
	// setup bridge traffic allow rule
	setupBridgeTrafficRule()
	// setup syste-mod traffic to container
	setupSystemModTrafficRule()
	// setup ip forwarding
	setupIpForward()
}

func createInitFile() {
	fd, err := os.Create(INIT_FLAG)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	if _, err := fd.Write([]byte("init complete")); err != nil {
		panic(err)
	}
}

func SetupEnvironment() {
	if !checkInitStatus() {
		// setup directory
		setupDirectory()
		// setup network environment
		setupNetworks()

		// change container status to stopped
		changeContainerStatusToStop()
		// setup module
		karakuri_mod.SetupModules()
		// create init file
		createInitFile()
	}
}
