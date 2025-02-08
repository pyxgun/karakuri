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

// request retrieve container id
func RequestContainerId(name string) (result bool, id string, message string) {
	url := SERVER + "/container/getid/" + name

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
