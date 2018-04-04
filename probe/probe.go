package probe

import (
	"fmt"
	"log"
    "time"
    "encoding/json"

    "github.com/maoge/ibsp-collectd/utils"
)

type Probe struct {
    RootUrl          string
    ServID           string

    ServType         string
    DBGenProbe       *DBServProbe
    MQGenProbe       *MQServProbe
    CacheGenProbe    *CacheServProbe

    DeployFlagMap    map[string]string

    Running          bool
}

type DBServProbe struct {
    TiDBSrvProbes    []*GeneralDataProbe
    PDSrvProbes      []*GeneralDataProbe
    TiKVSrvProbes    []*GeneralDataProbe
}

type MQServProbe struct {
    SwitchSrvProbes  []*GeneralDataProbe
    BrokerSrvProbes  []*GeneralDataProbe
}

type CacheServProbe struct {
    ProxySrvProbes   []*GeneralDataProbe
    CacheNodeProbes  []*GeneralDataProbe
}

func (pb *Probe) Init(url, servid string) bool {
    pb.RootUrl = url
    pb.ServID  = servid

    loadServTopoUrl := fmt.Sprintf("%s/%s/%s", url, "configsvr", "loadServiceTopoByInstID")

    params := fmt.Sprintf("%s=%s", "INST_ID", servid)
    result, httpErr := utils.HttpPost(loadServTopoUrl, params)
    if httpErr != nil {
        log.Fatalf("http post error: %s", httpErr)
        return false 
    }

    valid := json.Valid(result)
    if valid != true {
        log.Fatalf("result is not a valid JSON encoding: %s", result)
        return false
    }

    return pb.parse(result)
}

func (pb *Probe) parse(body []byte) bool {
    resultInfo := &utils.ResultInfo{}
    unmarshelErr := json.Unmarshal(body, &resultInfo)
    if unmarshelErr != nil {
        log.Fatalf("json unmarshel error: %s", unmarshelErr)
        return false
    }

    if resultInfo.RET_CODE != utils.REVOKE_OK {
        log.Fatalf("result nok: %s", string(body))
        return false
    }

    deployFlag := resultInfo.RET_INFO["DEPLOY_FLAG"]
    deployFlagArr, arrOk := deployFlag.([]interface{})
    if !arrOk {
        log.Fatalf("DEPLOY_FLAG reflex error!")
        return false
    }

    if !pb.parseDeployFlag(deployFlagArr) {
        return false
    }

    topo := resultInfo.RET_INFO[utils.DB_SERV_CONTAINER]
    var topoParse = false
    for k, _ := range resultInfo.RET_INFO {
        if k != "DEPLOY_FLAG" {            
            topoInfo, _ := topo.(map[string]interface{})

            switch k {
            case utils.DB_SERV_CONTAINER:
                pb.ServType = utils.SERV_TYPE_DB
                pb.DBGenProbe = new(DBServProbe)
                topoParse = pb.parseDB(topoInfo)
            case utils.MQ_SERV_CONTAINER:
                pb.ServType = utils.SERV_TYPE_MQ
                pb.MQGenProbe = new(MQServProbe)
                topoParse = pb.parseMQ(topoInfo)
            case utils.CACHE_SERV_CONTAINER:
                pb.ServType = utils.SERV_TYPE_CACHE
                pb.CacheGenProbe = new(CacheServProbe)
                topoParse = pb.parseCache(topoInfo)
            default:
                topoParse = false
            }
        }
    }

    return topoParse
}

func (pb *Probe) parseDeployFlag(arr []interface{}) bool {
    pb.DeployFlagMap = make(map[string]string)
    l := len(arr)
    for i := 0; i < l; i++ {
        m, mapOk := arr[i].(map[string]interface{})
        if !mapOk {
            log.Fatalf("DEPLOY_FLAG item reflex error!")
            return false
        }
        
        for k, v := range m {
            s, strOk := v.(string)
            if !strOk {
                log.Fatalf("DEPLOY_FLAG item value reflex error!")
                return false
            }
            pb.DeployFlagMap[k] = s
        }
    }

    return true
}

func (pb *Probe) parseDB(topo map[string]interface{}) bool {
    // TIDB
    tidbContainer := topo[utils.DB_TIDB_CONTAINER]
    tidbContainerMap, _ := tidbContainer.(map[string]interface{})

    tidb := tidbContainerMap[utils.DB_TIDB]
    tidbArr, _ := tidb.([]interface{})
    
    lenTiDB := len(tidbArr)
    pb.DBGenProbe.TiDBSrvProbes = make([]*GeneralDataProbe, lenTiDB)
    cnt := 0
    for i := 0; i < lenTiDB; i++ {
        m, _ := tidbArr[i].(map[string]interface{})

        id, _ := m["TIDB_ID"].(string)
        deployFlag := pb.DeployFlagMap[id]
        if deployFlag == "" || deployFlag == utils.NOT_DEPLOYED {
            continue
        }
        
        var genDataProbe *GeneralDataProbe = new(GeneralDataProbe)
        genDataProbe.ID             = id
        genDataProbe.Name,        _ = m["TIDB_NAME"].(string)
        genDataProbe.IP,          _ = m["IP"].(string)
        genDataProbe.Port,        _ = m["PORT"].(string)
        genDataProbe.StatPort,    _ = m["STAT_PORT"].(string)
        genDataProbe.User,        _ = m["OS_USER"].(string)
        genDataProbe.Passwd,      _ = m["OS_PWD"].(string)
        genDataProbe.UniFlag        = "\\-host " + genDataProbe.IP + " \\-P " + genDataProbe.Port  // TODO: change to TiDB ID later
        genDataProbe.InstallPath    = "tidb_deploy/" + genDataProbe.Port

        pb.DBGenProbe.TiDBSrvProbes[cnt] = genDataProbe
        cnt++
    }
    

    // PD
    pdContainer := topo[utils.DB_PD_CONTAINER]
    pdContainerMap, _ := pdContainer.(map[string]interface{})

    pd := pdContainerMap[utils.DB_PD]
    pdArr, _ := pd.([]interface{})

    lenPD := len(pdArr)
    pb.DBGenProbe.PDSrvProbes = make([]*GeneralDataProbe, lenPD)
    cnt = 0
    for i := 0; i < lenPD; i++ {
        m, _ := pdArr[i].(map[string]interface{})

        id, _ := m["PD_ID"].(string)
        deployFlag := pb.DeployFlagMap[id]
        if deployFlag == "" || deployFlag == utils.NOT_DEPLOYED {
            continue
        }

        var genDataProbe *GeneralDataProbe = new(GeneralDataProbe)
        genDataProbe.ID             = id
        genDataProbe.Name,        _ = m["PD_NAME"].(string)
        genDataProbe.IP,          _ = m["IP"].(string)
        genDataProbe.Port,        _ = m["PORT"].(string)
        genDataProbe.ClusterPort, _ = m["CLUSTER_PORT"].(string)
        genDataProbe.User,        _ = m["OS_USER"].(string)
        genDataProbe.Passwd,      _ = m["OS_PWD"].(string)
        genDataProbe.UniFlag        = "\\--name=" + genDataProbe.ID
        genDataProbe.InstallPath    = "pd_deploy/" + genDataProbe.Port
        
        pb.DBGenProbe.PDSrvProbes[cnt] = genDataProbe
        cnt++
    }

    // TIKV
    tikvContainer := topo[utils.DB_TIKV_CONTAINER]
    tikvContainerMap, _ := tikvContainer.(map[string]interface{})
    
    tikv := tikvContainerMap[utils.DB_TIKV]
    tikvArr, _ := tikv.([]interface{})

    lenTiKV := len(tikvArr)
    pb.DBGenProbe.TiKVSrvProbes = make([]*GeneralDataProbe, lenTiKV)
    cnt = 0
    for i := 0; i < lenTiKV; i++ {
        m, _ := tikvArr[i].(map[string]interface{})

        id, _ := m["TIKV_ID"].(string)
        deployFlag := pb.DeployFlagMap[id]
        if deployFlag == "" || deployFlag == utils.NOT_DEPLOYED {
            continue
        }

        var genDataProbe *GeneralDataProbe = new(GeneralDataProbe)
        genDataProbe.ID             = id
        genDataProbe.Name,        _ = m["TIKV_NAME"].(string)
        genDataProbe.IP,          _ = m["IP"].(string)
        genDataProbe.Port,        _ = m["PORT"].(string)
        genDataProbe.User,        _ = m["OS_USER"].(string)
        genDataProbe.Passwd,      _ = m["OS_PWD"].(string)
        genDataProbe.UniFlag        = "\\--addr " + genDataProbe.IP + ":" + genDataProbe.Port
        genDataProbe.InstallPath    = "tikv_deploy/" + genDataProbe.Port

        pb.DBGenProbe.TiKVSrvProbes[cnt] = genDataProbe
        cnt++
    }

    return true
}

func (pb *Probe) parseMQ(topo map[string]interface{}) bool {
    // TODO
    return false
}

func (pb *Probe) parseCache(topo map[string]interface{}) bool {
    // TODO
    return false
}

func (pb *Probe) doCollecting() {
    pb.Running = true

    for {
        if !pb.Running {
            break
        }

        switch pb.ServType {
        case utils.SERV_TYPE_DB:
            pb.collectDB()
        case utils.SERV_TYPE_MQ:
            pb.collectMQ()
        case utils.SERV_TYPE_CACHE:
            pb.collectCache()
        default:
            ;
        }
        
        time.Sleep(utils.COLLECT_INTERVAL)
    }
}

func (pb *Probe) collectDB() {
    dbprobe := pb.DBGenProbe
    if dbprobe == nil {
        return
    }

    tidbSrvArr := dbprobe.TiDBSrvProbes
    for _, m := range tidbSrvArr {
        if m == nil {
            continue
        }

        m.CollectGeneralData()
    }

    pdSrvArr := dbprobe.PDSrvProbes
    for _, m := range pdSrvArr {
        if m == nil {
            continue
        }

        m.CollectGeneralData()
    }

    tikvSrvArr := dbprobe.TiKVSrvProbes
    for _, m := range tikvSrvArr {
        if m == nil {
            continue
        }

        m.CollectGeneralData()
    }
}

func (pb *Probe) collectMQ() {
    // TODO
}

func (pb *Probe) collectCache() {
    // TODO
}

func (pb *Probe) GetCollectData() string {
    switch pb.ServType {
    case utils.SERV_TYPE_DB:
        return pb.getDBCollectData()
    case utils.SERV_TYPE_MQ:
        return pb.getMQCollectData()
    case utils.SERV_TYPE_CACHE:
        return pb.getCacheCollectData()
    default:
        return ""
    }
}

func (pb *Probe) getDBCollectData() string {
    var json string = ""
    json += "{"

    
    dbprobe := pb.DBGenProbe
    if dbprobe != nil {
        // tidb-server
        json += "\"DB_TIDB\":["
        tidbSrvArr := dbprobe.TiDBSrvProbes
        cnt := 0
        for _, m := range tidbSrvArr {
            if m == nil {
                continue
            }
            
            if cnt > 0 {
                json += ","
            }

            json += m.GetDataAsJson()

            cnt++
        }
        json += "]"


        // pd-server
        json += ",\"DB_PD\":["
        pdSrvArr := dbprobe.PDSrvProbes
        cnt = 0
        for _, m := range pdSrvArr {
            if m == nil {
                continue
            }

            if cnt > 0 {
                json += ","
            }

            json += m.GetDataAsJson()

            cnt++
        }
        json += "]"


        // tikv-server
        json += ",\"DB_TIKV\":["
        tikvSrvArr := dbprobe.TiKVSrvProbes
        cnt = 0
        for _, m := range tikvSrvArr {
            if m == nil {
                continue
            }

            if cnt > 0 {
                json += ","
            }

            json += m.GetDataAsJson()

            cnt++
        }
        json += "]"
    }

    json += "}"
    return json
}

func (pb *Probe) getMQCollectData() string {
    // TODO
    return ""
}

func (pb *Probe) getCacheCollectData() string {
    // TODO
    return ""
}

func (pb *Probe) Start() {
    go pb.doCollecting()
}

func (pb *Probe) Stop() {
    pb.Running = false
}
