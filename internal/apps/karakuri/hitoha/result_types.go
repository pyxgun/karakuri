package hitoha

import "karakuripkgs"

type ResponseContainerInfo struct {
	Result    string        `json:"result"`
	Message   string        `json:"message"`
	Container ContainerInfo `json:"container"`
}

type ResponseContainerList struct {
	Result        string        `json:"result"`
	ContainerList ContainerList `json:"containers"`
}

type ResponseContainerId struct {
	Result  string `json:"result"`
	Id      string `json:"id"`
	Message string `json:"message"`
}

type ResponseDeleteContainer struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

type ResponseImageList struct {
	Result    string    `json:"result"`
	ImageList ImageList `json:"images"`
}

type ResponseRunContainer struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

type TargetImage struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Os         string `json:"os"`
	Arch       string `json:"arch"`
}

type ResponsePullImage struct {
	Result      string      `json:"result"`
	ImageExists bool        `json:"inlocal"`
	Image       TargetImage `json:"image"`
}

type ResponsePushImage struct {
	Result  string `json:"result"`
	Image   string `json:"image"`
	Message string `json:"message"`
}

type ResponseStopContainer struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

type ResponseDeleteImage struct {
	Result string `json:"result"`
	Id     string `json:"id"`
}

type ResponseContainerSpec struct {
	Result string                  `json:"result"`
	Spec   karakuripkgs.ConfigSpec `json:"spec"`
}

type ResponseNamespaceList struct {
	Result    string        `json:"result"`
	Namespace NamespaceList `json:"namespace"`
}

type ResponseCreateNamespace struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type ResponseDeleteNamespace struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func createResponseContainerInfo(result string, container_info ContainerInfo, message string) ResponseContainerInfo {
	resp := ResponseContainerInfo{
		Result:    result,
		Message:   message,
		Container: container_info,
	}
	return resp
}

func createResponseContainerList(result string, container_list ContainerList) ResponseContainerList {
	resp := ResponseContainerList{
		Result:        result,
		ContainerList: container_list,
	}
	return resp
}

func createResponseContainerId(result string, id string, message string) ResponseContainerId {
	resp := ResponseContainerId{
		Result:  result,
		Message: message,
		Id:      id,
	}
	return resp
}

func createResponseDeleteContainer(result string, id string, message string) ResponseDeleteContainer {
	resp := ResponseDeleteContainer{
		Result:  result,
		Message: message,
		Id:      id,
	}
	return resp
}

func createResponseStopContainer(result string, id string, message string) ResponseStopContainer {
	resp := ResponseStopContainer{
		Result:  result,
		Message: message,
		Id:      id,
	}
	return resp
}

func createResponseRunContainer(result string, id string, message string) ResponseRunContainer {
	resp := ResponseRunContainer{
		Result:  result,
		Message: message,
		Id:      id,
	}
	return resp
}

func createResponseSpecContainer(result string, spec karakuripkgs.ConfigSpec) ResponseContainerSpec {
	resp := ResponseContainerSpec{
		Result: result,
		Spec:   spec,
	}
	return resp
}

func createResponseImageList(result string, image_list ImageList) ResponseImageList {
	resp := ResponseImageList{
		Result:    result,
		ImageList: image_list,
	}
	return resp
}

func createResponsePullImage(result string, inlocal bool, image string, tag string, os string, arch string) ResponsePullImage {
	resp := ResponsePullImage{
		Result:      result,
		ImageExists: inlocal,
		Image: TargetImage{
			Repository: image,
			Tag:        tag,
			Os:         os,
			Arch:       arch,
		},
	}
	return resp
}

func createResponsePushImage(result string, image string, message string) ResponsePushImage {
	resp := ResponsePushImage{
		Result:  result,
		Image:   image,
		Message: message,
	}
	return resp
}

func createResponseDeleteImage(result string, id string) ResponseDeleteImage {
	resp := ResponseDeleteImage{
		Result: result,
		Id:     id,
	}
	return resp
}

func createResponseNamespaceList(result string, namespace_list NamespaceList) ResponseNamespaceList {
	resp := ResponseNamespaceList{
		Result:    result,
		Namespace: namespace_list,
	}
	return resp
}

func createResponseCreateNamespace(result string, message string) ResponseCreateNamespace {
	resp := ResponseCreateNamespace{
		Result:  result,
		Message: message,
	}
	return resp
}

func createResponseDeleteNamespace(result string, message string) ResponseDeleteNamespace {
	resp := ResponseDeleteNamespace{
		Result:  result,
		Message: message,
	}
	return resp
}

// registry controller
type ResponseGetTargetRegistry struct {
	Result       string       `json:"result"`
	RegistryInfo RegistryInfo `json:"registry_info"`
}

type ResponseConnectRegistry struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type ResponseDisconnectRegistry struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type ResponseShowRepository struct {
	Result     string        `json:"result"`
	Message    string        `json:"message"`
	Repository RepogitryList `json:"repositories"`
}

type ResponseShowTag struct {
	Result  string  `json:"result"`
	Message string  `json:"message"`
	Tag     TagList `json:"tag"`
}

type ResponseDeleteManifest struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func createResponseGetTargetRegistry(result string, registry_info RegistryInfo) ResponseGetTargetRegistry {
	resp := ResponseGetTargetRegistry{
		Result:       result,
		RegistryInfo: registry_info,
	}
	return resp
}

func createResponseConnectRegistry(result, message string) ResponseConnectRegistry {
	resp := ResponseConnectRegistry{
		Result:  result,
		Message: message,
	}
	return resp
}

func createResponseDisconnectRegistry(result, message string) ResponseDisconnectRegistry {
	resp := ResponseDisconnectRegistry{
		Result:  result,
		Message: message,
	}
	return resp
}

func createResponseShowRepository(result, message string, repository_list RepogitryList) ResponseShowRepository {
	resp := ResponseShowRepository{
		Result:     result,
		Message:    message,
		Repository: repository_list,
	}
	return resp
}

func createResponseShowTag(result, message string, tag_list TagList) ResponseShowTag {
	resp := ResponseShowTag{
		Result:  result,
		Message: message,
		Tag:     tag_list,
	}
	return resp
}

func createResponseDeleteManifest(result, message string) ResponseDeleteManifest {
	resp := ResponseDeleteManifest{
		Result:  result,
		Message: message,
	}
	return resp
}
