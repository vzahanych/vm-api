package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vzahanych/vm-api/core"
)

func createVM(vmRequest core.VMCreationRequest) (*core.VMCreationResponse, error) {
	url := "http://localhost:8080/vms" // The endpoint URL for VM creation
	body, err := json.Marshal(vmRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling VM creation request: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error sending request to create VM: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResponse core.ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResponse)
		return nil, fmt.Errorf("failed to create VM: %s", errResponse.Error.Message)
	}

	var createVMResponse core.VMCreationResponse
	if err := json.NewDecoder(resp.Body).Decode(&createVMResponse); err != nil {
		return nil, fmt.Errorf("error decoding create VM response: %v", err)
	}

	return &createVMResponse, nil
}

func getVMStatus(vmID string) (*core.GetVMStatusResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/vms/%s/status", vmID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request to get VM status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse core.ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResponse)
		return nil, fmt.Errorf("failed to get VM status: %s", errResponse.Error.Message)
	}

	var statusResponse core.GetVMStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		return nil, fmt.Errorf("error decoding get VM status response: %v", err)
	}

	return &statusResponse, nil
}

func deleteVM(vmID string) (*core.VMDeletionResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/vms/%s", vmID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating delete request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse core.ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResponse)
		return nil, fmt.Errorf("failed to delete VM: %s", errResponse.Error.Message)
	}

	var deletionResponse core.VMDeletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&deletionResponse); err != nil {
		return nil, fmt.Errorf("error decoding delete VM response: %v", err)
	}

	return &deletionResponse, nil
}

func main() {
	// Example VM creation request
	vmRequest := core.VMCreationRequest{
		VCPUs:    2,
		Memory:   4096,
		DiskSize: 20,
		BaseImage: "/var/lib/libvirt/images/ubuntu24.04-2.qcow2",
		CPUPinning: &core.CPUPinning{
			Cores: []int{0, 1},
		},
	}

	// Create a VM
	createResponse, err := createVM(vmRequest)
	if err != nil {
		log.Fatalf("Error creating VM: %v", err)
	}
	fmt.Printf("VM created successfully: %+v\n", createResponse)

	// Get VM status
	vmID := createResponse.VMID
	statusResponse, err := getVMStatus(vmID)
	if err != nil {
		log.Fatalf("Error getting VM status: %v", err)
	}
	fmt.Printf("VM Status: %+v\n", statusResponse)

	// Delete VM
	deletionResponse, err := deleteVM(vmID)
	if err != nil {
		log.Fatalf("Error deleting VM: %v", err)
	}
	fmt.Printf("VM Deleted: %+v\n", deletionResponse)
}
