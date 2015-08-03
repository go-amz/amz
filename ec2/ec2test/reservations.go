// The ec2test package implements a fake EC2 provider with
// the capability of inducing errors on any given operation,
// and retrospectively determining what operations have been
// carried out.
//
// This file contains code handling AWS API around instance
// Reservations.
package ec2test

import "fmt"

// reservation holds a simulated ec2 reservation.
type reservation struct {
	id        string
	instances map[string]*Instance
	groups    []*securityGroup
}

func (srv *Server) newReservation(groups []*securityGroup) *reservation {
	r := &reservation{
		id:        fmt.Sprintf("r-%d", srv.reservationId.next()),
		instances: make(map[string]*Instance),
		groups:    groups,
	}

	srv.reservations[r.id] = r
	return r
}

func (r *reservation) hasRunningMachine() bool {
	for _, inst := range r.instances {
		if inst.state.Code != ShuttingDown.Code && inst.state.Code != Terminated.Code {
			return true
		}
	}
	return false
}
