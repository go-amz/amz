//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011-2015 Canonical Ltd.
//
// This file contains the Server type itself and a few other
// unexported methods, not fitting anywhere else.

package ec2test

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"gopkg.in/amz.v3/ec2"
)

// TODO possible other things:
// - some virtual time stamp interface, so a client
// can ask for all actions after a certain virtual time.

// Server implements an EC2 simulator for use in testing.
type Server struct {
	url             string
	listener        net.Listener
	mu              sync.Mutex
	reqs            []*Action
	createRootDisks bool

	attributes           map[string][]string       // attr name -> values
	instances            map[string]*Instance      // id -> instance
	reservations         map[string]*reservation   // id -> reservation
	groups               map[string]*securityGroup // id -> group
	zones                []availabilityZone
	vpcs                 map[string]*vpc                 // id -> vpc
	subnets              map[string]*subnet              // id -> subnet
	ifaces               map[string]*iface               // id -> iface
	networkAttachments   map[string]*interfaceAttachment // id -> attachment
	volumes              map[string]*volume              // id -> volume
	volumeAttachments    map[string]*volumeAttachment    // id -> volumeAttachment
	maxId                counter
	reqId                counter
	reservationId        counter
	groupId              counter
	vpcId                counter
	dhcpOptsId           counter
	subnetId             counter
	volumeId             counter
	ifaceId              counter
	attachId             counter
	initialInstanceState ec2.InstanceState
}

// NewServer returns a new server.
func NewServer() (*Server, error) {
	srv := &Server{
		attributes:           make(map[string][]string),
		instances:            make(map[string]*Instance),
		groups:               make(map[string]*securityGroup),
		vpcs:                 make(map[string]*vpc),
		subnets:              make(map[string]*subnet),
		ifaces:               make(map[string]*iface),
		networkAttachments:   make(map[string]*interfaceAttachment),
		volumes:              make(map[string]*volume),
		volumeAttachments:    make(map[string]*volumeAttachment),
		reservations:         make(map[string]*reservation),
		initialInstanceState: Pending,
	}

	// Add default security group.
	g := &securityGroup{
		name:        "default",
		description: "default group",
		id:          fmt.Sprintf("sg-%d", srv.groupId.next()),
	}
	g.perms = map[permKey]bool{
		permKey{
			protocol: "icmp",
			fromPort: -1,
			toPort:   -1,
			group:    g,
		}: true,
		permKey{
			protocol: "tcp",
			fromPort: 0,
			toPort:   65535,
			group:    g,
		}: true,
		permKey{
			protocol: "udp",
			fromPort: 0,
			toPort:   65535,
			group:    g,
		}: true,
	}
	srv.groups[g.id] = g

	// Add a default availability zone.
	var z availabilityZone
	z.Name = defaultAvailZone
	z.Region = "us-east-1"
	z.State = "available"
	srv.zones = []availabilityZone{z}

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("cannot listen on localhost: %v", err)
	}
	srv.listener = l

	srv.url = "http://" + l.Addr().String()

	// we use HandlerFunc rather than *Server directly so that we
	// can avoid exporting HandlerFunc from *Server.
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		srv.serveHTTP(w, req)
	}))
	return srv, nil
}

// Quit closes down the server.
func (srv *Server) Quit() {
	srv.listener.Close()
}

// URL returns the URL of the server.
func (srv *Server) URL() string {
	return srv.url
}
