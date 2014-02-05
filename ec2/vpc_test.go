//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011 Canonical Ltd.
//
// Written by Gustavo Niemeyer <gustavo.niemeyer@canonical.com>
//

package ec2_test

import (
	"launchpad.net/goamz/ec2"
	. "launchpad.net/gocheck"
)

// VPC tests with example responses

func (s *S) TestCreateVPCExample(c *C) {
	testServer.Response(200, nil, CreateVpcExample)

	resp, err := s.ec2.CreateVPC("10.0.0.0/16", ec2.DefaultTenancy)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateVpc"})
	c.Assert(req.Form["CidrBlock"], DeepEquals, []string{"10.0.0.0/16"})
	c.Assert(req.Form["InstanceTenancy"], DeepEquals, []string{ec2.DefaultTenancy})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	vpc := resp.VPC
	c.Check(vpc.Id, Equals, "vpc-1a2b3c4d")
	c.Check(vpc.State, Equals, ec2.PendingState)
	c.Check(vpc.CIDRBlock, Equals, "10.0.0.0/16")
	c.Check(vpc.DHCPOptionsId, Equals, "dopt-1a2b3c4d2")
	c.Check(vpc.Tags, HasLen, 0)
	c.Check(vpc.IsDefault, Equals, false)
	c.Check(vpc.InstanceTenancy, Equals, ec2.DefaultTenancy)
}

func (s *S) TestDeleteVPCExample(c *C) {
	testServer.Response(200, nil, DeleteVpcExample)

	resp, err := s.ec2.DeleteVPC("vpc-id")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteVpc"})
	c.Assert(req.Form["VpcId"], DeepEquals, []string{"vpc-id"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestVPCsExample(c *C) {
	testServer.Response(200, nil, DescribeVpcsExample)

	resp, err := s.ec2.VPCs([]string{"vpc-1a2b3c4d"}, nil)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeVpcs"})
	c.Assert(req.Form["VpcId.1"], DeepEquals, []string{"vpc-1a2b3c4d"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	c.Assert(resp.VPCs, HasLen, 1)
	vpc := resp.VPCs[0]
	c.Check(vpc.Id, Equals, "vpc-1a2b3c4d")
	c.Check(vpc.State, Equals, ec2.AvailableState)
	c.Check(vpc.CIDRBlock, Equals, "10.0.0.0/23")
	c.Check(vpc.DHCPOptionsId, Equals, "dopt-7a8b9c2d")
	c.Check(vpc.Tags, HasLen, 0)
	c.Check(vpc.IsDefault, Equals, false)
	c.Check(vpc.InstanceTenancy, Equals, ec2.DefaultTenancy)
}

// VPC tests run against either a local test server or live on EC2.

func (s *ServerTests) TestVPCs(c *C) {
	resp1, err := s.ec2.CreateVPC("10.0.0.0/16", "")
	c.Assert(err, IsNil)
	assertVPC(c, resp1.VPC, "", "10.0.0.0/16")
	id1 := resp1.VPC.Id

	resp2, err := s.ec2.CreateVPC("1.2.0.0/18", ec2.DefaultTenancy)
	c.Assert(err, IsNil)
	assertVPC(c, resp2.VPC, "", "1.2.0.0/18")
	id2 := resp2.VPC.Id

	list, err := s.ec2.VPCs(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(list.VPCs, HasLen, 2)
	if list.VPCs[0].Id != id1 {
		list.VPCs[0], list.VPCs[1] = list.VPCs[1], list.VPCs[0]
	}
	assertVPC(c, list.VPCs[0], id1, resp1.VPC.CIDRBlock)
	assertVPC(c, list.VPCs[1], id2, resp2.VPC.CIDRBlock)

	list, err = s.ec2.VPCs([]string{id1}, nil)
	c.Assert(err, IsNil)
	c.Assert(list.VPCs, HasLen, 1)
	assertVPC(c, list.VPCs[0], id1, resp1.VPC.CIDRBlock)

	f := ec2.NewFilter()
	f.Add("cidr", resp2.VPC.CIDRBlock)
	list, err = s.ec2.VPCs(nil, f)
	c.Assert(err, IsNil)
	c.Assert(list.VPCs, HasLen, 1)
	assertVPC(c, list.VPCs[0], id2, resp2.VPC.CIDRBlock)

	_, err = s.ec2.DeleteVPC(id1)
	c.Assert(err, IsNil)
	_, err = s.ec2.DeleteVPC(id2)
	c.Assert(err, IsNil)
}

func assertVPC(c *C, obtained ec2.VPC, expectId, expectCidr string) {
	if expectId != "" {
		c.Check(obtained.Id, Equals, expectId)
	} else {
		c.Check(obtained.Id, Matches, `^vpc-[0-9a-f]+$`)
	}
	c.Check(obtained.State, Matches, "("+ec2.AvailableState+"|"+ec2.PendingState+")")
	if expectCidr != "" {
		c.Check(obtained.CIDRBlock, Equals, expectCidr)
	} else {
		c.Check(obtained.CIDRBlock, Matches, `^\d+\.\d+\.\d+\.\d+/\d+$`)
	}
	c.Check(obtained.DHCPOptionsId, Matches, `^dopt-[0-9a-f]+$`)
	c.Check(obtained.IsDefault, Equals, false)
	c.Check(obtained.Tags, HasLen, 0)
	c.Check(
		obtained.InstanceTenancy,
		Matches,
		"("+ec2.DefaultTenancy+"|"+ec2.DedicatedTenancy+")",
	)
}
