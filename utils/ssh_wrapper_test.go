package utils

import (
	"testing"
    "fmt"
)

func TestConnect(t *testing.T) {
	client, connErr := Connect("mq", "amqp", "192.168.14.208", "22")
	if (connErr == nil) {
		defer client.Close()
		t.Log(client)
	} else {
		t.Log(connErr)
	}
}

func TestExecCmd(t *testing.T) {
	client, connErr := Connect("mq", "amqp", "192.168.14.208", "22")
	if connErr == nil {
		defer client.Close()

		out1, cmdErr1 := ExecCmd("echo $HOSTNAME", client)
		if cmdErr1 == nil {
			t.Log(out1)
		} else {
			t.Log(cmdErr1)
		}
	} else {
		t.Log(connErr)
	}
}

func TestGetPid(t *testing.T) {
    client, connErr := Connect("mq", "amqp", "192.168.14.206", "22")
    if connErr == nil {
        defer client.Close()

        out1, cmdErr1 := GetPid("addr 192.168.14.206:15011", client)
        if cmdErr1 == nil {
            fmt.Println(out1)
		} else {
			t.Log(cmdErr1)
		}
    } else {
		t.Log(connErr)
	}
}
