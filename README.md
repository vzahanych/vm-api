# VM Management API

This project provides a PoC for a RESTful API for managing KVM virtual machines (VMs). The API allows for the creation, deletion, lifecycle management, and querying of VM configurations. It also supports advanced features like CPU pinning, I/O rate limiting, and performance monitoring.

In addition to the RESTful API, we use Cobra to facilitate the creation of CLI tools for managing VMs, enabling easy extension of functionality. This allows developers and administrators to interact with the system directly from the command line, making it easier to automate tasks and integrate with other services.

## Features

1. **VM Provisioning**:
   - Create VMs with customizable vCPUs, memory, and disk sizes.
   - Advanced options like CPU pinning and I/O rate limiting are supported.

2. **Lifecycle Management**:
   - Start, stop (gracefully or forcefully), and reboot VMs.
   - Delete VMs, ensuring proper cleanup of resources.

3. **VM Information**:
   - Retrieve detailed configuration of specific VMs (vCPUs, memory, disk size, etc.).
   - Query real-time performance metrics like CPU and memory usage.

4. **Monitoring**:
   - Query performance metrics like CPU usage, memory usage, disk I/O, and network statistics.

## Prerequisites

Ensure that you have the following installed:

- Go 1.24 (Recommended version: 1.24 or later)
- `qemu-img` (for managing qcow2 images)
- libvirt installed and configured
- Make (for building and managing dependencies)

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/vzahanych/vm-api.git
   cd https://github.com/vzahanych/vm-api.git
   ```

2. Install the necessary Go dependencies:

   ```bash
   make deps
   ```

3. Run tests:

   ```bash
   make test
   ```

4. (Optional) Run static analysis and build the binary:

   ```bash
   make build
   ```

5. The API runs on `http://localhost:8080` by default.

## Base Image Assumptions

1. Base Image Format: The API assumes that the base images are in the QCOW2 format (.qcow2), which is commonly used with KVM for its features like snapshotting and compression.

2. Image Path: The base image should be located on the local file system of the KVM hypervisor, and the path to the base image must be provided when creating a VM.

3. Base Image Requirements: Ensure that the base image you provide has a valid operating system installation (such as Ubuntu, CentOS, etc.) that is compatible with your VM configuration. The API does not create the base image for you; it expects an existing image to be available at the specified path.

4. Storage Location: The default storage path for base images in this API is /var/lib/libvirt/images/. If you use a different storage path, ensure that the correct path to the base image is provided when creating the VM.

5. Permissions: The user running the API must have appropriate permissions to access and read the base image from the specified path. If the permissions are not set correctly, the VM creation process will fail.

6. Base Image Example: A sample base image could be "/var/lib/libvirt/images/ubuntu24.04-2.qcow2", but ensure that the actual base image used matches your setup.

## Create VM

To create a VM, send a `POST` request to `/vms` with the following example payload:

```bash
curl -X POST http://localhost:8080/vms \
    -H "Content-Type: application/json" \
    -d '{
        "vcpus": 2,
        "memory": 4096,
        "disk_size": 20,
        "base_image": "/var/lib/libvirt/images/ubuntu24.04-2.qcow2",
        "cpu_pinning": {
            "cores": [0, 1]
        },
        "io_limits": {
            "iops": 1000
        }
    }'
```

## Get VM Status

To retrieve the status of a VM, send a `GET` request to `/vms/{id}/status`:

```bash
curl -X GET http://localhost:8080/vms/c00b825f-630e-41df-86bb-e77efa314d7d/status
```

## Delete VM

To delete a VM, send a `DELETE` request to `/vms/{id}`:

```bash
curl -X DELETE http://localhost:8080/vms/c00b825f-630e-41df-86bb-e77efa314d7d
```

## Example

To demonstrate how to interact with the VM Management API, we have provided an example Go client in the file `./cmd/client/main.go`. This client showcases how to perform operations such as creating, deleting, and retrieving VM status via the API.

### Client Example (`./cmd/client/main.go`)

The Go client makes use of the `net/http` package to interact with the API endpoints. Below is a basic example that shows how to create a VM, get its status, and delete it.

#### Steps to run the Go client:

1. **Navigate to the client directory**:

   ```bash
   cd ./cmd/client
   ```

2. **Install necessary dependencies**:

   If you haven’t already installed the necessary Go modules for the client, run:

   ```bash
   go mod tidy
   ```

3. **Run the client**:

   To run the client, use the following command:

   ```bash
   go run main.go
   ```

This example demonstrates the following API operations:

- **Create a VM**: Sends a `POST` request to the `/vms` endpoint with the required parameters (vCPUs, memory, disk size, base image, etc.).
- **Get VM Status**: Sends a `GET` request to `/vms/{id}/status` to retrieve the current status of the created VM.
- **Delete the VM**: Sends a `DELETE` request to `/vms/{id}` to delete the VM after retrieving its status.

The Go client is a simple, illustrative example of how to interact with the API. It can be expanded with more features, such as error handling, retry logic, or more complex VM management tasks.


