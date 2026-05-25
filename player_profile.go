package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const playerProfilePreviewLimit = 10

type playerProfileSectionError struct {
	Section string `json:"section"`
	Error   string `json:"error"`
}

type playerProfileOnlineState struct {
	PlayerID int64  `json:"player_id"`
	Name     string `json:"name"`
	Map      string `json:"map"`
	Status   string `json:"status"`
	LastSeen string `json:"last_seen,omitempty"`
}

type playerProfileLocation struct {
	Map    string `json:"map"`
	Source string `json:"source,omitempty"`
}

type playerProfileInventorySummary struct {
	TotalItems      int        `json:"total_items"`
	TotalStackSize  int64      `json:"total_stack_size"`
	UniqueTemplates int        `json:"unique_templates"`
	PreviewItems    []itemInfo `json:"preview_items"`
}

type playerProfileCharXP struct {
	XP    int64 `json:"xp"`
	Level int   `json:"level"`
}

type playerProfileJourneySummary struct {
	TotalNodes     int           `json:"total_nodes"`
	CompleteNodes  int           `json:"complete_nodes"`
	RevealedNodes  int           `json:"revealed_nodes"`
	PendingRewards int           `json:"pending_rewards"`
	PreviewNodes   []journeyNode `json:"preview_nodes"`
}

type playerProfileResponse struct {
	PlayerID         int64                         `json:"player_id"`
	Identity         *playerInfo                   `json:"identity,omitempty"`
	OnlineState      *playerProfileOnlineState     `json:"online_state,omitempty"`
	Location         *playerProfileLocation        `json:"location,omitempty"`
	InventorySummary playerProfileInventorySummary `json:"inventory_summary"`
	Vehicles         []vehicleRow                  `json:"vehicles"`
	Currencies       []currencyRow                 `json:"currencies"`
	Factions         []factionRep                  `json:"factions"`
	Specializations  []specTrack                   `json:"specializations"`
	CharacterXP      *playerProfileCharXP          `json:"character_xp,omitempty"`
	JourneySummary   playerProfileJourneySummary   `json:"journey_summary"`
	RecentEvents     []gameEvent                   `json:"recent_events"`
	DungeonHistory   []dungeonRecord               `json:"dungeon_history"`
	SectionErrors    []playerProfileSectionError   `json:"section_errors"`
}

func newPlayerProfileResponse(playerID int64) playerProfileResponse {
	return playerProfileResponse{
		PlayerID: playerID,
		InventorySummary: playerProfileInventorySummary{
			PreviewItems: []itemInfo{},
		},
		Vehicles:        []vehicleRow{},
		Currencies:      []currencyRow{},
		Factions:        []factionRep{},
		Specializations: []specTrack{},
		JourneySummary: playerProfileJourneySummary{
			PreviewNodes: []journeyNode{},
		},
		RecentEvents:   []gameEvent{},
		DungeonHistory: []dungeonRecord{},
		SectionErrors:  []playerProfileSectionError{},
	}
}

func (p *playerProfileResponse) addSectionError(section string, err error) {
	if err == nil {
		return
	}
	p.addSectionErrorMessage(section, safePlayerProfileError(err))
}

func (p *playerProfileResponse) addSectionErrorMessage(section, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "section unavailable"
	}
	p.SectionErrors = append(p.SectionErrors, playerProfileSectionError{Section: section, Error: message})
}

func safePlayerProfileError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "not connected") {
		return "not connected"
	}
	if strings.Contains(msg, "internal error") {
		return "internal error"
	}
	return "section unavailable"
}

func handleGetPlayerProfile(w http.ResponseWriter, r *http.Request) {
	playerID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		jsonErr(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}

	profile, found := buildPlayerProfile(playerID)
	if !found {
		jsonErr(w, fmt.Errorf("player %d not found", playerID), http.StatusNotFound)
		return
	}
	jsonOK(w, profile)
}

func buildPlayerProfile(playerID int64) (playerProfileResponse, bool) {
	profile := newPlayerProfileResponse(playerID)
	identityLookupComplete := false

	if msg, ok := cmdFetchPlayers().(msgPlayers); !ok {
		profile.addSectionErrorMessage("identity", "internal error")
	} else if msg.err != nil {
		profile.addSectionError("identity", msg.err)
	} else {
		identityLookupComplete = true
		for _, row := range msg.rows {
			if row.ID == playerID {
				identity := row
				profile.Identity = &identity
				break
			}
		}
	}

	if identityLookupComplete && profile.Identity == nil {
		return profile, false
	}

	addProfileOnlineState(&profile, playerID)
	addProfileCurrency(&profile, playerID)
	addProfileInventory(&profile, playerID)
	addProfileVehicles(&profile)
	addProfileFactions(&profile, playerID)
	addProfileSpecializations(&profile, playerID)
	addProfileCharacterXP(&profile, playerID)
	addProfileJourney(&profile)
	addProfileEvents(&profile, playerID)
	addProfileDungeons(&profile, playerID)

	return profile, true
}

func addProfileOnlineState(profile *playerProfileResponse, playerID int64) {
	if profile.Identity != nil {
		profile.OnlineState = &playerProfileOnlineState{
			PlayerID: profile.Identity.ControllerID,
			Name:     profile.Identity.Name,
			Map:      profile.Identity.Map,
			Status:   profile.Identity.OnlineStatus,
		}
		if profile.Identity.Map != "" {
			profile.Location = &playerProfileLocation{Map: profile.Identity.Map, Source: "players"}
		}
	}

	msg, ok := cmdFetchOnlineState().(msgOnlineState)
	if !ok {
		profile.addSectionErrorMessage("online_state", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("online_state", msg.err)
		return
	}

	if row, found := selectProfileOnlineState(msg.rows, playerID, profile.Identity); found {
		profile.OnlineState = &row
		if row.Map != "" {
			profile.Location = &playerProfileLocation{Map: row.Map, Source: "online_state"}
		}
	}
}

func selectProfileOnlineState(rows []onlineStateRow, playerID int64, identity *playerInfo) (playerProfileOnlineState, bool) {
	ids := profileControllerIDs(playerID, identity)
	for _, row := range rows {
		if ids[row.PlayerID] {
			return playerProfileOnlineState{
				PlayerID: row.PlayerID,
				Name:     row.Name,
				Map:      row.Map,
				Status:   row.Status,
				LastSeen: row.LastSeen,
			}, true
		}
	}
	return playerProfileOnlineState{}, false
}

func addProfileCurrency(profile *playerProfileResponse, playerID int64) {
	msg, ok := cmdFetchCurrency().(msgCurrency)
	if !ok {
		profile.addSectionErrorMessage("currencies", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("currencies", msg.err)
		return
	}
	ids := profileControllerIDs(playerID, profile.Identity)
	for _, row := range msg.rows {
		if ids[row.PlayerID] {
			profile.Currencies = append(profile.Currencies, row)
		}
	}
}

func addProfileInventory(profile *playerProfileResponse, playerID int64) {
	msg, ok := cmdFetchInventory(playerID)().(msgInventory)
	if !ok {
		profile.addSectionErrorMessage("inventory_summary", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("inventory_summary", msg.err)
		return
	}
	profile.InventorySummary = summarizeProfileInventory(msg.rows)
}

func summarizeProfileInventory(rows []itemInfo) playerProfileInventorySummary {
	summary := playerProfileInventorySummary{PreviewItems: []itemInfo{}}
	seenTemplates := map[string]bool{}
	for _, item := range rows {
		summary.TotalItems++
		summary.TotalStackSize += item.StackSize
		seenTemplates[strings.ToLower(item.TemplateID)] = true
		if len(summary.PreviewItems) < playerProfilePreviewLimit {
			summary.PreviewItems = append(summary.PreviewItems, item)
		}
	}
	summary.UniqueTemplates = len(seenTemplates)
	return summary
}

func addProfileVehicles(profile *playerProfileResponse) {
	if profile.Identity == nil || profile.Identity.AccountID == 0 {
		profile.addSectionErrorMessage("vehicles", "account_id unavailable")
		return
	}
	msg, ok := cmdGetPlayerVehicles(profile.Identity.AccountID)().(msgVehicles)
	if !ok {
		profile.addSectionErrorMessage("vehicles", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("vehicles", msg.err)
		return
	}
	if msg.rows != nil {
		profile.Vehicles = msg.rows
	}
}

func addProfileFactions(profile *playerProfileResponse, playerID int64) {
	msg, ok := cmdFetchFactions().(msgFactions)
	if !ok {
		profile.addSectionErrorMessage("factions", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("factions", msg.err)
		return
	}
	ids := profileActorIDs(playerID, profile.Identity)
	for _, row := range msg.rows {
		if ids[row.ActorID] {
			profile.Factions = append(profile.Factions, row)
		}
	}
}

func addProfileSpecializations(profile *playerProfileResponse, playerID int64) {
	specPlayerID := playerID
	if profile.Identity != nil && profile.Identity.ControllerID != 0 {
		specPlayerID = profile.Identity.ControllerID
	}
	msg, ok := cmdFetchPlayerSpecs(specPlayerID)().(msgSpecs)
	if !ok {
		profile.addSectionErrorMessage("specializations", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("specializations", msg.err)
		return
	}
	if msg.rows != nil {
		profile.Specializations = msg.rows
	}
}

func addProfileCharacterXP(profile *playerProfileResponse, playerID int64) {
	msg, ok := cmdFetchCharXP(playerID)().(msgCharXP)
	if !ok {
		profile.addSectionErrorMessage("character_xp", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("character_xp", msg.err)
		return
	}
	profile.CharacterXP = &playerProfileCharXP{XP: msg.xp, Level: msg.level}
}

func addProfileJourney(profile *playerProfileResponse) {
	if profile.Identity == nil || profile.Identity.AccountID == 0 {
		profile.addSectionErrorMessage("journey_summary", "account_id unavailable")
		return
	}
	msg, ok := cmdFetchJourneyNodes(profile.Identity.AccountID)().(msgJourney)
	if !ok {
		profile.addSectionErrorMessage("journey_summary", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("journey_summary", msg.err)
		return
	}
	profile.JourneySummary = summarizeProfileJourney(msg.rows)
}

func summarizeProfileJourney(rows []journeyNode) playerProfileJourneySummary {
	summary := playerProfileJourneySummary{PreviewNodes: []journeyNode{}}
	for _, node := range rows {
		summary.TotalNodes++
		if node.IsComplete {
			summary.CompleteNodes++
		}
		if node.IsRevealed {
			summary.RevealedNodes++
		}
		if node.HasPendingReward {
			summary.PendingRewards++
		}
		if len(summary.PreviewNodes) < playerProfilePreviewLimit {
			summary.PreviewNodes = append(summary.PreviewNodes, node)
		}
	}
	return summary
}

func addProfileEvents(profile *playerProfileResponse, playerID int64) {
	msg, ok := cmdFetchEventLog(playerID)().(msgEvents)
	if !ok {
		profile.addSectionErrorMessage("recent_events", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("recent_events", msg.err)
		return
	}
	if msg.rows != nil {
		profile.RecentEvents = msg.rows
	}
}

func addProfileDungeons(profile *playerProfileResponse, playerID int64) {
	msg, ok := cmdFetchPlayerDungeons(playerID)().(msgDungeons)
	if !ok {
		profile.addSectionErrorMessage("dungeon_history", "internal error")
		return
	}
	if msg.err != nil {
		profile.addSectionError("dungeon_history", msg.err)
		return
	}
	if msg.rows != nil {
		profile.DungeonHistory = msg.rows
	}
}

func profileControllerIDs(playerID int64, identity *playerInfo) map[int64]bool {
	ids := map[int64]bool{playerID: true}
	if identity != nil && identity.ControllerID != 0 {
		ids[identity.ControllerID] = true
	}
	return ids
}

func profileActorIDs(playerID int64, identity *playerInfo) map[int64]bool {
	ids := map[int64]bool{playerID: true}
	if identity != nil && identity.ID != 0 {
		ids[identity.ID] = true
	}
	if identity != nil && identity.ControllerID != 0 {
		ids[identity.ControllerID] = true
	}
	return ids
}
