package main

import (
	"net/http"
	"strings"
)

type mutationSafetyClass struct {
	Action                   string   `json:"action"`
	Risk                     string   `json:"risk"`
	RequiresReason           bool     `json:"requires_reason"`
	ReasonEnforcementEnabled bool     `json:"reason_enforcement_enabled"`
	RequiresPreview          bool     `json:"requires_preview"`
	Destructive              bool     `json:"destructive"`
	RollbackHint             string   `json:"rollback_hint,omitempty"`
	OperatorWarnings         []string `json:"operator_warnings,omitempty"`
	RecommendedPath          string   `json:"recommended_path,omitempty"`
}

func classifyMutationRequest(method, path string) mutationSafetyClass {
	action := auditActionName(method, path)
	risk := mutationRiskForRequest(method, path)
	classification := mutationSafetyClass{
		Action:                   action,
		Risk:                     risk,
		RequiresReason:           risk == "high" || risk == "destructive",
		ReasonEnforcementEnabled: adminReasonEnforcementEnabled(),
		RequiresPreview:          risk == "high" || risk == "destructive",
		Destructive:              risk == "destructive",
	}

	lower := strings.ToLower(path)
	switch {
	case method == http.MethodDelete || strings.Contains(lower, "/players/item/"):
		classification.RollbackHint = "Capture the item row before deletion. Restore from audit-supported snapshot when Inventory Studio v2 is available."
		classification.OperatorWarnings = append(classification.OperatorWarnings, "Deleting inventory data is not live-safe and may require player relog or zone refresh.")
	case strings.Contains(lower, "/give-item"):
		classification.RecommendedPath = "Use Direct Inventory Write for graded or augmented items. Use Claim Rewards Queue only for plain live grants."
		classification.OperatorWarnings = append(classification.OperatorWarnings, "Direct inventory writes may not appear for online players until relog.")
	case strings.Contains(lower, "/grant-live"):
		classification.RecommendedPath = "Use only for plain template-plus-amount claim rewards. Do not use for graded or augmented items."
	case strings.Contains(lower, "/teleport"):
		classification.RollbackHint = "Record the player's prior partition/location before teleporting."
		classification.OperatorWarnings = append(classification.OperatorWarnings, "Teleport operations should be used for rescue/support, not routine travel.")
	case strings.Contains(lower, "/journey/") || strings.Contains(lower, "/wipe") || strings.Contains(lower, "/reset"):
		classification.RollbackHint = "Snapshot affected progression rows before mutation."
		classification.OperatorWarnings = append(classification.OperatorWarnings, "Progression edits can permanently alter quest/tutorial/codex state.")
	case strings.Contains(lower, "/storage/"):
		classification.OperatorWarnings = append(classification.OperatorWarnings, "Storage edits may require server zone restart before players see changes.")
	}
	return classification
}

func mutationSafetyForPath(method, path string) mutationSafetyClass {
	return classifyMutationRequest(method, path)
}
