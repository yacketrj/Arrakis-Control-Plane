package main

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	registerPublicRoutes(mux)
	registerDiscordAuthRoutes(mux)
	registerSelfServiceRoutes(mux)
	registerCoreAdminRoutes(mux)
	registerBattlegroupRoutes(mux)
	registerPlayerRoutes(mux)
	registerInventoryCoordinationRoutes(mux)
	registerDatabaseRoutes(mux)
	registerLogRoutes(mux)
	registerNotificationRoutes(mux)
	registerStorageRoutes(mux)
	registerBlueprintRoutes(mux)
}

func registerPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/public/status", handlePublicStatus)
}

func registerDiscordAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/auth/discord/login", handleDiscordLogin)
	mux.HandleFunc("GET /api/v1/auth/discord/callback", handleDiscordCallback)
	mux.HandleFunc("GET /api/v1/auth/discord/me", handleDiscordMe)
	mux.HandleFunc("POST /api/v1/auth/discord/logout", handleDiscordLogout)
	mux.HandleFunc("GET /api/v1/auth/discord/users", handleDiscordUsers)
	mux.HandleFunc("GET /api/v1/auth/discord/player-links", handleListDiscordPlayerLinks)
	mux.HandleFunc("POST /api/v1/auth/discord/player-links", handleUpsertDiscordPlayerLink)
	mux.HandleFunc("DELETE /api/v1/auth/discord/player-links/{discord_id}", handleDeleteDiscordPlayerLink)
}

func registerSelfServiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/self/player-link", handleSelfPlayerLink)
	mux.HandleFunc("GET /api/v1/self/player-card", handleSelfPlayerCard)
}

func registerCoreAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/status", handleStatus)
	mux.HandleFunc("POST /api/v1/reconnect", handleReconnect)
	mux.HandleFunc("GET /api/v1/connectivity/diagnostics", handleConnectivityDiagnostics)
	mux.HandleFunc("GET /api/v1/diagnostics/export", handleDiagnosticExport)
	mux.HandleFunc("GET /api/v1/audit/events", handleAdminAuditEvents)
	mux.HandleFunc("GET /api/v1/mutation-safety/classify", handleMutationSafetyClassify)
}

func registerBattlegroupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/battlegroup/status", handleBGStatus)
	mux.HandleFunc("GET /api/v1/battlegroup/health", handleBGHealth)
	mux.HandleFunc("POST /api/v1/battlegroup/exec", handleBGExec)
	mux.HandleFunc("GET /api/v1/battlegroup/pods", handleBGPods)
}

func registerPlayerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/players", handleGetPlayers)
	mux.HandleFunc("GET /api/v1/players/online", handleGetOnlineState)
	mux.HandleFunc("GET /api/v1/players/currency", handleGetCurrency)
	mux.HandleFunc("GET /api/v1/players/factions", handleGetFactions)
	mux.HandleFunc("GET /api/v1/players/specs", handleGetSpecs)
	mux.HandleFunc("GET /api/v1/players/templates", handleGetTemplates)
	mux.HandleFunc("POST /api/v1/players/templates/refresh", handleRefreshTemplates)
	mux.HandleFunc("GET /api/v1/players/{id}/profile", handleGetPlayerProfile)
	mux.HandleFunc("GET /api/v1/players/{id}/inventory", handleGetInventory)
	mux.HandleFunc("GET /api/v1/players/{id}/journey", handleGetJourney)
	mux.HandleFunc("POST /api/v1/players/give-item", handleGiveItems)
	mux.HandleFunc("POST /api/v1/players/give-currency", handleGiveCurrency)
	mux.HandleFunc("POST /api/v1/players/grant-live", handleGrantLive)
	mux.HandleFunc("POST /api/v1/players/give-faction-rep", handleGiveFactionRep)
	mux.HandleFunc("POST /api/v1/players/give-scrip", handleGiveScrip)
	mux.HandleFunc("POST /api/v1/players/award-xp", handleAwardXP)
	mux.HandleFunc("POST /api/v1/players/award-char-xp", handleAwardCharXP)
	mux.HandleFunc("POST /api/v1/players/award-intel", handleAwardIntel)
	mux.HandleFunc("POST /api/v1/players/kick", handleKick)
	mux.HandleFunc("DELETE /api/v1/players/item/{id}", handleDeleteItem)
	mux.HandleFunc("POST /api/v1/players/item/stack-size", handleSetItemStackSize)
	mux.HandleFunc("POST /api/v1/players/reset-spec", handleResetSpec)
	mux.HandleFunc("POST /api/v1/players/set-faction-tier", handleSetFactionTier)
	mux.HandleFunc("POST /api/v1/players/journey/complete", handleJourneyComplete)
	mux.HandleFunc("POST /api/v1/players/journey/reset", handleJourneyReset)
	mux.HandleFunc("POST /api/v1/players/journey/wipe", handleJourneyWipe)
	mux.HandleFunc("POST /api/v1/players/delete-tutorials", handleDeleteTutorials)
	mux.HandleFunc("POST /api/v1/players/wipe-codex", handleWipeCodex)
	mux.HandleFunc("GET /api/v1/players/{id}/char-xp", handleGetCharXP)
	mux.HandleFunc("GET /api/v1/players/{id}/specs", handleGetPlayerSpecs)
	mux.HandleFunc("POST /api/v1/players/set-spec-xp", handleSetSpecXP)
	mux.HandleFunc("GET /api/v1/players/{id}/vehicles", handleGetPlayerVehicles)
	mux.HandleFunc("POST /api/v1/players/repair-item", handleRepairItem)
	mux.HandleFunc("GET /api/v1/players/partitions", handleGetPartitions)
	mux.HandleFunc("POST /api/v1/players/teleport", handleTeleportPlayer)
	mux.HandleFunc("GET /api/v1/players/{id}/events", handleGetPlayerEvents)
	mux.HandleFunc("GET /api/v1/players/{id}/dungeons", handleGetPlayerDungeons)
}

func registerInventoryCoordinationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/inventory/requests", handleListInventoryRequests)
	mux.HandleFunc("POST /api/v1/inventory/requests", handleCreateInventoryRequest)
	mux.HandleFunc("PATCH /api/v1/inventory/requests/{id}", handleUpdateInventoryRequest)
	mux.HandleFunc("GET /api/v1/inventory/orders", handleListInventoryOrders)
	mux.HandleFunc("POST /api/v1/inventory/orders", handleCreateInventoryOrder)
	mux.HandleFunc("PATCH /api/v1/inventory/orders/{id}", handleUpdateInventoryOrder)
}

func registerDatabaseRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/database/tables", handleDBTables)
	mux.HandleFunc("GET /api/v1/database/describe", handleDBDescribe)
	mux.HandleFunc("GET /api/v1/database/sample", handleDBSample)
	mux.HandleFunc("GET /api/v1/database/search", handleDBSearch)
	mux.HandleFunc("GET /api/v1/database/functions", handleDBFunctions)
	mux.HandleFunc("GET /api/v1/database/functions/inspect", handleDBFunctionInspect)
	mux.HandleFunc("POST /api/v1/database/sql", handleDBSQL)
}

func registerLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/logs/pods", handleLogPods)
	mux.HandleFunc("POST /api/v1/logs/stream-ticket", handleIssueLogStreamTicket)
	mux.HandleFunc("GET /api/v1/logs/stream", handleLogStream)
	mux.HandleFunc("GET /api/v1/logs/cheats", handleGetCheatLog)
}

func registerNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/notify", handleNotify)
}

func registerStorageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/storage", handleListStorage)
	mux.HandleFunc("GET /api/v1/storage/{id}/items", handleGetStorageItems)
	mux.HandleFunc("POST /api/v1/storage/{id}/give-item", handleGiveItemToStorage)
}

func registerBlueprintRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/blueprints", handleListBlueprints)
	mux.HandleFunc("GET /api/v1/blueprints/{id}/export", handleExportBlueprint)
	mux.HandleFunc("POST /api/v1/blueprints/import", handleImportBlueprint)
}
