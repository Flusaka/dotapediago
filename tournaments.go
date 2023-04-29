package dotapediago

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

type TournamentStatus int8

const (
	tournamentsEndpoint = "https://liquipedia.net/dota2/api.php?action=parse&origin=*&format=json&page=Portal:Tournaments"

	TournamentStatusUpcoming TournamentStatus = 0
	TournamentStatusOngoing  TournamentStatus = 1
	TournamentStatusComplete TournamentStatus = 2
)

type Tournament struct {
	Name         string
	Status       TournamentStatus
	Tier         int
	StartDate    time.Time
	EndDate      time.Time
	PrizePool    string
	Participants int
	Location     string
}

func (client *Client) GetTournaments() ([]*Tournament, error) {
	doc, err := client.getDocument(tournamentsEndpoint)

	if err != nil {
		return nil, err
	}

	var tournaments []*Tournament
	tables := doc.Find(".divTable")
	if tables.Length() > 0 {
		// This code naively assumes we'll have 3 tables, might not be the best way :)

		// First table is upcoming tournaments
		upcomingTable := tables.Eq(0)
		upcomingTournaments := parseTournaments(upcomingTable, TournamentStatusUpcoming)

		tournaments = append(tournaments, upcomingTournaments...)

		// Second table is ongoing tournaments
		ongoingTable := tables.Eq(1)
		ongoingTournaments := parseTournaments(ongoingTable, TournamentStatusOngoing)

		tournaments = append(tournaments, ongoingTournaments...)

		// Third table is complete tournaments
		completeTable := tables.Eq(2)
		completeTournaments := parseTournaments(completeTable, TournamentStatusComplete)

		tournaments = append(tournaments, completeTournaments...)
	}

	return tournaments, nil
}

func parseTournaments(table *goquery.Selection, status TournamentStatus) []*Tournament {
	var tournaments []*Tournament
	table.Find(".divRow").Each(func(i int, selection *goquery.Selection) {
		tournament := parseTournament(selection, status)
		tournaments = append(tournaments, tournament)
	})
	return tournaments
}

func parseTournament(row *goquery.Selection, status TournamentStatus) *Tournament {
	// Get, and convert, the tier of the tournament as an integer
	tierString := row.Find(".Tier").Find("a").Text()
	tierStringParts := strings.Split(tierString, " ")
	tier, _ := strconv.Atoi(tierStringParts[1])

	// Get the tournament name
	tournamentName := row.Find(".Tournament").Find("b").Find("a").Text()

	// Get the start and end dates of the tournament

	// Get the prize pool
	prizePool := row.Find(".Prize").Text()

	// Get the number of participants
	participantsString := row.Find(".PlayerNumber").Text()

	// First we have to break any instances of nbsp, we can't rely on it being a "traditional" space
	participantsString = strings.ReplaceAll(participantsString, "\u00a0", " ")
	participantsStringParts := strings.Split(participantsString, " ")
	participants, _ := strconv.Atoi(participantsStringParts[0])

	// Get the location (could be a region, country or city)
	location := row.Find(".Location").Text()

	return &Tournament{
		tournamentName,
		status,
		tier,
		time.Now(),
		time.Now(),
		prizePool,
		participants,
		location,
	}
}
