package futaba

import (
	"os/exec"
	"strconv"

	"karakuripkgs"
)

func setupContainerNetwork(pid int, networ_spec karakuripkgs.SpecNetwork) {
	str_pid := strconv.Itoa(pid)
	// create veth
	cmd1 := exec.Command("ip", "link", "add", "name", "karakuri"+str_pid, "type", "veth", "peer", "name", karakuripkgs.HOST_NIC, "netns", str_pid)
	if err := cmd1.Run(); err != nil {
		panic(err)
	}
	// set eth0 ip address
	cmd2 := exec.Command("nsenter", "-t", str_pid, "-n", "ip", "address", "add", networ_spec.Address, "dev", karakuripkgs.HOST_NIC)
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
	// link device
	cmd4 := exec.Command("ip", "link", "set", "dev", "karakuri"+str_pid, "master", networ_spec.HostDevice)
	if err := cmd4.Run(); err != nil {
		panic(err)
	}
	// link up lo
	cmd5 := exec.Command("nsenter", "-t", str_pid, "-n", "ip", "link", "set", "up", "lo")
	if err := cmd5.Run(); err != nil {
		panic(err)
	}
	// link up eth0
	cmd6 := exec.Command("nsenter", "-t", str_pid, "-n", "ip", "link", "set", "up", karakuripkgs.HOST_NIC)
	if err := cmd6.Run(); err != nil {
		panic(err)
	}
	// link up veth
	cmd7 := exec.Command("ip", "link", "set", "up", "karakuri"+str_pid)
	if err := cmd7.Run(); err != nil {
		panic(err)
	}
	// set default route
	cmd8 := exec.Command("nsenter", "-t", str_pid, "-n", "ip", "route", "add", "default", "via", networ_spec.Gateway)
	if err := cmd8.Run(); err != nil {
		panic(err)
	}
}
