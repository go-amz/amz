//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011 Canonical Ltd.
//
// This file contains code handling AWS API around Subnets.

package ec2test

import (
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/amz.v3/ec2"
)

type subnet struct {
	ec2.Subnet
}

func (s *subnet) matchAttr(attr, value string) (ok bool, err error) {
	switch attr {
	case "cidr":
		return s.CIDRBlock == value, nil
	case "availability-zone":
		return s.AvailZone == value, nil
	case "state":
		return s.State == value, nil
	case "subnet-id":
		return s.Id == value, nil
	case "vpc-id":
		return s.VPCId == value, nil
	case "defaultForAz", "default-for-az":
		val, err := strconv.ParseBool(value)
		if err != nil {
			return false, fmt.Errorf("bad flag %q: %s", attr, value)
		}
		return s.DefaultForAZ == val, nil
	case "tag", "tag-key", "tag-value", "available-ip-address-count":
		return false, fmt.Errorf("%q filter not implemented", attr)
	}
	return false, fmt.Errorf("unknown attribute %q", attr)
}

// getDefaultSubnet returns the first default subnet for the AZ in the
// default VPC (if available).
func (srv *Server) getDefaultSubnet() *subnet {
	// We need to get the default VPC id and one of its subnets to use.
	defaultVPCId := ""
	for _, vpc := range srv.vpcs {
		if vpc.IsDefault {
			defaultVPCId = vpc.Id
			break
		}
	}
	if defaultVPCId == "" {
		// No default VPC, so nothing to do.
		return nil
	}
	for _, subnet := range srv.subnets {
		if subnet.VPCId == defaultVPCId && subnet.DefaultForAZ {
			return subnet
		}
	}
	return nil
}

func (srv *Server) calcSubnetAvailIPs(cidrBlock string) (int, error) {
	_, ipnet, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		return 0, err
	}
	// calculate the available IP addresses, removing the first 4 and
	// the last, which are reserved by AWS.
	maskOnes, maskBits := ipnet.Mask.Size()
	return 1<<uint(maskBits-maskOnes) - 5, nil
}

func (srv *Server) createSubnet(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	v := srv.vpc(req.Form.Get("VpcId"))
	cidrBlock := parseCidr(req.Form.Get("CidrBlock"))
	availZone := req.Form.Get("AvailabilityZone")
	if availZone == "" {
		// Assign one automatically as AWS does.
		availZone = "us-east-1b"
	}
	availIPs, err := srv.calcSubnetAvailIPs(cidrBlock)
	if err != nil {
		fatalf(400, "InvalidParameterValue", "calcSubnetAvailIPs(%q) failed: %v", cidrBlock, err)
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()
	s := &subnet{ec2.Subnet{
		Id:               fmt.Sprintf("subnet-%d", srv.subnetId.next()),
		VPCId:            v.Id,
		State:            "available",
		CIDRBlock:        cidrBlock,
		AvailZone:        availZone,
		AvailableIPCount: availIPs,
	}}
	srv.subnets[s.Id] = s
	var resp struct {
		XMLName xml.Name
		ec2.CreateSubnetResp
	}
	resp.XMLName = xml.Name{defaultXMLName, "CreateSubnetResponse"}
	resp.RequestId = reqId
	resp.Subnet = s.Subnet
	return resp
}

func (srv *Server) deleteSubnet(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	s := srv.subnet(req.Form.Get("SubnetId"))
	srv.mu.Lock()
	defer srv.mu.Unlock()

	delete(srv.subnets, s.Id)
	return &ec2.SimpleResp{
		XMLName:   xml.Name{defaultXMLName, "DeleteSubnetResponse"},
		RequestId: reqId,
		Return:    true,
	}
}

func (srv *Server) describeSubnets(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	idMap := parseIDs(req.Form, "SubnetId.")
	f := newFilter(req.Form)
	var resp struct {
		XMLName xml.Name
		ec2.SubnetsResp
	}
	resp.XMLName = xml.Name{defaultXMLName, "DescribeSubnetsResponse"}
	resp.RequestId = reqId
	for _, s := range srv.subnets {
		ok, err := f.ok(s)
		_, known := idMap[s.Id]
		if ok && (len(idMap) == 0 || known) {
			resp.Subnets = append(resp.Subnets, s.Subnet)
		} else if err != nil {
			fatalf(400, "InvalidParameterValue", "describe subnets: %v", err)
		}
	}
	return &resp
}

func (srv *Server) modifySubnetAttribute(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	id := req.Form.Get("SubnetId")
	s := srv.subnet(id)
	mapIp := strings.ToLower(req.Form.Get("MapPublicIpOnLaunch.Value")) == "true"
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if s == nil {
		fatalf(400, "InvalidSubnetID.NotFound", "no such subnet %v", id)
	}
	s.MapPublicIPOnLaunch = mapIp
	srv.subnets[id] = s

	return &ec2.SimpleResp{
		XMLName:   xml.Name{defaultXMLName, "ModifySubnetAttributeResponse"},
		RequestId: reqId,
		Return:    true,
	}
}

func (srv *Server) subnet(id string) *subnet {
	if id == "" {
		fatalf(400, "MissingParameter", "missing subnetId")
	}
	srv.mu.Lock()
	defer srv.mu.Unlock()
	s, found := srv.subnets[id]
	if !found {
		fatalf(400, "InvalidSubnetID.NotFound", "subnet %s not found", id)
	}
	return s
}
