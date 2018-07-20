// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package fields

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StrArray type for easy db->json->db convert string arrays
type StrArray []string

// Scan unmarshal json-field into []string
func (h *StrArray) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		return json.Unmarshal(v, h)
	case string:
		return json.Unmarshal([]byte(v), h)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

// Value marshals []string into json field
func (h StrArray) Value() (driver.Value, error) {
	return json.Marshal(h)
}
