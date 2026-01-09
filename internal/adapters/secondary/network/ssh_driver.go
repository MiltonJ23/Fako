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
}

func NewSSHDriver(host, user, privateKeyPath string, passphrase *domain.Secret) (*SSHDriver, error) {
	return &SSHDriver{
		Host:       host,
		User:       user,
		Port:       "22",
		PrivateKey: privateKeyPath,
		Passphrase: passphrase,
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
		s.Passphrase.Wipe()
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
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: implement a secure way out of this pile of dark shiet
		Timeout:         5 * time.Second,
	}

	// now we can return the dialed conversation
	return ssh.Dial("tcp", net.JoinHostPort(s.Host, s.Port), config)
}

func (s *SSHDriver) ApplyResource(ctx context.Context, r *domain.Resource) error {
	// TODO  implement a graceful shutdown here
	sshClient, clientConnectionError := s.connect()
	if clientConnectionError != nil {
		return fmt.Errorf("unable to establish a ssh connection, an error occured %v\n", clientConnectionError)
	}
	defer sshClient.Close() // the resource will be deleted at the end of the method execution

	// now let's open a new session
	sshSession, sshSessionCreationError := sshClient.NewSession()
	if sshSessionCreationError != nil {
		return fmt.Errorf("failed to open a new session, %v\n", sshSessionCreationError)
	}
	defer sshSession.Close()
	// reaching here means we were able to establish a connection and open a terminal session

	fmt.Printf("[SSH DRIVER] Connected to %s through a secure channel\n", s.Host)

	//We are then going to simulate, to test if we can properly execute a command, we are going to create a file to attest that we were there
	var b bytes.Buffer
	sshSession.Stdout = &b

	cmdToExecute := fmt.Sprintf("touch /tmp/fako_was_here_%s.txt", r.ID)
	runningCommandError := sshSession.Run(cmdToExecute)
	if runningCommandError != nil {
		return fmt.Errorf("failed to execute the command remotely: %v", runningCommandError)
	}
	// if the error is not triggered, we can safely admit the command was executed properly
	fmt.Printf("[SSH DRIVER] remote command executed successfully : %v\n", cmdToExecute)
	return nil
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
