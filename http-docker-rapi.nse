local http = require "http"
local shortport = require "shortport"
local stdnse = require "stdnse"


description = [[
An attacker can start containers with malicious images, allowing 
the execution of commands with root permissions.
"The Engine API is an HTTP API served by Docker Engine. It is the 
API the Docker client uses to communicate with the Engine, so 
everything the Docker client can do can be done with the API."

There is no default port for this service

https://pyperanger.github.io/2018/01/18/docker-api/
https://docs.docker.com/engine/api/v1.24/

]]

--@usage
-- nmap --script http-docker-rapi -sV <target>

--@output
--PORT     STATE SERVICE REASON         VERSION
--4243/tcp open  http    syn-ack ttl 64 Golang net/http server
--| docker: 
--|   references: 
--|     https://pyperanger.github.io/2018/01/18/docker-api/
--|     https://docs.docker.com/engine/api/v1.24/
--|   description: An attacker can start containers with malicious images, allowing 
--|   the execution of commands with root permissions.
--|   "The Engine API is an HTTP API served by Docker Engine. It is the 
--|   API the Docker client uses to communicate with the Engine, so 
--|   everything the Docker client can do can be done with the API."
--|     
--|   Server /version: {"Version":"1.12.6"
--|   "ApiVersion":"1.24"
--|   "GitCommit":"ae7d637/1.12.6"
--|   "GoVersion":"go1.7.6"
--|   "Os":"linux"
--|   "Arch":"amd64"
--|   "KernelVersion":"4.4"
--|   "BuildTime":"2017-07-18T16:18:12.179285019+00:00"
--|   "PkgVersion":"docker-common-1.12.6-7.gitae7d637.fc25.x86_64"}
--| 
--|   risk_factor: High
--|_  title: Docker API Remote Privilege Escalation



author = "pype"
license = "Same as Nmap--See https://nmap.org/book/man-legal.html"
categories = { "vuln", "safe" }

portrule = shortport.http

action = function(host, port)
  local response = http.generic_request(host, port,"OPTIONS", "/")

  version = {}
  local vuln = {
    title = "Docker API Remote Privilege Escalation",
    risk_factor = "High",
    description = [[
An attacker can start containers with malicious images, allowing 
  the execution of commands with root permissions.
  "The Engine API is an HTTP API served by Docker Engine. It is the 
  API the Docker client uses to communicate with the Engine, so 
  everything the Docker client can do can be done with the API."
    ]],
    references = {
      'https://pyperanger.github.io/2018/01/18/docker-api/',
      'https://docs.docker.com/engine/api/v1.24/'
    }
  }


  if response.status == 200 and string.match(response.header["server"], "Docker") then
  	gver = http.get(host, port, "/version")
  	version["Server /version"] =  gver.body:gsub(",","\n\t")

  else
    return
  end

  local res_unauth = http.get(host, port, "/images/search?term=ubuntu")
  
  if res_unauth.status == 200 then
    vuln["Server /version"] =  gver.body:gsub(",","\n\t")
    return vuln
  else
    return version

  end
  
  end
