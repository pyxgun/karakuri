# Karakuri Workflow
## Container creation process
```mermaid
sequenceDiagram
   participant karakuri
   participant hitoha
   participant futaba
   participant /etc/karakuri
   participant Registry

   Note over karakuri, hitoha: REST API

   karakuri->>+hitoha:http://localhost:9806/create
   alt image not exist in local
      hitoha->>Registry:pull image
      Registry-->>hitoha:Image Layer
      hitoha->>/etc/karakuri:add image list (/hitoha/images/image_list.json)
   end
   hitoha->>/etc/karakuri:create new container layer (/futaba/[container_id])
   hitoha->>futaba:call "futaba spec"
   futaba->>/etc/karakuri:create SpecFile (/futaba/[container_id]/config.json)
   futaba-->>hitoha: create result: success
   hitoha-->>-karakuri: 200 OK, {"result": "success"}
```
