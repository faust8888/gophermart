package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewConfig_WithDefaultValues(t *testing.T) {
	cfg, err := New()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, DatabaseURIDefaultValue, cfg.DatabaseURI)
	require.Equal(t, RunAddressDefaultValue, cfg.RunAddress)
	require.Equal(t, AccrualSystemDefaultValue, cfg.AccrualSystemAddress)
}

func TestNewConfig_WithFlagOverrideValues(t *testing.T) {
	runAddressOverrideValue := "localhost:9999"
	databaseURIOverrideValue := "postgres://localhost:5432"
	accrualSystemAddressOverrideValue := "localhost:7777/accrual"
	runAddressEnvironmentOverrideValue := "localhost:1212"

	// Override flags using command-line arguments
	os.Args = []string{
		GophermartFlagName,
		"--" + RunAddressFlag, runAddressOverrideValue,
		"--" + AccrualSystemAddressFlag, accrualSystemAddressOverrideValue,
		"--" + DatabaseURIFlag, databaseURIOverrideValue,
	}
	// environment variables have priority over flags
	os.Setenv("RUN_ADDRESS", runAddressEnvironmentOverrideValue)

	cfg, err := New()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, cfg.DatabaseURI, databaseURIOverrideValue)
	require.Equal(t, cfg.AccrualSystemAddress, accrualSystemAddressOverrideValue)
	require.Equal(t, cfg.RunAddress, runAddressEnvironmentOverrideValue)
}

func TestNewConfig_WithEnvOverrideValues(t *testing.T) {
	runAddressOverrideValue := "localhost:1111"
	databaseURIOverrideValue := "postgres://localhost:6544"
	accrualSystemAddressOverrideValue := "localhost:3333/accrual"

	os.Setenv("RUN_ADDRESS", runAddressOverrideValue)
	os.Setenv("DATABASE_URI", databaseURIOverrideValue)
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", accrualSystemAddressOverrideValue)

	cfg, err := New()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, cfg.DatabaseURI, databaseURIOverrideValue)
	require.Equal(t, cfg.RunAddress, runAddressOverrideValue)
	require.Equal(t, cfg.AccrualSystemAddress, accrualSystemAddressOverrideValue)
}
