// SPDX-License-Identifier: Apache-2.0

package vault

import "github.com/Cargill/vela-server/constants"

// Driver outputs the configured secret driver.
func (c *Client) Driver() string {
	return constants.DriverVault
}
