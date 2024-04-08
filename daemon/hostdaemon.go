package main

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/openshift/dpu-operator/daemon/plugin"
	"github.com/openshift/dpu-operator/dpu-cni/pkgs/cnitypes"
	pb "github.com/opiproject/opi-api/network/evpn-gw/v1alpha1/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ctrl "sigs.k8s.io/controller-runtime"
)

type HostDaemon struct {
	dev           bool
	log           logr.Logger
	conn          *grpc.ClientConn
	client        pb.BridgePortServiceClient
	vsp           plugin.VendorPlugin
	addr          string
	port          int32
	cniServerPath string
}

func (d *HostDaemon) CreateBridgePort(pf int, vf int, vlan int) (*pb.BridgePort, error) {
	d.connectWithRetry()
	createRequest := &pb.CreateBridgePortRequest{
		BridgePort: &pb.BridgePort{
			Name: fmt.Sprintf("%d-%d-%d", pf, vf, vlan),
			Spec: &pb.BridgePortSpec{
				Ptype:          1,
				MacAddress:     []byte{},
				LogicalBridges: []string{},
			},
		},
	}

	return d.client.CreateBridgePort(context.TODO(), createRequest)
}

func (d *HostDaemon) DeleteBridgePort(pf int, vf int, vlan int) error {
	d.connectWithRetry()
	req := &pb.DeleteBridgePortRequest{
		Name: fmt.Sprintf("%d-%d-%d", pf, vf, vlan),
	}

	_, err := d.client.DeleteBridgePort(context.TODO(), req)
	return err
}

func NewHostDaemon(vsp plugin.VendorPlugin) *HostDaemon {
	return &HostDaemon{
		vsp:           vsp,
		log:           ctrl.Log.WithName("HostDaemon"),
		cniServerPath: cnitypes.ServerSocketPath,
	}
}

func (d *HostDaemon) connectWithRetry() {
	if d.conn != nil {
		return
	}
	// Might want to change waitForReady to true to
	// block on connection. Currently, we connect
	// "just in time" so the grpc immediately after
	// the dial will fail if connection failed (after
	// retries)
	retryPolicy := `{
		"methodConfig": [{
		  "waitForReady": false,
		  "retryPolicy": {
			  "MaxAttempts": 40,
			  "InitialBackoff": "1s",
			  "MaxBackoff": "16s",
			  "BackoffMultiplier": 2.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", d.addr, d.port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retryPolicy))
	if err != nil {
		d.log.Error(err, "did not connect")
	}
	d.conn = conn
	d.client = pb.NewBridgePortServiceClient(conn)
}

func (d *HostDaemon) Start() {
	d.log.Info("starting HostDaemon", "devflag", d.dev)

	addr, port, err := d.vsp.Start()
	if err != nil {
		d.log.Error(err, "VSP init returned error")
	}
	d.addr = addr
	d.port = port
}
