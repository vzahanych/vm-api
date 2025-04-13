package core

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"libvirt.org/go/libvirt"
)

// DeleteVM deletes the virtual machine, stops it (gracefully or forced), undefines it, and deletes the associated disk file
func DeleteVM(c *gin.Context, vmID string, lq LibvirtQemu) (*VMDeletionResponse, error) {
	// Delete the associated disk file
	diskFile := fmt.Sprintf("/var/lib/libvirt/images/%s.qcow2", vmID)
	
	defer func() {
		// Check if the file exists
		if _, err := os.Stat(diskFile); err == nil {
			// File exists, attempt to remove it
			err = os.Remove(diskFile)
			if err != nil {
				// Log the error but still report success for the VM deletion (disk deletion failure is logged)
				log.Printf("Failed to delete the disk file %s: %v", diskFile, err)
			}
		}
	}()

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

	// Stop the VM gracefully
	err = lq.Shutdown(domain)
	if err != nil {
		// If graceful shutdown fails, forcefully stop the VM
		log.Printf("Graceful shutdown failed, forcing stop for VM %s", vmID)
		err = lq.Destroy(domain) // Force stop
		if err != nil {
			return nil, fmt.Errorf("failed to forcefully stop the domain: %v", err)
		}
	}

	// Step 2: Ensure the VM is stopped (check state)
	state, err := lq.GetState(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM state: %v", err)
	}

	// If VM is still running after forceful stop, return an error
	if state != libvirt.DOMAIN_SHUTOFF {
		err = lq.Destroy(domain) // Force stop
		if err != nil {
			return nil, fmt.Errorf("failed to forcefully stop the domain: %v", err)
		}
	}
	// Undefine the domain (remove from libvirt)
	err = lq.Undefine(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to undefine domain: %v", err)
	}

	// Return the VM deletion response
	response := &VMDeletionResponse{
		VMID:     vmID,
		Status:   "deleted",
		Message:  "VM successfully deleted",
		DiskFile: diskFile,
	}

	// Return response indicating success
	return response, nil
}
