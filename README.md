# Introduction

A simple Node DNS Server proxy based on [Pēteris Ņikiforovs tutorial](https://peteris.rocks/blog/dns-proxy-server-in-node-js-with-ui/)

# Features
* Create a static file with static address for each hostname that you want
* Create your docker containers with `HOSTNAMES` env, example: `HOSTNAMES=dev.mageddo.com, blog.mageddo.com` then access these hostnames in browser
* DNS try to solve the hosts from **docker** containers then from **records.json** file then from 3rd remote DNS servers one by one
* Cache for remote DNS increasing internet velocity
* List docker containers using [http://dns.mageddo:5380/containers](http://dns.mageddo:5380/containers)
* List cached hosts using [http://127.0.0.1:5380/cache](http://127.0.0.1:5380/cache)(without docker) or [http://dns.mageddo:5380/cache](http://dns.mageddo:5380/cache) (with docker)
* [Change DNS server ip](#changing-default-ports), or GUI editor IP


# Running

## Running without docker

you **can not** have other dns running on 53 port, in some tests I discover that the DNS not work as well on another port

Build the project 

	npm install

Starting the server 

	sudo npm start

setting this DNS as default DNS

```bash
sudo echo 'nameserver <127.0.0.1 or docker container ip>' > /etc/resolv.conf
```

## Running on docker

From **1.4.0** setup resolv.conf **is not** needled, is automatically

setup it

```bash
$ npm install # install dependencies
$ ./gradlew build # build docker image and run the container starting the app
```

# Testing if DNS is working

	$ host google.com 127.0.0.1
	Using domain server:
	Name: 127.0.0.1
	Address: 127.0.0.1#53
	Aliases:


# See binded hostnames with containers

	$ curl -X GET dns.mageddo:5380/containers
	container=/docker-dns-server, domain=dns.mageddo

# Adding manual DNS entries

you can create/edit `conf/records.json` (based on `records.samples.json`) manually or use the **Gui Editor**

# Gui editor

* To access from docker [http://dns.mageddo:5380](http://dns.mageddo:5380)
* To access without docker(running at npm) [http://127.0.0.1:5380](http://127.0.0.1:5380)


the password is `cat`

![](https://peteris.rocks/blog/dns-proxy-server-in-node-js-with-ui/dns.png)

# A entry example

records.json

```javascript
{
  "remoteDns": [
    { "address": "8.8.8.8", "port": 53, "type": "udp" }
  ],
  "entries" : [
    {
      "records": [
        {
          "type": "A",
          "address": "127.0.0.1",
          "ttl": 300,
          "name": "testing.mageddo.com"
        }
      ],
      "domain": "testing.mageddo.com"
    }
  ]
}
```

testing on terminal 

	$ host testing.mageddo.com 127.0.0.1
	Using domain server:
	Name: 127.0.0.1
	Address: 127.0.0.1#53
	Aliases: 

	testing3.mageddo.com has address 127.0.0.1
	testing3.mageddo.com has address 127.0.0.1
	testing3.mageddo.com has address 127.0.0.1
	
# Changing default ports

After [create your](#adding-manual-dns-entries) `records.json` file edit the variables `uiPort` and `dnsServerPort` then restart the DNS server

# License

This project is released under version 2.0 of the [Apache License][].
[Apache License]: http://www.apache.org/licenses/LICENSE-2.0

# Credits
* [Pēteris Ņikiforovs (project owner)](https://peteris.rocks/)
