package karakuri

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"karakuripkgs"
	"net"
	"net/http"
	"os"
	"os/exec"
)

type NodeInfo struct {
	Mode   string `json:"mode"`
	Target string `json:"target"`
	Status string `json:"connection_status"`
}

const (
	STAND_ALONE_MODE    = "stand-alone"
	REMOTE_CONTROL_MODE = "remote-control"
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

func setTargetNode(node string) {
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

	// set node
	node_info.Target = node
	node_info.Status = "disconnected"

	data, _ := json.MarshalIndent(node_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_NODECTL_NODEINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func setNodeStatus(status string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_NODECTL_NODEINFO)
	if err != nil {
		panic(err)
	}

	var node_info NodeInfo
	if err := json.Unmarshal(bytes, &node_info); err != nil {
		panic(err)
	}

	// set node
	node_info.Status = status

	data, _ := json.MarshalIndent(node_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_NODECTL_NODEINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func setNodeMode(mode string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_NODECTL_NODEINFO)
	if err != nil {
		panic(err)
	}

	var node_info NodeInfo
	if err := json.Unmarshal(bytes, &node_info); err != nil {
		panic(err)
	}

	// set node
	node_info.Mode = mode

	data, _ := json.MarshalIndent(node_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_NODECTL_NODEINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
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

func verifyConnectionToNode(auth_code string) (string, bool) {
	node_info := getTargetNode()
	node := node_info.Target
	if node == "" {
		return "failed to get target node", false
	}
	// connect test
	// request retrieve catalog
	url := "http://" + node + "/container/ls/default"
	req, _ := http.NewRequest("GET", url, nil)
	// set auth code if not localhost
	if node != karakuripkgs.SERVER {
		if auth_code == "" {
			return "authentication code required. please specify auth code using '--auth' option", false
		}
		req.Header.Set("Authorization", auth_code)
	}

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		setNodeStatus("connectoin_failed")
		return "node connection failed", false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		setNodeStatus("authentication_failed")
		return "failed to authentication", false
	}
	if resp.StatusCode != 200 {
		setNodeStatus("connectoin_failed")
		return "node connection failed", false
	}

	setNodeStatus("connected")
	return "node: " + node + " connection success", true
}

func checkConnectionStatus() bool {
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

	if node_info.Status != "connected" {
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

// authentication method
func generateAuthCode() string {
	var (
		r = rand.Reader
		b = make([]byte, 32)
	)

	_, err := io.ReadFull(r, b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

type AuthCode struct {
	AuthCode string `json:"auth_code"`
}

type RemoteAuthCode struct {
	Node     string `json:"node"`
	AuthCode string `json:"auth_code"`
}

type RemoteNodeList struct {
	List []RemoteAuthCode `json:"node_list"`
}

func storeAuthCode(authcode string) {
	var auth_code AuthCode
	auth_code.AuthCode = authcode
	data, _ := json.MarshalIndent(auth_code, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_NODECTL_AUTHCODE, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func checkRemoteNodeListFile() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_NODECTL_REMOTE_AUTHCODE); stat != nil {
		var remote_node_list RemoteNodeList
		data, _ := json.MarshalIndent(remote_node_list, "", "  ")
		if err := os.WriteFile(karakuripkgs.KARAKURI_NODECTL_REMOTE_AUTHCODE, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func storeRemoteNodeList(node, auth_code string) {
	checkRemoteNodeListFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_NODECTL_REMOTE_AUTHCODE)
	if err != nil {
		return
	}

	var remote_node_list RemoteNodeList
	if err := json.Unmarshal(bytes, &remote_node_list); err != nil {
		panic(err)
	}

	set_flag := false
	for i, entry := range remote_node_list.List {
		if entry.Node == node {
			remote_node_list.List[i].AuthCode = auth_code
			set_flag = true
		}
	}
	if !set_flag {
		remote_node_list.List = append(remote_node_list.List,
			RemoteAuthCode{
				Node:     node,
				AuthCode: auth_code,
			},
		)
	}

	data, _ := json.MarshalIndent(remote_node_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_NODECTL_REMOTE_AUTHCODE, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func getRemoteAuthCode(node string) string {
	checkRemoteNodeListFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_NODECTL_REMOTE_AUTHCODE)
	if err != nil {
		return ""
	}

	var remote_node_list RemoteNodeList
	if err := json.Unmarshal(bytes, &remote_node_list); err != nil {
		panic(err)
	}

	for _, entry := range remote_node_list.List {
		if entry.Node == node {
			return entry.AuthCode
		}
	}
	return ""
}

// -- called from main.go
func ShowTargetNode() {
	node_info := getTargetNode()
	printTargetNode(node_info)
}

func ConnectNode(node string, auth_code string) {
	if checkConnectionStatus() {
		node_info := getTargetNode()
		node := node_info.Target
		fmt.Println("already connected to node: " + node + ".\nplease execute `karakuri nodectl disconnect` before change connection node.")
		return
	}
	// set env
	var target_node string
	if node == "default" || node == "localhost" {
		target_node = "localhost:9806"
	} else {
		target_node = node + ":9816"
	}
	setTargetNode(target_node)
	// connection test
	if target_node != karakuripkgs.SERVER {
		if auth_code == "" {
			auth_code = getRemoteAuthCode(target_node)
			if auth_code == "" {
				fmt.Println("it appears that we have not yet connected to node: " + node + ". please specify authentication code with '--auth' option.")
				return
			}
		}
	}
	message, res := verifyConnectionToNode(auth_code)
	if !res {
		fmt.Println(message)
		return
	}

	// store remote node list
	storeRemoteNodeList(target_node, auth_code)
	fmt.Println(message)
}

func DisconnectNode() {
	if !checkConnectionStatus() {
		fmt.Println("still not connected any node.\nplease execute `karakuri nodectl connect --node {node}` first.")
	}
	node := getTargetNode().Target
	setNodeStatus("disconnected")
	fmt.Println("node: " + node + " dissconnected.")
}

func EnableRemoteControllMode() {
	node_mode := getTargetNode().Mode
	if node_mode == REMOTE_CONTROL_MODE {
		fmt.Println("already running in remote-control mode")
		return
	}
	setNodeMode(REMOTE_CONTROL_MODE)
	// accept
	device_address, _ := getDeviceIpAddress()
	iptables_cmd := exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "-i", karakuripkgs.HOST_NIC, "-d", device_address.String(), "--dport", "9816", "-j", "DROP")
	if err := iptables_cmd.Run(); err != nil {
		fmt.Println("failed to enable remote-control mode")
		return
	}
	// generate auth code
	auth_code := generateAuthCode()
	storeAuthCode(auth_code)
	fmt.Println("currently running in remote-control mode")
	fmt.Println("now you can connect from remote to this node. please execute the following command on controller node:")
	fmt.Println("  karakuri nodectl connect --node " + device_address.String() + " --auth " + auth_code)
}

func DisableRemoteControllMode() {
	node_mode := getTargetNode().Mode
	if node_mode == STAND_ALONE_MODE {
		fmt.Println("already running in stand-alone mode")
		return
	}
	setNodeMode(STAND_ALONE_MODE)
	// accept
	device_address, _ := getDeviceIpAddress()
	iptables_cmd := exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "-i", karakuripkgs.HOST_NIC, "-d", device_address.String(), "--dport", "9816", "-j", "DROP")
	if err := iptables_cmd.Run(); err != nil {
		fmt.Println("failed to disable remote-control mode")
		return
	}
	fmt.Println("currently running in stand-alone mode")
}
