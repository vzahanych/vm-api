package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/vzahanych/vm-api/core/mocks"
	"go.uber.org/mock/gomock"
	"libvirt.org/go/libvirt"
	"fmt"
)

// TestDeleteVM tests the DeleteVM function
func TestDeleteVM(t *testing.T) {
	// Step 1: Set up the gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Step 2: Create mocks for LibvirtQemu
	mockLibvirt := mocks.NewMockLibvirtQemu(ctrl)

	// VM ID for testing
	vmID := "123e4567-e89b-12d3-a456-426614174000"
	diskPath := fmt.Sprintf("/var/lib/libvirt/images/%s.qcow2", vmID)

	// Step 3: Define the expected behavior for the mock methods
	// Mock the behavior for LookupDomainByUUIDString
	mockLibvirt.EXPECT().
		LookupDomainByUUIDString(vmID).
		Return(&libvirt.Domain{}, nil).Times(1)

	// Mock graceful shutdown of the VM
	mockLibvirt.EXPECT().
		Shutdown(gomock.Any()).
		Return(nil).Times(1)

	// Mock the state retrieval to ensure the VM is shut down
	mockLibvirt.EXPECT().
		GetState(gomock.Any()).
		Return(libvirt.DOMAIN_SHUTOFF, nil).Times(1)

	// Mock undefining the domain
	mockLibvirt.EXPECT().
		Undefine(gomock.Any()).
		Return(nil).Times(1)

	// Step 4: Call the function under test
	response, err := DeleteVM(nil, vmID, mockLibvirt)

	// Step 6: Assert the expected results
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "deleted", response.Status)
	assert.Equal(t, diskPath, response.DiskFile)
	assert.Contains(t, response.Message, "VM successfully deleted")
}


