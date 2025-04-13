package core

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VMDeletionResponse represents the response structure for VM deletion
type VMDeletionResponse struct {
	VMID     string `json:"vm_id"`     // Unique UUID identifier for the deleted VM
	Status   string `json:"status"`    // Status of the VM deletion process (e.g., "deleted", "error")
	Message  string `json:"message"`   // Message describing the deletion result (success or failure)
	DiskFile string `json:"disk_file"` // Path to the deleted disk image file
}

func DeleteVMHandler(c *gin.Context) {
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

	logger = logger.With("endpoint", "delete vm")
	// Extract the VM ID from the URL path (e.g., /vms/{id})
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

	// Call the DeleteVM function from libvirt to stop, undefine, and delete the disk
	response, err := DeleteVM(c, vmID, lq)
	if err != nil {
		logger.Error("Failed to delete VM", "error", err)

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

	// Return successful deletion response
	c.JSON(http.StatusOK, response)
}
