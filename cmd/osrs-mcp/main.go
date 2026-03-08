package main

import (
	"flag"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/crichmond1989/osrs-mcp/internal/hiscores"
	"github.com/crichmond1989/osrs-mcp/internal/prices"
	"github.com/crichmond1989/osrs-mcp/internal/tools"
	"github.com/crichmond1989/osrs-mcp/internal/wiki"
	"github.com/crichmond1989/osrs-mcp/internal/wikisync"
)

func main() {
	addr := flag.String("addr", os.Getenv("ADDR"), "HTTP listen address (e.g. :8080); if empty, uses stdio")
	flag.Parse()

	s := server.NewMCPServer(
		"OSRS Wiki MCP",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	tools.RegisterAll(s, wiki.NewClient(), prices.NewClient(), hiscores.NewClient(), wikisync.NewClient())

	if *addr != "" {
		log.Printf("starting HTTP server on %s", *addr)
		if err := server.NewStreamableHTTPServer(s).Start(*addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}
}
