package karakuripkgs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ResponseContainerId struct {
	Result  string `json:"result"`
	Id      string `json:"id"`
	Message string `json:"message"`
}

type NodeInfo struct {
	Mode   string `json:"mode"`
	Target string `json:"target"`
	Status string `json:"connection_status"`
}

type RemoteAuthCode struct {
	Node     string `json:"node"`
	AuthCode string `json:"auth_code"`
}

type RemoteNodeList struct {
	List []RemoteAuthCode `json:"node_list"`
}

func checkTargetNodeFile() {
	if _, stat := os.Stat(KARAKURI_NODECTL_ROOT); stat != nil {
		if err := os.MkdirAll(KARAKURI_NODECTL_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(KARAKURI_NODECTL_NODEINFO); stat != nil {
		var node_info NodeInfo
		node_info.Mode = "stand-alone"
		node_info.Target = SERVER
		node_info.Status = "connected"
		data, _ := json.MarshalIndent(node_info, "", "  ")
		if err := os.WriteFile(KARAKURI_NODECTL_NODEINFO, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func getTargetNode() NodeInfo {
	checkTargetNodeFile()

	var bytes []byte
	bytes, err := os.ReadFile(KARAKURI_NODECTL_NODEINFO)
	if err != nil {
		panic(err)
	}

	var node_info NodeInfo
	if err := json.Unmarshal(bytes, &node_info); err != nil {
		panic(err)
	}

	return node_info
}

func checkRemoteNodeListFile() {
	if _, stat := os.Stat(KARAKURI_NODECTL_REMOTE_AUTHCODE); stat != nil {
		var remote_node_list RemoteNodeList
		data, _ := json.MarshalIndent(remote_node_list, "", "  ")
		if err := os.WriteFile(KARAKURI_NODECTL_REMOTE_AUTHCODE, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
func getRemoteAuthCode(node string) string {
	checkRemoteNodeListFile()

	var bytes []byte
	bytes, err := os.ReadFile(KARAKURI_NODECTL_REMOTE_AUTHCODE)
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

// request retrieve container id
func RequestContainerId(name string) (result bool, id string, message string) {
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/container/getid/" + name

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != SERVER {
		auth_code := getRemoteAuthCode(node)
		if auth_code == "" {
			fmt.Println("failed to retrieve auth code")
			os.Exit(1)
		}
		req.Header.Set("Authorization", auth_code)
	}

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	// check status
	if resp.StatusCode == 401 {
		fmt.Println("failed to authentication. please reconnect node: " + node)
		os.Exit(1)
	} else if resp.StatusCode != 200 {
		fmt.Println("failed to request daemon.")
		os.Exit(1)
	}

	byte_array, _ := io.ReadAll(resp.Body)

	var response ResponseContainerId
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Id, response.Message
	}
	return true, response.Id, response.Message
}

func RetrieveContainerId(id string, name string) string {
	if len(id) == 0 && name == "none" {
		fmt.Println("Must specify Container ID via --id or Name via --name.")
		return ""
	}

	if len(id) != 0 {
		return id
	} else {
		// retrieve container id
		if res, resp_id, message := RequestContainerId(name); !res {
			fmt.Println(message)
			return ""
		} else {
			return resp_id
		}
	}
}
