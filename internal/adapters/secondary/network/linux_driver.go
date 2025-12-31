package network

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

type LinuxDriver struct{}

func NewLinuxDriver() *LinuxDriver {
	return &LinuxDriver{}
}

func (l *LinuxDriver) ApplyResource(ctx context.Context, r *domain.Resource) error {
	// since this will only be applied to interface
	if r.Kind != "INTERFACE" {
		fmt.Println("The LinuxDriver is unable to handle resources others than Interface")
		return nil
	}
	fmt.Printf("-> [Linux Driver] Configuring Interface %s .......\n ", r.ID)
	// we will avoid command injections
	cmd := exec.CommandContext(ctx, "sudo", "ip", "link", "add", r.ID, "type", "dummy")
	output, StdError := cmd.CombinedOutput()
	if StdError != nil {
		return fmt.Errorf("command faile %s : %s ", StdError, string(output))
	}
	fmt.Printf(" [Linux Driver] Successfully created interface %s \n", r.ID)
	return nil
}

func (l *LinuxDriver) DeleteResource(ctx context.Context, r *domain.Resource) error {
	if r.Kind != "INTERFACE" {
		return nil
	}

	cmd := exec.CommandContext(ctx, "sudo", "ip", "link", "del", r.ID)
	output, StdError := cmd.CombinedOutput()
	if StdError != nil {
		return fmt.Errorf("command faile %s : %s ", StdError, string(output))
	}
	fmt.Printf("[Linux Driver] Delete Interface %s \n", r.ID)
	return nil
}
