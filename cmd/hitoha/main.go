package main

import (
	"encoding/json"
	"fmt"
	"hitoha"
	"karakuripkgs"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type AuthCode struct {
	AuthCode string `json:"auth_code"`
}

// retrieve authentication code
func getAuthCode() string {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_NODECTL_AUTHCODE)
	if err != nil {
		return ""
	}

	var auth_code AuthCode
	if err := json.Unmarshal(bytes, &auth_code); err != nil {
		panic(err)
	}
	return auth_code.AuthCode
}

// authentication middle-ware
func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate token
		token := r.Header.Get("Authorization")

		if token != getAuthCode() {
			http.Error(w, "Authentication Failed", http.StatusUnauthorized)
			return
		}

		// call next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("karakuri version " + karakuripkgs.KARAKURI_VERSION)
	// initial setup
	hitoha.SetupEnvironment()

	// local router
	router := mux.NewRouter()
	// Container
	// GET
	router.HandleFunc("/container/ls/{namespace}", hitoha.GetContainerList).Methods("GET")
	router.HandleFunc("/container/spec/{id}", hitoha.GetContainerSpec).Methods("GET")
	router.HandleFunc("/container/getid/{name}", hitoha.GetContainerId).Methods("GET")
	// POST
	router.HandleFunc("/container/create/{image}/{port}/{mount}/{cmd}/{registry}/{name}/{namespace}/{restart}", hitoha.PostCreateContainer).Methods("POST")
	router.HandleFunc("/container/start/{id}/{terminal}", hitoha.PostStartContainer).Methods("POST")
	router.HandleFunc("/container/run/{image}/{port}/{mount}/{cmd}/{registry}/{name}/{namespace}/{restart}/{terminal}", hitoha.PostRunContainer).Methods("POST")
	router.HandleFunc("/container/exec/{id}/{cmd}/{terminal}", hitoha.PostExecContainer).Methods("POST")
	router.HandleFunc("/container/kill/{id}", hitoha.PostKillContainer).Methods("POST")
	// DELETE
	router.HandleFunc("/container/delete/{id}", hitoha.DeleteDeleteContainer).Methods("DELETE")

	// Image
	// GET
	router.HandleFunc("/image/ls", hitoha.GetShowImages).Methods("GET")
	router.HandleFunc("/image/pull/{image-tag}/{os-arch}/{registry}", hitoha.GetPullImage).Methods("GET")
	// POST
	router.HandleFunc("/image/push/{image-tag}/{registry}", hitoha.PostPushImage).Methods("POST")
	// DELETE
	router.HandleFunc("/image/delete/{id}", hitoha.DeleteDeleteImage).Methods("DELETE")

	// namespcae
	// GET
	router.HandleFunc("/namespace/ls", hitoha.GetNamespaceList).Methods("GET")
	// POST
	router.HandleFunc("/namespace/create/{namespace}", hitoha.PostNamespace).Methods("POST")
	// DELETE
	router.HandleFunc("/namespace/delete/{namespace}", hitoha.DeleteNamespace).Methods("DELETE")

	// module
	// GET
	router.HandleFunc("/mod/list", hitoha.GetModuleList).Methods("GET")
	// POST
	router.HandleFunc("/mod/enable/{mod_name}", hitoha.PostEnableModule).Methods("POST")
	// DELETE
	router.HandleFunc("/mod/disable/{mod_name}", hitoha.DeleteDisableModule).Methods("DELETE")

	// registry controller
	// GET
	router.HandleFunc("/reg/target", hitoha.GetTargetRegistry).Methods("GET")
	router.HandleFunc("/reg/repository", hitoha.GetShowRepositories).Methods("GET")
	router.HandleFunc("/reg/tag/{repository}", hitoha.GetShowTags).Methods("GET")
	// POST
	router.HandleFunc("/reg/connect/{registry}", hitoha.PostConnectRegistry).Methods("POST")
	// DELETE
	router.HandleFunc("/reg/disconnect", hitoha.DeleteDisconnectRegistry).Methods("DELETE")
	router.HandleFunc("/reg/delete/{image-tag}", hitoha.DeleteDeleteManifest).Methods("DELETE")

	// remote router
	remote_router := mux.NewRouter()
	// Container
	// GET
	remote_router.Handle("/container/ls/{namespace}", authMiddleWare(http.HandlerFunc(hitoha.GetContainerList))).Methods("GET")
	remote_router.Handle("/container/spec/{id}", authMiddleWare(http.HandlerFunc(hitoha.GetContainerSpec))).Methods("GET")
	remote_router.Handle("/container/getid/{name}", authMiddleWare(http.HandlerFunc(hitoha.GetContainerId))).Methods("GET")
	// POST
	remote_router.Handle("/container/create/{image}/{port}/{mount}/{cmd}/{registry}/{name}/{namespace}/{restart}", authMiddleWare(http.HandlerFunc(hitoha.PostCreateContainer))).Methods("POST")
	remote_router.Handle("/container/start/{id}/{terminal}", authMiddleWare(http.HandlerFunc(hitoha.PostStartContainer))).Methods("POST")
	remote_router.Handle("/container/run/{image}/{port}/{mount}/{cmd}/{registry}/{name}/{namespace}/{restart}/{terminal}", authMiddleWare(http.HandlerFunc(hitoha.PostRunContainer))).Methods("POST")
	remote_router.Handle("/container/exec/{id}/{cmd}/{terminal}", authMiddleWare(http.HandlerFunc(hitoha.PostExecContainer))).Methods("POST")
	remote_router.Handle("/container/kill/{id}", authMiddleWare(http.HandlerFunc(hitoha.PostKillContainer))).Methods("POST")
	// DELETE
	remote_router.Handle("/container/delete/{id}", authMiddleWare(http.HandlerFunc(hitoha.DeleteDeleteContainer))).Methods("DELETE")

	// Image
	// GET
	remote_router.Handle("/image/ls", authMiddleWare(http.HandlerFunc(hitoha.GetShowImages))).Methods("GET")
	remote_router.Handle("/image/pull/{image-tag}/{os-arch}/{registry}", authMiddleWare(http.HandlerFunc(hitoha.GetPullImage))).Methods("GET")
	// POST
	remote_router.Handle("/image/push/{image-tag}/{registry}", authMiddleWare(http.HandlerFunc(hitoha.PostPushImage))).Methods("POST")
	// DELETE
	remote_router.Handle("/image/delete/{id}", authMiddleWare(http.HandlerFunc(hitoha.DeleteDeleteImage))).Methods("DELETE")

	// namespcae
	// GET
	remote_router.Handle("/namespace/ls", authMiddleWare(http.HandlerFunc(hitoha.GetNamespaceList))).Methods("GET")
	// POST
	remote_router.Handle("/namespace/create/{namespace}", authMiddleWare(http.HandlerFunc(hitoha.PostNamespace))).Methods("POST")
	// DELETE
	remote_router.Handle("/namespace/delete/{namespace}", authMiddleWare(http.HandlerFunc(hitoha.DeleteNamespace))).Methods("DELETE")

	// module
	// GET
	remote_router.Handle("/mod/list", authMiddleWare(http.HandlerFunc(hitoha.GetModuleList))).Methods("GET")
	// POST
	remote_router.Handle("/mod/enable/{mod_name}", authMiddleWare(http.HandlerFunc(hitoha.PostEnableModule))).Methods("POST")
	// DELETE
	remote_router.Handle("/mod/disable/{mod_name}", authMiddleWare(http.HandlerFunc(hitoha.DeleteDisableModule))).Methods("DELETE")

	// registry controller
	// GET
	remote_router.Handle("/reg/target", authMiddleWare(http.HandlerFunc(hitoha.GetTargetRegistry))).Methods("GET")
	remote_router.Handle("/reg/repository", authMiddleWare(http.HandlerFunc(hitoha.GetShowRepositories))).Methods("GET")
	remote_router.Handle("/reg/tag/{repository}", authMiddleWare(http.HandlerFunc(hitoha.GetShowTags))).Methods("GET")
	// POST
	remote_router.Handle("/reg/connect/{registry}", authMiddleWare(http.HandlerFunc(hitoha.PostConnectRegistry))).Methods("POST")
	// DELETE
	remote_router.Handle("/reg/disconnect", authMiddleWare(http.HandlerFunc(hitoha.DeleteDisconnectRegistry))).Methods("DELETE")
	remote_router.Handle("/reg/delete/{image-tag}", authMiddleWare(http.HandlerFunc(hitoha.DeleteDeleteManifest))).Methods("DELETE")

	// execute server
	// access from local controller
	go func() {
		fmt.Println("Listen on \"127.0.0.1:9806\" ...")
		http.ListenAndServe("127.0.0.1:9806", router)
	}()

	// access from external controller
	// port 9816 allows access from all addresses.
	// However, when running in standalone mode
	// access is denied by packet filters (iptables).
	// In remote control mode, authentication by Authorization header is required for access.
	go func() {
		fmt.Println("Listen on \"0.0.0.0:9816\" ...")
		http.ListenAndServe("0.0.0.0:9816", remote_router)
	}()

	select {}
}
