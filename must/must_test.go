// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package must

import "testing"

type data struct {
	Foo string
}

func TestUnmarshalFile(t *testing.T) {
	type args struct {
		dest interface{}
		src  string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			"",
			args{
				new(data),
				"../testdata/valid.json",
			},
			false,
		},
		{
			"",
			args{
				new(data),
				"../testdata/invalid.json",
			},
			true,
		},
		{
			"",
			args{
				new(data),
				"../testdata/notfound.json",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				// recover from panic if one occured. Set err to nil otherwise.
				if recover() != nil {
					if !tt.wantPanic {
						t.Errorf("UnmarshalFile() - want panic")
					}
				} else {
					if tt.wantPanic {
						t.Errorf("UnmarshalFile() - want panic")
					}
				}
			}()
			UnmarshalFile(tt.args.dest, tt.args.src)
		})
	}
}
