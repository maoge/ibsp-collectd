package probe

import (
    "fmt"
    "log"
    "strings"
    "strconv"

    "golang.org/x/crypto/ssh"
    "github.com/maoge/ibsp-collectd/utils"
)

type Broker struct {
    BrokerID     string
    BrokerName   string
    HostName     string
    ErlCookie    string
    IP           string
    Port         string
    MgrPort      string
    SyncPort     string
    OSUser       string
    OSPwd        string
    RootPwd      string

    GenProbe    *GeneralDataProbe
}

type VBroker struct {
    VBrokerID    string
    VBrokerName  string
    MasterID     string
    BrokerMap    map[string]*Broker
}

type CacheProxy struct {
    CacheProxyID   string
    CacheProxyName string
    OSUser         string
    OSPwd          string
    IP             string
    Port           string
    StatPort       string
    RWSep          string

    GenProbe    *GeneralDataProbe
}

type CacheNode struct {
    CacheNodeID    string
    CacheNodeName  string
    IP             string
    Port           string
    OSUser         string
    OSPwd          string

    GenProbe    *GeneralDataProbe
}

type CacheNodeCluster struct {
    CacheNodeClusterID   string
    CacheNodeClusterName string
    MasterID             string
    MaxMemory            string   // UNIT GB
    CacheSlot            string
                         
    CacheNodeMap  map[string]*CacheNode
}

type GeneralDataProbe struct {
    ID          string
    Name        string
    IP          string
    Port        string
    StatPort    string
    ClusterPort string

    User        string
    Passwd      string
    UniFlag     string
    InstallPath string

    Cpu         *utils.CpuInfo
    Mem         *utils.MemInfo
    Disk        *utils.DiskInfo
}

func (pb *GeneralDataProbe) Connect() *ssh.Client {
    conn, err := utils.Connect(pb.User, pb.Passwd, pb.IP, utils.SSH_PORT)
    if err != nil {
        fmt.Println(err)
        log.Fatal("ssh connect fatel error:", err, "\n")
        return nil
    }

    return conn;
}

func (pb *GeneralDataProbe) GeneralInfo() (*utils.CpuInfo, *utils.MemInfo, *utils.DiskInfo) {
    return pb.Cpu, pb.Mem, pb.Disk
}

func (pb *GeneralDataProbe) CollectGeneralData() {
    sshClient := pb.Connect()

    if sshClient == nil {
        return
    }

    pid := pb.GetPid(sshClient)

    pb.CPU(pid, sshClient)
    pb.MEM(pid, sshClient)
    pb.DISK(pb.InstallPath, sshClient)

    pb.Close(sshClient)
}

func (pb *GeneralDataProbe) GetPid(client *ssh.Client) string {
    pid, err := utils.GetPid(pb.UniFlag, client)
    if err != nil {
        log.Fatal(err)
        return ""
    } else {
        return pid
    }
}

func (pb *GeneralDataProbe) CPU(pid string, client *ssh.Client) {
    //cmd := fmt.Sprintf("ps -eo 'pid,pcpu' | grep %s | grep -v grep", pid)
    cmd := fmt.Sprintf("top -p %s -bcn 1 | grep %s", pid, pid)
    result, err := utils.ExecCmd(cmd, client)
    var cpu float32 = 0
    if err == nil {
        arr := strings.Fields(string(result))
        f, _ := strconv.ParseFloat(arr[8], 32)
        cpu = float32(f)
    }

    if pb.Cpu == nil {
        pb.Cpu = new(utils.CpuInfo)
    }
    pb.Cpu.Used = cpu
}

func (pb *GeneralDataProbe) MEM(pid string, client *ssh.Client) {
    cmd := fmt.Sprintf("ps -eo 'pid,rsz' | grep %s | grep -v grep", pid)
    result, err := utils.ExecCmd(cmd, client)
    mem := 0
    if err == nil {
        arr := strings.Fields(string(result))
        mem, _ = strconv.Atoi(arr[1])
    } else {
        log.Fatal(err)
    }

    if pb.Mem == nil {
        pb.Mem = new(utils.MemInfo)
    }
    pb.Mem.Used = mem/1024
}

func (pb *GeneralDataProbe) DISK(path string, client *ssh.Client) {
    used := 0
    total := 0
    available := 0
    
    // du -s -m pd_deploy/15001
    // 192     pd_deploy/15001   
    cmd1 := fmt.Sprintf("du -s -m %s", path)
    result1, err1 := utils.ExecCmd(cmd1, client)
    if err1 == nil {
        arr := strings.Fields(string(result1))

        used, _ = strconv.Atoi(arr[0])
    }

    // df -h -m pd_deploy/15001
    // Filesystem               1M-blocks  Used   Available Use%  Mounted on
    // /dev/mapper/centos-home  205144     46415  158730    23%   /home
    cmd2 := fmt.Sprintf("df -hm --output='size,avail' %s", path)
    result2, err2 := utils.ExecCmd(cmd2, client)
    if err2 == nil {        
        // find the new line
        l := len(result2)
        start := 0
        for i:=0; i < l; i++ {
            if result2[i] == '\n' {
                start = i + 1
                break
            }
        }
        
        s := string(result2[start:l])
        arr := strings.Fields(s)

        total, _ = strconv.Atoi(arr[0])
        available, _ = strconv.Atoi(arr[1])
    }

    if pb.Disk == nil {
        pb.Disk = new(utils.DiskInfo)
    }
    pb.Disk.Total = total
    pb.Disk.Used = used
    pb.Disk.Available = available
}

func (pb *GeneralDataProbe) Close(client *ssh.Client) {
    if client != nil {
        client.Close()
    }
}

func (pb *GeneralDataProbe) GetDataAsJson() string {
    return fmt.Sprintf("{\"ID\":\"%s\",\"CPU\":{\"Used\":%f},\"MEM\":{\"Used\":%d},\"DISK\":{\"Total\":%d,\"Used\":%d,\"Available\":%d}}",
        pb.ID, pb.Cpu.Used, pb.Mem.Used, pb.Disk.Total, pb.Disk.Used, pb.Disk.Available)
}
