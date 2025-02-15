package karakuri

import (
	"encoding/json"
	"errors"
	"fmt"
	"karakuripkgs"
	"net"
	"net/http"
	"os"
	"os/exec"
)

type ClusterInfo struct {
	Mode   string `json:"mode"`
	Target string `json:"target"`
	Status string `json:"connection_status"`
}

const (
	STAND_ALONE_MODE = "stand-alone"
	CLUSTER_MODE     = "cluster-mode"
)

func checkTargetClusterFile() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_CLSCTL_ROOT); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_CLSCTL_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(karakuripkgs.KARAKURI_CLSCTL_CLSINFO); stat != nil {
		var cluster_info ClusterInfo
		cluster_info.Mode = STAND_ALONE_MODE
		cluster_info.Target = karakuripkgs.SERVER
		cluster_info.Status = "connected"
		data, _ := json.MarshalIndent(cluster_info, "", "  ")
		if err := os.WriteFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func setTargetCluster(cluster string) {
	checkTargetClusterFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO)
	if err != nil {
		panic(err)
	}

	var cluster_info ClusterInfo
	if err := json.Unmarshal(bytes, &cluster_info); err != nil {
		panic(err)
	}

	// set cluster
	cluster_info.Target = cluster
	cluster_info.Status = "disconnected"

	data, _ := json.MarshalIndent(cluster_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func setClusterStatus(status string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO)
	if err != nil {
		panic(err)
	}

	var cluster_info ClusterInfo
	if err := json.Unmarshal(bytes, &cluster_info); err != nil {
		panic(err)
	}

	// set cluster
	cluster_info.Status = status

	data, _ := json.MarshalIndent(cluster_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func setClusterMode(mode string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO)
	if err != nil {
		panic(err)
	}

	var cluster_info ClusterInfo
	if err := json.Unmarshal(bytes, &cluster_info); err != nil {
		panic(err)
	}

	// set cluster
	cluster_info.Mode = mode

	data, _ := json.MarshalIndent(cluster_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func getTargetCluster() ClusterInfo {
	checkTargetClusterFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO)
	if err != nil {
		panic(err)
	}

	var cluster_info ClusterInfo
	if err := json.Unmarshal(bytes, &cluster_info); err != nil {
		panic(err)
	}

	return cluster_info
}

func verifyConnectionToCluster() (string, bool) {
	cluster_info := getTargetCluster()
	cluster := cluster_info.Target
	if cluster == "" {
		return "failed to get target cluster", false
	}
	// connect test
	// request retrieve catalog
	url := "http://" + cluster + "/container/ls/default"
	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		setClusterStatus("connectoin_failed")
		return "cluster connection failed", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		setClusterStatus("connectoin_failed")
		return "cluster connection failed", false
	}

	setClusterStatus("connected")
	return "cluster: " + cluster + " connection success", true
}

func checkConnectionStatus() bool {
	checkTargetClusterFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_CLSCTL_CLSINFO)
	if err != nil {
		panic(err)
	}

	var cluster_info ClusterInfo
	if err := json.Unmarshal(bytes, &cluster_info); err != nil {
		panic(err)
	}

	if cluster_info.Status != "connected" {
		return false
	}
	return true
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

// -- called from endpoint.go
func ShowTargetCluster() {
	cluster_info := getTargetCluster()
	printTargetCluster(cluster_info)
}

func ConnectCluster(cluster string) {
	if checkConnectionStatus() {
		cluster_info := getTargetCluster()
		cluster := cluster_info.Target
		fmt.Println("already connected to cluster: " + cluster + ".\nplease execute `karakuri clsctl disconnect` before change connection cluster.")
		return
	}
	// set env
	var target_cluster string
	if cluster == "default" || cluster == "localhost" {
		target_cluster = "localhost:9806"
	} else {
		target_cluster = cluster + ":9816"
	}
	setTargetCluster(target_cluster)
	// connection test
	message, res := verifyConnectionToCluster()
	if !res {
		fmt.Println("failed to connect cluster: " + cluster + ". please verify the cluster address is correct.\nOr make sure the target cluster is running in 'cluster-mode'")
		return
	}
	fmt.Println(message)
}

func DisconnectCluster() {
	if !checkConnectionStatus() {
		fmt.Println("still not connected any cluster.\nplease execute `karakuri clsctl connect --cluster {cluster}` first.")
	}
	setClusterStatus("disconnected")
	fmt.Println("cluster dissconnected.")
}

func EnableClusterMode() {
	cluster_mode := getTargetCluster().Mode
	if cluster_mode == CLUSTER_MODE {
		fmt.Println("already running in cluster-mode")
		return
	}
	setClusterMode(CLUSTER_MODE)
	// accept
	device_address, _ := getDeviceIpAddress()
	iptables_cmd := exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "-i", karakuripkgs.HOST_NIC, "-d", device_address.String(), "--dport", "9816", "-j", "DROP")
	if err := iptables_cmd.Run(); err != nil {
		fmt.Println("failed to enable cluster-mode")
		return
	}
	fmt.Println("currently running in cluster-mode")
}

func DisableClusterMode() {
	cluster_mode := getTargetCluster().Mode
	if cluster_mode == STAND_ALONE_MODE {
		fmt.Println("already running in stand-alone-mode")
		return
	}
	setClusterMode(STAND_ALONE_MODE)
	// accept
	device_address, _ := getDeviceIpAddress()
	iptables_cmd := exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "-i", karakuripkgs.HOST_NIC, "-d", device_address.String(), "--dport", "9816", "-j", "DROP")
	if err := iptables_cmd.Run(); err != nil {
		fmt.Println("failed to enable cluster-mode")
		return
	}
	fmt.Println("currently running in stand-alone-mode")
}
