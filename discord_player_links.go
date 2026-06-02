package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const discordPlayerLinkStoreDefaultPath = "discord-player-links.json"

var discordPlayerLinkStoreMu sync.Mutex

type discordPlayerLink struct {
	DiscordID         string `json:"discord_id"`
	PlayerID          int64  `json:"player_id"`
	PlayerName        string `json:"player_name,omitempty"`
	Notes             string `json:"notes,omitempty"`
	LinkedByDiscordID string `json:"linked_by_discord_id,omitempty"`
	LinkedByAuthType   string `json:"linked_by_auth_type,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type discordPlayerLinkStore struct {
	Links []discordPlayerLink `json:"links"`
}

type discordPlayerLinkPayload struct {
	DiscordID  string `json:"discord_id"`
	PlayerID   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	Notes      string `json:"notes"`
}

type selfPlayerCardSummary struct {
	DiscordID        string                         `json:"discord_id"`
	PlayerID         int64                          `json:"player_id"`
	PlayerName       string                         `json:"player_name,omitempty"`
	Class            string                         `json:"class,omitempty"`
	Map              string                         `json:"map,omitempty"`
	OnlineStatus     string                         `json:"online_status,omitempty"`
	Location         *playerProfileLocation         `json:"location,omitempty"`
	InventorySummary playerProfileInventorySummary  `json:"inventory_summary"`
	VehicleCount     int                            `json:"vehicle_count"`
	Currencies       []currencyRow                  `json:"currencies"`
	Factions         []factionRep                   `json:"factions"`
	Specializations  []specTrack                    `json:"specializations"`
	CharacterXP      *playerProfileCharXP           `json:"character_xp,omitempty"`
	JourneySummary   playerProfileJourneySummary    `json:"journey_summary"`
	SectionErrors    []playerProfileSectionError    `json:"section_errors"`
}

func discordPlayerLinkStorePath() string {
	if path := strings.TrimSpace(os.Getenv("DISCORD_PLAYER_LINK_STORE")); path != "" {
		return path
	}
	return discordPlayerLinkStoreDefaultPath
}

func loadDiscordPlayerLinkStore() (discordPlayerLinkStore, error) {
	store := discordPlayerLinkStore{}
	data, err := os.ReadFile(discordPlayerLinkStorePath())
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return store, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return store, nil
	}
	if err := json.Unmarshal(data, &store); err != nil {
		return store, fmt.Errorf("decode Discord player link store: %w", err)
	}
	return store, nil
}

func saveDiscordPlayerLinkStore(store discordPlayerLinkStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(discordPlayerLinkStorePath(), data, 0o600)
}

func cleanDiscordPlayerLinkText(value string, maxLen int, field string, required bool) (string, error) {
	trimmed := strings.TrimSpace(value)
	if required && trimmed == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	if trimmed != "" && containsUnsafeControl(trimmed) {
		return "", fmt.Errorf("%s contains unsupported control characters", field)
	}
	if len(trimmed) > maxLen {
		return "", fmt.Errorf("%s exceeds %d characters", field, maxLen)
	}
	return trimmed, nil
}

func normalizeDiscordPlayerLinkPayload(payload discordPlayerLinkPayload, existing *discordPlayerLink, linkedByDiscordID string, linkedByAuthType string, now time.Time) (discordPlayerLink, error) {
	discordID, err := cleanDiscordPlayerLinkText(payload.DiscordID, 128, "discord_id", true)
	if err != nil {
		return discordPlayerLink{}, err
	}
	if payload.PlayerID <= 0 {
		return discordPlayerLink{}, fmt.Errorf("player_id must be greater than zero")
	}
	playerName, err := cleanDiscordPlayerLinkText(payload.PlayerName, 160, "player_name", false)
	if err != nil {
		return discordPlayerLink{}, err
	}
	notes, err := cleanDiscordPlayerLinkText(payload.Notes, 1000, "notes", false)
	if err != nil {
		return discordPlayerLink{}, err
	}
	linkedByDiscordID, err = cleanDiscordPlayerLinkText(linkedByDiscordID, 128, "linked_by_discord_id", false)
	if err != nil {
		return discordPlayerLink{}, err
	}
	linkedByAuthType, err = cleanDiscordPlayerLinkText(linkedByAuthType, 64, "linked_by_auth_type", false)
	if err != nil {
		return discordPlayerLink{}, err
	}
	stamp := now.UTC().Format(time.RFC3339Nano)
	createdAt := stamp
	if existing != nil && existing.CreatedAt != "" {
		createdAt = existing.CreatedAt
	}
	return discordPlayerLink{
		DiscordID:         discordID,
		PlayerID:          payload.PlayerID,
		PlayerName:        playerName,
		Notes:             notes,
		LinkedByDiscordID: linkedByDiscordID,
		LinkedByAuthType:   linkedByAuthType,
		CreatedAt:         createdAt,
		UpdatedAt:         stamp,
	}, nil
}

func findDiscordPlayerLink(store discordPlayerLinkStore, discordID string) (discordPlayerLink, bool) {
	discordID = strings.TrimSpace(discordID)
	for _, link := range store.Links {
		if link.DiscordID == discordID {
			return link, true
		}
	}
	return discordPlayerLink{}, false
}

func upsertDiscordPlayerLink(store *discordPlayerLinkStore, link discordPlayerLink) {
	for i := range store.Links {
		if store.Links[i].DiscordID == link.DiscordID {
			store.Links[i] = link
			return
		}
	}
	store.Links = append(store.Links, link)
}

func deleteDiscordPlayerLink(store *discordPlayerLinkStore, discordID string) bool {
	for i := range store.Links {
		if store.Links[i].DiscordID == discordID {
			store.Links = append(store.Links[:i], store.Links[i+1:]...)
			return true
		}
	}
	return false
}

func requestAuthAttribution(r *http.Request) (string, string) {
	if session, ok := discordSessionFromRequest(r); ok {
		return session.DiscordID, "discord"
	}
	if bearerToken(r.Header.Get("Authorization")) != "" || r.Header.Get("X-Admin-Token") != "" {
		return "", "admin-token"
	}
	return "", "unknown"
}

func handleListDiscordPlayerLinks(w http.ResponseWriter, r *http.Request) {
	discordPlayerLinkStoreMu.Lock()
	defer discordPlayerLinkStoreMu.Unlock()

	store, err := loadDiscordPlayerLinkStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	links := append([]discordPlayerLink{}, store.Links...)
	sort.Slice(links, func(i, j int) bool { return links[i].DiscordID < links[j].DiscordID })
	jsonOK(w, links)
}

func handleUpsertDiscordPlayerLink(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	defer r.Body.Close()
	var payload discordPlayerLinkPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	linkedByDiscordID, linkedByAuthType := requestAuthAttribution(r)

	discordPlayerLinkStoreMu.Lock()
	defer discordPlayerLinkStoreMu.Unlock()

	store, err := loadDiscordPlayerLinkStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	var existing *discordPlayerLink
	if found, ok := findDiscordPlayerLink(store, payload.DiscordID); ok {
		existing = &found
	}
	link, err := normalizeDiscordPlayerLinkPayload(payload, existing, linkedByDiscordID, linkedByAuthType, time.Now())
	if err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	upsertDiscordPlayerLink(&store, link)
	if err := saveDiscordPlayerLinkStore(store); err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, link)
}

func handleDeleteDiscordPlayerLink(w http.ResponseWriter, r *http.Request) {
	discordID, err := cleanDiscordPlayerLinkText(r.PathValue("discord_id"), 128, "discord_id", true)
	if err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}

	discordPlayerLinkStoreMu.Lock()
	defer discordPlayerLinkStoreMu.Unlock()

	store, err := loadDiscordPlayerLinkStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	if !deleteDiscordPlayerLink(&store, discordID) {
		jsonErr(w, fmt.Errorf("Discord player link not found"), http.StatusNotFound)
		return
	}
	if err := saveDiscordPlayerLinkStore(store); err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"ok": "deleted"})
}

func currentDiscordPlayerLink(r *http.Request) (discordPlayerLink, bool, error) {
	session, ok := discordSessionFromRequest(r)
	if !ok {
		return discordPlayerLink{}, false, fmt.Errorf("Discord session required")
	}
	discordPlayerLinkStoreMu.Lock()
	defer discordPlayerLinkStoreMu.Unlock()

	store, err := loadDiscordPlayerLinkStore()
	if err != nil {
		return discordPlayerLink{}, false, err
	}
	link, found := findDiscordPlayerLink(store, session.DiscordID)
	return link, found, nil
}

func handleSelfPlayerLink(w http.ResponseWriter, r *http.Request) {
	link, found, err := currentDiscordPlayerLink(r)
	if err != nil {
		jsonErr(w, err, http.StatusUnauthorized)
		return
	}
	if !found {
		jsonErr(w, fmt.Errorf("Discord identity is not linked to a player"), http.StatusNotFound)
		return
	}
	jsonOK(w, link)
}

func handleSelfPlayerCard(w http.ResponseWriter, r *http.Request) {
	link, found, err := currentDiscordPlayerLink(r)
	if err != nil {
		jsonErr(w, err, http.StatusUnauthorized)
		return
	}
	if !found {
		jsonErr(w, fmt.Errorf("Discord identity is not linked to a player"), http.StatusNotFound)
		return
	}
	profile, profileFound := buildPlayerProfile(link.PlayerID)
	if !profileFound {
		jsonErr(w, fmt.Errorf("linked player was not found"), http.StatusNotFound)
		return
	}
	summary := selfPlayerCardSummary{
		DiscordID:        link.DiscordID,
		PlayerID:         link.PlayerID,
		PlayerName:       link.PlayerName,
		Location:         profile.Location,
		InventorySummary: profile.InventorySummary,
		VehicleCount:     len(profile.Vehicles),
		Currencies:       profile.Currencies,
		Factions:         profile.Factions,
		Specializations:  profile.Specializations,
		CharacterXP:      profile.CharacterXP,
		JourneySummary:   profile.JourneySummary,
		SectionErrors:    profile.SectionErrors,
	}
	if profile.Identity != nil {
		summary.PlayerName = profile.Identity.Name
		summary.Class = profile.Identity.Class
		summary.Map = profile.Identity.Map
		summary.OnlineStatus = profile.Identity.OnlineStatus
	}
	jsonOK(w, summary)
}

func parseLinkedPlayerID(raw string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("player_id must be greater than zero")
	}
	return value, nil
}
