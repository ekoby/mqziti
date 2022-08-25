package main

import (
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/system"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/openziti/sdk-golang/ziti/edge"
	"log"
)

type zitiListener struct {
	ztx     ziti.Context
	service string
	server  edge.Listener
	config  *listeners.Config
}

func (z *zitiListener) SetConfig(config *listeners.Config) {
	z.config = config
}

func (z *zitiListener) Listen(s *system.Info) (e error) {
	z.server, e = z.ztx.Listen(z.service)
	return e
}

func (z *zitiListener) Serve(establishFunc listeners.EstablishFunc) {

	for {
		clt, err := z.server.Accept()

		edge_clt := clt.(edge.Conn)
		log.Printf("ziti client: %v", edge_clt.SourceIdentifier())
		if err != nil {
			log.Fatal("failed to accept ", err)
		}

		go func() {
			_ = establishFunc(z.service, clt, z.config.Auth)
		}()
	}
}

func (z *zitiListener) ID() string {
	return z.service
}

func (z *zitiListener) Close(closeFunc listeners.CloseFunc) {

	closeFunc(z.service)

	if z.server != nil {
		_ = z.server.Close()
	}
}

func NewZitiListener(identity string, service string) listeners.Listener {
	var ztx ziti.Context
	if identity != "" {
		if cfg, err := config.NewFromFile(identity); err == nil {
			ztx = ziti.NewContextWithConfig(cfg)
		} else {
			log.Fatalf("failed to load ziti: %s", err.Error())
		}
	} else {
		ztx = ziti.NewContext()
	}

	return &zitiListener{
		ztx:     ztx,
		service: service,
	}
}
