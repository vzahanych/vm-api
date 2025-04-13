package core

import (
	"fmt"
	"os/exec"

	"libvirt.org/go/libvirt"
)

//go:generate mockgen -source=libvirt-qemu.go -destination=mocks/libvirt_qemu_mock.go -package=mocks

// LibvirtQemu defines the interface for libvirt/qemu managment
type LibvirtQemu interface {
	NewConnect(uri string) error
	LookupDomainByUUIDString(uuid string) (*libvirt.Domain, error)
	CloneAndResizeDisk(baseImage string, newDiskPath string, diskSizeGB int, shrink bool) error
	DomainDefineXML(xmlConfig string) (*libvirt.Domain, error)
	Create(domain *libvirt.Domain) error
	GetName(domain *libvirt.Domain) (string, error)
	GetState(domain *libvirt.Domain) (libvirt.DomainState, error)
	Shutdown(domain *libvirt.Domain) error
	Destroy(domain *libvirt.Domain) error
	Undefine(domain *libvirt.Domain) error
}

type LibvirtQemuImpl struct {
	conn *libvirt.Connect
}

// NewConnect creates a connection to the libvirt daemon at the given URI
func (l *LibvirtQemuImpl) NewConnect(uri string) (err error) {
	l.conn, err = libvirt.NewConnect(uri)
	if err != nil {
		return fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	return nil
}

func (l *LibvirtQemuImpl) DomainDefineXML(xmlConfig string) (*libvirt.Domain, error) {
	return l.conn.DomainDefineXML(xmlConfig)
}

func (l *LibvirtQemuImpl) LookupDomainByUUIDString(id string) (*libvirt.Domain, error) {
	return l.conn.LookupDomainByUUIDString(id)
}

// CloneAndResizeDisk clones a base image and resizes the cloned disk
func (l *LibvirtQemuImpl) CloneAndResizeDisk(baseImage string, newDiskPath string, diskSizeGB int, shrink bool) error {
	// Clone the base image using qemu-img
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", "-F", "qcow2", "-b", baseImage, newDiskPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone the base image: %v", err)
	}

	// Resize the cloned disk
	var resizeCmd *exec.Cmd
	if shrink {
		resizeCmd = exec.Command("qemu-img", "resize", "--shrink", newDiskPath, fmt.Sprintf("%dG", diskSizeGB))
	} else {
		resizeCmd = exec.Command("qemu-img", "resize", newDiskPath, fmt.Sprintf("%dG", diskSizeGB))
	}

	if err := resizeCmd.Run(); err != nil {
		return fmt.Errorf("failed to resize the disk: %v", err)
	}

	return nil
}

func (l *LibvirtQemuImpl) Create(domain *libvirt.Domain) error {
	err := domain.Create()
	if err != nil {
		return fmt.Errorf("failed to create domain: %v", err)
	}
	return nil
}

// GetName retrieves the name of the domain (VM)
func (l *LibvirtQemuImpl) GetName(domain *libvirt.Domain) (string, error) {
	name, err := domain.GetName()
	if err != nil {
		return "", fmt.Errorf("failed to get the domain name: %v", err)
	}
	return name, nil
}

// GetState retrieves the state of the domain (VM)
func (l *LibvirtQemuImpl) GetState(domain *libvirt.Domain) (libvirt.DomainState, error) {
	state, _, err := domain.GetState()
	if err != nil {
		return -1, fmt.Errorf("failed to get the domain state: %v", err)
	}
	return state, nil
}

// Shutdown gracefully shuts down the domain (VM)
func (l *LibvirtQemuImpl) Shutdown(domain *libvirt.Domain) error {
	err := domain.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to gracefully shut down the domain: %v", err)
	}
	return nil
}

// Destroy forcefully stops the domain (VM)
func (l *LibvirtQemuImpl) Destroy(domain *libvirt.Domain) error {
	err := domain.Destroy()
	if err != nil {
		return fmt.Errorf("failed to forcefully stop the domain: %v", err)
	}
	return nil
}

// Undefine removes the domain (VM) definition from libvirt
func (l *LibvirtQemuImpl) Undefine(domain *libvirt.Domain) error {
	err := domain.Undefine()
	if err != nil {
		return fmt.Errorf("failed to undefine the domain: %v", err)
	}
	return nil
}
