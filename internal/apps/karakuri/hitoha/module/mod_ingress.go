package karakuri_mod

import (
	"karakuripkgs"
	"os"
	"os/exec"
	"strings"
)

func createIngressModuleKarakurifile() {
	// create Karakurifile
	karakurifile := `FROM nginx:alpine

COPY /etc/karakuri/modules/ingress/script /script
COPY /etc/karakuri/modules/ingress/cert /cert
COPY /etc/karakuri/modules/ingress/conf /conf

RUN apk update
RUN apk add openssl
RUN openssl genrsa 2048 > /cert/server.key
RUN openssl req -new -key /cert/server.key -config /cert/ssl.cnf -out /cert/server.csr
RUN openssl x509 -req -days 3650 -in /cert/server.csr -signkey /cert/server.key -out /cert/server.crt -extfile /cert/san.txt
RUN cp /conf/server.conf /etc/nginx/conf.d/

CMD ["sh", "/script/entrypoint.sh"]
`

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/Karakurifile"); stat != nil {
		fd, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/Karakurifile")
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		bytes := []byte(karakurifile)
		if _, err := fd.Write(bytes); err != nil {
			panic(err)
		}
	}
}

func createIngressModuleConfig() {
	// conf
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/conf"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS+"/conf", os.ModePerm); err != nil {
			panic(err)
		}
	}

	server_conf := `server{
    server_name    ingress.karakuri.container;

    listen 443 ssl;
    ssl_certificate /cert/server.crt;
    ssl_certificate_key /cert/server.key;

    client_max_body_size 5G;
    keepalive_timeout  130;
    send_timeout 130;
    client_body_timeout 130;
    client_header_timeout 130;
    proxy_send_timeout 130;
    proxy_read_timeout 130;
}
`
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/conf/server.conf"); stat != nil {
		fd_1, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/conf/server.conf")
		if err != nil {
			panic(err)
		}
		defer fd_1.Close()
		bytes_1 := []byte(server_conf)
		if _, err := fd_1.Write(bytes_1); err != nil {
			panic(err)
		}
	}

	// cert
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/cert"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS+"/cert", os.ModePerm); err != nil {
			panic(err)
		}
	}

	ssl_conf := `[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext
prompt = no
[ req_distinguished_name ]
C = JP
CN = karakuri ingress
[ req_ext ]
subjectAltName = @alt_names
[ alt_names ]
DNS.1 = *.karakuri.container
`
	san_txt := "subjectAltName = DNS:*.karakuri.container"

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/cert/ssl.cnf"); stat != nil {
		fd_2, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/cert/ssl.cnf")
		if err != nil {
			panic(err)
		}
		defer fd_2.Close()
		bytes_2 := []byte(ssl_conf)
		if _, err := fd_2.Write(bytes_2); err != nil {
			panic(err)
		}
	}
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/cert/san.txt"); stat != nil {
		fd_3, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/cert/san.txt")
		if err != nil {
			panic(err)
		}
		defer fd_3.Close()
		bytes_3 := []byte(san_txt)
		if _, err := fd_3.Write(bytes_3); err != nil {
			panic(err)
		}
	}

	// script
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/script"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS+"/script", os.ModePerm); err != nil {
			panic(err)
		}
	}
	entrypoint := `#!/bin/sh

/docker-entrypoint.sh nginx -g "daemon off;"
`
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/script/entrypoint.sh"); stat != nil {
		fd_4, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/script/entrypoint.sh")
		if err != nil {
			panic(err)
		}
		defer fd_4.Close()
		bytes_4 := []byte(entrypoint)
		if _, err := fd_4.Write(bytes_4); err != nil {
			panic(err)
		}
	}
}

func setupIngressModule() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS, os.ModePerm); err != nil {
			panic(err)
		}
	}
	// Karakurifile
	createIngressModuleKarakurifile()
	// ingress config
	createIngressModuleConfig()
}

func buildIngressImage(mod_info ModInfo) {
	build_args := []string{"build", "--name", mod_info.ImageName, "--buildpath", mod_info.Path}
	build := exec.Command("karakuri", build_args...)
	if err := build.Run(); err != nil {
		panic(err)
	}
}

func createIngressContainer(mod_info ModInfo) {
	create_args := []string{
		"create",
		"--name", mod_info.Name,
		"--image", mod_info.ImageName,
		"--port", "4443:443:tcp",
		"--restart", "on-boot",
		"--ns", "system-mod",
	}
	create := exec.Command("karakuri", create_args...)
	if err := create.Run(); err != nil {
		panic(err)
	}
}

func startIngressContainer(mod_info ModInfo) {
	start_args := []string{"start", "--name", mod_info.Name}
	start := exec.Command("karakuri", start_args...)
	if err := start.Run(); err != nil {
		panic(err)
	}
}

func enableIngressModule(mod_info ModInfo) {
	image_info := strings.Split(mod_info.ImageName, ":")
	if !isImageExists(image_info[0], image_info[1]) {
		buildIngressImage(mod_info)
	}
	createIngressContainer(mod_info)
	startIngressContainer(mod_info)
}

func disableIngressModule() {
	// remove container
	removeModuleContainer("ingress")
}

func EditIngressEntry() {
	edit_args := []string{"exec", "--it", "--name", "ingress", "--cmd", "vi,/etc/nginx/conf.d/server.conf"}
	edit := exec.Command("karakuri", edit_args...)

	edit.Stdin = os.Stdin
	edit.Stdout = os.Stdout
	edit.Stderr = os.Stderr

	if err := edit.Start(); err != nil {
		panic(err)
	}
	edit.Wait()
	// restart
	restart_args := []string{"restart", "--name", "ingress"}
	restart := exec.Command("karakuri", restart_args...)
	if err := restart.Run(); err != nil {
		panic(err)
	}
}
