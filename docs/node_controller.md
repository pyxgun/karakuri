# Node Controller
`karakuri` implements Node Controller for remote management.  
It is possible to connect to `karakuri` running on other node and manage containers remotely.  
By using Node Controller, you can operated separately as a node for controllers and a node for workers.

## Change mode
`karakuri` operates in stand-alone mode by default.  
This is the mode in which the system is not managed remotely, only locally.  

You can check which mode you are currently running in with the following command.
```
$ sudo karakuri nodectl info

Mode : stand-alone
```

On the node that you want to manage remotely, execute the following command to change remote-control mode.
```
$ sudo karakuri nodectl mode remote-control
```
If the remote-control mode is successfully activated, the following message is displayed.
```
currently running in remote-control mode
now you can connect from remote to this node. please execute the following command on controller node:
  karakuri nodectl connect --node [NODE_ADDRESS] --auth [AUTHCODE]
```
An authorization code is required to connect to a system operating in remote control mode.
Copy the displayed command and proceed next step.

(Note) To switch to stand-alone mode, execute the following command:
```
$ sudo karakuri nodectl mode stand-alone
```

## Connect to Remote node
By default, `karakuri` connected to local.  
To switch the node to be connected, first disconnect the current connection.
```
$ sudo karakuri nodectl disconnect
```
Then connect to remote node using the command (authentication code) obtained in the previous step.
```
$ sudo karakuri nodectl connect --node [NODE_ADDRESS] --auth [AUTHCODE]

node: [NODE_ADDRESS] connection success
```
You can check which node you are currently connected to with the following command.
```
$ sudo karakuri nodectl info

Mode : stand-alone

Target Node : [NODE_ADDRESS]
Status      : connected
```

(Note) The authorization code is only required for the first connection or if the authorization code is changed.  
The `--auth` option is not required when alternating between local and remote, which is explained in the next step.


## Connect to Local
To switch from a remote node to a local `karakuri`, execute the following command:
```
$ sudo karakuri nodectl disconnect
$ sudo karakuri nodectl connect --default
```
