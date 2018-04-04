package utils

import (
    "time"
)

const REVOKE_OK            = 0
const REVODE_NOK           = -1

const DEPLOYED             = "1"
const NOT_DEPLOYED         = "0"

const SSH_PORT             = "22"

const SERV_TYPE_MQ         = "MQ"
const SERV_TYPE_CACHE      = "CACHE"
const SERV_TYPE_DB         = "DB"

const MQ_SERV_CONTAINER    = "MQ_SERV_CONTAINER"
const CACHE_SERV_CONTAINER = "CACHE_SERV_CONTAINER"
const DB_SERV_CONTAINER    = "DB_SERV_CONTAINER"

const DB_TIDB_CONTAINER    = "DB_TIDB_CONTAINER"
const DB_TIKV_CONTAINER    = "DB_TIKV_CONTAINER"
const DB_PD_CONTAINER      = "DB_PD_CONTAINER"

const DB_TIDB              = "DB_TIDB"
const DB_TIKV              = "DB_TIKV"
const DB_PD                = "DB_PD"

const COLLECT_INTERVAL     = time.Duration(10)*time.Second

type WorkerPool struct {
	
}

type Instance struct {
    ID          string
    NAME        string
    OS_USER     string
    OS_PWD      string
    IP          string
    PORT        string
    STAT_PORT   string
}

type ResultInfo struct {
    RET_CODE    int
    RET_INFO    map[string]interface{}
}

type CpuInfo struct {
    Used      float32
}

type MemInfo struct {
    Used      int    // UNIT M
}

type DiskInfo struct {
    Total     int    // UNIT M
    Used      int    // UNIT M
    Available int    // UNIT M
}