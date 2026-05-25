package main

import "testing"

func TestSummarizeProfileInventory(t *testing.T) {
	items := []itemInfo{
		{ID: 1, TemplateID: "tpl-a", StackSize: 5},
		{ID: 2, TemplateID: "TPL-A", StackSize: 3},
		{ID: 3, TemplateID: "tpl-b", StackSize: 7},
	}

	got := summarizeProfileInventory(items)

	if got.TotalItems != 3 {
		t.Fatalf("expected 3 total items, got %d", got.TotalItems)
	}
	if got.TotalStackSize != 15 {
		t.Fatalf("expected total stack size 15, got %d", got.TotalStackSize)
	}
	if got.UniqueTemplates != 2 {
		t.Fatalf("expected 2 unique templates, got %d", got.UniqueTemplates)
	}
	if len(got.PreviewItems) != 3 {
		t.Fatalf("expected all preview items, got %d", len(got.PreviewItems))
	}
}

func TestSummarizeProfileInventoryPreviewLimit(t *testing.T) {
	items := make([]itemInfo, playerProfilePreviewLimit+5)
	for i := range items {
		items[i] = itemInfo{ID: int64(i + 1), TemplateID: "tpl", StackSize: 1}
	}

	got := summarizeProfileInventory(items)

	if got.TotalItems != playerProfilePreviewLimit+5 {
		t.Fatalf("expected total item count to include all rows, got %d", got.TotalItems)
	}
	if len(got.PreviewItems) != playerProfilePreviewLimit {
		t.Fatalf("expected preview limit %d, got %d", playerProfilePreviewLimit, len(got.PreviewItems))
	}
}

func TestSummarizeProfileJourney(t *testing.T) {
	nodes := []journeyNode{
		{NodeID: "a", IsComplete: true, IsRevealed: true, HasPendingReward: true},
		{NodeID: "b", IsComplete: false, IsRevealed: true, HasPendingReward: false},
		{NodeID: "c", IsComplete: true, IsRevealed: false, HasPendingReward: true},
	}

	got := summarizeProfileJourney(nodes)

	if got.TotalNodes != 3 {
		t.Fatalf("expected 3 total nodes, got %d", got.TotalNodes)
	}
	if got.CompleteNodes != 2 {
		t.Fatalf("expected 2 complete nodes, got %d", got.CompleteNodes)
	}
	if got.RevealedNodes != 2 {
		t.Fatalf("expected 2 revealed nodes, got %d", got.RevealedNodes)
	}
	if got.PendingRewards != 2 {
		t.Fatalf("expected 2 pending rewards, got %d", got.PendingRewards)
	}
	if len(got.PreviewNodes) != 3 {
		t.Fatalf("expected all preview nodes, got %d", len(got.PreviewNodes))
	}
}

func TestSelectProfileOnlineStateUsesControllerID(t *testing.T) {
	identity := &playerInfo{ID: 100, ControllerID: 200, Name: "Player", Map: "OldMap", OnlineStatus: "Offline"}
	rows := []onlineStateRow{
		{PlayerID: 300, Name: "Other", Map: "OtherMap", Status: "Online", LastSeen: "later"},
		{PlayerID: 200, Name: "Player", Map: "CurrentMap", Status: "Online", LastSeen: "now"},
	}

	got, found := selectProfileOnlineState(rows, 100, identity)

	if !found {
		t.Fatalf("expected online state match")
	}
	if got.PlayerID != 200 || got.Map != "CurrentMap" || got.Status != "Online" || got.LastSeen != "now" {
		t.Fatalf("unexpected online state: %#v", got)
	}
}

func TestProfileIDHelpers(t *testing.T) {
	identity := &playerInfo{ID: 100, ControllerID: 200}

	controllerIDs := profileControllerIDs(100, identity)
	if !controllerIDs[100] || !controllerIDs[200] {
		t.Fatalf("expected player and controller IDs in controller map: %#v", controllerIDs)
	}

	actorIDs := profileActorIDs(100, identity)
	if !actorIDs[100] || !actorIDs[200] {
		t.Fatalf("expected actor and controller IDs in actor map: %#v", actorIDs)
	}
}

func TestSafePlayerProfileError(t *testing.T) {
	cases := map[string]string{
		"not connected":                   "not connected",
		"internal error while aggregating": "internal error",
		"password=secret host=internal":    "section unavailable",
	}
	for input, want := range cases {
		got := safePlayerProfileError(testError(input))
		if got != want {
			t.Fatalf("safePlayerProfileError(%q) = %q, want %q", input, got, want)
		}
	}
}

type testError string

func (e testError) Error() string { return string(e) }
