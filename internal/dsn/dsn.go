// Package dsn provides datasourse name handling..
package dsn

import (
	"fmt"
	"net/url"
	"strings"
)

var noDSN = &DSN{}

// DSN is a database connection config.
type DSN struct {
	driver   string
	cleanDSN url.URL
}

// Parse returns DSN based on connection URI.
func Parse(uri string) (*DSN, error) {
	var connection *DSN

	parsed, err := url.Parse(uri)
	if err != nil {
		return noDSN, err
	}

	switch parsed.Scheme {
	case "sqlite":
		query := url.Values{
			// modern defaults
			"_cache_size": {"-20000"},    // 20MB (default: 2MB)
			"_fk":         {"1"},         // enable foreign key constraints (default: 0)
			"_journal":    {"WAL"},       // enable journal mode (default: DELETE)
			"_sync":       {"NORMAL"},    // less synchisation events
			"_timeout":    {"5000"},      // prevent SQLITE_BUSY exceptions
			"_txlock":     {"immediate"}, // prevent SQLITE_BUSY exceptions (default: deferred)
		}

		if strings.TrimSpace(parsed.Host) != "" && parsed.Host != "localhost" {
			return noDSN, fmt.Errorf("only localhost is supported")
		}

		if strings.TrimSpace(parsed.Opaque) == "" && strings.TrimSpace(parsed.Path) == "" {
			return noDSN, fmt.Errorf("missing path")
		}

		if parsed.User != nil {
			query.Add("_auth", "")
			query.Add("_auth_user", parsed.User.Username())
			password, ok := parsed.User.Password()
			if ok {
				query.Add("_auth_pass", password)
			}
			query.Add("_auth_crypt", "SHA256")
			query.Add("_auth_salt", "SSHA256")
		}

		parsedQuery, err := url.ParseQuery(parsed.RawQuery)
		if err != nil {
			return noDSN, err
		}
		for key, values := range parsedQuery {
			for _, value := range values {
				query.Set(key, value)
			}
		}

		cleanDSN := url.URL{
			Scheme:   "file",
			Opaque:   parsed.Opaque,
			Path:     parsed.Path,
			RawPath:  parsed.RawPath,
			Fragment: parsed.Fragment,
		}
		if strings.TrimSpace(cleanDSN.Opaque) != "" || strings.TrimSpace(cleanDSN.Path) != "" {
			cleanDSN.RawQuery = query.Encode()
		}

		connection = &DSN{
			driver:   "sqlite3",
			cleanDSN: cleanDSN,
		}

	default:
		return noDSN, fmt.Errorf("unsupported database scheme: %s", parsed.Scheme)
	}

	return connection, nil
}

// Driver returns database driver name.
func (dsn DSN) Driver() string {
	return dsn.driver
}

// String returns connection URI specific to database driver.
func (dsn DSN) String() string {
	return dsn.cleanDSN.String()
}
