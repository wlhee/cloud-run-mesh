package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/costinm/krun/pkg/hbone"
	"github.com/costinm/krun/pkg/k8s"
)


var initDebug func(run *k8s.KRun)

func main() {
	kr := &k8s.KRun{
		StartTime: time.Now(),
	}

	err := kr.InitK8SClient()
	if err != nil {
		log.Fatal("Failed to connect to GKE ", time.Since(kr.StartTime), kr, os.Environ(), err)
	}

	kr.LoadConfig()

	if len(os.Args) == 1 {
		// Default gateway label for now, we can customize with env variables.
		kr.Gateway = "ingressgateway"
		log.Println("Starting in gateway mode", os.Args)
	}

	kr.Refresh()

	if kr.XDSAddr == "" {
		kr.FindXDSAddr()
	}

	if kr.XDSAddr != "-" {
		proxyConfig := fmt.Sprintf(`{"discoveryAddress": "%s"}`, kr.XDSAddr)
		kr.StartIstioAgent(proxyConfig)
	}

	kr.StartApp()


	if InitDebug != nil {
		// Split for conditional compilation (to compile without ssh dep)
		InitDebug(kr)
	}


	// TODO: wait for app and proxy ready
	if kr.XDSAddr != "-" {
		hb := &hbone.HBone{
		}
		err = hb.Init()
		if err != nil {
			log.Println("Failed to init hbone", err)
		} else {

			err = hb.Start(":14009")
			if err != nil {
				panic(err)
			}
		}
	}
	select{}
}
