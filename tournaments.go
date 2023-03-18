package dotapediago

const (
	tournamentsEndpoint = "https://liquipedia.net/dota2/Portal:Tournaments"
)

type Tournament struct {
}

func (client *Client) GetTournaments() []*Tournament {
	return []*Tournament{}
}
