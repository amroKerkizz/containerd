/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package images

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/containerd/containerd/v2/plugins"
)

// Configuration represents the configuration settings for images.
type Configuration struct {
    CRIPluginSettings map[string]interface{}
}

// NewConfiguration creates a new instance of Configuration.
func NewConfiguration() *Configuration {
    return &Configuration{
        CRIPluginSettings: make(map[string]interface{}),
    }
}

// SetSandboxImage sets the default sandbox image.
func (c *Configuration) SetSandboxImage(image string) {
    c.CRIPluginSettings["sandbox_image"] = image
}

// SetRegistryConfigPath sets the registry configuration path.
func (c *Configuration) SetRegistryConfigPath(path string) {
    if c.CRIPluginSettings["registry"] == nil {
        c.CRIPluginSettings["registry"] = make(map[string]interface{})
    }
    c.CRIPluginSettings["registry"].(map[string]interface{})["config_path"] = path
}

// VerifySandboxImageMigration tests the migration process of sandbox images into a new configuration format.
// It sets up a default sandbox image, performs configuration migration, and verifies the migration result.
func VerifySandboxImageMigration(t *testing.T) {
    config := NewConfiguration()
    config.SetSandboxImage("rancher/mirrored-pause:3.9-amd64")

    // Perform configuration migration
    configMigration(context.Background(), 2, config.CRIPluginSettings)

    // Verify migration result
    extractedImages, exists := config.CRIPluginSettings[string(plugins.CRIServicePlugin)+".images"].(map[string]interface{})
    assert.True(t, exists, "expected extractedImages to exist")
    specificImages, exists := extractedImages["pinned_images"].(map[string]interface{})
    assert.True(t, exists, "expected specificImages to exist")
    migratedSandboxImage, exists := specificImages["sandbox"].(string)
    assert.True(t, exists, "expected migratedSandboxImage to exist")
    assert.Equal(t, "rancher/mirrored-pause:3.9-amd64", migratedSandboxImage)
}

// VerifyRegistryConfigurationMigration tests the migration of registry configuration to ensure path integrity.
// It sets up an expected path for the registry configuration, performs configuration migration, and verifies the migration result.
func VerifyRegistryConfigurationMigration(t *testing.T) {
    config := NewConfiguration()
    config.SetRegistryConfigPath("/etc/containerd/certs.d")

    // Perform configuration migration
    configMigration(context.Background(), 2, config.CRIPluginSettings)

    // Verify migration result
    extractedImages, exists := config.CRIPluginSettings[string(plugins.CRIServicePlugin)+".images"].(map[string]interface{})
    assert.True(t, exists, "expected extractedImages to exist")
    registrySettings, exists := extractedImages["registry"].(map[string]interface{})
    assert.True(t, exists, "expected registrySettings to exist")
    actualPath, exists := registrySettings["config_path"].(string)
    assert.True(t, exists, "expected actualPath to exist")
    assert.Equal(t, "/etc/containerd/certs.d", actualPath)
}
