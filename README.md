# Introduction

A DNS Server that works as a simple DNS and:

* [Resolves automatically Docker containers IPs to hostnames that you set](#resolves-automatically-docker-containers-ips-to-hostnames-that-you-set)

* [Have manually entries(config file or web gui editor) to you setup IPs for hostnames that you want ](#add-manual-dns-entries)


# Running

## Running without docker

you **can not** have other dns running on 53 port, in some tests I discover that the DNS not work as well on another port

Build the project 

	npm install

Starting the server 

	npm start

setting this DNS as default DNS, (on **1.4.0** setup resolv.conf **is not** needled, is automatically)

```bash
sudo echo 'nameserver <127.0.0.1 or docker container ip>' > /etc/resolv.conf
```

## Running on docker
setup it

```bash
$ npm install # install dependencies
$ gradle build-dev # build docker image and run the container starting the app
```

# Testing if DNS is working

	$ host google.com 127.0.0.1
	Using domain server:
	Name: 127.0.0.1
	Address: 127.0.0.1#53
	Aliases:


    
# Features

## Resolves automatically [Docker](https://docker.com) containers IPs to hostnames that you set

**Example**, running a apache server on docker

```bash
$ docker run --rm --hostname my.google.com eboraas/apache
```

with this **DNS server** to access this container, you **have not** to know the docker container IP(that is dynamically and changes on every run), or publish the container port to host port. Simply open your browser and type http://my.google.com and have fun.

**OBS:** `--hostname` is a docker native flag

**Add more than one hostname to the same docker container**
If you need more than one hostname resolving to the same docker container you can use as follows

```bash
$ docker run --rm --hostname my.google.com --env HOSTNAMES="my2.google.com,my3.google.com" eboraas/apache
```


## Add manual DNS entries

you can edit `records.json` manually or use the **Gui Editor**

#### Gui editor

Access http://localhost/ or http://dns.mageddo/ (when running on docker)

the **password** is `cat`
![](http://i.imgur.com/Zf5nlla.jpg)

#### A manual entry example

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


# License

This project is released under version 2.0 of the [Apache License](http://www.apache.org/licenses/LICENSE-2.0).

# Credits
* Elvis de Freitas
* [Pēteris Ņikiforovs (Project Owner)](https://peteris.rocks/)
