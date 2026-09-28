package terminal

import (
	"slices"
	"testing"

	"github.com/go-delve/delve/pkg/config"
)

func TestConfigureSetAliasDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		aliases map[string][]string
		remove  string
		want    map[string][]string
	}{
		{
			name:    "unique alias",
			aliases: map[string][]string{"print": {"a", "b", "c"}},
			remove:  "a",
			want:    map[string][]string{"print": {"b", "c"}},
		},
		{
			name:    "alias duplicated back to back",
			aliases: map[string][]string{"print": {"a", "a"}},
			remove:  "a",
			want:    map[string][]string{"print": {}},
		},
		{
			name:    "alias duplicated with another alias in between",
			aliases: map[string][]string{"print": {"a", "b", "a"}},
			remove:  "a",
			want:    map[string][]string{"print": {"b"}},
		},
		{
			name:    "alias duplicated with other aliases in between",
			aliases: map[string][]string{"print": {"a", "b", "a", "c", "a"}},
			remove:  "a",
			want:    map[string][]string{"print": {"b", "c"}},
		},
		{
			name: "alias of more than one command",
			aliases: map[string][]string{
				"print": {"a", "b"},
				"locals": {"a", "c"},
			},
			remove: "a",
			want: map[string][]string{
				"print":  {"b"},
				"locals": {"c"},
			},
		},
		{
			name:    "only alias",
			aliases: map[string][]string{"print": {"a"}},
			remove:  "a",
			want:    map[string][]string{"print": {}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			term := &Term{conf: &config.Config{Aliases: tc.aliases}, cmds: DebugCommands(nil)}
			term.cmds.Merge(term.conf.Aliases)

			if err := configureSet(term, `alias "`+tc.remove+`"`); err != nil {
				t.Fatalf("configureSet: %v", err)
			}

			for cmd, want := range tc.want {
				got := term.conf.Aliases[cmd]
				if !slices.Equal(got, want) {
					t.Errorf("aliases of %q: got %q, want %q", cmd, got, want)
				}
			}

			// the deleted alias must no longer resolve to a command
			if cmd := term.cmds.Find(tc.remove, noPrefix); cmd.aliases[0] != "nocmd" {
				t.Errorf("alias %q still resolves to %q after being removed", tc.remove, cmd.aliases[0])
			}
		})
	}
}
