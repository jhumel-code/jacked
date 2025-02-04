package test

import (
	"testing"

	"github.com/carbonetes/jacked/internal/config"
	"github.com/carbonetes/jacked/internal/events"
	jacked "github.com/carbonetes/jacked/internal/model"
	"github.com/carbonetes/jacked/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/model"
)

var (
	tr                            = bus.GetBus()
	testChannelManagerChannelName = "jacked-test-channel"
	testChannelManager            = tr.GetChannelManager()
	testImageName                 = "nginx"
	arguments                     = jacked.Arguments{
		DisableFileListing:  new(bool),
		SecretContentRegex:  new(string),
		DisableSecretSearch: new(bool),
		Image:               new(string),
		Dir:                 new(string),
		Tar:                 new(string),
		Quiet:               new(bool),
		OutputFile:          new(string),
		ExcludedFilenames:   &[]string{},
		EnabledParsers:      &[]string{},
		RegistryURI:         new(string),
		RegistryUsername:    new(string),
		RegistryPassword:    new(string),
		RegistryToken:       new(string),
	}
)

// Basic tests for transport-go package based on their github repository: see https://github.com/vmware/transport-go
func TestChannelManagerBoot(t *testing.T) {
	assert.Len(t, testChannelManager.GetAllChannels(), 0)
}

func TestCreateChannel(t *testing.T) {
	testChannelManager.CreateChannel(testChannelManagerChannelName)

	assert.Len(t, testChannelManager.GetAllChannels(), 1)
	fetchedChannel, _ := testChannelManager.GetChannel(testChannelManagerChannelName)

	assert.NotNil(t, fetchedChannel)
	assert.True(t, testChannelManager.CheckChannelExists(testChannelManagerChannelName))
}

func TestDestroyChannel(t *testing.T) {
	testChannelManager.CreateChannel(testChannelManagerChannelName)
	testChannelManager.DestroyChannel(testChannelManagerChannelName)

	fetchedChannel, err := testChannelManager.GetChannel(testChannelManagerChannelName)
	assert.Len(t, testChannelManager.GetAllChannels(), 0)
	assert.NotNil(t, err)
	assert.Nil(t, fetchedChannel)
}

func TestSubscribeChannelHandler(t *testing.T) {
	testChannelManager.CreateChannel(testChannelManagerChannelName)

	handler := func(*model.Message) {}
	uuid, err := testChannelManager.SubscribeChannelHandler(testChannelManagerChannelName, handler, false)
	assert.Nil(t, err)
	assert.NotNil(t, uuid)
	channel, _ := testChannelManager.GetChannel(testChannelManagerChannelName)
	assert.True(t, channel.ContainsHandlers())
}

// Scan test for diggity package
func TestDiggityScan(t *testing.T) {
	var cfg config.Configuration
	var packages []jacked.Package
	var secrets jacked.SecretResults

	cfg.SetDefault()

	arguments.Image = &testImageName
	arguments.DisableSecretSearch = &cfg.SecretConfig.Disabled
	arguments.SecretContentRegex = &cfg.SecretConfig.SecretRegex
	arguments.SecretMaxFileSize = cfg.SecretConfig.MaxFileSize
	arguments.EnabledParsers = &cfg.EnabledParsers
	arguments.DisableFileListing = &cfg.DisableFileListing
	arguments.RegistryURI = &cfg.Registry.URI
	arguments.RegistryToken = &cfg.Registry.Token
	arguments.RegistryUsername = &cfg.Registry.Username
	arguments.RegistryPassword = &cfg.Registry.Password

	sbom := events.RequestSBOMAnalysis(&arguments)
	assert.NotNil(t, sbom)

	parser.ParseSBOM(&sbom, &packages, &secrets)

	t.Logf("%v found packages", len(packages))

	assert.NotNil(t, packages)

}
