package probe

import (
	"testing"
    "fmt"
)

func TestMem(t *testing.T) {
    var genProbe GeneralDataProbe = GeneralDataProbe{"192.168.14.206",
                                                     "22",
                                                     "mq2",
                                                     "amqp2",
                                                     "name=c05b9b3c-d6fb-f127-5aa3-776e6223b1a4",
                                                     "pd_deploy/15001"}

    client := genProbe.Connect()
    if client == nil {
        t.Log("connect fail!")
        return
    }

    pid := genProbe.GetPid(client)
    fmt.Println("pid:", pid)
    if pid == "" {
        t.Log("pid get null string!")
        return
    }

    mem := genProbe.Mem(pid, client)
    fmt.Println("Mem info:", mem)
}

func TestDisk(t *testing.T) {
    var genProbe GeneralDataProbe = GeneralDataProbe{"192.168.14.206",
                                                     "22",
                                                     "mq2",
                                                     "amqp2",
                                                     "name=c05b9b3c-d6fb-f127-5aa3-776e6223b1a4",
                                                     "pd_deploy/15001"}

    client := genProbe.Connect()
    if client == nil {
        t.Log("connect fail!")
        return
    }

    disk := genProbe.Disk(genProbe.InstallPath, client)
    fmt.Println("Disk info:", disk)
}

func TestCPU(t *testing.T) {
    var genProbe GeneralDataProbe = GeneralDataProbe{"192.168.14.206",
                                                     "22",
                                                     "mq2",
                                                     "amqp2",
                                                     "name=c05b9b3c-d6fb-f127-5aa3-776e6223b1a4",
                                                     "pd_deploy/15001"}

    client := genProbe.Connect()
    if client == nil {
        t.Log("connect fail!")
        return
    }

    pid := genProbe.GetPid(client)
    fmt.Println("pid:", pid)
    if pid == "" {
        t.Log("pid get null string!")
        return
    }

    cpu := genProbe.CPU(pid, client)
    fmt.Println("CPU info:", cpu)
}
