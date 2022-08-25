package main

import (
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func zitiConnector(ztx ziti.Context) mqtt.OpenConnectionFunc {
	return func(uri *url.URL, options mqtt.ClientOptions) (net.Conn, error) {
		return ztx.Dial(uri.Path)
	}
}

func main() {
	identity := flag.String("identity", "id", "Ziti identity file")
	service := flag.String("service", "", "Ziti MQTT service name")
	topic := flag.String("topic", "", "MQTT topic")
	pub := flag.Bool("pub", false, "Publish")

	flag.Parse()
	var ztx ziti.Context
	if cfg, err := config.NewFromFile(*identity); err == nil {
		ztx = ziti.NewContextWithConfig(cfg)
	} else {
		ztx = ziti.NewContext()
	}

	id := uuid.New()
	mqttOpts := &mqtt.ClientOptions{
		ClientID: id.String(),
		Servers: []*url.URL{
			{Scheme: "mqziti", Path: *service},
		},
		CustomOpenConnectionFn: zitiConnector(ztx),
	}
	c := mqtt.NewClient(mqttOpts)
	connected := c.Connect()
	connected.Wait()
	if connected.Error() != nil {
		log.Fatalf("failed to connect: %v", connected.Error())
	}

	if *pub {
		msg := strings.Join(flag.Args(), " ")
		tok := c.Publish(*topic, 1, true, msg)
		tok.Wait()
		if tok.Error() != nil {
			log.Fatalf("failed to publish = %v\n", tok.Error())
		}
	} else {
		tok := c.Subscribe(*topic, 1, func(client mqtt.Client, message mqtt.Message) {
			fmt.Printf(">>> %v\n", string(message.Payload()))
		})
		tok.Wait()

		if tok.Error() != nil {
			fmt.Printf("subscription error = %v\n", tok.Error())
		} else {
			sigs := make(chan os.Signal, 1)
			done := make(chan bool)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigs
				done <- true
			}()
			fmt.Println("wainting for messages...")
			<-done
		}

	}
}
