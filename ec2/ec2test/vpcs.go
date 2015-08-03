//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011 Canonical Ltd.
//
// This file contains code handling AWS API around VPCs.

package ec2test

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"gopkg.in/amz.v3/ec2"
)

type vpc struct {
	ec2.VPC
}

func (v *vpc) matchAttr(attr, value string) (ok bool, err error) {
	switch attr {
	case "cidr":
		return v.CIDRBlock == value, nil
	case "state":
		return v.State == value, nil
	case "vpc-id":
		return v.Id == value, nil
	case "tag", "tag-key", "tag-value", "dhcp-options-id", "isDefault":
		return false, fmt.Errorf("%q filter is not implemented", attr)
	}
	return false, fmt.Errorf("unknown attribute %q", attr)
}

func (srv *Server) createVpc(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	cidrBlock := parseCidr(req.Form.Get("CidrBlock"))
	tenancy := req.Form.Get("InstanceTenancy")
	if tenancy == "" {
		tenancy = "default"
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()
	v := &vpc{ec2.VPC{
		Id:              fmt.Sprintf("vpc-%d", srv.vpcId.next()),
		State:           "available",
		CIDRBlock:       cidrBlock,
		DHCPOptionsId:   fmt.Sprintf("dopt-%d", srv.dhcpOptsId.next()),
		InstanceTenancy: tenancy,
	}}
	srv.vpcs[v.Id] = v
	var resp struct {
		XMLName xml.Name
		ec2.CreateVPCResp
	}
	resp.XMLName = xml.Name{defaultXMLName, "CreateVpcResponse"}
	resp.RequestId = reqId
	resp.VPC = v.VPC
	return resp
}

func (srv *Server) deleteVpc(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	v := srv.vpc(req.Form.Get("VpcId"))
	srv.mu.Lock()
	defer srv.mu.Unlock()

	delete(srv.vpcs, v.Id)
	return &ec2.SimpleResp{
		XMLName:   xml.Name{defaultXMLName, "DeleteVpcResponse"},
		RequestId: reqId,
		Return:    true,
	}
}

func (srv *Server) describeVpcs(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	idMap := parseIDs(req.Form, "VpcId.")
	f := newFilter(req.Form)
	var resp struct {
		XMLName xml.Name
		ec2.VPCsResp
	}
	resp.XMLName = xml.Name{defaultXMLName, "DescribeVpcsResponse"}
	resp.RequestId = reqId
	for _, v := range srv.vpcs {
		ok, err := f.ok(v)
		_, known := idMap[v.Id]
		if ok && (len(idMap) == 0 || known) {
			resp.VPCs = append(resp.VPCs, v.VPC)
		} else if err != nil {
			fatalf(400, "InvalidParameterValue", "describe VPCs: %v", err)
		}
	}
	return &resp
}

func (srv *Server) vpc(id string) *vpc {
	if id == "" {
		fatalf(400, "MissingParameter", "missing vpcId")
	}
	srv.mu.Lock()
	defer srv.mu.Unlock()
	v, found := srv.vpcs[id]
	if !found {
		fatalf(400, "InvalidVpcID.NotFound", "VPC %s not found", id)
	}
	return v
}
