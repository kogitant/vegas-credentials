package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigPath(t *testing.T) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Could not resolve homedir")
	}

	want := fmt.Sprintf("%s%s", homedir, "/.aws/config")

	output, err := resolveConfigPath()
	if output != want || err != nil {
		t.Fatalf("Got %s, want %s", output, want)
	}
}

func TestLoadConfig(t *testing.T) {

	configPath := getTestdataFilePath("valid-minimal-config.ini")
	config, err := loadConfig(configPath)

	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}

	want := "arn:aws:iam::123456789012:role/Demo"
	output := config.Section("profile my-profile").Key("assume_role_arn").String()

	if output != want {
		t.Fatalf("Got %s, want %s", output, want)
	}
}

func TestLoadProfileMinimal(t *testing.T) {
	configPath := getTestdataFilePath("valid-minimal-config.ini")
	config, err := loadConfig(configPath)

	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}

	profileName := "my-profile"

	profile, err := loadProfile(config, profileName)
	if err != nil {
		t.Fatalf("Error loading profile: %s", err.Error())
	}

	want := Profile{
		AssumeRoleArn:   "arn:aws:iam::123456789012:role/Demo",
		SourceProfile:   "default",
		MfaSerial:       "arn:aws:iam::123456789012:mfa/JaneDoeMFA",
		DurationSeconds: 3600,
	}

	if profile != want {
		t.Fatalf(`Got %q, want %q`, profile, want)
	}
}

func TestLoadProfileFull(t *testing.T) {
	configPath := getTestdataFilePath("valid-full-config.ini")
	config, err := loadConfig(configPath)

	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}

	profileName := "my-profile"

	profile, err := loadProfile(config, profileName)
	if err != nil {
		t.Fatalf("Error loading profile: %s", err.Error())
	}

	want := Profile{
		AssumeRoleArn:   "arn:aws:iam::123456789012:role/Demo",
		SourceProfile:   "default",
		MfaSerial:       "arn:aws:iam::123456789012:mfa/JaneDoeMFA",
		DurationSeconds: 900,
		ExternalID:      "extid123",
		RoleSessionName: "my-session",
		Region:          "eu-west-1",
		YubikeySerial:   "123456",
		YubikeyLabel:    "foobar",
	}

	if profile != want {
		t.Fatalf(`Got %q, want %q`, profile, want)
	}
}

func TestLoadProfileInvalid(t *testing.T) {
	configPath := getTestdataFilePath("invalid-config.ini")
	config, err := loadConfig(configPath)

	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}

	profileName := "my-profile"

	profile, err := loadProfile(config, profileName)

	if err == nil {
		t.Fatalf("Expected error, but got %q", profile)
	}

	want := fmt.Sprintf("Missing assume_role_arn from profile %s config", profileName)

	if err.Error() != want {
		t.Fatalf("Error loading profile: %s", err.Error())
	}
}

func TestLoadProfileMissing(t *testing.T) {
	configPath := getTestdataFilePath("missing-profile.ini")
	config, err := loadConfig(configPath)

	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}

	profileName := "my-profile"

	profile, err := loadProfile(config, profileName)

	if err == nil {
		t.Fatalf("Expected error, but got %q", profile)
	}

	want := fmt.Sprintf(`section "profile %s" does not exist`, profileName)

	if err.Error() != want {
		t.Fatalf("Error loading profile: %s", err.Error())
	}
}

func getCurrentDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func getTestdataFilePath(file string) string {
	cwd := getCurrentDirectory()
	return filepath.Join(cwd, "testdata", file)
}
