package karakuri

import (
	"encoding/json"
	"fmt"
	"hitoha"
	"io"
	"karakuri_mod"
	"karakuripkgs"
	"net/http"
	"os"
	"strings"
)

type RequestCreateContainer struct {
	Image     string
	Name      string
	Namespace string
	Port      string
	Mount     string
	Cmd       string
	Registry  string
	Restart   string
}

type RequestStartContainer struct {
	Id       string
	Name     string
	Terminal bool
}

type RequestRunContainer struct {
	Name      string
	Namespace string
	Image     string
	Port      string
	Mount     string
	Terminal  bool
	Cmd       string
	Registry  string
	Restart   string
	Remove    bool
}

type RequestExecContainer struct {
	Id       string
	Name     string
	Terminal bool
	Cmd      string
}

type RequsetRestartContainer struct {
	Id       string
	Name     string
	Terminal bool
}

type RequestStopContainer struct {
	Id   string
	Name string
}

type RequestDeleteContainer struct {
	Id   string
	Name string
}

type RequestShowContainerSpec struct {
	Id   string
	Name string
}

// ----------------------
// container requests
// request create container
func requestCreateContainer(request_param RequestCreateContainer) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	new_image := strings.Replace(request_param.Image, "/", "!", -1)
	new_mount := strings.Replace(request_param.Mount, "/", "-", -1)
	new_command := strings.Replace(request_param.Cmd, "/", "!", -1)

	url := "http://" + node +
		"/container/create/" +
		new_image + "/" +
		request_param.Port + "/" +
		new_mount + "/" +
		new_command + "/" +
		request_param.Registry + "/" +
		request_param.Name + "/" +
		request_param.Namespace + "/" +
		request_param.Restart

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseContainerInfo
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

// request start container
func requestStartContainer(id, terminal string) (result bool, meessage string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}
	// Check if it is an available option
	if node != karakuripkgs.SERVER {
		if terminal == "true" {
			fmt.Println("'--it' option is not available for remote node")
			os.Exit(1)
		}
	}

	url := "http://" + node + "/container/start/" + id + "/" + terminal

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseContainerInfo
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

// func requestRunContainer(image string, port string, mount string, cmd string, registry string) (bool, string) {
func requestRunContainer(request_param RequestRunContainer, terminal string) (bool, string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}
	// Check if it is an available option
	if node != karakuripkgs.SERVER {
		if terminal == "true" {
			fmt.Println("'--it' option is not available for remote node")
			os.Exit(1)
		}
	}

	new_image := strings.Replace(request_param.Image, "/", "!", -1)
	new_mount := strings.Replace(request_param.Mount, "/", "-", -1)
	new_command := strings.Replace(request_param.Cmd, "/", "-", -1)
	url := "http://" + node +
		"/container/run/" +
		new_image + "/" +
		request_param.Port + "/" +
		new_mount + "/" +
		new_command + "/" +
		request_param.Registry + "/" +
		request_param.Name + "/" +
		request_param.Namespace + "/" +
		request_param.Restart + "/" +
		terminal

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseRunContainer
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, ""
	}
	return true, response.Id
}

// request exec container
func requestExecContainer(id, terminal, cmd string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}
	// Check if it is an available option
	if node != karakuripkgs.SERVER {
		if terminal == "true" {
			fmt.Println("'--it' option is not available for remote node")
			os.Exit(1)
		}
	}

	new_command := strings.Replace(cmd, "/", "-", -1)
	url := "http://" + node + "/container/exec/" + id + "/" + new_command + "/" + terminal

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseContainerInfo
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func requestShowContainer(namespace string) (string, hitoha.ContainerList) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/container/ls/" + namespace

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseContainerList
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return response.Result, hitoha.ContainerList{}
	}
	return response.Result, response.ContainerList
}

func requestShowContainerSpec(id string) (bool, karakuripkgs.ConfigSpec) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}
	// Check if it is an available option
	if node != karakuripkgs.SERVER {
		fmt.Println("'karakuri spec' command is not available for remote node")
		os.Exit(1)
	}

	url := "http://" + node + "/container/spec/" + id

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseContainerSpec
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, karakuripkgs.ConfigSpec{}
	}

	return true, response.Spec
}

func requestStopContainer(id string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/container/kill/" + id

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseStopContainer
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func requestDeleteContainer(id string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/container/delete/" + id

	req, _ := http.NewRequest("DELETE", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseDeleteContainer
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

// ----------------------------------

func CreateContainer(request_param RequestCreateContainer) {
	if res, message := requestCreateContainer(request_param); !res {
		fmt.Println(message)
	} else {
		fmt.Println(message)
	}
}

func StartContainer(request_param RequestStartContainer) {
	container_id := karakuripkgs.RetrieveContainerId(request_param.Id, request_param.Name)
	if container_id == "" {
		return
	}
	terminal := "false"
	if request_param.Terminal {
		terminal = "true"
	}
	if res, message := requestStartContainer(container_id, terminal); !res {
		fmt.Println(message)
	} else {
		// if terminal is true, execute from karakuri
		if request_param.Terminal {
			// execute runtime: start
			karakuripkgs.RuntimeStart(container_id, request_param.Terminal)
			hitoha.UpdateContainerStatus(container_id, "stopped")
		} else {
			fmt.Println("container: " + container_id + " start success.")
		}
	}
}

func RunContainer(request_param RequestRunContainer) {
	terminal := "false"
	if request_param.Terminal {
		terminal = "true"
	}
	res, container_id := requestRunContainer(request_param, terminal)
	if !res {
		fmt.Println("ERR: failed to run container")
	} else {
		if request_param.Terminal {
			// execute runtime: start
			karakuripkgs.RuntimeStart(container_id, request_param.Terminal)

			hitoha.UpdateContainerStatus(container_id, "stopped")
			// delete container
			if request_param.Remove {
				DeleteContainer(RequestDeleteContainer{
					Id:   container_id,
					Name: "none",
				})
			}
		}
	}
}

func ExecContainer(request_param RequestExecContainer) {
	container_id := karakuripkgs.RetrieveContainerId(request_param.Id, request_param.Name)
	if container_id == "" {
		return
	}
	terminal := "false"
	if request_param.Terminal {
		terminal = "true"
	}
	if res, message := requestExecContainer(container_id, terminal, request_param.Cmd); !res {
		fmt.Println(message)
		return
	} else {
		if request_param.Terminal {
			// execute runtime: exec
			karakuripkgs.RuntimeExec(container_id, request_param.Terminal, request_param.Cmd)
		}
	}
}

func StopContainer(request_param RequestStopContainer) {
	container_id := karakuripkgs.RetrieveContainerId(request_param.Id, request_param.Name)
	if container_id == "" {
		return
	}
	if res, message := requestStopContainer(container_id); !res {
		fmt.Println(message)
	} else {
		fmt.Println(message)
	}
}

func StopAllContaier(namespace string) {
	res, container_list := requestShowContainer(namespace)
	if res != "success" {
		fmt.Println("Namespace: \"" + namespace + "\" is not exists.")
		return
	}

	for _, entry := range container_list.List {
		StopContainer(RequestStopContainer{
			Id:   entry.Id,
			Name: "none",
		})
	}
}

func RestartContainer(request_param RequsetRestartContainer) {
	container_id := karakuripkgs.RetrieveContainerId(request_param.Id, request_param.Name)
	if container_id == "" {
		return
	}
	// stop container
	StopContainer(RequestStopContainer{
		Id:   container_id,
		Name: "none",
	})
	// start container
	StartContainer(RequestStartContainer(request_param))

	fmt.Println("container: " + container_id + " restart success")
}

func ShowContainerList(namespace string) {
	res, container_list := requestShowContainer(namespace)
	if res != "success" {
		fmt.Println("Namespace: \"" + namespace + "\" is not exists.")
		return
	}

	printContainerList(container_list, namespace)
}

func ShowContainerSpec(request_param RequestShowContainerSpec) {
	container_id := karakuripkgs.RetrieveContainerId(request_param.Id, request_param.Name)
	if container_id == "" {
		return
	}
	res, spec := requestShowContainerSpec(container_id)
	if !res {
		fmt.Println("ERR: failed to get spec")
	}

	printContainerSpec(spec)
}

func DeleteContainer(request_param RequestDeleteContainer) {
	container_id := karakuripkgs.RetrieveContainerId(request_param.Id, request_param.Name)
	if container_id == "" {
		return
	}
	if res, message := requestDeleteContainer(container_id); !res {
		fmt.Println(message)
	} else {
		fmt.Println(message)
	}
}

// ----------------------

// ----------------------
// image request
// show image
func requestShowImage() (string, hitoha.ImageList) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/image/ls"

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseImageList
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return response.Result, hitoha.ImageList{}
	}
	return response.Result, response.ImageList
}

// pull image
func requestPullImage(image_tag string, os_arch string, registry string) (result bool, inlocal bool) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	new_image_tag := strings.Replace(image_tag, "/", "!", -1)
	url := "http://" + node + "/image/pull/" + new_image_tag + "/" + os_arch + "/" + registry

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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
		fmt.Println("Cannot find " + image_tag + " from registry.")
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

	var response hitoha.ResponsePullImage
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, false
	}

	if response.ImageExists {
		return true, true
	}
	return true, false
}

// push image
func requestPushImage(image_tag string, registry string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	new_image_tag := strings.Replace(image_tag, "/", "!", -1)
	url := "http://" + node + "/image/push/" + new_image_tag + "/" + registry

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponsePushImage
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}

	return true, response.Message
}

// delete image
func requestDeleteImage(id string) bool {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/image/delete/" + id

	req, _ := http.NewRequest("DELETE", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseDeleteImage
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false
	}
	return true
}

func ShowImage() {
	res, image_list := requestShowImage()
	if res != "success" {
		fmt.Println("ERR: failed to get image list")
		return
	}

	printImageList(image_list)
}

func PullImage(image_tag string, os_arch string, registry string) {
	fmt.Println("Pulling image, " + image_tag + " ...")
	result, is_exist := requestPullImage(image_tag, os_arch, registry)
	if !result {
		fmt.Println("ERR: failed to pull image")
		return
	}
	if is_exist {
		fmt.Println("\"" + image_tag + "\" is already exists in local")
	} else {
		fmt.Println("Pull \"" + image_tag + "\" completed")
	}
}

func PushImage(image_tag string, registry string) {
	fmt.Println("Pushing image, " + image_tag + " ...")
	result, message := requestPushImage(image_tag, registry)
	if !result {
		fmt.Println(message)
		return
	}
	fmt.Println("Push completed")
}

func DeleteImage(id string) {
	if !requestDeleteImage(id) {
		fmt.Println("ERR: failed to delete image")
	}
}

// ----------------------
// namespace request
func requestShowNamespace() (result bool, namespace_list hitoha.NamespaceList) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/namespace/ls"

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseNamespaceList
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, hitoha.NamespaceList{}
	}
	return true, response.Namespace
}

func requestCreateNamespace(namespace string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/namespace/create/" + namespace

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseCreateNamespace
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func requestDeleteNamespace(namespace string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/namespace/delete/" + namespace

	req, _ := http.NewRequest("DELETE", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseDeleteNamespace
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func ShowNamespace() {
	res, namespace_list := requestShowNamespace()
	if !res {
		fmt.Println("ERR: failed to get namespace list")
		return
	}
	printNamespaceList(namespace_list)
}

func CreateNamespace(namespace string) {
	_, message := requestCreateNamespace(namespace)
	fmt.Println(message)
}

func DeleteNamespace(namespace string) {
	_, message := requestDeleteNamespace(namespace)
	fmt.Println(message)
}

// ----------------------
// module request
func requestEnableModule(mod_name string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/mod/enable/" + mod_name

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response karakuri_mod.ResponseEnableModule
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func requestDisableModule(mod_name string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/mod/disable/" + mod_name

	req, _ := http.NewRequest("DELETE", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response karakuri_mod.ResponseEnableModule
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func requestShowModule() (result bool, mod_list karakuri_mod.ModList) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/mod/list"

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response karakuri_mod.ResponseModuleList
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.List
	}
	return true, response.List
}

func EnableModule(mod_name string) {
	_, message := requestEnableModule(mod_name)
	fmt.Println(message)
}

func DisableModule(mod_name string) {
	_, message := requestDisableModule(mod_name)
	fmt.Println(message)
}

func ShowModuleList() {
	res, module_list := requestShowModule()
	if !res {
		fmt.Println("Failed to retrieve module list")
		return
	}
	printModuleList(module_list)
}

// ----------------------
// registry controller
// connect registry
func requestConnectRegistry(registry string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/reg/connect/" + registry

	req, _ := http.NewRequest("POST", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseConnectRegistry
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func requestTargetRegistry() (result bool, registry_info hitoha.RegistryInfo) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/reg/target"

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseGetTargetRegistry
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.RegistryInfo
	}
	return true, response.RegistryInfo
}

func requestGetRepository() (result bool, message string, repository_list hitoha.RepogitryList) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/reg/repository"

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseShowRepository
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message, response.Repository
	}
	return true, response.Message, response.Repository
}

func requestGetTag(repository string) (result bool, message string, tag_list hitoha.TagList) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	new_repository := strings.Replace(repository, "/", "!", -1)
	url := "http://" + node + "/reg/tag/" + new_repository

	req, _ := http.NewRequest("GET", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseShowTag
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message, response.Tag
	}
	return true, response.Message, response.Tag
}

func requestDisconnectRegistry() (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	url := "http://" + node + "/reg/disconnect"

	req, _ := http.NewRequest("DELETE", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseDisconnectRegistry
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func requestDeleteImageManifest(image_tag string) (result bool, message string) {
	// node
	node_info := getTargetNode()
	node := node_info.Target
	if node_info.Status != "connected" {
		fmt.Println("node: " + node + " is not connected.")
		os.Exit(1)
	}

	new_image_tag := strings.Replace(image_tag, "/", "!", -1)
	url := "http://" + node + "/reg/delete/" + new_image_tag

	req, _ := http.NewRequest("DELETE", url, nil)
	// set auth_code
	if node != karakuripkgs.SERVER {
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

	var response hitoha.ResponseDeleteManifest
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	if response.Result != "success" {
		return false, response.Message
	}
	return true, response.Message
}

func ConnectRegistry(registry string) {
	_, message := requestConnectRegistry(registry)
	fmt.Println(message)
}

func DisconnectRegistry() {
	_, message := requestDisconnectRegistry()
	fmt.Println(message)
}

func ShowTargetRegistry() {
	_, registry_info := requestTargetRegistry()
	printTargetRegistry(registry_info)
}

func ShowRepository() {
	res, message, repository_list := requestGetRepository()
	if !res {
		fmt.Println(message)
		return
	}
	printRepository(repository_list)
}

func ShowTag(repository string) {
	res, message, tag_list := requestGetTag(repository)
	if !res {
		fmt.Println(message)
		return
	}
	printTag(repository, tag_list)
}

func DeleteImageManifest(image_tag string) {
	res, message := requestDeleteImageManifest(image_tag)
	if !res {
		fmt.Println(message)
		return
	}
	fmt.Println(message)
}
