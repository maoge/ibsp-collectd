package utils

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

func sshCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func Connect(user, password, ip, port string) (*ssh.Client, error) {
	ip_port := fmt.Sprintf("%s:%s", ip, port)
	auth := []ssh.AuthMethod{ssh.Password(password)}
	conf := ssh.ClientConfig{User: user, Auth: auth, HostKeyCallback: sshCallback}
	return ssh.Dial("tcp", ip_port, &conf)
}

func ExecCmd(cmd string, client *ssh.Client) ([]byte, error) {
	session, err := client.NewSession()
	if err == nil {
		defer session.Close()
		return session.Output(cmd)
	} else {
		return nil, err
	}
}

func GetPid(UniFlag string, client *ssh.Client) (string, error) {
    session, err := client.NewSession()
    if err == nil {
        defer session.Close()
        
        cmd := fmt.Sprintf("ps -eo 'pid,comm,args' | grep \"%s \" | grep -v grep", UniFlag)
        result, err1 := session.CombinedOutput(cmd)
        if err1 == nil {
            start, end := find(result, 32)
            if start != -1 && end != -1 {
                pid := string(result[start:end])
                return pid, nil
            } else {
                return "", nil
            }
        } else {
            return "", err1
        }
    } else {
        return "", err
    }
}

func find(src []byte, b byte) (int, int) {
    l := len(src)
    start := -1
    end := -1
    for i:=0; i < l; i++ {
        if (start == -1) {
            if src[i] != 32 {
                start = i;
            }
        } else {
            if src[i] == 32 {
                end = i
                break
            }
        }
    }

    return start, end
}
