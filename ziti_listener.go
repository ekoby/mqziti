package main

import (
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/system"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/edge"
	"log"
)

type ZitiListener struct {
	ztx     ziti.Context
	service string
	server  edge.Listener
	config  *listeners.Config
}

func (z *ZitiListener) SetConfig(config *listeners.Config) {
	z.config = config
}

func (z *ZitiListener) Listen(s *system.Info) (e error) {
	z.server, e = z.ztx.Listen(z.service)
	return e
}

func (z *ZitiListener) Serve(establishFunc listeners.EstablishFunc) {

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

func (z *ZitiListener) ID() string {
	return z.service
}

func (z ZitiListener) Close(closeFunc listeners.CloseFunc) {

	closeFunc(z.service)

	if z.server != nil {
		_ = z.server.Close()
	}
}

func NewZitiListener(service string) listeners.Listener {
	return &ZitiListener{
		ztx:     ziti.NewContext(),
		service: service,
	}
}
