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
	Repositry string
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
	Repositry string
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
	new_mount := strings.Replace(request_param.Mount, "/", "-", -1)
	new_command := strings.Replace(request_param.Cmd, "/", "!", -1)
	url := karakuripkgs.SERVER +
		"/container/create/" +
		request_param.Image + "/" +
		request_param.Port + "/" +
		new_mount + "/" +
		new_command + "/" +
		request_param.Repositry + "/" +
		request_param.Name + "/" +
		request_param.Namespace

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
func requestStartContainer(id string) (result bool, meessage string) {
	url := karakuripkgs.SERVER + "/container/start/" + id

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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

// func requestRunContainer(image string, port string, mount string, cmd string, repositry string) (bool, string) {
func requestRunContainer(request_param RequestRunContainer) (bool, string) {
	new_mount := strings.Replace(request_param.Mount, "/", "-", -1)
	new_command := strings.Replace(request_param.Cmd, "/", "-", -1)
	url := karakuripkgs.SERVER +
		"/container/run/" +
		request_param.Image + "/" +
		request_param.Port + "/" +
		new_mount + "/" +
		new_command + "/" +
		request_param.Repositry + "/" +
		request_param.Name + "/" +
		request_param.Namespace

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
func requestExecContainer(id string) (result bool, message string) {
	url := karakuripkgs.SERVER + "/container/exec/" + id

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/container/ls/" + namespace

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/container/spec/" + id

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/container/kill/" + id

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/container/delete/" + id

	req, _ := http.NewRequest("DELETE", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	if res, message := requestStartContainer(container_id); !res {
		fmt.Println(message)
	} else {
		// setup port forward
		config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + container_id)
		SetupPortForwarding("add", config_spec.Network)

		// execute runtime: start
		karakuripkgs.RuntimeStart(container_id, request_param.Terminal)

		// update status
		if request_param.Terminal {
			hitoha.UpdateContainerStatus(container_id, "stopped")
			// delete port forward
			SetupPortForwarding("delete", config_spec.Network)
		} else {
			hitoha.UpdateContainerStatus(container_id, "running")
			fmt.Println("container: " + container_id + " start success.")
		}
	}
}

func RunContainer(request_param RequestRunContainer) {
	res, container_id := requestRunContainer(request_param)
	if !res {
		fmt.Println("ERR: failed to run container")
	}

	// setup port forward
	config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + container_id)
	SetupPortForwarding("add", config_spec.Network)

	// execute runtime: start
	karakuripkgs.RuntimeStart(container_id, request_param.Terminal)

	if request_param.Terminal {
		hitoha.UpdateContainerStatus(container_id, "stopped")
		// delete port forward
		SetupPortForwarding("delete", config_spec.Network)
		// delete container
		if request_param.Remove {
			DeleteContainer(RequestDeleteContainer{
				Id:   container_id,
				Name: "none",
			})
		}
	} else {
		hitoha.UpdateContainerStatus(container_id, "running")
	}
}

func ExecContainer(request_param RequestExecContainer) {
	container_id := karakuripkgs.RetrieveContainerId(request_param.Id, request_param.Name)
	if container_id == "" {
		return
	}
	if res, message := requestExecContainer(container_id); !res {
		fmt.Println(message)
		return
	} else {
		// execute runtime: start
		karakuripkgs.RuntimeExec(container_id, request_param.Terminal, request_param.Cmd)
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
		// setup port forward
		config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + container_id)
		// delete port forward
		SetupPortForwarding("delete", config_spec.Network)

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
	url := karakuripkgs.SERVER + "/image/ls"

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
func requestPullImage(image_tag string, os_arch string, repositry string) (result bool, inlocal bool) {
	new_image_tag := strings.Replace(image_tag, "/", "!", -1)
	url := karakuripkgs.SERVER + "/image/pull/" + new_image_tag + "/" + os_arch + "/" + repositry

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot find " + image_tag + " from registry.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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

// delete image
func requestDeleteImage(id string) bool {
	url := karakuripkgs.SERVER + "/image/delete/" + id

	req, _ := http.NewRequest("DELETE", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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

func PullImage(image_tag string, os_arch string, repositry string) {
	fmt.Println("Pulling image, " + image_tag + " ...")
	result, is_exist := requestPullImage(image_tag, os_arch, repositry)
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

func DeleteImage(id string) {
	if !requestDeleteImage(id) {
		fmt.Println("ERR: failed to delete image")
	}
}

// ----------------------
// namespace request
func requestShowNamespace() (result bool, namespace_list hitoha.NamespaceList) {
	url := karakuripkgs.SERVER + "/namespace/ls"

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/namespace/create/" + namespace

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/namespace/delete/" + namespace

	req, _ := http.NewRequest("DELETE", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/mod/enable/" + mod_name

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/mod/disable/" + mod_name

	req, _ := http.NewRequest("DELETE", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
	url := karakuripkgs.SERVER + "/mod/list"

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Cannot connect to the Karakuri daemon. Please start the karakuri daemon.")
		os.Exit(1)
	}
	defer resp.Body.Close()
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
