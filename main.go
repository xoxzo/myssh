package myssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Get default location of a private key
func privateKeyPath() string {
	return os.Getenv("HOME") + "/.ssh/id_rsa"
}

// Get private key for ssh authentication
func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, _ := ioutil.ReadFile(keyPath)
	return ssh.ParsePrivateKey(buff)
}

func makeSshConfig(user string) (*ssh.ClientConfig, error) {
	/*key, err := parsePrivateKey(privateKeyPath())
	  if err != nil {
	    return nil, err
	  }*/

	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, err
	}
	agentClient := agent.NewClient(conn)

	config := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			//ssh.PublicKeys(key),
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return &config, nil
}

func Run(command string, hosts []string) {
	for _, host := range hosts {
		fmt.Printf("Running %s at %s\n", command, host)
		RunHost(command, host)
	}
}

func RunHost(command string, host string) {
	config, err := makeSshConfig("kamal")
	if err != nil {
		log.Fatalf("[ERR] %v", err)
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatalf("[ERR] %v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("[ERR] %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatalf("[ERR] Failed running %s %v", command, err)
	}
	fmt.Printf("%s: %s\n", host, b.String())
}
