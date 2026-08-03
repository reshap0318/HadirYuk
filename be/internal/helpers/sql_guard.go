package helpers

import (
	"fmt"
	"regexp"
	"strings"
)

// blacklistedTables are fully off-limits to the AI chat (§6.2 of the FSD) — no
// business relevance for a CS chatbot, or contain sensitive tokens.
var blacklistedTables = map[string]bool{
	"password_resets":      true,
	"role_has_permissions": true,
	"user_has_roles":       true,
	"roles":                true,
	"permissions":          true,
}

// usersSensitiveColumns are the users columns the AI chat must never see
// (only id/name are allowed). Checked as whole words rather than proper
// column parsing.
//
// ponytail: regex heuristic, not a real SQL parser — good enough as a second
// layer behind the DB-level column GRANT (§6.1), which is the real boundary.
// Upgrade to a proper SQL parser (e.g. vitess/sqlparser) if this ever needs
// to guard more than one sensitive table.
var usersSensitiveColumns = regexp.MustCompile(`(?i)\b(password|email|avatar)\b`)

var (
	selectOnlyRegex    = regexp.MustCompile(`(?i)^select\b`)
	forbiddenKeywords  = regexp.MustCompile(`(?i)\b(insert|update|delete|drop|alter|truncate|create|replace|grant|revoke|exec|execute|call)\b`)
	tableRefRegex      = regexp.MustCompile("(?i)\\b(?:from|join)\\s+`?([a-zA-Z_][a-zA-Z0-9_]*)`?")
	limitRegex         = regexp.MustCompile(`(?i)\blimit\s+\d+`)
	usersTableRefRegex = regexp.MustCompile("(?i)\\busers\\b")
	starFromUsersRegex = regexp.MustCompile("(?i)select\\s+\\*\\s+from\\s+`?users`?")
)

// ValidateAndPrepareSQL enforces SELECT-only, table/column blacklist (§6.2),
// and injects LIMIT maxRows if the query doesn't already have one. Returns
// the (possibly rewritten) query, or an error if the query fails validation.
func ValidateAndPrepareSQL(sql string, maxRows int) (string, error) {
	trimmed := strings.TrimSpace(sql)
	trimmed = strings.TrimSuffix(strings.TrimSpace(trimmed), ";")

	if strings.Contains(trimmed, ";") {
		return "", fmt.Errorf("multiple statements are not allowed")
	}
	if !selectOnlyRegex.MatchString(trimmed) {
		return "", fmt.Errorf("only SELECT queries are allowed")
	}
	if forbiddenKeywords.MatchString(trimmed) {
		return "", fmt.Errorf("query contains a disallowed keyword")
	}

	for _, m := range tableRefRegex.FindAllStringSubmatch(trimmed, -1) {
		if blacklistedTables[strings.ToLower(m[1])] {
			return "", fmt.Errorf("access to table %q is not allowed", m[1])
		}
	}

	if usersTableRefRegex.MatchString(trimmed) {
		if starFromUsersRegex.MatchString(trimmed) {
			return "", fmt.Errorf("select * from users is not allowed, only id/name")
		}
		if usersSensitiveColumns.MatchString(trimmed) {
			return "", fmt.Errorf("query references a restricted users column")
		}
	}

	if !limitRegex.MatchString(trimmed) {
		trimmed = fmt.Sprintf("%s LIMIT %d", trimmed, maxRows)
	}

	return trimmed, nil
}
