package host

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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

func (h *Host) connect() error {
	config := &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", h.Ipaddr, h.Port), config); err != nil {
		return errors.New(fmt.Sprintf("connect to %s:%s failed, %s", h.Ipaddr, h.Port, err.Error()))
	} else {
		h.SSHClient = client
	}
	return nil
}

func (h *Host) openSession() error {
	if err := h.connect(); err != nil {
		return err
	}
	if session, err := h.SSHClient.NewSession(); err != nil {
		return errors.New(fmt.Sprintf("open session failed, %s", err.Error()))
	} else {
		h.SSHSession = session
	}
	return nil
}

func (h *Host) Cmd(cmd string) error {
	if err := h.openSession(); err != nil {
		return err
	}
	if err := h.SSHSession.Run(cmd); err != nil {
		return errors.New(fmt.Sprintf("run %s failed, %s", cmd, err.Error()))
	}
	return nil
}

func (h *Host) CmdGet(cmd string) (string, error) {
	if err := h.openSession(); err != nil {
		return "", err
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	h.SSHSession.Stderr = &stderr
	h.SSHSession.Stdout = &stdout
	if err := h.SSHSession.Run(cmd); err != nil {
		return strings.Trim(stderr.String(), "\n"), errors.New(fmt.Sprintf("run %s failed, %s", cmd, err.Error()))
	} else {
		return strings.Trim(stdout.String(), "\n"), nil
	}
}

func (h *Host) OpenSftp() error {
	if err := h.connect(); err != nil {
		return err
	}
	if sftpClient, err := sftp.NewClient(h.SSHClient); err != nil {
		return errors.New(fmt.Sprintf("open sftp failed, %s", err.Error()))
	} else {
		h.SFTPClient = sftpClient
	}
	return nil
}

func (h *Host) Put(local, remote string) (int64, error) {
	localFile, err := os.Open(local)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Open local file %s failed: %s", local, err.Error()))
	}
	filename := path.Base(local)
	remotePath := path.Join(remote, filename)
	remoteFile, err := h.SFTPClient.Create(remotePath)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("[%s]: Create remote file %s failed: %s", h.Hostname, remotePath, err.Error()))
	}
	size, err := io.Copy(remoteFile, localFile)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("[%s]: Upload file to %s failed: %s", h.Hostname, remotePath, err.Error()))
	}
	return size, nil
}

func (h *Host) Get(local, remote string) (int64, error) {
	filename := path.Base(remote)
	localPath := path.Join(local, filename)
	localFile, err := os.Create(localPath)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Create local file %s failed: %s", localPath, err.Error()))
	}
	remoteFile, err := h.SFTPClient.Open(remote)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Open remote file %s failed: %s", remote, err.Error()))
	}
	size, err := io.Copy(localFile, remoteFile)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("[%s]: Download file to %s failed: %s", h.Hostname, localPath, err.Error()))
	}
	return size, nil
}
