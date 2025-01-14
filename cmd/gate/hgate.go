package main

import (
	"log"
	"os"
	"time"

	"github.com/costinm/cloud-run-mesh/pkg/gcp"
	_ "github.com/costinm/cloud-run-mesh/pkg/gcp"
	"github.com/costinm/cloud-run-mesh/pkg/k8s"
	"github.com/costinm/cloud-run-mesh/pkg/snigate"
)

// Based on krun, start pilot-agent to get the certs and create the XDS proxy, and implement
// a SNI to H2 proxy - similar with the current multi-net gateway protocol from Istio side.
//
// This has a dependency on k8s - will auto-update the WorkloadInstance for H2R.
//
// However it does not depend directly on Istio or XDS - the certificates can be mounted or generated with
// krun+pilot-agent.
func main() {
	kr := k8s.New()

	kr.VendorInit = gcp.InitGCP

	_, err := snigate.InitSNIGate(kr, ":15443", ":15441")
	if err != nil {
		log.Fatal("Failed to connect to GKE ", time.Since(kr.StartTime), kr, os.Environ(), err)
	}
	log.Println("Started SNIGate", os.Environ())

	select {}

}
