// The ec2test package implements a fake EC2 provider with
// the capability of inducing errors on any given operation,
// and retrospectively determining what operations have been
// carried out.
//
// This file contains internals of the EC2 test server, including the
// map of supported actions.
package ec2test

import (
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/amz.v3/ec2"
)

const (
	ownerId = "9876"
	// defaultAvailZone is the availability zone to use by default.
	defaultAvailZone = "us-east-1c"
)

// Action represents a request that changes the ec2 state.
type Action struct {
	RequestId string

	// Request holds the requested action as a url.Values instance
	Request url.Values

	// If the action succeeded, Response holds the value that
	// was marshalled to build the XML response for the request.
	Response interface{}

	// If the action failed, Err holds an error giving details of the failure.
	Err *ec2.Error
}

var actions = map[string]func(*Server, http.ResponseWriter, *http.Request, string) interface{}{
	"RunInstances":                  (*Server).runInstances,
	"TerminateInstances":            (*Server).terminateInstances,
	"DescribeInstances":             (*Server).describeInstances,
	"CreateSecurityGroup":           (*Server).createSecurityGroup,
	"DescribeAvailabilityZones":     (*Server).describeAvailabilityZones,
	"DescribeSecurityGroups":        (*Server).describeSecurityGroups,
	"DeleteSecurityGroup":           (*Server).deleteSecurityGroup,
	"AuthorizeSecurityGroupIngress": (*Server).authorizeSecurityGroupIngress,
	"RevokeSecurityGroupIngress":    (*Server).revokeSecurityGroupIngress,
	"CreateVpc":                     (*Server).createVpc,
	"DeleteVpc":                     (*Server).deleteVpc,
	"DescribeVpcs":                  (*Server).describeVpcs,
	"CreateSubnet":                  (*Server).createSubnet,
	"DeleteSubnet":                  (*Server).deleteSubnet,
	"DescribeSubnets":               (*Server).describeSubnets,
	"ModifySubnetAttribute":         (*Server).modifySubnetAttribute,
	"CreateNetworkInterface":        (*Server).createIFace,
	"DeleteNetworkInterface":        (*Server).deleteIFace,
	"DescribeNetworkInterfaces":     (*Server).describeIFaces,
	"AttachNetworkInterface":        (*Server).attachIFace,
	"DetachNetworkInterface":        (*Server).detachIFace,
	"DescribeAccountAttributes":     (*Server).accountAttributes,
	"AssignPrivateIpAddresses":      (*Server).assignPrivateIP,
	"UnassignPrivateIpAddresses":    (*Server).unassignPrivateIP,
	"CreateVolume":                  (*Server).createVolume,
	"DeleteVolume":                  (*Server).deleteVolume,
	"DescribeVolumes":               (*Server).describeVolumes,
	"AttachVolume":                  (*Server).attachVolume,
	"DetachVolume":                  (*Server).detachVolume,
	"ModifyInstanceAttribute":       (*Server).modifyInstanceAttribute,
	"CreateTags":                    (*Server).createTags,
}

// newAction allocates a new action and adds it to the
// recorded list of server actions.
func (srv *Server) newAction() *Action {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	a := new(Action)
	srv.reqs = append(srv.reqs, a)
	return a
}

// serveHTTP serves the EC2 protocol.
func (srv *Server) serveHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	a := srv.newAction()
	a.RequestId = fmt.Sprintf("req%d", srv.reqId.next())
	a.Request = req.Form

	// Methods on Server that deal with parsing user data
	// may fail. To save on error handling code, we allow these
	// methods to call fatalf, which will panic with an *ec2.Error
	// which will be caught here and returned
	// to the client as a properly formed EC2 error.
	defer func() {
		switch err := recover().(type) {
		case *ec2.Error:
			a.Err = err
			err.RequestId = a.RequestId
			writeError(w, err)
		case nil:
		default:
			panic(err)
		}
	}()

	f := actions[req.Form.Get("Action")]
	if f == nil {
		fatalf(400, "InvalidParameterValue", "Unrecognized Action")
	}

	response := f(srv, w, req, a.RequestId)
	a.Response = response

	w.Header().Set("Content-Type", `xml version="1.0" encoding="UTF-8"`)
	xmlMarshal(w, response)
}

func (srv *Server) notImplemented(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	fatalf(500, "InternalError", "not implemented")
	panic("not reached")
}
