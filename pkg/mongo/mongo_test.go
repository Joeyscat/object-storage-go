package mongo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCli(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test InitCli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitCli()
		})
	}
}

func TestPutMetadata(t *testing.T) {
	InitCli()

	type args struct {
		name string
		hash string
		size uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test 1",
			args:    args{name: "xxx1", size: 10000, hash: "xxx"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PutMetadata(tt.args.name, tt.args.hash, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("PutMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddVersion(t *testing.T) {
	InitCli()

	type args struct {
		name    string
		hash    string
		version uint64
		size    uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test AddVersion",
			args:    args{name: "xxx1", size: 10000, hash: "xxx", version: 2},
			wantErr: false,
		},
	}
	for i := 0; i < 10; i++ {
		tests = append(tests, struct {
			name    string
			args    args
			wantErr bool
		}{
			name:    "Test AddVersion",
			args:    args{name: "xxx1", size: 10000, hash: "xxx", version: tests[0].args.version + uint64(i) + 1},
			wantErr: false,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddVersion(tt.args.name, tt.args.hash, tt.args.version, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("AddVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchAllVersions(t *testing.T) {
	InitCli()

	type args struct {
		name string
		from int64
		size int64
	}
	tests := []struct {
		name      string
		args      args
		wantMetas []*Metadata
		wantErr   bool
	}{
		{
			name:      "Test SearchAllVersions",
			args:      args{name: "xxx1", from: 0, size: 3},
			wantMetas: []*Metadata{{Version: 1}, {Version: 2}, {Version: 3}},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMetas, err := SearchAllVersions(tt.args.name, tt.args.from, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchAllVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, len(gotMetas), len(tt.wantMetas))
			for i, meta := range tt.wantMetas {
				t.Log(gotMetas[i])
				assert.Equal(t, meta.Version, gotMetas[i].Version)
			}
		})
	}
}

func TestSearchLatestVersion(t *testing.T) {
	InitCli()

	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantMeta *Metadata
		wantErr  bool
	}{
		{
			name:     "TestSearchLatestVersion",
			args:     args{name: "xxx1"},
			wantMeta: &Metadata{Name: "xxx1", Version: 12, Size: 10000, Hash: "xxx"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMeta, err := SearchLatestVersion(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchLatestVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMeta, tt.wantMeta) {
				t.Errorf("SearchLatestVersion() gotMeta = %v, want %v", gotMeta, tt.wantMeta)
			}
		})
	}
}

func TestGetMetadata(t *testing.T) {
	InitCli()

	type args struct {
		name    string
		version int
	}
	tests := []struct {
		name     string
		args     args
		wantMeta *Metadata
		wantErr  bool
	}{
		{
			name:     "TestGetMetadata_OK",
			args:     args{name: "xxx1", version: 1},
			wantMeta: &Metadata{Name: "xxx1", Version: 1, Size: 10000, Hash: "xxx"},
			wantErr:  false,
		},
		{
			name:     "TestGetMetadata_OK",
			args:     args{name: "xxx1", version: 3},
			wantMeta: &Metadata{Name: "xxx1", Version: 3, Size: 10000, Hash: "xxx"},
			wantErr:  false,
		},
		{
			name:     "TestGetMetadata_VERSION_NOT_FOUND",
			args:     args{name: "xxx1", version: 13},
			wantMeta: &Metadata{Name: "xxx1", Version: 0, Size: 0, Hash: ""},
			wantErr:  false,
		},
		{
			name:     "TestGetMetadata_METADATA_NOT_FOUND",
			args:     args{name: "xxx2", version: 1},
			wantMeta: &Metadata{Name: "xxx2", Version: 0, Size: 0, Hash: ""},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMeta, err := GetMetadata(tt.args.name, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMeta, tt.wantMeta) {
				t.Errorf("GetMetadata() gotMeta = %v, want %v", gotMeta, tt.wantMeta)
			}
		})
	}
}

func TestDelMetadata(t *testing.T) {
	type args struct {
		name    string
		version int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DelMetadata(tt.args.name, tt.args.version)
		})
	}
}

func TestSearchVersionStatus(t *testing.T) {
	type args struct {
		minDocCount int
	}
	tests := []struct {
		name    string
		args    args
		want    []Bucket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchVersionStatus(tt.args.minDocCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchVersionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchVersionStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
