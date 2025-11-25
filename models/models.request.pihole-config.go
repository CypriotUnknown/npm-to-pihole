package models

import "gitlab.com/CypriotUnknown/npm-to-pihole/constants"

type PiholeConfigRequest struct {
	Cname string
	Type  constants.PiholeRequestType
}
