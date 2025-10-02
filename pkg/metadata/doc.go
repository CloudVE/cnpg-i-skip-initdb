// Package metadata contains the metadata of this plugin
package metadata

import "github.com/cloudnative-pg/cnpg-i/pkg/identity"

// PluginName is the name of the plugin
const PluginName = "cnpg-i-skip-initdb.leonardoce.github.com"

// Data is the metadata of this plugin
var Data = identity.GetPluginMetadataResponse{
	Name:          PluginName,
	Version:       "0.0.1",
	DisplayName:   "Plugin feature showcase",
	ProjectUrl:    "https://github.com/leonardoce/cnpg-i-skip-initdb",
	RepositoryUrl: "https://github.com/leonardoce/cnpg-i-skip-initdb",
	License:       "Proprietary",
	LicenseUrl:    "https://github.com/leonardoce/cnpg-i-skip-initdb/LICENSE",
	Maturity:      "alpha",
}
