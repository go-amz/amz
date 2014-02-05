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
	"fmt"
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

const (
	// Common status values for network interfaces / attachments.
	AvailableStatus = "available"
	AttachingStatus = "attaching"
	AttachedStatus  = "attached"
	PendingStatus   = "pending"
	InUseStatus     = "in-use"
	DetachingStatus = "detaching"
	DetachedStatus  = "detached"
)

// PrivateIP describes a private IP address of a network interface.
//
// See http://goo.gl/jtuQEJ for more details.
type PrivateIP struct {
	Address   string `xml:"privateIpAddress"`
	DNSName   string `xml:"privateDnsName"`
	IsPrimary bool   `xml:"primary"`
}

// NetworkInterface describes a network interface for AWS VPC.
//
// See http://goo.gl/G63OQL for more details.
type NetworkInterface struct {
	Id               string                     `xml:"networkInterfaceId"`
	SubnetId         string                     `xml:"subnetId"`
	VPCId            string                     `xml:"vpcId"`
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
	PrivateIPs       []PrivateIP                `xml:"privateIpAddressesSet>item"`
}

// NetworkInterfaceOptions encapsulates options for the
// CreateNetworkInterface call.
//
// Only the SubnetId is required, the rest are optional.
//
// You can specify a primary private IP address by setting
// PrivateIPAddress or by using PrivateIPs slice, to set more than one
// IP. Only one of the given IPs can be set as primary.
//
// If you don't specify a private IP address, EC2 selects one for you
// from the subnet range.
//
// SecondaryPrivateIPsCount is the number of secondary private IP
// addresses to assign to the network interface. When you specify a
// number of secondary IP addresses, Amazon EC2 selects these IP
// addresses within the subnet range. The number of IP addresses you
// can assign to a network interface varies by instance type
type NetworkInterfaceOptions struct {
	SubnetId                 string
	PrivateIPAddress         string
	PrivateIPs               []PrivateIP
	SecondaryPrivateIPsCount int
	Description              string
	SecurityGroupIds         []string
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
func (ec2 *EC2) CreateNetworkInterface(opts NetworkInterfaceOptions) (resp *CreateNetworkInterfaceResp, err error) {
	params := makeParamsVPC("CreateNetworkInterface")
	params["SubnetId"] = opts.SubnetId
	var ips []PrivateIP
	if opts.PrivateIPAddress != "" {
		ips = append(ips, PrivateIP{
			Address:   opts.PrivateIPAddress,
			IsPrimary: true,
		})
	}
	ips = append(ips, opts.PrivateIPs...)
	for i, ip := range ips {
		prefix := fmt.Sprintf("PrivateIpAddresses.%d.", i+1)
		params[prefix+"PrivateIpAddress"] = ip.Address
		params[prefix+"Primary"] = strconv.FormatBool(ip.IsPrimary)
	}
	if opts.Description != "" {
		params["Description"] = opts.Description
	}
	if opts.SecondaryPrivateIPsCount > 0 {
		count := strconv.Itoa(opts.SecondaryPrivateIPsCount)
		params["SecondaryPrivateIpAddressCount"] = count
	}
	for i, groupId := range opts.SecurityGroupIds {
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
	params := makeParamsVPC("DeleteNetworkInterface")
	params["NetworkInterfaceId"] = id
	resp = &SimpleResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NetworkInterfacesResp is the response to a NetworkInterfaces
// request.
//
// See http://goo.gl/2LcXtM for more details.
type NetworkInterfacesResp struct {
	RequestId  string             `xml:"requestId"`
	Interfaces []NetworkInterface `xml:"networkInterfaceSet>item"`
}

// NetworkInterfaces describes one or more network interfaces. Both
// parameters are optional, and if specified will limit the returned
// interfaces to the matching ids or filtering rules.
//
// See http://goo.gl/2LcXtM for more details.
func (ec2 *EC2) NetworkInterfaces(ids []string, filter *Filter) (resp *NetworkInterfacesResp, err error) {
	params := makeParamsVPC("DescribeNetworkInterfaces")
	for i, id := range ids {
		params["NetworkInterfaceId."+strconv.Itoa(i+1)] = id
	}
	filter.addParams(params)

	resp = &NetworkInterfacesResp{}
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
	params := makeParamsVPC("AttachNetworkInterface")
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
	params := makeParamsVPC("DetachNetworkInterface")
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
