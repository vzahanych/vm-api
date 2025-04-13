# KVM Management API

## Create VM endpoint

Requirements:

1. VM Provisioning: Create new VMs specifying vCPUs, memory (MB), and desired root disk size (GB). The service must handle the underlying storage creation (cloning from a base image and resizing).

2. Advanced Resource Control: Allow optional specification during VM creation for: Pinning VM vCPUs to specific host physical CPU cores (for performance-critical workloads). Applying I/O rate limits (e.g., IOPS) to VM disks (for predictable storage performance).

```yaml
paths:
  /vms:
    post:
      summary: Create a new virtual machine (VM)
      description: |
        This endpoint allows you to create a new VM by specifying vCPUs, memory, root disk size, and the base image.
        The service will handle the underlying storage creation by cloning from the base image and resizing the disk.
        It also generates a unique MAC address and defines the VM in libvirt.
        Optionally, advanced resource control such as CPU pinning or I/O limits can be specified.
      operationId: createVM
      tags:
        - VM Provisioning
      requestBody:
        description: Parameters for creating a new virtual machine, including resource allocation and disk size.
        content:
          application/json:
            schema:
              type: object
              properties:
                vcpus:
                  type: integer
                  example: 2
                  description: Number of virtual CPUs to be assigned to the new VM.
                memory:
                  type: integer
                  example: 4096
                  description: Amount of memory (in MB) to be allocated to the new VM.
                disk_size:
                  type: integer
                  example: 20
                  description: The desired root disk size for the VM in GB.
                base_image:
                  type: string
                  example: "/var/lib/libvirt/images/ubuntu-base.qcow2"
                  description: Path to the base image that will be cloned for the VM.
                cpu_pinning:
                  type: object
                  properties:
                    cores:
                      type: array
                      items:
                        type: integer
                      example: [0, 1]
                      description: List of physical CPU cores to pin the VM's vCPUs to.
                  description: Optional CPU pinning configuration.
                io_limits:
                  type: object
                  properties:
                    iops:
                      type: integer
                      example: 1000
                      description: The I/O operations per second (IOPS) limit for the VM's disk.
                  description: Optional I/O tuning for limiting the disk I/O.
              required:
                - vcpus
                - memory
                - disk_size
                - base_image
      responses:
        '201':
          description: VM successfully created
          content:
            application/json:
              schema:
                type: object
                properties:
                  type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                    description: Unique UUID identifier for the newly created VM.
                  status:
                    type: string
                    example: "created"
                    description: Status of the VM creation process.
                  vcpus:
                    type: integer
                    example: 2
                    description: Number of virtual CPUs assigned to the VM.
                  memory:
                    type: integer
                    example: 4096
                    description: Amount of memory allocated to the VM.
                  disk_size:
                    type: integer
                    example: 20
                    description: Disk size in GB allocated to the VM.
                  disk_file:
                    type: string
                    example: "/var/lib/libvirt/images/vm-12345.qcow2"
                    description: Path to the root disk image file.
                  mac_address:
                    type: string
                    example: "00:16:3e:2b:8a:9d"
                    description: The unique MAC address generated for the VM's network interface.
                  message:
                    type: string
                    example: "VM successfully created and storage cloned."
                    description: Confirmation message about the VM creation status.
        '400':
          description: Bad Request - Invalid parameters or missing required fields.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                      message:
                        type: string
                        example: "Invalid parameter: disk_size must be greater than 0."
        '404':
          description: Not Found - The specified base image does not exist.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                      message:
                        type: string
                        example: "Base image not found: /var/lib/libvirt/images/ubuntu-base.qcow2."
        '500':
          description: Internal Server Error - Failure in VM creation process (e.g., libvirt error).
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                      message:
                        type: string
                        example: "Failed to clone base image: /path/to/base-image.qcow2"
```

**Request Example** (POST /vms)

```json
{
  "vcpus": 2,
  "memory": 4096,
  "disk_size": 20,
  "base_image": "/var/lib/libvirt/images/ubuntu-base.qcow2",
  "cpu_pinning": {
    "cores": [0, 1]
  },
  "io_limits": {
    "iops": 1000
  }
}
```

**Response Example** (POST /vms)

```json
{
  "type": "uuid",
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "created",
  "vcpus": 2,
  "memory": 4096,
  "disk_size": 20,
  "disk_file": "/var/lib/libvirt/images/vm-12345.qcow2",
  "mac_address": "00:16:3e:2b:8a:9d",
  "message": "VM successfully created and storage cloned."
}
```

**Error Responses**:

```json
{
  "error": {
    "code": 400,
    "message": "Invalid parameter: disk_size must be greater than 0."
  }
}
```

```json
{
  "error": {
    "code": 404,
    "message": "Base image not found: /var/lib/libvirt/images/ubuntu-base.qcow2."
  }
}
```

```json
{
  "error": {
    "code": 500,
    "message": "Failed to clone base image: /path/to/base-image.qcow2"
  }
}
```

Comments:
For this VM creation endpoint, the key design choices are rooted in flexibility, clarity, and standard RESTful principles.

1. **UUID as ID**: Choosing a UUID for the VM’s ID ensures that each VM is uniquely identifiable, especially in distributed environments. Unlike auto-increment integers, UUIDs are globally unique and avoid collisions, which is essential when scaling across multiple systems or services. This also helps when clients or systems generate the ID externally, ensuring uniqueness.

2. **Resource Structure**: The VM creation request is structured around essential attributes: vcpus, memory, disk_size, and base_image. These are the most common and fundamental parameters for provisioning a VM. The additional parameters like cpu_pinning and io_limits are optional, providing flexibility for more advanced configurations based on user needs without overwhelming the API. Optional fields are nested under objects to clearly separate them and maintain a clean structure.

3. **Actions vs. State Updates**: The design follows standard RESTful practices where actions like VM creation (POST) are distinguished from state updates. In this case, creating a VM involves sending a POST request to /vms with required parameters, while any changes to an existing VM (e.g., starting, stopping) would be handled by other actions like PUT or POST to specific endpoints (e.g., /vms/{id}/start). The state of the VM after creation (e.g., its running status) is not part of the creation process and would be managed separately.

4. **Advanced Resource Parameters**: Advanced parameters like cpu_pinning and io_limits are optional because not all users will need them. They are nested as objects to provide clear separation and better extensibility for future resource controls. By structuring them this way, we allow for easy expansion of additional resource controls without cluttering the request body.

5. **Error Handling**: Clear and specific error codes (400 Bad Request, 404 Not Found, 500 Internal Server Error) ensure that users understand what went wrong. The use of structured error responses, with code and message, helps clients handle failures appropriately, ensuring that validation and backend issues are communicated effectively.

## Update VM configuration endpoints 

Usage of dedicated endpoints to update a VM  allows more granular controle over the whole update proccess.  

###  Update CPU Configuration

```yaml
paths:
  /vms/{id}/cpu:
    put:
      summary: Update the CPU configuration of a VM.
      description: |
        This endpoint updates the CPU configuration for a VM, such as changing the number of vCPUs and optionally pinning vCPUs to specific physical CPU cores.
      operationId: updateVMCpu
      tags:
        - VM Configuration
      parameters:
        - name: id
          in: path
          required: true
          description: The UUID of the VM to update.
          schema:
            type: string
            format: uuid
      requestBody:
        description: CPU configuration parameters for the VM.
        content:
          application/json:
            schema:
              type: object
              properties:
                vcpus:
                  type: integer
                  example: 4
                  description: The number of virtual CPUs to assign to the VM.
                cpu_pinning:
                  type: object
                  properties:
                    cores:
                      type: array
                      items:
                        type: integer
                      example: [0, 1]
                      description: List of physical CPU cores to pin the VM's vCPUs to.
      responses:
        '200':
          description: CPU configuration updated successfully.
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                    description: The UUID of the updated VM.
                  vcpus:
                    type: integer
                    example: 4
                    description: The number of virtual CPUs assigned to the VM.
                  cpu_pinning:
                    type: object
                    properties:
                      cores:
                        type: array
                        items:
                          type: integer
                        example: [0, 1]
                        description: The physical CPU cores that the vCPUs are pinned to.
                  status:
                    type: string
                    example: "updated"
                    description: The status of the CPU update operation.
                  message:
                    type: string
                    example: "CPU configuration successfully updated."
                    description: A message confirming the successful update.
        '400':
          description: |
            Bad request - invalid parameters.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 400).
            - `message`: A description of the error (e.g., "Invalid parameter: vcpus must be greater than 0").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Invalid parameter: vcpus must be greater than 0."
                        description: A message describing the error.
        '404':
          description: |
            VM not found - The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error.
        '500':
          description: |
            Internal error during update - The server encountered an error while attempting to update the VM's CPU configuration.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error during CPU update").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error during CPU update."
                        description: A message describing the error
```

Request Example (PUT /vms/{id}/cpu)

```json
{
  "vcpus": 4,
  "cpu_pinning": {
    "cores": [0, 1, 2, 3]
  }
}

```

Response Example (200 OK)

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "vcpus": 4,
  "cpu_pinning": {
    "cores": [0, 1, 2, 3]
  },
  "status": "updated",
  "message": "CPU configuration successfully updated."
}
```

Error Responses:

```json
{
  "error": {
    "code": 400,
    "message": "Invalid parameter: vcpus must be greater than 0."
  }
}
```

```json
{
  "error": {
    "code": 404,
    "message": "VM not found: 123e4567-e89b-12d3-a456-426614174000."
  }
}
```

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during CPU configuration update."
  }
}
```

###  Update Memory Configuration endpoint

```yaml
paths:
  /vms/{id}/memory:
    put:
      summary: Update the memory configuration of a VM.
      description: |
        This endpoint updates the memory configuration for a VM, such as increasing or decreasing the allocated memory (in MB).
      operationId: updateVMMemory
      tags:
        - VM Configuration
      parameters:
        - name: id
          in: path
          required: true
          description: The UUID of the VM to update.
          schema:
            type: string
            format: uuid
      requestBody:
        description: Memory configuration for the VM.
        content:
          application/json:
            schema:
              type: object
              properties:
                memory:
                  type: integer
                  example: 8192
                  description: The amount of memory (in MB) to assign to the VM.
      responses:
        '200':
          description: Memory configuration updated successfully.
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                    description: The UUID of the updated VM.
                  memory:
                    type: integer
                    example: 8192
                    description: The new amount of memory (in MB) allocated to the VM.
                  status:
                    type: string
                    example: "updated"
                    description: The status of the memory update operation.
                  message:
                    type: string
                    example: "Memory configuration successfully updated."
                    description: A message confirming the successful update of the memory configuration.
        '400':
          description: |
            Bad request - invalid parameters.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 400).
            - `message`: A description of the error (e.g., "Invalid parameter: memory must be greater than 0").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Invalid parameter: memory must be greater than 0."
                        description: A message describing the error.
        '404':
          description: |
            VM not found - The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error.
        '500':
          description: |
            Internal error during update - The server encountered an error while attempting to update the VM's memory configuration.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error during memory update").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error during memory update."
                        description: A message describing the error.
```

Request Example (PUT /vms/{id}/memory)

```json
{
  "memory": 8192
}
```

Response Example (200 OK)

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "memory": 8192,
  "status": "updated",
  "message": "Memory configuration successfully updated."
}
```

Error Responses

```json
{
  "error": {
    "code": 400,
    "message": "Invalid parameter: memory must be greater than 0."
  }
}
```

```json
{
  "error": {
    "code": 404,
    "message": "VM not found: 123e4567-e89b-12d3-a456-426614174000."
  }
}
```

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during memory configuration update."
  }
}
```

###  Update Disk Configuration endpoint

```yaml
paths:
  /vms/{id}/disk:
    put:
      summary: Resize the disk of a VM.
      description: |
        This endpoint resizes the root disk of a VM to the specified size (in GB). The disk will be resized, and the VM must be properly handled to ensure the disk resizing does not affect its functionality.
      operationId: resizeVMDisk
      tags:
        - VM Configuration
      parameters:
        - name: id
          in: path
          required: true
          description: The UUID of the VM to update.
          schema:
            type: string
            format: uuid
      requestBody:
        description: Disk resizing configuration for the VM.
        content:
          application/json:
            schema:
              type: object
              properties:
                disk_size:
                  type: integer
                  example: 50
                  description: The new size for the root disk in GB.
      responses:
        '200':
          description: Disk resized successfully.
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                    description: The UUID of the updated VM.
                  disk_size:
                    type: integer
                    example: 50
                    description: The new disk size in GB allocated to the VM.
                  status:
                    type: string
                    example: "updated"
                    description: The status of the disk resize operation.
                  message:
                    type: string
                    example: "Disk resized successfully."
                    description: A message confirming the successful disk resize.
        '400':
          description: |
            Bad request - invalid parameters.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 400).
            - `message`: A description of the error (e.g., "Invalid parameter: disk_size must be greater than 0").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Invalid parameter: disk_size must be greater than 0."
                        description: A message describing the error.
        '404':
          description: |
            VM not found - The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error.
        '500':
          description: |
            Internal error during disk resize - The server encountered an error while attempting to resize the VM's disk.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error during disk resize").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error during disk resize."
                        description: A message describing the error.
```

Request Example (PUT /vms/{id}/disk)

```json
{
  "disk_size": 50
}
```

Response Example (200 OK)

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "disk_size": 50,
  "status": "updated",
  "message": "Disk resized successfully."
}
```

Error Responses:

```json
{
  "error": {
    "code": 400,
    "message": "Invalid parameter: disk_size must be greater than 0."
  }
}
```

```json
{
  "error": {
    "code": 404,
    "message": "VM not found: 123e4567-e89b-12d3-a456-426614174000."
  }
}
```

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during disk resize operation."
  }
}
```

Comments:

1. The decision to create separate endpoints for CPU, memory, and disk configurations (/vms/{id}/cpu, /vms/{id}/memory, /vms/{id}/disk) was made to follow a modular and focused approach. This is in line with best practices observed in systems like ProxMox and VMware, where each resource type (CPU, memory, disk) can be updated independently. This provides clear separation of concerns, enabling flexible operations and reducing unnecessary complexity.
2. By having dedicated endpoints for each resource, users can update only the specific components they need without triggering unnecessary changes or consuming more resources.
3. The operations in these endpoints focus on state updates (e.g., modifying the VM’s memory or resizing the disk), not on triggering actions. Since the state of the VM (its CPU count, memory size, and disk size) is being changed, the PUT method is appropriate, as it is meant to replace the current resource state with a new one. These are idempotent operations, meaning multiple identical requests will yield the same result (no duplication or unexpected behavior), which aligns with the semantics of PUT.
4. For actions (like starting or stopping the VM), separate endpoints using POST would be used (as seen in the lifecycle management endpoints like /vms/{id}/start or /vms/{id}/stop). This distinction helps keep actions and state updates logically separate, maintaining a clean API structure.
5. For the CPU update endpoint, we allow optional CPU pinning as an advanced feature. This is useful for performance-critical workloads that need fine-grained control over which physical CPUs the virtual machine's vCPUs will be pinned to. By nesting CPU pinning under an object (e.g., "cpu_pinning": { "cores": [0, 1] }), we can cleanly extend this feature in the future, allowing users to configure additional parameters like CPU frequency or affinity.
6. Similarly, we could extend disk resizing and memory updates to include more advanced resource controls (e.g., setting I/O limits for disk operations), keeping the API extensible without cluttering the core resource definitions.
7. The design follows the principle of modularity to allow future scalability. For example, if the need arises to configure new resource controls (e.g., NUMA nodes, memory ballooning), these can easily be added to the existing endpoint without disrupting the current structure. Modularity ensures that each resource can evolve independently, which is essential for maintaining backward compatibility as the API grows.
8. If users require updates to other resources (such as network interfaces, disk partitioning, etc.), they can be handled in similarly structured endpoints (e.g., /vms/{id}/network, /vms/{id}/disks), following the same clear and granular structure.
9. The endpoints are named consistently to follow RESTful best practices, which improves readability and ease of use. For example:

/vms/{id}/cpu for CPU updates,

/vms/{id}/memory for memory updates,

/vms/{id}/disk for disk resizing.

10.The design is built with extensibility in mind. If new VM resource parameters are required (such as additional storage controls or network configurations), they can be easily added to these endpoints in the future without breaking existing functionality.

For example, the disk resize endpoint could later support additional parameters like disk type (e.g., SSD vs. HDD) or partitioning options by adding optional fields in the request body, while maintaining backward compatibility with the existing structure.


## Lifecycle Management endpoints

Requirements:
Lifecycle Management: Start, stop (gracefully preferred), and reboot VMs. Delete VMs, ensuring complete cleanup of associated resources (definition and storage).

### Start VM endpoint

```yaml
paths:
  /vms/{id}/start:
    post:
      summary: Start the specified VM.
      description: |
        Starts the VM. If the VM is already running, no action is taken (idempotent).
      operationId: startVM
      tags:
        - Lifecycle Management
      parameters:        
      responses:
        '200':
          description: VM successfully started, or already running.
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                    description: The UUID of the VM.
                  status:
                    type: string
                    example: "running"
                    description: The current status of the VM after attempting to start it.
                  message:
                    type: string
                    example: "VM successfully started, or already running."
                    description: A message confirming that the VM has been started or is already in a running state.
        '404':
          description: |
            VM not found - The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error.
        '500':
          description: |
            Internal error during start - The server encountered an error while attempting to start the VM.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error during VM start").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error during VM start."
                        description: A message describing the error.

```

Request Example (POST /vms/{id}/start)

```json
{}
```

Response Example (200 OK)

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "running",
  "message": "VM successfully started, or already running."
}
```

Error Responses:

```json
{
  "error": {
    "code": 404,
    "message": "VM not found: 123e4567-e89b-12d3-a456-426614174000."
  }
}
```

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during VM start operation."
  }
}
```

Comments:

1. Starting a VM is a state change in the VM's lifecycle, but it is more of an action than an update to the resource. POST is semantically more appropriate for operations that trigger actions rather than updating resource attributes.
2. URI Path: The URI /vms/{id}/start adheres to the common pattern of operating on specific resources (in this case, a VM) with the identifier {id}. The start action is directly related to the VM and is an appropriate sub-resource action for the specific VM identified by its UUID.
3. Idempotency: If the VM is already running, the operation should still return 200 OK as it is an idempotent action. This means the result is the same whether the VM was started or already running. This is common in many systems where certain actions (e.g., "start" or "reboot") can be safely repeated without causing side effects.
4. Consistency: The 200 OK response ensures that the operation has been processed correctly, even if the VM was already running, aligning with the idea that the "start" action does not create an error when no action is required.


### Stop VM endpoint

```yaml
paths:
  /vms/{id}/stop:
    post:
      summary: Gracefully stop the specified VM.
      description: |
        Gracefully stops the VM. If it cannot be stopped gracefully, a forced shutdown is attempted.
      operationId: stopVM
      tags:
        - Lifecycle Management
      parameters:
      responses:
        '200':
          description: VM successfully stopped, or already stopped (idempotent).
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                    description: The UUID of the VM that was stopped.
                  status:
                    type: string
                    example: "stopped"
                    description: The current status of the VM after attempting to stop it.
                  message:
                    type: string
                    example: "VM successfully stopped, or already stopped."
                    description: A message confirming that the VM has been stopped or was already stopped.
        '400':
          description: Bad Request - Invalid UUID format.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                      message:
                        type: string
                        example: "Invalid UUID format."
        '404':
          description: |
            VM not found - The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error.
        '500':
          description: |
            Internal error during start - The server encountered an error while attempting to start the VM.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error during VM start").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error during VM start."
                        description: A message describing the error.
```

Request Example (POST /vms/{id}/stop)

```json
{}
```

Response Example (200 OK)

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "stopped",
  "message": "VM successfully stopped."
}
```

Error Responses:

```json
{
  "error": {
    "code": 404,
    "message": "VM not found: 123e4567-e89b-12d3-a456-426614174000."
  }
}
```

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during VM stop operation."
  }
}
```

Comments:

1. Action vs. State Update: Stopping a VM is a state change in the VM's lifecycle but is considered more of an action than a configuration update to the resource itself. Using POST is semantically more appropriate for operations that trigger an action, such as stopping a VM, rather than directly modifying its attributes.

2. URI Path: The URI /vms/{id}/stop follows the common pattern for operating on specific resources, in this case, a VM, identified by its unique UUID ({id}). The stop action is a sub-resource action related to the specific VM, aligning with RESTful design principles.

3. Graceful Stop and Forced Stop: The POST request initiates the graceful shutdown of the VM, and if that fails, it will attempt a forced stop. This approach ensures that the operation gracefully handles various VM states, and it is idempotent in nature—repeated requests (such as stopping an already stopped VM) should not lead to errors but should return a successful response with no further action.

4. Idempotency and Consistency: If the VM is already stopped, the POST request should still return 200 OK. The operation is idempotent, meaning that attempting to stop a VM that is already stopped does not result in an error. This is consistent with common practices, where actions like "stop" can be safely repeated without causing additional side effects.

### Reboot VM endpoint

```yaml
paths:
  /vms/{id}/reboot:
    post:
      summary: Reboot the specified VM.
      description: |
        Gracefully reboots the VM. If necessary, performs a forced reboot.
      operationId: rebootVM
      tags:
        - Lifecycle Management
      parameters:
      responses:
        '200':
          description: VM successfully rebooted, or already rebooting.
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                    description: The UUID of the VM being rebooted.
                  status:
                    type: string
                    example: "rebooting"
                    description: The current status of the VM after the reboot operation.
                  message:
                    type: string
                    example: "VM successfully rebooted."
                    description: A message confirming that the VM has been rebooted or is in the process of rebooting.
        '404':
          description: |
            VM not found - The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error..
        '409':
          description: VM is not running.
        '500':
          description: |
            Internal error during start - The server encountered an error while attempting to start the VM.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error during VM start").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error during VM start."
                        description: A message describing the error.
```

Request Example (POST /vms/{id}/reboot)

```json
{}
```

Response Example (200 OK)

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "rebooting",
  "message": "VM successfully rebooted."
}
```

Error Responses:

```json
{
  "error": {
    "code": 404,
    "message": "VM not found: 123e4567-e89b-12d3-a456-426614174000."
  }
}
```

```json
{
  "error": {
    "code": 409,
    "message": "VM is not running, cannot reboot."
  }
}
```

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during VM reboot operation."
  }
}
```

Comments:

1. Action vs. State Update: Rebooting a VM is a state change in its lifecycle and is best classified as an action rather than an update to the VM's attributes. Using POST is appropriate because rebooting is an operation that triggers a backend process (the reboot) rather than directly modifying the configuration of the VM. The POST method reflects that this is an action to be performed on the resource, not an update to its state.

2. URI Path: The URI /vms/{id}/reboot adheres to the RESTful design pattern of performing actions on specific resources (the VM) identified by its unique UUID ({id}). The reboot action is clearly tied to the individual VM, making this an appropriate sub-resource operation for triggering the reboot process.

3. Graceful Reboot and Forced Reboot: The operation initiates a graceful reboot of the VM, and if that fails, a forced reboot is triggered. This ensures that the system makes an attempt to preserve the state of the VM by gracefully rebooting it, falling back to a forced reboot only if necessary. Like other lifecycle actions, rebooting is idempotent—the operation can be safely repeated without side effects, such as rebooting an already rebooting or running VM.

4. Idempotency and Consistency: If the VM is already in a state where a reboot is not possible (e.g., if the VM is not running), the operation should return a 409 Conflict to indicate that the action cannot be completed. Otherwise, if the reboot is initiated successfully, a 200 OK is returned, confirming the request has been processed. This behavior aligns with common practices in APIs where actions like "reboot" or "restart" can be retried without causing adverse effects.

5. Clear Success Response: The 200 OK status code is used when the reboot operation succeeds, providing clear feedback to the client that the VM has been successfully rebooted. If the VM is not running, a 409 Conflict ensures that the client knows the reboot cannot be performed. This clear status reporting helps clients understand the result of their requests.

### Delete VM endpoint

```yaml
paths:
  /vms/{id}:
    delete:
      summary: Delete the specified VM and its resources.
      description: |
        Deletes the VM, ensuring that it is stopped first and all associated resources (e.g., storage) are cleaned up.
      operationId: deleteVM
      tags:
        - Lifecycle Management
      parameters:
        - name: id
          in: path
          required: true
          description: The UUID of the VM to delete.
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: VM successfully deleted.
          content:
            application/json:
              schema:
                type: object
                properties:
                  vm_id:
                    type: string
                    description: The UUID of the deleted VM.
                  status:
                    type: string
                    description: The status of the VM deletion process (e.g., "deleted").
                  message:
                    type: string
                    description: Message describing the result (e.g., "VM successfully deleted").
                  disk_file:
                    type: string
                    description: Path to the deleted disk image file.
        '400':
          description: Bad Request - Invalid UUID format.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                      message:
                        type: string
                        example: "Invalid UUID format."
        '404':
          description: VM not found.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                      message:
                        type: string
                        example: "VM with ID {id} not found."
        '500':
          description: Internal error during deletion.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                      message:
                        type: string
                        example: "Internal error during VM deletion.".
```


`200` - VM successfully deleted:

```json
{
  "vm_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "deleted",
  "message": "VM successfully deleted.",
  "disk_file": "/var/lib/libvirt/images/123e4567-e89b-12d3-a456-426614174000.qcow2"
}
```

`400` - Bad Request (Invalid UUID format):

```json
{
  "error": {
    "code": 400,
    "message": "Invalid UUID format."
  }
}
```

`404` - VM Not Found:

```json
{
  "error": {
    "code": 404,
    "message": "VM with ID 123e4567-e89b-12d3-a456-426614174000 not found."
  }
}
```

`500` - Internal Error During Deletion:

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during VM deletion."
  }
}
```

Comments for Delete VM Endpoint:

1. Action vs. State Update: Deleting a VM is an action that removes the VM's definition and its associated resources (such as storage). The DELETE HTTP method is the correct choice because it clearly represents the operation of removing a resource. This is in line with RESTful principles, where DELETE signifies the complete removal of a resource from the system.

2. URI Path: The URI /vms/{id} targets the specific VM identified by its UUID ({id}). This follows the common RESTful design pattern of performing operations on specific resources, in this case, the VM, and aligns with the structure of the API.

3. Idempotency: the DELETE operation should be considered idempotent as long as the VM is stopped. If the VM has already been deleted, a repeated DELETE request should return 404 Not Found, indicating that the resource is no longer available to delete. If the VM is already stopped, but not deleted, the operation will succeed and the resource will be removed.

## VM Information

Requirements:

Retrieve configuration details (vCPUs, memory, disk size, assigned resource controls like pinning) and current status (e.g., running, stopped) for specific VMs. List existing VMs on the hypervisor.


### List All VMs endpoint

```yaml
paths:
  /vms:
    get:
      summary: List all virtual machines on the hypervisor.
      description: |
        This endpoint retrieves a list of all virtual machines (VMs) on the hypervisor, including basic information
        such as their UUID, name, and current status.
      operationId: listVMs
      tags:
        - VM Information
      responses:
        '200':
          description: Successfully retrieved the list of VMs.
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: string
                      format: uuid
                      example: "123e4567-e89b-12d3-a456-426614174000"
                      description: The unique identifier of the VM.
                    name:
                      type: string
                      example: "vm1"
                      description: The name of the VM.
                    status:
                      type: string
                      example: "running"
                      description: Current status of the VM (e.g., "running", "stopped").
        '500':
          description: |
            Internal error during start - The server encountered an error while attempting to start the VM.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error during VM start").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error during VM start."
                        description: A message describing the error.

```

Here are the JSON response examples for the `/vms` GET endpoint:

`200` - Successfully retrieved the list of VMs:

```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "vm1",
    "status": "running"
  },
  {
    "id": "987e6543-e89b-12d3-a456-426614174001",
    "name": "vm2",
    "status": "stopped"
  },
  {
    "id": "234e5678-e89b-12d3-a456-426614174002",
    "name": "vm3",
    "status": "paused"
  }
]
```

`500` - Internal error during the process:

```json
{
  "error": {
    "code": 500,
    "message": "Internal server error while retrieving VMs."
  }
}
```
Comments:

1. The endpoint /vms represents the collection of all VMs in the system. This is in line with RESTful principles, where collections of resources are accessed via pluralized URIs (e.g., /vms for the collection). The design follows the common practice of listing resources through a dedicated endpoint that returns the most basic information required for further actions or queries.

2. This endpoint is designed to return essential information for all VMs at once (e.g., UUID, name, and status). This high-level overview is useful in scenarios where the user or system requires a quick summary of the VM state without needing to retrieve detailed configuration data. For example, a client might want to quickly assess how many VMs are running or stopped on a hypervisor without retrieving full details for each VM.

3. By providing a simple list of VMs with high-level data, the endpoint allows the API to remain modular and scalable. The detailed configuration or resource-specific information (such as CPU pinning or disk I/O limits) is separated into other endpoints, like /vms/{id}/config. This design reduces unnecessary complexity and makes it easier to scale the API to support future resource details or new operations.

4. Each VM in the list includes basic but consistent fields (ID, name, and status), making the response uniform across all entries. This standardization makes it easier for consumers of the API to process the data, ensuring that every VM entry has the same structure, and can be iterated or manipulated in a predictable way.

5. The endpoint follows standard REST practices by using the GET method to list the resources (VMs) and by leveraging a clear, simple URI structure (/vms). This design is intuitive and easy to use, ensuring that it aligns with industry standards and is easily adoptable by developers familiar with RESTful APIs.

6. While this current design simply lists the VMs with minimal information, there is room to extend it for future needs. For instance, adding support for filtering, sorting, and pagination could help accommodate situations where the number of VMs is large. Pagination ensures that large datasets are delivered in manageable chunks, improving both the API's performance and the user experience.

### Get Configuration Details endpoint

```yaml
paths:
  /vms/{id}/config:
    get:
      summary: Retrieve configuration details of a specific VM.
      description: |
        This endpoint retrieves the configuration details for a specific VM, including its vCPUs, memory, disk size, 
        and any advanced resource controls (e.g., CPU pinning, I/O limits). The configuration details are set during VM creation.
      operationId: getVMConfig
      tags:
        - VM Information
      parameters:
        - name: id
          in: path
          required: true
          description: The unique identifier (UUID) of the VM.
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Successfully retrieved the VM configuration details.
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                  vcpus:
                    type: integer
                    example: 2
                  memory:
                    type: integer
                    example: 4096
                  disk_size:
                    type: integer
                    example: 20
                  cpu_pinning:
                    type: object
                    properties:
                      cores:
                        type: array
                        items:
                          type: integer
                        example: [0, 1]
                  io_limits:
                    type: object
                    properties:
                      iops:
                        type: integer
                        example: 1000
        '400':
          description: Bad Request - Invalid UUID format.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                      message:
                        type: string
                        example: "Invalid UUID format."
        '404':
          description: VM not found.
        '500':
          description: Internal server error while retrieving VM configuration details.

```

`200` - Successfully retrieved the VM configuration details:

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "vcpus": 2,
  "memory": 4096,
  "disk_size": 20,
  "cpu_pinning": {
    "cores": [0, 1]
  },
  "io_limits": {
    "iops": 1000
  }
}
```
`400` - Bad Request (Invalid UUID format):

```json
{
  "error": {
    "code": 400,
    "message": "Invalid UUID format."
  }
}
```

`404` - VM Not Found:

```json
{
  "error": {
    "code": 404,
    "message": "VM with ID 123e4567-e89b-12d3-a456-426614174000 not found."
  }
}
```

`500` - Internal Server Error:

```json
{
  "error": {
    "code": 500,
    "message": "Internal server error while retrieving VM configuration details."
  }
}
```

Comments:

1. GET is the appropriate HTTP method because it is used to retrieve data without causing any changes to the resource. This endpoint allows the user to access essential configuration information about a VM, such as vCPUs, memory, disk size, and resource controls like CPU pinning or I/O limits.
2. The path /vms/{id}/config targets a specific VM identified by its UUID ({id}) and retrieves the configuration details. By placing /config as a sub-resource of /vms/{id}, this structure clearly distinguishes configuration details from other operational data (e.g., status or lifecycle actions), in line with common RESTful practices. This naming convention provides clarity and reflects that the endpoint specifically deals with the configuration of the VM.
3. The configuration details include advanced parameters such as CPU pinning and I/O limits, which are optional but valuable for performance-sensitive workloads. These parameters are included in the response as part of the VM configuration to provide a complete and detailed specification of how the VM is set up. By structuring advanced resource parameters like cpu_pinning and io_limits as objects, the API allows for easy extensibility in the future—new parameters can be added in a structured way without disrupting the existing design.
4. The GET operation for retrieving VM configuration details is idempotent. This means that no matter how many times the client requests the same details, the result will be the same and will not cause any changes to the VM. This is important for consistency and reliability in APIs—clients can call the endpoint repeatedly without side effects or unexpected changes to the resource. Additionally, if the VM does not exist, a 404 Not Found status is returned, indicating that the resource is absent, ensuring that clients can handle errors appropriately.

### Get VM Status endpoint

```yaml
paths:
  /vms/{id}/status:
    get:
      summary: Retrieve the current status of a specific VM.
      description: |
        This endpoint retrieves the current status of a specific VM, such as whether it is running, stopped, paused, etc.
        The status represents the VM's current lifecycle state.
      operationId: getVMStatus
      tags:
        - VM Information
      parameters:
        - name: id
          in: path
          required: true
          description: The unique identifier (UUID) of the VM.
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Successfully retrieved the VM's status.
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                    example: "123e4567-e89b-12d3-a456-426614174000"
                  name:
                    type: string
                    example: "vm1"
                    description: The name of the VM.
                  status:
                    type: string
                    example: "running"
                    description: Current status of the VM (e.g., "running", "stopped").
                  message:
                    type: string
                    example: "VM is running"
                    description: A message describing the current state of the VM.
      '400':
        description: |
          Bad Request - Invalid UUID format for the VM ID.
          The error response will include a `code` and a `message` field:
          - `code`: The error code (e.g., 400).
          - `message`: A description of the error (e.g., "Invalid UUID format").
        content:
          application/json:
            schema:
              type: object
              properties:
                error:
                  type: object
                  properties:
                    code:
                      type: integer
                      example: 400
                      description: The error code for the error.
                    message:
                      type: string
                      example: "Invalid UUID format"
                      description: A message describing the error.
        '404':
          description: |
            The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404 but it might be custom error code).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error.
        '500':
          description: |
            An internal server error occurred while processing the request.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 500).
            - `message`: A description of the error (e.g., "Internal server error").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 500
                        description: The error code for the error.
                      message:
                        type: string
                        example: "Internal server error"
                        description: A message describing the error.
```

Here are the corresponding JSON responses for the specified endpoints (`200`, `400`, `404`, and `500`):

`200` - Successfully retrieved the VM's status:

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "vm1",
  "status": "running",
  "message": "VM is running"
}
```

`400` - Bad Request - Invalid UUID format:

```json
{
  "error": {
    "code": 400,
    "message": "Invalid UUID format"
  }
}
```

`404` - VM not found:

```json
{
  "error": {
    "code": 404,
    "message": "VM not found"
  }
}
```

`500` - Internal server error:

```json
{
  "error": {
    "code": 500,
    "message": "Internal server error"
  }
}
```

These JSON responses match the descriptions and provide appropriate error or success information for each status code. Let me know if you need further modifications or additional details!

Comments for Get VM Status Endpoint:

1. The GET method for retrieving the VM status is a read-only operation that provides the current state of the VM without making any changes to the VM itself. The GET HTTP method is appropriate because it adheres to RESTful principles, where GET is used to fetch data without causing side effects. The operation is designed solely to return the current operational state of the VM, such as whether it is running, stopped, paused, or in another state.

2. The path /vms/{id}/status is used to target the status of a specific VM identified by its UUID ({id}). The use of /status as a sub-resource indicates that the endpoint focuses on querying the VM’s current operational state, distinguishing it from other endpoints that deal with the configuration or lifecycle of the VM. This clear and concise URI structure follows RESTful conventions and ensures that users can easily identify the purpose of the endpoint.

3. By isolating status from configuration, this design adheres to the principle of separation of concerns. VM configuration details (e.g., vCPUs, memory, disk size) are managed separately through other endpoints like /vms/{id}/config, while the status endpoint focuses purely on the current state of the VM. This clean separation improves API clarity and helps prevent confusion between static configuration and dynamic operational states.

4. The status of the VM is represented by a simple string value (e.g., "running", "stopped"). This makes the status response lightweight and easy to process. Providing the status as a single string makes the response predictable, allowing clients to quickly check the state of the VM without needing to parse complex data structures.

5. The GET operation is idempotent. Regardless of how many times the client queries the status of a VM, the result will remain the same unless the VM’s status changes. This is a fundamental property of GET requests, ensuring that the state of the resource remains unchanged during the request. If the VM exists and the status is successfully retrieved, the client will always receive the same response unless the state of the VM has been modified externally.


## Monitoring

Monitoring: Provide a way to query real-time performance metrics for a running VM (specifically CPU usage and memory usage).

### VM performance endpoint

```yaml
paths:
  /vms/{id}/performance:
    get:
      summary: Retrieve comprehensive real-time performance metrics for the specified VM.
      description: |
        This endpoint provides a detailed set of performance metrics for a running VM, including CPU usage, memory usage, disk I/O statistics, network statistics, and more. The metrics are returned as a JSON object for easy parsing and integration.
      operationId: getVMPerformanceMetrics
      tags:
        - Monitoring
      parameters:
        - name: id
          in: path
          required: true
          description: The UUID of the VM for which performance metrics are being queried.
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Real-time performance metrics successfully retrieved.
          content:
            application/json:
              schema:
                type: object
                properties:
                  cpu_usage:
                    type: number
                    format: float
                    example: 45.5
                    description: The current CPU usage percentage.
                  memory_usage:
                    type: integer
                    example: 2048
                    description: The current memory usage in MB.
                  disk_read_bytes:
                    type: integer
                    example: 1048576
                    description: Total bytes read from disk.
                  disk_write_bytes:
                    type: integer
                    example: 524288
                    description: Total bytes written to disk.
                  disk_read_iops:
                    type: integer
                    example: 150
                    description: Disk read operations per second.
                  disk_write_iops:
                    type: integer
                    example: 75
                    description: Disk write operations per second.
                  network_in_bytes:
                    type: integer
                    example: 1024000
                    description: Total incoming network traffic in bytes.
                  network_out_bytes:
                    type: integer
                    example: 512000
                    description: Total outgoing network traffic in bytes.
                  network_in_packets:
                    type: integer
                    example: 1200
                    description: Incoming network packets per second.
                  network_out_packets:
                    type: integer
                    example: 800
                    description: Outgoing network packets per second.
        '400':
          description: Bad Request - Invalid UUID format.
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 400
                      message:
                        type: string
                        example: "Invalid UUID format."
        '404':
          description: |
            The specified virtual machine (VM) was not found.
            The error response will include a `code` and a `message` field:
            - `code`: The error code (e.g., 404 but it might be custom error code).
            - `message`: A description of the error (e.g., "VM not found").
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: object
                    properties:
                      code:
                        type: integer
                        example: 404
                        description: The error code for the error.
                      message:
                        type: string
                        example: "VM not found"
                        description: A message describing the error.
        '500':
          description: Internal error during performance metrics retrieval.
```

`200` - Real-time performance metrics successfully retrieved:

```json
{
  "cpu_usage": 45.5,
  "memory_usage": 2048,
  "disk_read_bytes": 1048576,
  "disk_write_bytes": 524288,
  "disk_read_iops": 150,
  "disk_write_iops": 75,
  "network_in_bytes": 1024000,
  "network_out_bytes": 512000,
  "network_in_packets": 1200,
  "network_out_packets": 800
}
```

`400` - Bad Request (Invalid UUID format):

```json
{
  "error": {
    "code": 400,
    "message": "Invalid UUID format."
  }
}
```
`404` - VM Not Found:

```json
{
  "error": {
    "code": 404,
    "message": "VM not found"
  }
}
```

`500` - Internal Server Error:

```json
{
  "error": {
    "code": 500,
    "message": "Internal error during performance metrics retrieval."
  }
}
```

Comments:

1. The GET method for retrieving the VM performance metrics is a read-only operation that provides real-time data about the VM’s resource utilization without making any changes to the VM itself. The GET HTTP method is appropriate because it adheres to RESTful principles, where GET is used to fetch data without causing side effects. This operation focuses on returning the current performance metrics such as CPU usage, memory usage, disk I/O, and network I/O, providing users with a comprehensive snapshot of the VM’s health.

2. The path /vms/{id}/performance is used to target the performance metrics of a specific VM identified by its UUID ({id}). The use of /performance as a sub-resource indicates that this endpoint is dedicated to querying dynamic performance data rather than static configuration or lifecycle management. This clear and concise URI structure follows RESTful conventions, making it easy for users to identify the endpoint’s purpose — retrieving performance data for a specific VM.

3. By aggregating multiple performance metrics into a single endpoint, this design aims to simplify the process of retrieving essential data. Users can get all key performance indicators (KPIs) — CPU, memory, disk, and network statistics — in one request. This reduces the need for multiple API calls, improving both efficiency and usability. Additionally, providing all metrics in a single response ensures that users have a comprehensive view of the VM's performance in real time.

4. The response is structured to provide each metric in a predictable and easy-to-process format, such as cpu_usage, memory_usage, disk_read_bytes, network_in_bytes, etc. This structure makes it easy for clients to parse and analyze the data. Each metric is represented with an appropriate unit of measurement (e.g., bytes, percentage), and the use of clear labels ensures that the information is easily understood without requiring further processing or documentation lookup.

5. The GET operation is idempotent. Regardless of how many times the client queries the performance metrics, the result will remain the same unless the VM’s performance data changes due to system workload fluctuations. This is a fundamental property of GET requests, ensuring that repeated requests will consistently yield the same information, as long as the VM’s performance does not change externally. This idempotency guarantees that users can safely query the performance data multiple times without unexpected side effects.

6. The endpoint is designed to be scalable and extensible. As new performance metrics become relevant or as the system evolves, it is straightforward to add additional fields to the response (e.g., disk latency, load averages) without requiring changes to the overall API structure. This design allows the endpoint to grow with the evolving needs of the system while maintaining backward compatibility for existing users.

7. Error handling is implemented with standard HTTP status codes, ensuring that clients can easily understand the outcome of their requests. A 200 OK response indicates successful retrieval of performance data, while a 404 Not Found error is returned if the specified VM is not found or not running. A 500 Internal Server Error response helps indicate any server-side issues during the performance data retrieval process.

8. This design keeps the performance metrics request lightweight and focused. Instead of offering unnecessary details, the endpoint provides only the most relevant performance data for a VM, avoiding the overhead of sending excessive information that might not be needed for typical use cases. This makes it efficient for monitoring and real-time analysis purposes, aligning with the needs of developers and system administrators who require quick access to key metrics.