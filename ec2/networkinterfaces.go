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

// NetworkInterfaceAttachment describes a network interface
// attachment.
//
// See http://goo.gl/KtiKuV for more details.
type NetworkInterfaceAttachment struct {
	Id                  string `xml:"attachmentId"`
	InstanceId          string `xml:"instanceId"`
	InstanceOwnerId     string `xml:"instanceOwnerId"`
	DeviceIndex         int    `xml:"deviceIndex"`
	Status              string `xml:"status"`
	AttachTime          string `xml:"attachTime"`
	DeleteOnTermination bool   `xml:"deleteOnTermination"`
}

// NetworkInterface describes a network interface for AWS VPC.
//
// See http://goo.gl/G63OQL for more details.
type NetworkInterface struct {
	Id               string                     `xml:"networkInterfaceId"`
	SubnetId         string                     `xml:"subnetId"`
	VpcId            string                     `xml:"vpcId"`
	AvailZone        string                     `xml:"availabilityZone"`
	Description      string                     `xml:"description"`
	OwnerId          string                     `xml:"ownerId"`
	RequesterId      string                     `xml:"requesterId"`
	RequesterManaged bool                       `xml:"requesterManaged"`
	Status           string                     `xml:"status"`
	MACAddress       string                     `xml:"macAddress"`
	PrivateIPAddress string                     `xml:"privateIpAddress"`
	PrivateDNSName   string                     `xml:"privateDnsName"`
	SourceDestCheck  bool                       `xml:"sourceDestCheck"`
	Groups           []SecurityGroup            `xml:"groupSet>item"`
	Attachment       NetworkInterfaceAttachment `xml:"attachment"`
	Tags             []Tag                      `xml:"tagSet>item"`
}

// NetworkInterfaceOptions encapsulates options for the
// CreateNetworkInterface call.
//
// Only the SubnetId is required, the rest are optional.
type NetworkInterfaceOptions struct {
	SubnetId         string
	PrivateIPAddress string
	Description      string
	SecurityGroupIds []string
}

// CreateNetworkInterfaceResp is the response to a
// CreateNetworkInterface request.
//
// See http://goo.gl/ze3VhA for more details.
type CreateNetworkInterfaceResp struct {
	RequestId        string           `xml:"requestId"`
	NetworkInterface NetworkInterface `xml:"networkInterface"`
}

// CreateNetworkInterface creates a network interface in the specified
// subnet.
//
// See http://goo.gl/ze3VhA for more details.
func (ec2 *EC2) CreateNetworkInterface(options NetworkInterfaceOptions) (resp *CreateNetworkInterfaceResp, err error) {
	params := makeParams("CreateNetworkInterface")
	params["SubnetId"] = options.SubnetId
	if options.PrivateIPAddress != "" {
		params["PrivateIpAddress"] = options.PrivateIPAddress
	}
	if options.Description != "" {
		params["Description"] = options.Description
	}
	for i, groupId := range options.SecurityGroupIds {
		params["SecurityGroupId."+strconv.Itoa(i+1)] = groupId
	}
	resp = &CreateNetworkInterfaceResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteNetworkInterface deletes the specified network interface.
// You must detach the network interface before you can delete it.
//
// See http://goo.gl/MC1yOj for more details.
func (ec2 *EC2) DeleteNetworkInterface(id string) (resp *SimpleResp, err error) {
	params := makeParams("DeleteNetworkInterface")
	params["NetworkInterfaceId"] = id
	resp = &SimpleResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DescribeNetworkInterfacesResp is the response to a
// DescribeNetworkInterfaces request.
//
// See http://goo.gl/2LcXtM for more details.
type DescribeNetworkInterfacesResp struct {
	RequestId  string             `xml:"requestId"`
	Interfaces []NetworkInterface `xml:"networkInterfaceSet>item"`
}

// DescribeNetworkInterfaces describes one or more network
// interfaces. Both parameters are optional, and if specified will
// limit the returned interfaces to the matching ids or filtering
// rules.
//
// See http://goo.gl/2LcXtM for more details.
func (ec2 *EC2) DescribeNetworkInterfaces(ids []string, filter *Filter) (resp *DescribeNetworkInterfacesResp, err error) {
	params := makeParams("DescribeNetworkInterfaces")
	for i, id := range ids {
		params["NetworkInterfaceId."+strconv.Itoa(i+1)] = id
	}
	filter.addParams(params)

	resp = &DescribeNetworkInterfacesResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AttachNetworkInterfaceResp is the response to a
// AttachNetworkInterface request.
//
// See http://goo.gl/rEbSii for more details.
type AttachNetworkInterfaceResp struct {
	RequestId    string `xml:"requestId"`
	AttachmentId string `xml:"attachmentId"`
}

// AttachNetworkInterface attaches a network interface to an instance.
//
// See http://goo.gl/rEbSii for more details.
func (ec2 *EC2) AttachNetworkInterface(interfaceId, instanceId string, deviceIndex int) (resp *AttachNetworkInterfaceResp, err error) {
	params := makeParams("AttachNetworkInterface")
	params["NetworkInterfaceId"] = interfaceId
	params["InstanceId"] = instanceId
	params["DeviceIndex"] = strconv.Itoa(deviceIndex)
	resp = &AttachNetworkInterfaceResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DetachNetworkInterface detaches a network interface from an
// instance.
//
// See http://goo.gl/0Xc1px for more details.
func (ec2 *EC2) DetachNetworkInterface(attachmentId string, force bool) (resp *SimpleResp, err error) {
	params := makeParams("DetachNetworkInterface")
	params["AttachmentId"] = attachmentId
	if force {
		// Force is optional.
		params["Force"] = "true"
	}
	resp = &SimpleResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
