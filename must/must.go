// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package must

import (
	"encoding/json"
	"os"
)

// UnmarshalFile - unmarshal json file or panic
func UnmarshalFile(dest interface{}, src string) {
	f, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(f).Decode(dest)
	if err != nil {
		panic(err)
	}
}
