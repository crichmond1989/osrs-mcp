package tools

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/crich/osrs-mcp/internal/hiscores"
	"github.com/crich/osrs-mcp/internal/prices"
	"github.com/crich/osrs-mcp/internal/wiki"
	"github.com/crich/osrs-mcp/internal/wikisync"
)

// RegisterAll wires all MCP tools onto the server.
// Accepts interfaces so tests can substitute fakes.
func RegisterAll(s *server.MCPServer, wikiClient wiki.Client, pricesClient prices.Client, hiscoresClient hiscores.Client, wikiSyncClient wikisync.Client) {
	ws := NewWikiSearchTool(wikiClient)
	s.AddTool(ws.Definition(), ws.Handler)

	wsf := NewWikiSearchFullTool(wikiClient)
	s.AddTool(wsf.Definition(), wsf.Handler)

	wp := NewWikiPageTool(wikiClient)
	s.AddTool(wp.Definition(), wp.Handler)

	pl := NewPricesLatestTool(pricesClient)
	s.AddTool(pl.Definition(), pl.Handler)

	pm := NewPricesMappingTool(pricesClient)
	s.AddTool(pm.Definition(), pm.Handler)

	pt := NewPricesTimeSeriesTool(pricesClient)
	s.AddTool(pt.Definition(), pt.Handler)

	ps := NewPlayerStatsTool(hiscoresClient)
	s.AddTool(ps.Definition(), ps.Handler)

	pq := NewPlayerQuestsTool(wikiSyncClient)
	s.AddTool(pq.Definition(), pq.Handler)

	qi := NewQuestInfoTool(wikiClient)
	s.AddTool(qi.Definition(), qi.Handler)
}
