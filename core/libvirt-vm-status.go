package core

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"libvirt.org/go/libvirt"
)

func getVMstatus(c *gin.Context, vmID string, lq LibvirtQemu) (*GetVMStatusResponse, error) {
	// Lookup the domain (VM) by UUID
	domain, err := lq.LookupDomainByUUIDString(vmID)
	if err != nil {
		if er, ok := err.(libvirt.Error); ok {
            if er.Code == libvirt.ERR_NO_DOMAIN {
				return nil, NewNotFoundError(vmID)
			}			
		}
		return nil, err
	}

	// Retrieve the VM's name and state
	name, err := lq.GetName(domain)
	if err != nil {
		return nil, err
	}

	// Retrieve the VM's current state (running, stopped, etc.)
	state, err := lq.GetState(domain)
	if err != nil {
		return nil, err
	}

	// Convert the state to a string (e.g., "running", "stopped")
	var status string
	switch state {
	case libvirt.DOMAIN_RUNNING:
		status = "running"
	case libvirt.DOMAIN_SHUTOFF:
		status = "stopped"
	case libvirt.DOMAIN_PAUSED:
		status = "paused"
	default:
		status = "unknown"
	}

	// Prepare the response
	return &GetVMStatusResponse{
		VMID:    vmID,
		Name:    name,
		Status:  status,
		Message: fmt.Sprintf("VM %s is %s", name, status),
	}, nil
}
