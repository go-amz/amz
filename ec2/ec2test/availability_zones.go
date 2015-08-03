//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011-2015 Canonical Ltd.
//
// This file contains code handling AWS API for availability zones.

package ec2test

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"gopkg.in/amz.v3/ec2"
)

// SetAvailabilityZones initializes the availability zones the test
// server is returning.
func (srv *Server) SetAvailabilityZones(zones []ec2.AvailabilityZoneInfo) {
	srv.mu.Lock()
	srv.zones = make([]availabilityZone, len(zones))
	for i, z := range zones {
		srv.zones[i] = availabilityZone{z}
	}
	srv.mu.Unlock()
}

type availabilityZone struct {
	ec2.AvailabilityZoneInfo
}

func (z *availabilityZone) matchAttr(attr, value string) (ok bool, err error) {
	switch attr {
	case "message":
		for _, m := range z.MessageSet {
			if m == value {
				return true, nil
			}
		}
		return false, nil
	case "region-name":
		return z.Region == value, nil
	case "state":
		switch value {
		case "available", "impaired", "unavailable":
			return z.State == value, nil
		}
		return false, fmt.Errorf("invalid state %q", value)
	case "zone-name":
		return z.Name == value, nil
	}
	return false, fmt.Errorf("unknown attribute %q", attr)
}

func (srv *Server) describeAvailabilityZones(w http.ResponseWriter, req *http.Request, reqId string) interface{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	f := newFilter(req.Form)
	var resp struct {
		XMLName xml.Name
		ec2.AvailabilityZonesResp
	}
	resp.XMLName = xml.Name{defaultXMLName, "DescribeAvailabilityZonesResponse"}
	resp.RequestId = reqId
	for _, zone := range srv.zones {
		ok, err := f.ok(&zone)
		if ok {
			resp.Zones = append(resp.Zones, zone.AvailabilityZoneInfo)
		} else if err != nil {
			fatalf(400, "InvalidParameterValue", "describe availability zones: %v", err)
		}
	}
	return &resp
}
