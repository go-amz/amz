//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011 Canonical Ltd.
//
// Written by Gustavo Niemeyer <gustavo.niemeyer@canonical.com>
//

package ec2

import (
	"strconv"
)

// Vpc describes an Amazon Virtual Private Cloud (VPC).
//
// See http://goo.gl/Uy6ZLL for more details.
type Vpc struct {
	Id            string `xml:"vpcId"`
	State         string `xml:"state"`
	CidrBlock     string `xml:"cidrBlock"`
	DhcpOptionsId string `xml:"dhcpOptionsId"`
	Tags          []Tag  `xml:"tagSet>item"`
}

// CreateVpcResp is the response to a CreateVpc request.
//
// See http://goo.gl/nkwjvN for more details.
type CreateVpcResp struct {
	RequestId string `xml:"requestId"`
	Vpc       Vpc    `xml:"vpc"`
}

// CreateVpc creates a VPC with the specified CIDR block.
//
// The smallest VPC you can create uses a /28 netmask (16 IP
// addresses), and the largest uses a /16 netmask (65,536 IP
// addresses).
//
// See http://goo.gl/nkwjvN for more details.
func (ec2 *EC2) CreateVpc(cidrBlock string) (resp *CreateVpcResp, err error) {
	params := makeParams("CreateVpc")
	params["CidrBlock"] = cidrBlock
	resp = &CreateVpcResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteVpc deletes the specified VPC. You must detach or delete all
// gateways and resources that are associated with the VPC before you
// can delete it. For example, you must terminate all instances
// running in the VPC, delete all security groups associated with the
// VPC (except the default one), delete all route tables associated
// with the VPC (except the default one), and so on.
//
// See http://goo.gl/bcxtbf for more details.
func (ec2 *EC2) DeleteVpc(id string) (resp *SimpleResp, err error) {
	params := makeParams("DeleteVpc")
	params["VpcId"] = id
	resp = &SimpleResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DescribeVpcsResp is the response to a DescribeVpcs request.
//
// See http://goo.gl/Y5kHqG for more details.
type DescribeVpcsResp struct {
	RequestId string `xml:"requestId"`
	Vpcs      []Vpc  `xml:"vpcSet>item"`
}

// DescribeVpcs describes one or more VPCs. Both parameters are
// optional, and if specified will limit the returned VPCs to the
// matching ids or filtering rules.
//
// See http://goo.gl/Y5kHqG for more details.
func (ec2 *EC2) DescribeVpcs(ids []string, filter *Filter) (resp *DescribeVpcsResp, err error) {
	params := makeParams("DescribeVpcs")
	for i, id := range ids {
		params["VpcId."+strconv.Itoa(i+1)] = id
	}
	filter.addParams(params)

	resp = &DescribeVpcsResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
