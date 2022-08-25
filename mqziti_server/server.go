package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/logrusorgru/aurora"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

func main() {
	identity := flag.String("identity", "",
		"Ziti Idenity file (optional if ZITI_SDK_CONFIG environment var is set)")
	service := flag.String("service", "", "Ziti MQTT service name (required)")

	flag.Parse()
	fmt.Println(aurora.Magenta("Mochi MQTT Server initializing..."), aurora.Cyan("Ziti"))

	// An example of configuring various server options...
	options := &mqtt.Options{
		BufferSize:      0, // Use default values
		BufferBlockSize: 0, // Use default values
	}

	server := mqtt.NewServer(options)

	zl := NewZitiListener(*identity, *service)

	err := server.AddListener(zl, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println(aurora.BgMagenta("  Started!  "))

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	<-done
	fmt.Println(aurora.BgRed("  Caught Signal  "))

	_ = server.Close()
	fmt.Println(aurora.BgGreen("  Finished  "))

}
