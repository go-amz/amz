//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2014 Canonical Ltd.
//

package ec2_test

import (
	"time"

	. "gopkg.in/check.v1"

	"gopkg.in/amz.v4-unstable/aws"
	"gopkg.in/amz.v4-unstable/ec2"
	"strconv"
)

// Volume tests with example responses

func (s *S) TestCreateVolumeExample(c *C) {
	testServer.Response(200, nil, CreateVolumeExample)

	volumeToCreate := ec2.CreateVolume{
		AvailZone:  "us-east-1a",
		VolumeType: "ssd",
		VolumeSize: 10,
		IOPS:       3000,
		Encrypted:  true,
	}
	resp, err := s.ec2.CreateVolume(volumeToCreate)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateVolume"})
	c.Assert(req.Form["AvailabilityZone"], DeepEquals, []string{"us-east-1a"})
	c.Assert(req.Form["VolumeType"], DeepEquals, []string{"ssd"})
	c.Assert(req.Form["Size"], DeepEquals, []string{"10"})
	c.Assert(req.Form["Iops"], DeepEquals, []string{"3000"})
	c.Assert(req.Form["Encrypted"], DeepEquals, []string{"true"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
	volume := resp.Volume
	c.Check(volume.Id, Equals, "vol-1a2b3c4d")
	c.Check(volume.AvailZone, Equals, "us-east-1a")
	c.Check(volume.Status, Equals, "creating")
	c.Check(volume.VolumeType, Equals, "standard")
	c.Check(volume.Size, Equals, 80)
	c.Check(volume.IOPS, Equals, int64(3000))
	c.Check(volume.Encrypted, Equals, true)
	c.Check(volume.Tags, HasLen, 0)
}

func (s *S) TestDeleteVolumeExample(c *C) {
	testServer.Response(200, nil, DeleteVolumeExample)

	resp, err := s.ec2.DeleteVolume("volume-id")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteVolume"})
	c.Assert(req.Form["VolumeId"], DeepEquals, []string{"volume-id"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
}

func (s *S) TestVolumesExample(c *C) {
	testServer.Response(200, nil, DescribeVolumesExample)

	ids := []string{"vol-1a2b3c4d"}
	resp, err := s.ec2.Volumes(ids, nil)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeVolumes"})
	c.Assert(req.Form["VolumeId.1"], DeepEquals, []string{ids[0]})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
	c.Check(resp.Volumes, HasLen, 1)
	volume := resp.Volumes[0]
	c.Check(volume.Id, Equals, "vol-1a2b3c4d")
	c.Check(volume.AvailZone, Equals, "us-east-1a")
	c.Check(volume.Status, Equals, "in-use")
	c.Check(volume.VolumeType, Equals, "standard")
	c.Check(volume.Size, Equals, 80)
	c.Check(volume.IOPS, Equals, int64(3000))
	c.Check(volume.Encrypted, Equals, true)
	c.Check(volume.Tags, HasLen, 0)
	attachments := volume.Attachments
	c.Check(attachments, HasLen, 1)
	attachment := attachments[0]
	c.Check(attachment.VolumeId, Equals, ids[0])
	c.Check(attachment.Status, Equals, "attached")
	c.Check(attachment.Device, Equals, "/dev/sdh")
	c.Check(attachment.DeleteOnTermination, Equals, false)
}

func (s *S) TestAttachVolumeExample(c *C) {
	testServer.Response(200, nil, AttachVolumeExample)

	resp, err := s.ec2.AttachVolume("volume-id", "instance-id", "device")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"AttachVolume"})
	c.Assert(req.Form["VolumeId"], DeepEquals, []string{"volume-id"})
	c.Assert(req.Form["InstanceId"], DeepEquals, []string{"instance-id"})
	c.Assert(req.Form["Device"], DeepEquals, []string{"device"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
	c.Check(resp.VolumeId, Equals, "vol-1a2b3c4d")
	c.Check(resp.Device, Equals, "/dev/sdh")
	c.Check(resp.InstanceId, Equals, "i-1a2b3c4d")
	c.Check(resp.Status, Equals, "attaching")
}

func (s *S) TestDetachVolumeExample(c *C) {
	testServer.Response(200, nil, DetachVolumeExample)

	resp, err := s.ec2.DetachVolume("volume-id", "instance-id", "device", true)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DetachVolume"})
	c.Assert(req.Form["VolumeId"], DeepEquals, []string{"volume-id"})
	c.Assert(req.Form["InstanceId"], DeepEquals, []string{"instance-id"})
	c.Assert(req.Form["Device"], DeepEquals, []string{"device"})
	c.Assert(req.Form["Force"], DeepEquals, []string{"true"})

	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
	c.Check(resp.VolumeId, Equals, "vol-1a2b3c4d")
	c.Check(resp.Device, Equals, "/dev/sdh")
	c.Check(resp.InstanceId, Equals, "i-1a2b3c4d")
	c.Check(resp.Status, Equals, "detaching")
}

// Volume tests run against either a local test server or live on EC2.

func (s *ServerTests) TestVolumes(c *C) {
	vol1 := ec2.CreateVolume{
		AvailZone:  "us-east-1b",
		VolumeType: "standard",
		VolumeSize: 20,
	}
	resp1, err := s.ec2.CreateVolume(vol1)
	c.Assert(err, IsNil)
	assertVolume(c, resp1.Volume, "", vol1.VolumeType, vol1.AvailZone, 20, vol1.IOPS)
	id1 := resp1.Volume.Id

	vol2 := ec2.CreateVolume{
		AvailZone:  "us-east-1b",
		VolumeType: "io1",
		VolumeSize: 101,
		IOPS:       3030,
	}
	resp2, err := s.ec2.CreateVolume(vol2)
	c.Assert(err, IsNil)
	assertVolume(c, resp2.Volume, "", vol2.VolumeType, vol2.AvailZone, 101, vol2.IOPS)
	id2 := resp2.Volume.Id

	// We only check for the Volumes we just created, because the user
	// might have others in his account (when testing against the EC2
	// servers). In some cases it takes a short while until both Volumes
	// are created, so we need to retry a few times to make sure.
	var list *ec2.VolumesResp
	done := false
	testAttempt := aws.AttemptStrategy{
		Total: 2 * time.Minute,
		Delay: 5 * time.Second,
	}
	for a := testAttempt.Start(); a.Next(); {
		c.Logf("waiting for %v to be created", []string{id1, id2})
		list, err = s.ec2.Volumes(nil, nil)
		if err != nil {
			c.Logf("retrying; Volumes returned: %v", err)
			continue
		}
		found := 0
		for _, vol := range list.Volumes {
			c.Logf("found Volume %v", vol)
			switch vol.Id {
			case id1:
				assertVolume(c, vol, id1, vol1.VolumeType, vol1.AvailZone, 20, vol1.IOPS)
				found++
			case id2:
				assertVolume(c, vol, id2, vol2.VolumeType, vol2.AvailZone, 101, vol2.IOPS)
				found++
			}
			if found == 2 {
				done = true
				break
			}
		}
		if done {
			c.Logf("all Volumes were created")
			break
		}
	}
	if !done {
		c.Fatalf("timeout while waiting for Volumes %v", []string{id1, id2})
	}

	list, err = s.ec2.Volumes([]string{id1}, nil)
	c.Assert(err, IsNil)
	c.Assert(list.Volumes, HasLen, 1)
	assertVolume(c, list.Volumes[0], id1, vol1.VolumeType, vol1.AvailZone, 20, vol1.IOPS)

	f := ec2.NewFilter()
	f.Add("size", strconv.Itoa(resp2.Volume.Size))
	list, err = s.ec2.Volumes(nil, f)
	c.Assert(err, IsNil)
	c.Assert(list.Volumes, HasLen, 1)
	assertVolume(c, list.Volumes[0], id2, vol2.VolumeType, vol2.AvailZone, 101, vol2.IOPS)

	_, err = s.ec2.DeleteVolume(id1)
	c.Assert(err, IsNil)
	_, err = s.ec2.DeleteVolume(id2)
	c.Assert(err, IsNil)
}

func assertVolume(c *C, obtained ec2.Volume, expectId, expectType, availZone string, expectSize int, expectIOPS int64) {
	if expectId != "" {
		c.Check(obtained.Id, Equals, expectId)
	} else {
		c.Check(obtained.Id, Matches, `^vol-[0-9a-f]+$`)
	}
	c.Check(obtained.VolumeType, Equals, expectType)
	c.Check(obtained.AvailZone, Equals, availZone)
	if expectSize > 0 {
		c.Check(obtained.Size, Equals, expectSize)
	}
	if expectIOPS > 0 {
		c.Check(obtained.IOPS, Equals, expectIOPS)
	}
	c.Check(obtained.Status, Matches, "(creating|available)")
	c.Check(obtained.Encrypted, Equals, false)
	c.Check(obtained.Tags, HasLen, 0)
}

// Volume Attachment tests run against either a local test server or live on EC2.

func (s *ServerTests) TestVolumeAttachments(c *C) {
	vol1 := ec2.CreateVolume{
		AvailZone:  "us-east-1d",
		VolumeType: "standard",
		VolumeSize: 20,
	}
	resp1, err := s.ec2.CreateVolume(vol1)
	c.Assert(err, IsNil)
	volId := resp1.Id

	// Create an instance to attach the volume to.
	instList, err := s.ec2.RunInstances(&ec2.RunInstances{
		ImageId:      imageId,
		InstanceType: "m1.medium",
		AvailZone:    "us-east-1d",
	})
	c.Assert(err, IsNil)
	inst := instList.Instances[0]
	c.Assert(inst, NotNil)
	instId := inst.InstanceId
	defer terminateInstances(c, s.ec2, []string{instId})

	// Instance needs to be running before attaching volume.
	testAttempt := aws.AttemptStrategy{
		Total: 5 * time.Minute,
		Delay: 5 * time.Second,
	}
	var resp2 *ec2.VolumeAttachmentResp
	for a := testAttempt.Start(); a.Next(); {
		resp2, err = s.ec2.AttachVolume(volId, instId, "/dev/sdb")
		if err != nil {
			c.Logf("AttachVolume returned: %v; retrying...", err)
			continue
		}
		if resp2 != nil {
			break
		}
	}
	if resp2 == nil {
		c.Fatalf("timeout while waiting for the instance to be running")
	}
	assertVolumeAttachment(c, resp2, volId, instId, "/dev/sdb")

	_, err = s.ec2.DetachVolume(volId, "", "", false)
	c.Assert(err, IsNil)
}

func assertVolumeAttachment(c *C, obtained *ec2.VolumeAttachmentResp, volId, instanceId, device string) {
	c.Check(obtained.VolumeId, Equals, volId)
	c.Check(obtained.InstanceId, Equals, instanceId)
	c.Check(obtained.Device, Equals, device)
	c.Check(obtained.Status, Matches, "(attaching|attached)")
}
