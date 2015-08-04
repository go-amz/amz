//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2015 Canonical Ltd.
//

package ec2_test

import (
	. "gopkg.in/check.v1"
)

// Internet Gateway tests with example responses

func (s *S) TestInternetGatewaysExample(c *C) {
	testServer.Response(200, nil, DescribeInternetGatewaysExample)

	ids := []string{"igw-eaad4883EXAMPLE"}
	resp, err := s.ec2.InternetGateways(ids, nil)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeInternetGateways"})
	c.Assert(req.Form["InternetGatewayId.1"], DeepEquals, ids)

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
	c.Assert(resp.InternetGateways, HasLen, 1)
	igw := resp.InternetGateways[0]
	c.Check(igw.Id, Equals, "igw-eaad4883EXAMPLE")
	c.Check(igw.VPCId, Equals, "vpc-11ad4878")
	c.Check(igw.AttachmentState, Equals, "available")
	c.Check(igw.Tags, HasLen, 0)
}
