package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/valyala/fasthttp"
    "github.com/maoge/ibsp-collectd/probe"
    "github.com/maoge/ibsp-collectd/routing"
)

type Global struct {
    probe      *probe.Probe
}

var (
    name     = flag.String("name",    "",      "universal unique id")
    addr     = flag.String("addr",    ":8080", "TCP address to listen to")
    compress = flag.Bool("compress",  false,   "Whether to enable transparent response compression")
    rooturl  = flag.String("rooturl", "",      "Where to fetch meta data")
    servid   = flag.String("servid",  "",      "Specify the id of service to collect")

    global     Global
)

func main() {
    flag.Parse()

    bootstrap()

    waitExit()

    destroy()
}

func bootstrap() bool {
    go bootHttpService()

    probe, ok := initProbe()
    if ok {
        global.probe = probe
    } else {
        return false
    }

    return true
}

func initProbe() (*probe.Probe, bool) {
    var probe *probe.Probe = new(probe.Probe)
    ok := probe.Init(*rooturl, *servid)
    if ok {
        probe.Start()
        return probe, true 
    } else {
        return nil, false
    }
}

func bootHttpService() bool {
    //h := requestHandler
    //if *compress {
    //	h = fasthttp.CompressHandler(h)
    //}

    //if err := fasthttp.ListenAndServe(*addr, h); err != nil {
    //	log.Fatalf("Error in ListenAndServe: %s", err)
    //    return false
    //}

    router := routing.New()
    router.To("GET,POST", "/test", func(c *routing.Context) error {
        fmt.Fprintf(c, "Hello, FastHttp!")
        return nil
    })

    router.To("GET,POST", "/getCollectData", getCollectData)

    panic(fasthttp.ListenAndServe(*addr, router.HandleRequest))

    return true
}

func getCollectData(c *routing.Context) error {
    if global.probe != nil {
        fmt.Fprintf(c, global.probe.GetCollectData())
    } else {
        fmt.Fprintf(c, "{}")
    }

    return nil
}

//func initGorutinePool() bool, WorkerPool {
//}

func requestHandler(ctx *fasthttp.RequestCtx) {
	//fmt.Fprintf(ctx, "Hello, world!\n\n")

	//fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	//fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	//fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	//fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	//fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	//fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	//fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	//fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	//fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	//fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	//fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	//ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	//var c fasthttp.Cookie
	//c.SetKey("cookie-name")
	//c.SetValue("cookie-value")
	//ctx.Response.Header.SetCookie(&c)
}

func waitExit() {
    exitChan := make(chan int)
    signalChan := make(chan os.Signal, 1)
    go func() {
        <-signalChan
        exitChan <- 1
    }()
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    <-exitChan

    return
}

func destroy() {
    if global.probe != nil {
        global.probe.Stop()
    }
}
