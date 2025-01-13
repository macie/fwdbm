package dsn

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		connectionURI string
		wantDriver    string
		wantString    string
		wantErr       bool
	}{
		{
			name:          "PostgreSQL undefined",
			connectionURI: "postgres://",
			wantErr:       true,
		},
		{
			name:          "MySQL undefined",
			connectionURI: "mysql://",
			wantErr:       true,
		},
		{
			name:          "SQLite undefined",
			connectionURI: "sqlite://",
			wantErr:       true,
		},
		{
			name:          "SQLite encrypted",
			connectionURI: "sqlite://root:secret@localhost/default.db",
			wantDriver:    "sqlite3",
			wantString:    "file:///default.db?_auth=&_auth_crypt=SHA256&_auth_pass=secret&_auth_salt=SSHA256&_auth_user=root&_cache_size=-20000&_fk=1&_journal=WAL&_sync=NORMAL&_timeout=5000&_txlock=immediate",
		},
		{
			name:          "SQLite encrypted invalid domain",
			connectionURI: "sqlite://root:secret@example.com/default.db",
			wantErr:       true,
		},
		{
			name:          "SQLite encrypted without domain",
			connectionURI: "sqlite://root:secret/default.db",
			wantErr:       true,
		},
		{
			name:          "SQLite no path",
			connectionURI: "sqlite:?_fk=0",
			wantErr:       true,
		},
		{
			name:          "SQLite overridden parameters",
			connectionURI: "sqlite:default.db?_fk=0&_timeout=1000",
			wantDriver:    "sqlite3",
			wantString:    "file:default.db?_cache_size=-20000&_fk=0&_journal=WAL&_sync=NORMAL&_timeout=1000&_txlock=immediate",
		},
		{
			name:          "SQLite duplicate parameters",
			connectionURI: "sqlite:/default.db?_fk=1&_timeout=1000&_timeout=2000",
			wantDriver:    "sqlite3",
			wantString:    "file:///default.db?_cache_size=-20000&_fk=1&_journal=WAL&_sync=NORMAL&_timeout=2000&_txlock=immediate",
		},
		{
			name:          "SQLite invalid parameter",
			connectionURI: "sqlite:/dir/default.db?_journal=OFF&_invalid=2000",
			wantDriver:    "sqlite3",
			wantString:    "file:///dir/default.db?_cache_size=-20000&_fk=1&_invalid=2000&_journal=OFF&_sync=NORMAL&_timeout=5000&_txlock=immediate",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.connectionURI)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDB() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !strings.EqualFold(got.Driver(), tt.wantDriver) {
				t.Errorf("parseDB().Driver() got = %v, want %v", got.Driver(), tt.wantDriver)
			}
			if !strings.EqualFold(got.String(), tt.wantString) {
				t.Errorf("parseDB().String() got = %v, want %v", got.String(), tt.wantString)
			}
		})
	}
}
