// Copyright 2014 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)

package web

import (
	"encoding/json"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMatchesApi(t *testing.T) {
	web := setupTestWeb(t)

	match1 := model.Match{Type: "qualification", DisplayName: "1", Time: time.Unix(0, 0), Red1: 1, Red2: 2, Red3: 3,
		Blue1: 4, Blue2: 5, Blue3: 6, Blue1IsSurrogate: true, Blue2IsSurrogate: true, Blue3IsSurrogate: true}
	match2 := model.Match{Type: "qualification", DisplayName: "2", Time: time.Unix(600, 0), Red1: 7, Red2: 8, Red3: 9,
		Blue1: 10, Blue2: 11, Blue3: 12, Red1IsSurrogate: true, Red2IsSurrogate: true, Red3IsSurrogate: true}
	match3 := model.Match{Type: "practice", DisplayName: "1", Time: time.Now(), Red1: 6, Red2: 5, Red3: 4,
		Blue1: 3, Blue2: 2, Blue3: 1}
	web.arena.Database.CreateMatch(&match1)
	web.arena.Database.CreateMatch(&match2)
	web.arena.Database.CreateMatch(&match3)
	matchResult1 := model.BuildTestMatchResult(match1.Id, 1)
	web.arena.Database.CreateMatchResult(matchResult1)

	recorder := web.getHttpResponse("/api/matches/qualification")
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "application/json", recorder.HeaderMap["Content-Type"][0])
	var matchesData []MatchWithResult
	err := json.Unmarshal([]byte(recorder.Body.String()), &matchesData)
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(matchesData)) {
		assert.Equal(t, match1.Id, matchesData[0].Match.Id)
		assert.Equal(t, *matchResult1, matchesData[0].Result.MatchResult)
		assert.Equal(t, match2.Id, matchesData[1].Match.Id)
		assert.Nil(t, matchesData[1].Result)
	}
}

func TestRankingsApi(t *testing.T) {
	web := setupTestWeb(t)

	// Test that empty rankings produces an empty array.
	recorder := web.getHttpResponse("/api/rankings")
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "application/json", recorder.HeaderMap["Content-Type"][0])
	rankingsData := struct {
		Rankings           []RankingWithNickname
		TeamNicknames      map[string]string
		HighestPlayedMatch string
	}{}
	err := json.Unmarshal([]byte(recorder.Body.String()), &rankingsData)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(rankingsData.Rankings))
	assert.Equal(t, "", rankingsData.HighestPlayedMatch)

	ranking1 := RankingWithNickname{*game.TestRanking2(), "Simbots"}
	ranking2 := RankingWithNickname{*game.TestRanking1(), "ChezyPof"}
	web.arena.Database.CreateRanking(&ranking1.Ranking)
	web.arena.Database.CreateRanking(&ranking2.Ranking)
	web.arena.Database.CreateMatch(&model.Match{Type: "qualification", DisplayName: "29", Status: "complete"})
	web.arena.Database.CreateMatch(&model.Match{Type: "qualification", DisplayName: "30"})
	web.arena.Database.CreateTeam(&model.Team{Id: 254, Nickname: "ChezyPof"})
	web.arena.Database.CreateTeam(&model.Team{Id: 1114, Nickname: "Simbots"})

	recorder = web.getHttpResponse("/api/rankings")
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "application/json", recorder.HeaderMap["Content-Type"][0])
	err = json.Unmarshal([]byte(recorder.Body.String()), &rankingsData)
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(rankingsData.Rankings)) {
		assert.Equal(t, ranking1, rankingsData.Rankings[1])
		assert.Equal(t, ranking2, rankingsData.Rankings[0])
	}
	assert.Equal(t, "29", rankingsData.HighestPlayedMatch)
}

func TestSponsorSlides(t *testing.T) {
	web := setupTestWeb(t)

	slide1 := model.SponsorSlide{1, "subtitle", "line1", "line2", "image", 2}
	slide2 := model.SponsorSlide{2, "Chezy Sponsaur", "Teh", "Chezy Pofs", "ejface.jpg", 54}
	web.arena.Database.CreateSponsorSlide(&slide1)
	web.arena.Database.CreateSponsorSlide(&slide2)

	recorder := web.getHttpResponse("/api/sponsor_slides")
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "application/json", recorder.HeaderMap["Content-Type"][0])
	var sponsorSlides []model.SponsorSlide
	err := json.Unmarshal([]byte(recorder.Body.String()), &sponsorSlides)
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(sponsorSlides)) {
		assert.Equal(t, slide1, sponsorSlides[0])
		assert.Equal(t, slide2, sponsorSlides[1])
	}
}