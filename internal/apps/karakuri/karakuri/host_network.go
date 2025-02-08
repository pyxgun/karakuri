package karakuri

import (
	"karakuripkgs"
	"os/exec"
	"strconv"
	"strings"
)

func setIpTables(mode string, host_port int, target_ip string, target_port int, protocol string) {
	source_port := strconv.Itoa(host_port)
	container_ip := (strings.Split(target_ip, "/"))[0]
	container_port := strconv.Itoa(target_port)

	var mode_option string
	if mode == "add" {
		mode_option = "-A"
	} else {
		mode_option = "-D"
	}

	cmd1 := exec.Command("iptables", "-t", "nat", mode_option, "PREROUTING", "-p", protocol, "--dport", source_port, "-j", "DNAT", "--to-destination", container_ip+":"+container_port)
	if err := cmd1.Run(); err != nil {
		panic(err)
	}
	cmd2 := exec.Command("iptables", "-t", "nat", mode_option, "OUTPUT", "-p", protocol, "--dport", source_port, "-j", "DNAT", "--to-destination", container_ip+":"+container_port)
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
	cmd3 := exec.Command("iptables", "-t", "nat", mode_option, "POSTROUTING", "-p", protocol, "-d", container_ip, "--dport", container_port, "-j", "MASQUERADE")
	if err := cmd3.Run(); err != nil {
		panic(err)
	}
	cmd4 := exec.Command("iptables", mode_option, "FORWARD", "-p", protocol, "-d", container_ip, "--dport", container_port, "-j", "ACCEPT")
	if err := cmd4.Run(); err != nil {
		panic(err)
	}
	cmd5 := exec.Command("iptables", mode_option, "FORWARD", "-p", protocol, "-s", container_ip, "--sport", container_port, "-j", "ACCEPT")
	if err := cmd5.Run(); err != nil {
		panic(err)
	}
}

func SetupPortForwarding(mode string, network karakuripkgs.SpecNetwork) {
	container_ip := network.Address
	for _, entry := range network.Port {
		host_port := entry.HostPort
		container_port := entry.TargetPort
		protocol := entry.Protocol
		setIpTables(mode, host_port, container_ip, container_port, protocol)
	}
}
