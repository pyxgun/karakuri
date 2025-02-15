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

type ClusterInfo struct {
	Mode   string `json:"mode"`
	Target string `json:"target"`
	Status string `json:"connection_status"`
}

func checkTargetClusterFile() {
	if _, stat := os.Stat(KARAKURI_CLSCTL_ROOT); stat != nil {
		if err := os.MkdirAll(KARAKURI_CLSCTL_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(KARAKURI_CLSCTL_CLSINFO); stat != nil {
		var cluster_info ClusterInfo
		cluster_info.Mode = "stand-alone"
		cluster_info.Target = SERVER
		cluster_info.Status = "connected"
		data, _ := json.MarshalIndent(cluster_info, "", "  ")
		if err := os.WriteFile(KARAKURI_CLSCTL_CLSINFO, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func getTargetCluster() ClusterInfo {
	checkTargetClusterFile()

	var bytes []byte
	bytes, err := os.ReadFile(KARAKURI_CLSCTL_CLSINFO)
	if err != nil {
		panic(err)
	}

	var cluster_info ClusterInfo
	if err := json.Unmarshal(bytes, &cluster_info); err != nil {
		panic(err)
	}

	return cluster_info
}

// request retrieve container id
func RequestContainerId(name string) (result bool, id string, message string) {
	cluster_info := getTargetCluster()
	cluster := cluster_info.Target
	if cluster_info.Status != "connected" {
		fmt.Println("cluster: " + cluster + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + cluster + "/container/getid/" + name

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
