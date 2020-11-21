package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	"github.com/linde12/gowol"

	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":8080", "TCP address to listen to")
)

func main() {
	flag.Parse()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/verify":
			verifyResponse(ctx)
		case "/ping":
			pingResponse(ctx)
		case "/wake":
			wakeUp(ctx)
		default:
			noFunction(ctx)
		}
	}

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
	port, err := strconv.Atoi(string(ctx.QueryArgs().Peek("port")))
	if err != nil {
		wrongResponseString(ctx, "Port is not defined")
		return
	}
	if port < 0 || port > 65535 {
		wrongResponseString(ctx, "Invalid port")
		return
	}
	redirectstring := fmt.Sprintf("/?host=%s&ip=%s&address=%s&port=%s", host, ip, address, fmt.Sprint(port))
	ctx.Redirect(redirectstring, fasthttp.StatusTemporaryRedirect)
}

func wrongResponseError(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusUnauthorized)
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
	ctx.SetStatusCode(fasthttp.StatusUnauthorized)
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

	host := string(ctx.QueryArgs().Peek("ip"))
	ip := net.ParseIP(host)
	if ip == nil {
		wrongResponseString(ctx, "Invalid IP")
		return
	}
	address := string(ctx.QueryArgs().Peek("address"))
	port, err := strconv.Atoi(string(ctx.QueryArgs().Peek("port")))
	if err != nil {
		wrongResponseString(ctx, "Port is not defined")
		return
	}
	if port < 0 || port > 65535 {
		wrongResponseString(ctx, "Invalid port")
		return
	}

	packet, err := gowol.NewMagicPacket(address)
	if err != nil {
		wrongResponseError(ctx, err)
		return
	}
	err = packet.SendPort(host, fmt.Sprint(port))
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
