//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011-2015 Canonical Ltd.
//
// This file contains code handling AWS account attributes
// discovery API.

package ec2test

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"gopkg.in/amz.v3/ec2"
)

// SetInitialAttributes sets the given account attributes on the server.
func (srv *Server) SetInitialAttributes(attrs map[string][]string) {
	for attrName, values := range attrs {
		srv.attributes[attrName] = values
		if attrName == "default-vpc" {
			// The default-vpc attribute was provided, so create the
			// respective VPCs and their subnets.
			for _, vpcId := range values {
				srv.vpcs[vpcId] = &vpc{ec2.VPC{
					Id:              vpcId,
					State:           "available",
					CIDRBlock:       "10.0.0.0/16",
					DHCPOptionsId:   fmt.Sprintf("dopt-%d", srv.dhcpOptsId.next()),
					InstanceTenancy: "default",
					IsDefault:       true,
				}}
				subnetId := fmt.Sprintf("subnet-%d", srv.subnetId.next())
				cidrBlock := "10.10.0.0/20"
				availIPs, _ := srv.calcSubnetAvailIPs(cidrBlock)
				srv.subnets[subnetId] = &subnet{ec2.Subnet{
					Id:               subnetId,
					VPCId:            vpcId,
					State:            "available",
					CIDRBlock:        cidrBlock,
					AvailZone:        "us-east-1b",
					AvailableIPCount: availIPs,
					DefaultForAZ:     true,
				}}
			}
		}
	}
}

func (srv *Server) accountAttributes(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	attrsMap := parseIDs(req.Form, "AttributeName.")
	var resp struct {
		XMLName xml.Name
		ec2.AccountAttributesResp
	}
	resp.XMLName = xml.Name{defaultXMLName, "DescribeAccountAttributesResponse"}
	resp.RequestId = reqId
	for attrName, _ := range attrsMap {
		vals, ok := srv.attributes[attrName]
		if !ok {
			fatalf(400, "InvalidParameterValue", "describe attrs: not found %q", attrName)
		}
		resp.Attributes = append(resp.Attributes, ec2.AccountAttribute{attrName, vals})
	}
	return &resp
}
