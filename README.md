MQTT => MQZiti
==============


This project shows how to easily Zitify(run dark on zero trust overlay network) your Golang MQTT server and clients.

## Setup
1. Get your code -- clone this repo
1. Get yourself an OpenZiti network and Ziti identities
   * follow [quickstart](https://openziti.github.io/ziti/quickstarts/quickstart-overview.html) or, use [Ziti Edge Developer Sandbox](https://zeds.openziti.org)
   * create a Ziti service that will be used to for MQTT communication 
   * add/enroll Ziiti identities
1. build this project
   ```console
   $ mkdir build && go build -o build ./...
   ```
   this creates `mqziti_server` and `mqziti_client` in your build directory


This following is assumed for the rest of this document:
* 'mqziti' - the name of the service we are going to use
* `server.json` - Ziti identity file for the server
* `client.json` - Ziti identity for the client

## Server

We use [Mochi MQTT](https://github.com/mochi-co/mqtt) as base and implement a 
[`Listener`](https://github.com/mochi-co/mqtt/blob/v1.3.1/server/listeners/listeners.go#L46) that binds to the Ziti service.

Run the server
```console
$ ./build/mqziti_server -identity server.json -service mqziti

```

You can check that the process has no listening sockets. This means that you need to open your firewall.

## Client

We use [Paho MQTT](https://github.com/eclipse/paho.mqtt.golang) and implement a connector that connects to Ziti service.

You will need to run `mqziti_client` twice for this test: one instance to subscribe, and one instance to publish.

Subscriber:
```console
$ ./build/mqziti_client -identity client.json -service mqziti -topic /openziti
```

Publish something:
```console
$ ./build/mqziti_client -identity client.json -service mqziti -topic /openziti -pub "Hello OpenZiti!"
```

You should see the message printed on the subscriber console.


## Links

* Follow our [Blog](https://openziti.io/)
* Join [Discussion](https://openziti.discourse.group)
* [Development](https://github.com/openziti)
* [Documentation](https://openziti.github.io) 
