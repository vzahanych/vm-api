package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/vzahanych/vm-api/core/mocks"
	"go.uber.org/mock/gomock"
	"libvirt.org/go/libvirt"
)

// TestGetVMStatus tests the getVMstatus function
func TestGetVMStatus(t *testing.T) {
	// Step 1: Set up the gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Step 2: Create a mock of the LibvirtQemu interface
	mockLibvirt := mocks.NewMockLibvirtQemu(ctrl)

	// Step 3: Define the expected behavior of the mock for domain lookup and state retrieval
	vmID := "123e4567-e89b-12d3-a456-426614174000"
	mockLibvirt.EXPECT().
		LookupDomainByUUIDString(vmID).
		Return(&libvirt.Domain{}, nil). // Mock successful domain lookup
		Times(1)

	mockLibvirt.EXPECT().
		GetName(gomock.Any()).
		Return("testVM", nil). // Mock the VM name retrieval
		Times(1)

	mockLibvirt.EXPECT().
		GetState(gomock.Any()).
		Return(libvirt.DOMAIN_RUNNING, nil). // Mock the VM state as running
		Times(1)

	// Step 4: Call the function under test
	response, err := getVMstatus(nil, vmID, mockLibvirt)

	// Step 5: Assert the expected results
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response.VMID)
	assert.Equal(t, "testVM", response.Name)
	assert.Equal(t, "running", response.Status)
	assert.Contains(t, response.Message, "VM testVM is running")
}
