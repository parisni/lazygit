package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleFollowRenames = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by the path of a renamed file without following renames, then turn following on from the filtering menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.FollowRenamesWhenFiltering = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("oldFile", "a\nb\nc\n")
		shell.Commit("add old file")
		shell.UpdateFileAndAdd("oldFile", "x\nb\nc\n")
		shell.Commit("update old file")
		shell.RenameFileInGit("oldFile", "newFile")
		shell.Commit("rename file")
		shell.UpdateFileAndAdd("newFile", "y\nb\nc\n")
		shell.Commit("update new file")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter path to filter by")).
			Confirm()
		t.ExpectPopup().Prompt().
			Title(Equals("Enter path:")).
			Type("newFile").
			Confirm()

		t.Views().Information().Content(Equals("Filtering by 'newFile' (Reset)"))

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("update new file").IsSelected(),
				Contains("rename file"),
			).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("[ ] Follow renames")).
			Confirm()

		t.Views().Information().Content(Equals("Filtering by 'newFile' (following renames) (Reset)"))

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("update new file").IsSelected(),
				Contains("rename file"),
				Contains("update old file"),
				Contains("add old file"),
			).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("[✓] Follow renames")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("update new file").IsSelected(),
				Contains("rename file"),
			)
	},
})
