package core

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetVMStatusResponse represents the response structure for retrieving a specific VM's status
type GetVMStatusResponse struct {
	VMID    string `json:"vm_id"`   // Unique UUID identifier of the VM
	Name    string `json:"name"`    // The name of the VM
	Status  string `json:"status"`  // The current state of the VM (e.g., "running", "stopped")
	Message string `json:"message"` // Message describing the current state of the VM (e.g., "VM is running")
}

// GetVMStatus handles retrieving the status of a specific VM based on its ID
func GetVMStatus(c *gin.Context) {
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

	logger = logger.With("endpoint", "vm status")

	// Extract the VM ID from the URL path (e.g., /vms/{id}/status)
	vmID := c.Param("id")

	// Validate that the VM ID is a valid UUID
	_, err := uuid.Parse(vmID)
	if err != nil {
		// If not a valid UUID, return a 400 Bad Request error
		// Log the error
		logger.Error("Failed to parse id", "error", err)

		// Validation failed, return 400 Bad Request
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
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

	response, err := getVMstatus(c, vmID, lq)
	if err != nil {
		logger.Error("Failed to get VM status", "error", err)

		if _, ok := err.(*NotFoundError); ok {
            c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetails{
					Code:    http.StatusNotFound,
					Message: "VM not found",
				},
			})
			return			
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			},
		})
		return
	}

	// Return the response
	c.JSON(http.StatusOK, response)
}
