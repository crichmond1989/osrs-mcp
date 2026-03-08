package wikisync

// Quest state constants returned by the WikiSync API.
const (
	QuestNotStarted = 0
	QuestInProgress = 1
	QuestComplete   = 2
)

// WikiSyncResponse holds the data returned by the WikiSync player endpoint.
// Quests maps OSRS internal quest IDs (as strings) to their completion state.
type WikiSyncResponse struct {
	Quests map[string]int `json:"quests"`
}
