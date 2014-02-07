//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2014 Canonical Ltd.
//
// Written by Gustavo Niemeyer <gustavo.niemeyer@canonical.com>
//

package ec2

import (
	"strconv"
)

// AssignPrivateIPAddresses assigns one or more secondary private IP
// addresses to the specified network interface. One or more specific
// secondary IP addresses can be given explicitly, or a count of
// secondary IP addresses can be given to be automatically assigned
// within the subnet's CIDR block range. The total number of secondary
// IP addresses that can be assigned to an instance varies by instance
// type.
//
// If ipAddresses is non-empty, secondaryIPsCount is ignored (they are
// mutually exclusive).
//
// allowReassignment specifies whether to allow reassignment of
// addresses currently assigned to a different network interface.
//
// See http://goo.gl/MoeH0L more details.
func (ec2 *EC2) AssignPrivateIPAddresses(interfaceId string, ipAddresses []string, secondaryIPsCount int, allowReassignment bool) (resp *SimpleResp, err error) {
	params := makeParamsVPC("AssignPrivateIpAddresses")
	params["NetworkInterfaceId"] = interfaceId
	if secondaryIPsCount > 0 {
		params["SecondaryPrivateIpAddressCount"] = strconv.Itoa(secondaryIPsCount)
	} else {
		for i, ip := range ipAddresses {
			// PrivateIpAddress is zero indexed.
			n := strconv.Itoa(i)
			params["PrivateIpAddress."+n] = ip
		}
	}
	if allowReassignment {
		params["AllowReassignment"] = "true"
	}
	resp = &SimpleResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UnassignPrivateIPAddresses unassigns one or more secondary private
// IP addresses from a network interface.
//
// See http://goo.gl/RjGZdB for more details.
func (ec2 *EC2) UnassignPrivateIPAddresses(interfaceId string, ipAddresses []string) (resp *SimpleResp, err error) {
	params := makeParamsVPC("UnassignPrivateIpAddresses")
	params["NetworkInterfaceId"] = interfaceId
	for i, ip := range ipAddresses {
		// PrivateIpAddress is zero indexed.
		n := strconv.Itoa(i)
		params["PrivateIpAddress."+n] = ip
	}
	resp = &SimpleResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
