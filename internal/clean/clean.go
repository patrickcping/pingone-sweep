package clean

import (
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

type CleanEnvironmentConfig struct {
	Client        *management.APIClient
	EnvironmentID string
	DryRun        bool
}

// Sample populations and users
//
// Populations: [Sample Users, More Sample Users]
// Groups: [Sample Group, Another Sample Group]

// Disable System Applications
//
// Applications: [PingOne Application Portal]

// SOP
//
// Policies: [Multi_Factor, Single_Factor (if not default)]

// Password Policies
//
// Policies: [Standard (if not default), Basic, Passphrase]

// MFA Policies
//
// Policies: [Default MFA Policy (if not default)]

// FIDO Policies
//
// Policies: [Passkeys (if not default), Security Keys]

// Risk Policies
//
// Policies: [Default Risk Policy (if not default)]

// Verify Policies
//
// Policies: [Default Verify Policy (if not default)]

// Authorize decision endpoints
//
// Endpoints: [DEV, TEST, PROD]

// Notification Templates
//
// Unknown

// Notification Policies
//
// Policies: [Default Notification Policy (if not default)]

// Self-Service Selection
//
// All Self-Service capabilities

// Languages
//
// All that can be

// Key Rotation Policy
//
// KRP: TBC
