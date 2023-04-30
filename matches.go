package dotapediago

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

type MatchStatus uint8

const (
	matchesEndpoint   = "https://liquipedia.net/dota2/api.php?action=parse&origin=*&format=json&page=Liquipedia:Upcoming_and_ongoing_matches"
	streamEndpointFmt = "https://liquipedia.net/dota2/Special:Stream/twitch/%s"

	MatchStatusUpcoming MatchStatus = 0
	MatchStatusOngoing  MatchStatus = 1
	MatchStatusComplete MatchStatus = 2
)

type Match struct {
	TeamOne        *BaseTeam
	TeamTwo        *BaseTeam
	BestOf         int
	Status         MatchStatus
	StartTime      time.Time
	Stream         string
	TournamentName string
}

func (client *Client) GetMatches() ([]*Match, error) {
	// Add a year for the default GetMatches call
	return client.GetMatchesUntil(time.Now().Add(time.Hour * 24 * 365))
}

func (client *Client) GetMatchesUntil(maxStartTime time.Time) ([]*Match, error) {
	doc, err := client.getDocument(matchesEndpoint)

	if err != nil {
		return nil, err
	}

	// This is a little awkward, and works with how the current HTML structure is, but might need reworking, need to investigate data queries
	matchesGroups := doc.Find(".matches-list").Children().Last().Children()
	// First one is _all_ matches of any "prestige"
	matches := client.parseMatches(matchesGroups.Eq(0), maxStartTime)
	return matches, nil
}

func (client *Client) parseMatches(table *goquery.Selection, maxStartTime time.Time) []*Match {
	var matches []*Match
	table.Find(".infobox_matches_content").Each(func(i int, selection *goquery.Selection) {
		if match, shouldAdd := client.parseMatch(selection, maxStartTime); shouldAdd {
			matches = append(matches, match)
		}
	})
	return matches
}

func (client *Client) parseMatch(row *goquery.Selection, maxStartTime time.Time) (*Match, bool) {
	// Parse the teams first
	status := MatchStatusUpcoming

	teamOne := parseTeam(row.Find(".team-left"))
	teamTwo := parseTeam(row.Find(".team-right"))

	timerObject := row.Find(".timer-object")

	var startTime time.Time
	startTimestampStr, exists := timerObject.Attr("data-timestamp")
	if exists {
		startTimestamp, _ := strconv.ParseInt(startTimestampStr, 10, 64)
		startTime = time.Unix(startTimestamp, 0)
		if startTime.After(maxStartTime) {
			return nil, false
		}
		if startTime.Before(time.Now()) {
			status = MatchStatusOngoing
		}
	} else {
		// If the timestamp doesn't exist for some reason, return false for whether it should be added to the matches array
		return nil, false
	}

	var streamUrl = ""
	twitchStream, exists := timerObject.Attr("data-stream-twitch")
	if exists {
		streamUrl, _ = client.GetStreamURL(fmt.Sprintf(streamEndpointFmt, twitchStream))
	}

	var bestOf = 0
	boContainer := row.Find(".versus abbr")
	boText := boContainer.Text()
	if boText != "" {
		boText = strings.TrimPrefix(boText, "Bo")
		bestOf, _ = strconv.Atoi(boText)
	}

	tournamentName := row.Find(".league-icon-small-image > a").AttrOr("title", "")

	return &Match{
		TeamOne:        teamOne,
		TeamTwo:        teamTwo,
		Stream:         streamUrl,
		StartTime:      startTime,
		Status:         status,
		BestOf:         bestOf,
		TournamentName: tournamentName,
	}, true
}

func parseTeam(teamSelector *goquery.Selection) *BaseTeam {
	teamContainer := teamSelector.Children().First()
	teamShortName := teamContainer.Find(".team-template-text > a").Text()
	teamFullName := teamContainer.AttrOr("data-highlightingclass", teamShortName)
	return &BaseTeam{
		ShortName: teamShortName,
		FullName:  teamFullName,
	}
}
