// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package fields

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestStrArray_Scan(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		h       *StrArray
		args    args
		want    *StrArray
		wantErr bool
	}{
		{
			"success",
			&StrArray{},
			args{
				val: []byte(`["1", "2"]`),
			},
			&StrArray{"1", "2"},
			false,
		},
		{
			"success",
			&StrArray{},
			args{
				val: `["1", "2"]`,
			},
			&StrArray{"1", "2"},
			false,
		},
		{
			"error",
			&StrArray{},
			args{
				val: 1000,
			},
			&StrArray{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.Scan(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("StrArray.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.h, tt.want) {
				t.Errorf("got - %+v, want - %+v", tt.h, tt.want)
			}
		})
	}
}

func TestStrArray_Value(t *testing.T) {
	tests := []struct {
		name    string
		h       StrArray
		want    driver.Value
		wantErr bool
	}{
		{
			"success",
			StrArray{"1","3","8"},
			[]byte(`["1","3","8"]`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("StrArray.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrArray.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
