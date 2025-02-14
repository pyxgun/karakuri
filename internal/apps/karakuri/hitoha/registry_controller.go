package hitoha

import (
	"encoding/json"
	"io"
	"karakuripkgs"
	"net/http"
	"os"
	"sort"
	"strings"
)

type RegistryInfo struct {
	Target string `json:"target"`
	Status string `json:"connection_status"`
}

func checkTargetRegistryFile() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_REGCTL_ROOT); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_REGCTL_ROOT, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(karakuripkgs.KARAKURI_REGCTL_REGINFO); stat != nil {
		var registry_info RegistryInfo
		data, _ := json.MarshalIndent(registry_info, "", "  ")
		if err := os.WriteFile(karakuripkgs.KARAKURI_REGCTL_REGINFO, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func setTargetRegistry(registry string) {
	checkTargetRegistryFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_REGCTL_REGINFO)
	if err != nil {
		panic(err)
	}

	var registry_info RegistryInfo
	if err := json.Unmarshal(bytes, &registry_info); err != nil {
		panic(err)
	}

	// set registry
	registry_info.Target = registry
	registry_info.Status = "disconnected"

	data, _ := json.MarshalIndent(registry_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_REGCTL_REGINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func setRegistryStatus(status string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_REGCTL_REGINFO)
	if err != nil {
		panic(err)
	}

	var registry_info RegistryInfo
	if err := json.Unmarshal(bytes, &registry_info); err != nil {
		panic(err)
	}

	// set registry
	registry_info.Status = status

	data, _ := json.MarshalIndent(registry_info, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_REGCTL_REGINFO, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func getTargetRegistry() RegistryInfo {
	checkTargetRegistryFile()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_REGCTL_REGINFO)
	if err != nil {
		panic(err)
	}

	var registry_info RegistryInfo
	if err := json.Unmarshal(bytes, &registry_info); err != nil {
		panic(err)
	}

	return registry_info
}

func verifyConnectionToRegistry() (string, bool) {
	registry_info := getTargetRegistry()
	registry := registry_info.Target
	if registry == "" {
		return "failed to get target registry", false
	}
	// connect test
	// request retrieve catalog
	url := "http://" + registry + "/v2/_catalog"
	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		setRegistryStatus("connectoin_failed")
		return "registry connection failed", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		setRegistryStatus("connectoin_failed")
		return "registry connection failed", false
	}

	setRegistryStatus("connected")
	return "registry: " + registry + " connection success", true
}

func checkConnectionStatus() bool {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_REGCTL_REGINFO)
	if err != nil {
		panic(err)
	}

	var registry_info RegistryInfo
	if err := json.Unmarshal(bytes, &registry_info); err != nil {
		panic(err)
	}

	if registry_info.Status != "connected" {
		return false
	}
	return true
}

// Registry API
// get repogitry list
type RepogitryList struct {
	Repository []string `json:"repositories"`
}

func getRepositoryList() (RepogitryList, string, bool) {
	registry_info := getTargetRegistry()
	registry := registry_info.Target
	if registry == "" {
		return RepogitryList{}, "failed to get target registry", false
	}
	// status check
	if !checkConnectionStatus() {
		return RepogitryList{}, "still not connected any registry.\nplease execute `karakuri regctl connect --registry {REGISTRY}` first.", false
	}

	// request retrieve repository
	url := "http://" + registry + "/v2/_catalog"
	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		return RepogitryList{}, "registry connection failed", false
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return RepogitryList{}, "failed to get repository list", false
	}

	var response RepogitryList
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	// sort a-z
	array_slice := response.Repository[:]
	sort.Strings(array_slice)
	response.Repository = array_slice

	return response, "success", true
}

// get tag list
type TagList struct {
	Name string   `json:"name"`
	Tag  []string `json:"tags"`
}

func getTagList(repogitory string) (TagList, string, bool) {
	registry_info := getTargetRegistry()
	registry := registry_info.Target
	if registry == "" {
		return TagList{}, "failed to get target registry", false
	}
	// status check
	if !checkConnectionStatus() {
		return TagList{}, "still not connected any registry.\nplease execute `karakuri regctl connect --registry {REGISTRY}` first.", false
	}

	// request retrieve repository
	url := "http://" + registry + "/v2/" + repogitory + "/tags/list"
	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		return TagList{}, "registry connection failed", false
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return TagList{}, "no such repository: " + repogitory, false
	}

	var response TagList
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	// sort a-z
	array_slice := response.Tag[:]
	sort.Strings(array_slice)
	response.Tag = array_slice

	return response, "success", true
}

func getManifest(image, tag string) (string, bool) {
	registry_info := getTargetRegistry()
	registry := registry_info.Target
	if registry == "" {
		return "failed to get target registry", false
	}
	// status check
	if !checkConnectionStatus() {
		return "still not connected any registry.\nplease execute `karakuri regctl connect --registry {REGISTRY}` first.", false
	}

	// request retrieve manifest
	url := "http://" + registry + "/v2/" + image + "/manifests/" + tag
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		return "no such image: " + image + ":" + tag, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "no such image: " + image + ":" + tag, false
	}

	return resp.Header.Get("Docker-Content-Digest"), true
}

func deleteManifest(image, digest string) (string, bool) {
	registry_info := getTargetRegistry()
	registry := registry_info.Target
	if registry == "" {
		return "failed to get target registry", false
	}
	// status check
	if !checkConnectionStatus() {
		return "still not connected any registry.\nplease execute `karakuri regctl connect --registry {REGISTRY}` first.", false
	}

	// request retrieve manifest
	url := "http://" + registry + "/v2/" + image + "/manifests/" + digest
	req, _ := http.NewRequest("DELETE", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		return "failed to delete image: " + image, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return "no such image: " + image, false
	}
	return "delete success", true
}

// -- called from endpoint.go
func ShowTargetRegistry() ResponseGetTargetRegistry {
	registry_info := getTargetRegistry()
	return createResponseGetTargetRegistry("success", registry_info)
}

func ConnectRegistry(registry string) ResponseConnectRegistry {
	if checkConnectionStatus() {
		registry_info := getTargetRegistry()
		registry := registry_info.Target
		return createResponseConnectRegistry("error", "already connected to registry: "+registry+".\nplease execute `karakuri regctl disconnect` before change connection registry.")
	}
	// set env
	setTargetRegistry(registry)
	// connection test
	message, res := verifyConnectionToRegistry()
	if !res {
		return createResponseConnectRegistry("error", message)
	}
	return createResponseConnectRegistry("success", message)
}

func DisconnectRegistry() ResponseDisconnectRegistry {
	if !checkConnectionStatus() {
		return createResponseDisconnectRegistry("error", "still not connected any registry.\nplease execute `karakuri regctl connect --registry {REGISTRY}` first.")
	}
	setRegistryStatus("disconnected")
	return createResponseDisconnectRegistry("success", "registry dissconnected.")
}

func ShowRepositories() ResponseShowRepository {
	// retrieve repository list
	repository_list, message, res := getRepositoryList()
	if !res {
		return createResponseShowRepository("error", message, RepogitryList{})
	}
	return createResponseShowRepository("success", message, repository_list)
}

func ShowTags(repository string) ResponseShowTag {
	// retrieve repository list
	tag_list, message, res := getTagList(repository)
	if !res {
		return createResponseShowTag("error", message, TagList{})
	}
	return createResponseShowTag("success", message, tag_list)
}

func DeleteImageManifest(image_tag string) ResponseDeleteManifest {
	// parse image:tag
	image_tag_info := strings.Split(image_tag, ":")
	image := image_tag_info[0]
	tag := "latest"
	if len(image_tag_info) == 2 {
		tag = image_tag_info[1]
	}

	// get content-digest
	content_digest, res := getManifest(image, tag)
	if !res {
		return createResponseDeleteManifest("error", content_digest)
	}
	message, delete_result := deleteManifest(image, content_digest)
	if !delete_result {
		return createResponseDeleteManifest("error", message)
	}
	return createResponseDeleteManifest("success", "delete "+image+":"+tag+" success.")
}
