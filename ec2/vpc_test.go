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

func (s *S) TestCreateVpcExample(c *C) {
	testServer.Response(200, nil, CreateVpcExample)

	resp, err := s.ec2.CreateVpc("10.0.0.0/16")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateVpc"})
	c.Assert(req.Form["CidrBlock"], DeepEquals, []string{"10.0.0.0/16"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	vpc := resp.Vpc
	c.Check(vpc.Id, Equals, "vpc-1a2b3c4d")
	c.Check(vpc.State, Equals, "pending")
	c.Check(vpc.CidrBlock, Equals, "10.0.0.0/16")
	c.Check(vpc.DhcpOptionsId, Equals, "dopt-1a2b3c4d2")
	c.Check(vpc.Tags, HasLen, 0)
}

func (s *S) TestDeleteVpcExample(c *C) {
	testServer.Response(200, nil, DeleteVpcExample)

	resp, err := s.ec2.DeleteVpc("vpc-id")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteVpc"})
	c.Assert(req.Form["VpcId"], DeepEquals, []string{"vpc-id"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestDescribeVpcsExample(c *C) {
	testServer.Response(200, nil, DescribeVpcsExample)

	resp, err := s.ec2.DescribeVpcs([]string{"vpc-1a2b3c4d"}, nil)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeVpcs"})
	c.Assert(req.Form["VpcId.1"], DeepEquals, []string{"vpc-1a2b3c4d"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	c.Check(resp.Vpcs, HasLen, 1)
	vpc := resp.Vpcs[0]
	c.Check(vpc.Id, Equals, "vpc-1a2b3c4d")
	c.Check(vpc.State, Equals, "available")
	c.Check(vpc.CidrBlock, Equals, "10.0.0.0/23")
	c.Check(vpc.DhcpOptionsId, Equals, "dopt-7a8b9c2d")
	c.Check(vpc.Tags, HasLen, 0)
}

// VPC tests run against either a local test server or live on EC2.

func (s *ServerTests) TestVpcs(c *C) {
	resp1, err := s.ec2.CreateVpc("10.0.0.0/16")
	c.Assert(err, IsNil)
	assertVpc(c, resp1.Vpc, "", "10.0.0.0/16")
	id1 := resp1.Vpc.Id

	resp2, err := s.ec2.CreateVpc("1.2.0.0/18")
	c.Assert(err, IsNil)
	assertVpc(c, resp2.Vpc, "", "1.2.0.0/18")
	id2 := resp2.Vpc.Id

	list, err := s.ec2.DescribeVpcs(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(list.Vpcs, HasLen, 2)
	if list.Vpcs[0].Id != id1 {
		list.Vpcs[0], list.Vpcs[1] = list.Vpcs[1], list.Vpcs[0]
	}
	assertVpc(c, list.Vpcs[0], id1, resp1.Vpc.CidrBlock)
	assertVpc(c, list.Vpcs[1], id2, resp2.Vpc.CidrBlock)

	list, err = s.ec2.DescribeVpcs([]string{id1}, nil)
	c.Assert(err, IsNil)
	c.Assert(list.Vpcs, HasLen, 1)
	assertVpc(c, list.Vpcs[0], id1, resp1.Vpc.CidrBlock)

	f := ec2.NewFilter()
	f.Add("cidr", resp2.Vpc.CidrBlock)
	list, err = s.ec2.DescribeVpcs(nil, f)
	c.Assert(err, IsNil)
	c.Assert(list.Vpcs, HasLen, 1)
	assertVpc(c, list.Vpcs[0], id2, resp2.Vpc.CidrBlock)

	_, err = s.ec2.DeleteVpc(id1)
	c.Assert(err, IsNil)
	_, err = s.ec2.DeleteVpc(id1)
	c.Assert(err, ErrorMatches, `.*\(InvalidVpcID.NotFound\)`)
	_, err = s.ec2.DeleteVpc("invalid-id")
	c.Assert(err, ErrorMatches, `.*\(InvalidVpcID.NotFound\)`)
	_, err = s.ec2.DeleteVpc(id2)
	c.Assert(err, IsNil)
}

func assertVpc(c *C, obtained ec2.Vpc, expectId, expectCidr string) {
	if expectId != "" {
		c.Check(obtained.Id, Equals, expectId)
	} else {
		c.Check(obtained.Id, Matches, `^vpc-[0-9a-f]+$`)
	}
	c.Check(obtained.State, Matches, "(available|pending)")
	if expectCidr != "" {
		c.Check(obtained.CidrBlock, Equals, expectCidr)
	} else {
		c.Check(obtained.CidrBlock, Matches, `^\d+\.\d+\.\d+\.\d+/\d+$`)
	}
	c.Check(obtained.DhcpOptionsId, Matches, `^dopt-[0-9a-f]+$`)
}
