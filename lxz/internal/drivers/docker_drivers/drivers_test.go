package docker_drivers

import "testing"

func TestListContainers(t *testing.T) {
	// Initialize Docker client
	if err := InitDockerClient(); err != nil {
		t.Fatalf("Failed to initialize Docker client: %v", err)
	}

	// List containers
	containers, err := ListContainers()
	if err != nil {
		t.Fatalf("Failed to list containers: %v", err)
	}

	// Check if containers are returned
	if len(containers) == 0 {
		t.Error("Expected to find at least one container, but found none")
	} else {
		t.Logf("Found %d containers", len(containers))
	}
}
