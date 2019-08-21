package host

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type Host struct {
	Hostname   string
	Username   string
	Password   string
	Ipaddr     string
	Port       string
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
	SFTPClient *sftp.Client
}

func (h *Host) connect() {
	config := &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", h.Ipaddr, h.Port), config)
	if err != nil {
		panic(fmt.Sprintf("[%s]: Failed to dial: %s", h.Hostname, err))
	}
	h.SSHClient = client
}

func (h *Host) openSession() {
	session, err := h.SSHClient.NewSession()
	if err != nil {
		panic(fmt.Sprintf("[%s]: Failed to create session: %s", h.Hostname, err))
	}
	h.SSHSession = session
}

func (h *Host) Cmd(cmd string) {
	h.connect()
	h.openSession()
	var stderr bytes.Buffer
	h.SSHSession.Stderr = &stderr
	if err := h.SSHSession.Run(cmd); err != nil {
		log.Println(stderr.String())
		panic(fmt.Sprintf("[%s]: Failed to run: %s", h.Hostname, err.Error()))
	}
}

func (h *Host) CmdGet(cmd string) string {
	h.connect()
	h.openSession()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	h.SSHSession.Stderr = &stderr
	h.SSHSession.Stdout = &stdout
	if err := h.SSHSession.Run(cmd); err != nil {
		log.Println(stderr.String())
		panic(fmt.Sprintf("[%s]: Failed to run: %s", h.Hostname, err.Error()))
	}
	return strings.Trim(stdout.String(), "\n")
}

func (h *Host) OpenSftp() {
	h.connect()
	sftpClient, err := sftp.NewClient(h.SSHClient)
	if err != nil {
		panic(fmt.Sprintf("[%s]: Failed to open sftp: %s", h.Hostname, err))
	}
	h.SFTPClient = sftpClient
}

func (h *Host) Put(local, remote string) int64 {
	localFile, err := os.Open(local)
	if err != nil {
		panic(fmt.Sprintf("Open local file %s: %s", local, err))
	}
	filename := path.Base(local)
	remotePath := path.Join(remote, filename)
	remoteFile, err := h.SFTPClient.Create(path.Join(remotePath))
	if err != nil {
		panic(fmt.Sprintf("[%s]: Create remote file %s: %s", h.Hostname, remote, err))
	}
	size, err := io.Copy(remoteFile, localFile)
	if err != nil {
		panic(fmt.Sprintf("[%s]: Upload file to %s: %s", h.Hostname, remote, err))
	}
	return size
}

func (h *Host) Get(local, remote string) int64 {
	filename := path.Base(remote)
	localFath := path.Join(local, filename)
	localFile, err := os.Create(localFath)
	if err != nil {
		panic(fmt.Sprintf("Create local file %s: %s", local, err))
	}
	remoteFile, err := h.SFTPClient.Open(remote)
	if err != nil {
		panic(fmt.Sprintf("Open remote file %s: %s", remote, err))
	}
	size, err := io.Copy(localFile, remoteFile)
	if err != nil {
		panic(fmt.Sprintf("[%s]: Download file to %s: %s", h.Hostname, remote, err))
	}
	return size
}
