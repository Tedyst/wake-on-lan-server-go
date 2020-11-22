package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-ping/ping"
	"github.com/linde12/gowol"

	"github.com/valyala/fasthttp"
)

var (
	addr             = flag.String("addr", ":8080", "TCP address to listen to")
	externalAddress  = flag.String("externaladdress", "https://wakeup.stoicatedy.ovh", "External address of this server")
	checkingInterval = flag.Int("timeout", 60, "Interval in which to not check a host")
	dir              = flag.String("dir", "frontend/build/", "Directory to serve static files from")
	wolPort          = flag.String("wolport", "9", "Port to send WoL packet")
	broadcastAddress = flag.String("broadcast", "192.168.10.255", "IP for the broadcast")
	debug            = flag.Bool("debug", false, "Debug Mode")
)

func main() {
	flag.Parse()

	fs := &fasthttp.FS{
		Root:       *dir,
		IndexNames: []string{"index.html"},
	}
	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		if *debug {
			log.Printf("IP: %s, RequestURI:%s", ctx.RemoteIP(), ctx.RequestURI())
		}
		switch string(ctx.Path()) {
		case "/verify":
			verifyResponse(ctx)
		case "/ping":
			pingResponse(ctx)
		case "/wake":
			wakeUp(ctx)
		default:
			fsHandler(ctx)
		}
	}

	log.Print("Starting server")
	if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func noFunction(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "React app should be here\n")
}

func verifyResponse(ctx *fasthttp.RequestCtx) {
	host := string(ctx.QueryArgs().Peek("ip"))
	ip := net.ParseIP(host)
	if ip == nil {
		wrongResponseString(ctx, "Invalid IP")
		return
	}
	address := string(ctx.QueryArgs().Peek("address"))
	redirectURL := string(ctx.QueryArgs().Peek("redirectURL"))

	if isRecentlyPinged(host) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
	pinger := ping.New(host)
	pinger.Count = 1
	pinger.Timeout = time.Millisecond * 1000
	err := pinger.Run()
	if err == nil {
		if pinger.Statistics().PacketsRecv == 1 {
			addHost(host, pinger.Statistics().AvgRtt)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}
	}

	redirectstring := fmt.Sprintf("%s/?ip=%s&address=%s&redirectURL=%s", *externalAddress, host, address, redirectURL)
	log.Print(redirectstring)
	ctx.Redirect(redirectstring, fasthttp.StatusTemporaryRedirect)
}

func wrongResponseError(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	response := &jsonResponse{
		Success: false,
		Error:   fmt.Sprintf("%s", err),
	}
	str, _ := json.Marshal(response)
	fmt.Fprintf(ctx, "%s", str)
	return
}

func wrongResponseString(ctx *fasthttp.RequestCtx, err string) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	response := &jsonResponse{
		Success: false,
		Error:   err,
	}
	str, _ := json.Marshal(response)
	fmt.Fprintf(ctx, "%s", str)
	return
}

func pingResponse(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	host := string(ctx.QueryArgs().Peek("ip"))
	ip := net.ParseIP(host)

	if ip == nil {
		wrongResponseString(ctx, "Invalid IP")
		return
	}

	if isRecentlyPinged(host) {
		response := &jsonResponse{
			Success: true,
		}
		str, _ := json.Marshal(response)
		fmt.Fprintf(ctx, "%s", str)
		return
	}
	pinger := ping.New(host)
	pinger.Count = 1
	pinger.Timeout = time.Millisecond * 200
	err := pinger.Run()
	if err != nil {
		wrongResponseError(ctx, err)
		return
	}
	if pinger.Statistics().PacketsRecv == 0 {
		wrongResponseString(ctx, "System not started yet")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
	response := &jsonResponse{
		Success: true,
	}
	str, _ := json.Marshal(response)
	fmt.Fprintf(ctx, "%s", str)
}

func wakeUp(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	address := string(ctx.QueryArgs().Peek("address"))

	packet, err := gowol.NewMagicPacket(address)
	if err != nil {
		wrongResponseError(ctx, err)
		return
	}
	err = packet.SendPort(*broadcastAddress, *wolPort)
	if err != nil {
		wrongResponseError(ctx, err)
		return
	}

	response := &jsonResponse{
		Success: true,
		Error:   "",
	}
	str, _ := json.Marshal(response)
	fmt.Fprintf(ctx, "%s", str)
}
