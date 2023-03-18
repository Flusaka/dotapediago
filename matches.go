package dotapediago

const (
	matchesEndpoint = "https://liquipedia.net/dota2/api.php?action=parse&origin=*&format=json&page=Liquipedia:Upcoming_and_ongoing_matches"
)

type Match struct {
}

func (client *Client) GetMatches() []*Match {
	return []*Match{}
}
