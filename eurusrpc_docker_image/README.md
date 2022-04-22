# Install Eurus Node
##### This guide is for Ubuntu 20.04 

Requirements: basic Linux command. 

A VPS with 2cpu 8gb ram 
HD 300GB

## Make sure your host network port is allow income package for below services

30303 UDP & TCP for p2p

8545 TCP for RPC

8546 TCP for websocket 

8547 TCP for graphql

9545 TCP for metrics



## Prerequisites

Install docker
https://docs.docker.com/engine/install/ubuntu/

Install docker composer 
https://docs.docker.com/compose/install/

## Setup the environment
Copy data to your root folder /

### Change your IP address 
Change p2p-host=" YOUR EXTERNAL IP Address" 

on /data/config.toml

Copy docker-composer.yml to /home/ubuntu/

Run this command under /home/ubuntu/ 

```bash
docker-composer up -d
```
install the service from the yaml file

## Get the Enode URL
when service up . Go to docker container Log and find a message like below . 

Sample as below

```bash
"level":"INFO","thread":"main","class":"DefaultP2PNetwork","message":"Enode URL enode://6645ee02b748c5ce6c5267b8372902157db5ef547ce3c6a43c788eb27256d45e9983dd533636a9e034d74e1b35115168596717250b59d84ab471b8a3974702ba@52.74.36.240:30303","throwable":""}
```

Copy your Enode URL and sent us a email for approval  
