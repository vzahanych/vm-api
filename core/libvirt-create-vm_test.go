package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vzahanych/vm-api/core/mocks"
	"go.uber.org/mock/gomock"
	"libvirt.org/go/libvirt"
)

// TestCreateVM tests the createVM function
func TestCreateVM(t *testing.T) {
	// Step 1: Setup gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Step 2: Create a mock of the LibvirtQemu interface
	mockLibvirt := mocks.NewMockLibvirtQemu(ctrl)

	// Step 3: Define the expected behavior for CloneAndResizeDisk
	baseImage := "/var/lib/libvirt/images/ubuntu-base.qcow2"
	diskSizeGB := 20
	mockLibvirt.EXPECT().
		CloneAndResizeDisk(baseImage, gomock.Any(), diskSizeGB, true).
		Return(nil) // Simulate successful cloning and resizing

	// Step 4: Define the expected behavior for DomainDefineXML
	// We will now use a regular assertion inside the Do method to check if the XML contains the necessary fields
	mockLibvirt.EXPECT().
		DomainDefineXML(gomock.Any()). // We expect any XML here, we will check it dynamically
		Do(func(xmlConfig string) {
			// Validate that the XML contains the necessary tags
			assert.Contains(t, xmlConfig, "<name>")
			assert.Contains(t, xmlConfig, "<uuid>")
			assert.Contains(t, xmlConfig, "<mac address='00:16:3e:") // Check MAC address format
		}).
		Return(&libvirt.Domain{}, nil).Times(1) // Simulate successful domain definition

	// Step 5: Define the expected behavior for Create
	// We expect the Create method to be called after the VM is defined
	mockLibvirt.EXPECT().
		Create(gomock.Any()). // Expect the Create method to be called with any domain
		Return(nil).          // Simulate successful creation
		Times(1)              // Expect it to be called exactly once

	// Step 6: Prepare the VMCreationRequest with a valid body
	request := &VMCreationRequest{
		VCPUs:    2,
		Memory:   4096,
		DiskSize: 20,
		BaseImage: "/var/lib/libvirt/images/ubuntu-base.qcow2",
		CPUPinning: &CPUPinning{
			Cores: []int{0, 1},
		},
	}

	// Step 7: Call the function to test
	vmResponse, err := createVM(nil, request, mockLibvirt)

	// Step 8: Assert the results
	assert.Nil(t, err)
	assert.NotNil(t, vmResponse)
	assert.Equal(t, "running", vmResponse.Status)
	assert.Equal(t, 2, vmResponse.VCPUs)
	assert.Equal(t, 4096, vmResponse.Memory)
	assert.Equal(t, 20, vmResponse.DiskSize)
	assert.Contains(t, vmResponse.DiskFile, "qcow2") // Check that the disk file path contains ".qcow2"
	assert.NotEmpty(t, vmResponse.MacAddress) // Ensure the MAC address is generated
}




