package network

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MiltonJ23/Fako/internal/core/domain"
	"golang.org/x/crypto/ssh"
)

type SSHDriver struct {
	Host       string
	User       string
	Port       string
	PrivateKey string
	Passphrase *domain.Secret
	Mapper     domain.CommandMapper
	DryRun     bool
}

func NewSSHDriver(host, user, privateKeyPath string, passphrase *domain.Secret, mapper domain.CommandMapper, dryRun bool) (*SSHDriver, error) {
	return &SSHDriver{
		Host:       host,
		User:       user,
		Port:       "22",
		PrivateKey: privateKeyPath,
		Passphrase: passphrase,
		Mapper:     mapper,
		DryRun:     dryRun,
	}, nil
}

// Let's write the connect method responsible for establishing a secure tunnel
func (s *SSHDriver) connect() (*ssh.Client, error) {
	// First thing first, we are loading the key
	key, loadingKeyError := os.ReadFile(s.PrivateKey)
	if loadingKeyError != nil {
		return nil, fmt.Errorf("unable to read the private key, an error happened %v\n", loadingKeyError)
	}
	//We are going to create the signer
	var signer ssh.Signer
	var signerCreationError error

	if s.Passphrase != nil {
		// First, reveal the secret bytes
		secretBytes := s.Passphrase.Reveal()
		// Being done here, let's use it to create our signer
		signer, signerCreationError = ssh.ParsePrivateKeyWithPassphrase(key, secretBytes)

		// now let's wipe the secret
		// s.Passphrase.Wipe()
	} else {
		signer, signerCreationError = ssh.ParsePrivateKey(key)
	}

	if signerCreationError != nil {
		return nil, fmt.Errorf("unable to parse private key, an error happened %v\n", signerCreationError)
	}

	//Now let's go to the configuration of the SSH client

	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: implement a secure way out of this pile of dark shit
		Timeout:         5 * time.Second,
	}

	// now we can return the dialed conversation
	return ssh.Dial("tcp", net.JoinHostPort(s.Host, s.Port), config)
}

func (s *SSHDriver) ApplyResource(ctx context.Context, r *domain.Resource) error {
	// TODO  implement a graceful shutdown here
	cmds, commandGenerationErrors := s.Mapper.GenerateApplyCommands(r)
	if commandGenerationErrors != nil {
		return fmt.Errorf("error while generating commands: %v", commandGenerationErrors)
	}

	return s.RunCommands(ctx, cmds)
}

// now let's go with DeleteResource, but as with ApplyResource, we are going to delete the file we created

func (s *SSHDriver) DeleteResource(ctx context.Context, r *domain.Resource) error {
	sshClient, clientConnectionError := s.connect()
	if clientConnectionError != nil {
		return fmt.Errorf("unable to establish a ssh connection, an error occured %v\n", clientConnectionError)
	}

	defer sshClient.Close()
	sshSession, sshSessionCreationError := sshClient.NewSession()
	if sshSessionCreationError != nil {
		return fmt.Errorf("failed to open a new session, %v\n", sshSessionCreationError)
	}
	defer sshSession.Close()

	var buf bytes.Buffer
	sshSession.Stdout = &buf
	cmdToExecute := fmt.Sprintf("rm /tmp/fako_was_here_%s.txt", r.ID)
	return sshSession.Run(cmdToExecute)
}

func (s *SSHDriver) RunCommands(ctx context.Context, commands []domain.RemoteCommand) error {
	// Ensure we are not going with a dry run before running the code
	if s.DryRun {
		fmt.Println("\n--- DRY RUN MODE (No changes applied) ---")
		for i, cmd := range commands {
			fmt.Printf("[%d] Desc: %s\n", i+1, cmd.Description)
			fmt.Printf("    Cmd : %s\n", cmd.Cmd)
		}
		fmt.Println("-----------------------------------------")
		return nil
	}

	// let's configure the connection
	sshClient, clientConnectionError := s.connect()
	if clientConnectionError != nil {
		return fmt.Errorf("unable to establish a ssh connection, an error occured %v\n", clientConnectionError)
	}
	defer sshClient.Close()

	sshSession, sshSessionCreationError := sshClient.NewSession()
	if sshSessionCreationError != nil {
		return fmt.Errorf("failed to open a new session, %v\n", sshSessionCreationError)
	}
	defer sshSession.Close()

	// We are going to fuse all  the commands into an executable file to avoid overload on the ssh connection
	var executableScript bytes.Buffer
	executableScript.WriteString("set -e\n") // Ensure the transaction is stopped if an error occur

	for _, command := range commands {
		executableScript.WriteString(fmt.Sprintf("echo Executing: %s\n", command.Description))
		executableScript.WriteString(command.Cmd + "\n")
	}

	// Then we capture the standard output and err
	var stdout, stderr bytes.Buffer
	sshSession.Stdout = &stdout
	sshSession.Stderr = &stderr

	fmt.Printf("-> [SSH TRANSACTION] Sending %d commands to %s .....\n", len(commands), s.Host)

	ScriptExecutionError := sshSession.Run(executableScript.String())
	if ScriptExecutionError != nil {
		return fmt.Errorf("transaction failed on host %s:\nERROR: %v\nSTDERR: %s\nLAST OUTPUT: %s", s.Host, ScriptExecutionError, stderr.String(), stdout.String())
	}

	fmt.Printf("-> [SSH TRANSACTION] Success on Host %s \n ", s.Host)
	return nil
}

// ExecuteCommand will launch a command and fetch its output (stdout)
func (s *SSHDriver) ExecuteCommand(ctx context.Context, command *domain.RemoteCommand) (string, error) {
	// let's make a dry-run read in-case
	if s.DryRun {
		fmt.Printf("[DRY-RUN READ] Executing: %s\n\n", command.Description)
		return "[ ]", nil
	}

	// now , manage the connexion and session launch
	// let's configure the connexion
	sshClient, clientConnectionError := s.connect()
	if clientConnectionError != nil {
		return "[ ]", fmt.Errorf("unable to establish a ssh connection, an error occured %v\n", clientConnectionError)
	}
	defer sshClient.Close()

	sshSession, sshSessionCreationError := sshClient.NewSession()
	if sshSessionCreationError != nil {
		return "[ ]", fmt.Errorf("failed to open a new session, %v\n", sshSessionCreationError)
	}
	defer sshSession.Close()

	// now let's capture the output
	var stdout bytes.Buffer
	sshSession.Stdout = &stdout

	runningCommandError := sshSession.Run(command.Cmd)
	if runningCommandError != nil {
		return "[ ]", fmt.Errorf("failed to run the command: %v\n", runningCommandError.Error())
	}

	return stdout.String(), nil
}
