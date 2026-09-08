package prd

import (
	"strings"
	"testing"
)

func TestStructuredPRDConstruction(t *testing.T) {
	svc := NewPRDService()
	req := GeneratePRDRequest{
		ProjectIdea:           "A modern cloud inventory and POS application",
		TargetUsers:           "Retail store owners and cashiers",
		BusinessGoals:         "Automate inventory tracking and reduce transaction time",
		MainFeatures:          "Realtime stock management\nBarcode scanner support\nReceipt printing",
		TechnicalRequirements: "Offline first with local sync",
	}

	result := svc.constructStructuredPRD("InventoryPro", req, "Laravel / Vue / MySQL")

	if !strings.Contains(result, "InventoryPro") {
		t.Errorf("Expected project name in PRD, but got: %s", result)
	}

	// Verify all 17 PRD sections are present
	sections := []string{
		"1. Product Overview",
		"2. Problem Statement",
		"3. Goals",
		"4. Target Users",
		"5. User Roles",
		"6. Functional Requirements",
		"7. Modules",
		"8. User Flows",
		"9. Database Requirements",
		"10. API Requirements",
		"11. UI Requirements",
		"12. Non Functional Requirements",
		"13. Security Requirements",
		"14. Performance Requirements",
		"15. Deployment Requirements",
		"16. Acceptance Criteria",
		"17. Development Phases",
	}

	for _, sec := range sections {
		if !strings.Contains(result, sec) {
			t.Errorf("Missing expected section: %s", sec)
		}
	}
}
