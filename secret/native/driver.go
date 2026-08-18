// SPDX-License-Identifier: Apache-2.0

package native

import "github.com/Cargill/vela-server/constants"

// Driver outputs the configured secret driver.
func (c *Client) Driver() string {
	return constants.DriverNative
}
