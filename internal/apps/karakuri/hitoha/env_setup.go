package hitoha

import (
	"encoding/json"
	"errors"
	"futaba"
	"karakuri_mod"
	"karakuripkgs"
	"net"
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

type NodeInfo struct {
	Mode   string `json:"mode"`
	Target string `json:"target"`
	Status string `json:"connection_status"`
}

const (
	STAND_ALONE_MODE = "stand-alone"
)

func checkTargetNodeFile() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_NODECTL_ROOT); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_NODECTL_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(karakuripkgs.KARAKURI_NODECTL_NODEINFO); stat != nil {
		var node_info NodeInfo
		node_info.Mode = STAND_ALONE_MODE
		node_info.Target = karakuripkgs.SERVER
		node_info.Status = "connected"
		data, _ := json.MarshalIndent(node_info, "", "  ")
		if err := os.WriteFile(karakuripkgs.KARAKURI_NODECTL_NODEINFO, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func isPrivateIP(ip net.IP) bool {
	var prvMasks []*net.IPNet

	for _, cidr := range []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	} {
		_, mask, _ := net.ParseCIDR(cidr)
		prvMasks = append(prvMasks, mask)
	}

	for _, mask := range prvMasks {
		if mask.Contains(ip) {
			return true
		}
	}
	return false
}

func getDeviceIpAddress() (net.IP, error) {
	ift, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifi := range ift {
		addrs, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}
	return nil, errors.New("no IP")
}

func getTargetNode() NodeInfo {
	checkTargetNodeFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_NODECTL_NODEINFO)
	if err != nil {
		panic(err)
	}

	var node_info NodeInfo
	if err := json.Unmarshal(bytes, &node_info); err != nil {
		panic(err)
	}

	return node_info
}

func setupRemoteController() {
	node_mode := getTargetNode().Mode
	if node_mode == STAND_ALONE_MODE {
		device_address, _ := getDeviceIpAddress()
		iptables_cmd := exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "-i", karakuripkgs.HOST_NIC, "-d", device_address.String(), "--dport", "9816", "-j", "DROP")
		if err := iptables_cmd.Run(); err != nil {
			panic(err)
		}
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
	// setup remote controller
	setupRemoteController()
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

		// start restart: always container
		autoStartContainer()

		// create init file
		createInitFile()
	}
}
