package core

import (
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "libvirt.org/go/libvirt"
)

// CreateVM creates a virtual machine based on the VMCreationRequest
func createVM(c *gin.Context, request *VMCreationRequest, lq LibvirtQemu) (*VMCreationResponse, error) {
	// Generate a new UUID for the VM
	vmID := uuid.New().String()
	diskPath := fmt.Sprintf("/var/lib/libvirt/images/%s.qcow2", vmID)

	// Clone the base image and resize it
	err := lq.CloneAndResizeDisk(request.BaseImage, diskPath, request.DiskSize, true)
	if err != nil {
		return  nil, fmt.Errorf("failed to clone and resize disk: %v", err)
	}

	// Initialize the CPU pinning XML section as an empty string by default
	cpuPinningXML := ""

	// If CPU pinning is specified, modify the CPU section
	if request.CPUPinning != nil {
		// Create the CPU pinning section
		pinningXML := "<cpu mode='host-passthrough' check='full'><numactrls>"
		for _, core := range request.CPUPinning.Cores {
			pinningXML += fmt.Sprintf("<numactrl id='%d'/>", core)
		}
		pinningXML += "</numactrls></cpu>"

		// Set the CPU pinning section
		cpuPinningXML = pinningXML
	}

	// Set the VM's XML configuration with a placeholder for CPU pinning
	xmlConfig := fmt.Sprintf(`
<domain type='kvm'>
  <name>%s</name>
  <uuid>%s</uuid>
  <memory unit='KiB'>%d</memory>
  <vcpu placement='static'>%d</vcpu>
  <os>
    <type arch='x86_64' machine='pc-i440fx-2.9'>hvm</type>
    <boot dev='hd'/>
  </os>
  <devices>
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2'/>
      <source file='%s'/>
      <target dev='vda' bus='virtio'/>
      <address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x0'/>
    </disk>
    <interface type='network'>
      <mac address='%s'/>
      <source network='default'/>
      <model type='virtio'/>
      <address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
    </interface>
  </devices>
  %s
</domain>`, vmID, vmID, request.Memory*1024, request.VCPUs, diskPath, generateMACAddress(), cpuPinningXML)

	// Create the VM from the generated XML configuration
	domain, err := lq.DomainDefineXML(xmlConfig)
	if err != nil {
		return  nil, fmt.Errorf("failed to define the domain: %v", err)
	}

	// Start the VM
	if err := lq.Create(domain); err != nil {
		return  nil, fmt.Errorf("failed to start the domain: %v", err)
	}

	// If I/O limits are specified, apply them
	if request.IOLimits != nil {
		// Implement I/O limits logic using libvirt's API (optional step, based on your libvirt's capabilities)
		// This can be done by adjusting the disk's I/O limits after the VM creation
	}

	// Create the response object with the relevant details
	response := &VMCreationResponse{
		VMID:       vmID,
		Status:     "running",
		VCPUs:      request.VCPUs,
		Memory:     request.Memory,
		DiskSize:   request.DiskSize,
		DiskFile:   diskPath,
		MacAddress: generateMACAddress(),
		Message:    "VM successfully created and storage cloned",
	}

	// Return the response object and nil error (success)
	return response, nil
}

// generateMACAddress generates a unique MAC address from the UUID
// TODO: generate it from private range according to the requirements
func generateMACAddress() string {
	// Generate a new UUID
	uuid := uuid.New()

	// MAC address format: 00:16:3e:XX:XX:XX, where XX comes from the UUID
	return fmt.Sprintf("00:16:3e:%02x:%02x:%02x", uuid[0], uuid[1], uuid[2])
}

// cloneAndResizeDisk clones the base image and resizes the cloned disk
func cloneAndResizeDisk(baseImage string, newDiskPath string, diskSizeGB int, shrink bool) error {
	// Step 1: Clone the base image into a new disk file
	// We use qemu-img to create a QCOW2 file from the base image
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", "-F", "qcow2", "-b", baseImage, newDiskPath)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to clone the base image: %v", err)
	}

	// Step 2: Resize the disk to the desired size
	var resizeCmd *exec.Cmd
	if shrink {
		// If we need to shrink the disk, add the --shrink flag
		resizeCmd = exec.Command("qemu-img", "resize", "--shrink", newDiskPath, fmt.Sprintf("%dG", diskSizeGB))
	} else {
		// Otherwise, just resize without shrinking
		resizeCmd = exec.Command("qemu-img", "resize", newDiskPath, fmt.Sprintf("%dG", diskSizeGB))
	}

	err := resizeCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to resize the disk: %v", err)
	}

	return nil
}
