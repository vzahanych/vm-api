package core

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// VMCreationRequest represents the body of a request to create a new VM.
type VMCreationRequest struct {
	VCPUs      int         `json:"vcpus"`                 // Number of virtual CPUs to be assigned to the new VM.
	Memory     int         `json:"memory"`                // Amount of memory (in MB) to be allocated to the new VM.
	DiskSize   int         `json:"disk_size"`             // The desired root disk size for the VM in GB.
	BaseImage  string      `json:"base_image"`            // Path to the base image that will be cloned for the VM.
	CPUPinning *CPUPinning `json:"cpu_pinning,omitempty"` // Optional CPU pinning configuration.
	IOLimits   *IOLimits   `json:"io_limits,omitempty"`   // Optional I/O tuning for limiting disk I/O.
}

// CPUPinning represents the optional CPU pinning configuration for the VM.
type CPUPinning struct {
	Cores []int `json:"cores"` // List of physical CPU cores to pin the VM's vCPUs to.
}

// IOLimits represents the optional I/O limits configuration for the VM's disk.
type IOLimits struct {
	IOPS int `json:"iops"` // The I/O operations per second (IOPS) limit for the VM's disk.
}

// VMCreationResponse represents the response structure for VM creation
type VMCreationResponse struct {
	VMID       string `json:"vm_id"`       // Unique UUID identifier for the created VM
	Status     string `json:"status"`      // Status of the VM creation process (e.g., "created")
	VCPUs      int    `json:"vcpus"`       // Number of virtual CPUs assigned to the VM
	Memory     int    `json:"memory"`      // Amount of memory (in MB) allocated to the VM
	DiskSize   int    `json:"disk_size"`   // Disk size in GB allocated to the VM
	DiskFile   string `json:"disk_file"`   // Path to the root disk image file
	MacAddress string `json:"mac_address"` // The unique MAC address generated for the VM's network interface
	Message    string `json:"message"`     // Confirmation message about the VM creation status
}

// Handler to create VM
func CreateVMHandler(c *gin.Context) {
	var request VMCreationRequest

	// Get the logger from the Gin context
	l, ok := c.Get("logger")
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			},
		})
	}

	logger, ok := l.(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			},
		})
	}

	logger = logger.With("endpoint", "create vm")

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&request); err != nil {
		// Log the error
		logger.Error("Failed to bind request", "error", err)

		// Validation failed, return 400 Bad Request
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// Check if the base image exists
	if _, err := os.Stat(request.BaseImage); os.IsNotExist(err) {
		// Base image not found, return 404 Not Found
		logger.Error("Base image not found", "base_image", request.BaseImage)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Base image not found: %s", request.BaseImage),
			},
		})
		return
	}

	lq := &LibvirtQemuImpl{}
	if err := lq.NewConnect("qemu:///system"); err != nil {
		logger.Error("Fail to connect to libvirt", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			},
		})
		return
	}

	var res *VMCreationResponse
	var err error
	if res, err = createVM(c, &request, lq); err != nil {
		// Failure in VM creation process, return 500 Internal Server Error
		logger.Error("Failed to create VM", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}
