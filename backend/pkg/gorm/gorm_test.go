package gorm

import "testing"

func TestSplitSQLStatements(t *testing.T) {
	sql := `
-- create the first table
CREATE TABLE one (
	id BIGINT PRIMARY KEY,
	note VARCHAR(255) COMMENT 'keep; semicolon'
);

/* create the second table */
CREATE TABLE two (
	id BIGINT PRIMARY KEY,
	name VARCHAR(64) DEFAULT 'it''s ok'
);
`

	statements := splitSQLStatements(sql)
	if len(statements) != 2 {
		t.Fatalf("expected 2 statements, got %d: %#v", len(statements), statements)
	}

	if statements[0] == "" || statements[1] == "" {
		t.Fatalf("expected non-empty statements: %#v", statements)
	}
}
