package locate

import (
    "reflect"
    "testing"
)

func TestLocate(t *testing.T) {
    type args struct {
        name string
    }
    tests := []struct {
        name           string
        args           args
        wantLocateInfo map[int]string
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if gotLocateInfo := Locate(tt.args.name); !reflect.DeepEqual(gotLocateInfo, tt.wantLocateInfo) {
                t.Errorf("Locate() = %v, want %v", gotLocateInfo, tt.wantLocateInfo)
            }
        })
    }
}
