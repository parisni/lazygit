package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleShowTags = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hide the tags of the commits by config, then show them from the log menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.ShowTagsInCommits = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.CreateLightweightTag("some-tag", "HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").DoesNotContain("some-tag").IsSelected(),
				Contains("one"),
			)

		t.Views().Main().Content(DoesNotContain("some-tag"))

		t.Views().Commits().
			Press(keys.Commits.OpenLogMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Commit Log Options")).
			Select(Contains("[ ] Show tags")).
			Confirm()

		t.Views().Commits().
			Lines(
				MatchesRegexp(`some-tag.*two`).IsSelected(),
				Contains("one"),
			)

		t.Views().Main().Content(Contains("tag: some-tag"))

		t.Views().Commits().
			Press(keys.Commits.OpenLogMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Commit Log Options")).
			Select(Contains("[✓] Show tags")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("two").DoesNotContain("some-tag").IsSelected(),
				Contains("one"),
			)

		t.Views().Main().Content(DoesNotContain("some-tag"))
	},
})
