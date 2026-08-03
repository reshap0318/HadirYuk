package helpers

import "testing"

func TestValidateAndPrepareSQL(t *testing.T) {
	cases := []struct {
		name    string
		sql     string
		wantErr bool
		wantSQL string
	}{
		{"valid select passes", "SELECT id, name FROM users WHERE id = 1 LIMIT 10", false, "SELECT id, name FROM users WHERE id = 1 LIMIT 10"},
		{"dml rejected", "UPDATE users SET name = 'x'", true, ""},
		{"drop rejected", "DROP TABLE attendances", true, ""},
		{"multiple statements rejected", "SELECT 1; DROP TABLE attendances", true, ""},
		{"sensitive users column rejected", "SELECT password FROM users", true, ""},
		{"select star from users rejected", "SELECT * FROM users", true, ""},
		{"allowed users columns pass", "SELECT id, name FROM users", false, "SELECT id, name FROM users LIMIT 100"},
		{"blacklisted table rejected", "SELECT * FROM password_resets", true, ""},
		{"missing limit gets injected", "SELECT * FROM attendances", false, "SELECT * FROM attendances LIMIT 100"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ValidateAndPrepareSQL(tc.sql, 100)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got none (result: %q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantSQL {
				t.Fatalf("expected %q, got %q", tc.wantSQL, got)
			}
		})
	}
}
