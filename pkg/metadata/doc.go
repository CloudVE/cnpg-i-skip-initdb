// Package metadata contains the metadata of this plugin
package metadata

import "github.com/cloudnative-pg/cnpg-i/pkg/identity"

// PluginName is the name of the plugin
const PluginName = "cnpg-i-skip-initdb.leonardoce.github.com"

// Data is the metadata of this plugin
var Data = identity.GetPluginMetadataResponse{
	Name:          PluginName,
	Version:       "0.1.0",
	DisplayName:   "CNPG Skip InitDB plugin",
	ProjectUrl:    "https://github.com/CloudVE/cnpg-i-skip-initdb",
	RepositoryUrl: "https://github.com/CloudVE/cnpg-i-skip-initdb",
	License:       "Apache License 2.0",
	LicenseUrl:    "https://github.com/CloudVE/cnpg-i-skip-initdb/LICENSE",
	Maturity:      "beta",
}
