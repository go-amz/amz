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

// Subnet describes an Amazon VPC subnet.
//
// See http://goo.gl/CdkvO2 for more details.
type Subnet struct {
	Id                  string `xml:"subnetId"`
	State               string `xml:"state"`
	VPCId               string `xml:"vpcId"`
	CIDRBlock           string `xml:"cidrBlock"`
	AvailableIPCount    int    `xml:"availableIpAddressCount"`
	AvailZone           string `xml:"availabilityZone"`
	DefaultForAZ        bool   `xml:"defaultForAz"`
	MapPublicIPOnLaunch bool   `xml:"mapPublicIpOnLaunch"`
	Tags                []Tag  `xml:"tagSet>item"`
}

// CreateSubnetResp is the response to a CreateSubnet request.
//
// See http://goo.gl/wLPhfI for more details.
type CreateSubnetResp struct {
	RequestId string `xml:"requestId"`
	Subnet    Subnet `xml:"subnet"`
}

// CreateSubnet creates a subnet in an existing VPC.
//
// When you create each subnet, you provide the VPC ID and the CIDR
// block you want for the subnet. After you create a subnet, you can't
// change its CIDR block. The subnet's CIDR block can be the same as
// the VPC's CIDR block (assuming you want only a single subnet in the
// VPC), or a subset of the VPC's CIDR block. If you create more than
// one subnet in a VPC, the subnets' CIDR blocks must not overlap. The
// smallest subnet (and VPC) you can create uses a /28 netmask (16 IP
// addresses), and the largest uses a /16 netmask (65,536 IP
// addresses).
//
// availZone can be empty, in which case Amazon EC2 selects one for
// you (recommended).
//
// See http://goo.gl/wLPhfI for more details.
func (ec2 *EC2) CreateSubnet(vpcId, cidrBlock, availZone string) (resp *CreateSubnetResp, err error) {
	params := makeParamsVPC("CreateSubnet")
	params["VpcId"] = vpcId
	params["CidrBlock"] = cidrBlock
	if availZone != "" {
		params["AvailabilityZone"] = availZone
	}
	resp = &CreateSubnetResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteSubnet deletes the specified subnet. You must terminate all
// running instances in the subnet before you can delete the subnet.
//
// See http://goo.gl/KmhcBM for more details.
func (ec2 *EC2) DeleteSubnet(id string) (resp *SimpleResp, err error) {
	params := makeParamsVPC("DeleteSubnet")
	params["SubnetId"] = id
	resp = &SimpleResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SubnetsResp is the response to a Subnets request.
//
// See http://goo.gl/NTKQVI for more details.
type SubnetsResp struct {
	RequestId string   `xml:"requestId"`
	Subnets   []Subnet `xml:"subnetSet>item"`
}

// Subnets describes one or more of your subnets. Both parameters are
// optional, and if specified will limit the returned subnets to the
// matching ids or filtering rules.
//
// See http://goo.gl/NTKQVI for more details.
func (ec2 *EC2) Subnets(ids []string, filter *Filter) (resp *SubnetsResp, err error) {
	params := makeParamsVPC("DescribeSubnets")
	for i, id := range ids {
		params["SubnetId."+strconv.Itoa(i+1)] = id
	}
	filter.addParams(params)

	resp = &SubnetsResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
