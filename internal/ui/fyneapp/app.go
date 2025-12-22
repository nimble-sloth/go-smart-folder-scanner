package fyneapp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	dlg "github.com/sqweek/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

import "github.com/nimble-sloth/go-smart-folder-scanner/internal/scanner"

func RunApp() {
	a := app.New()
	// Keep it simple / “Windows standard” feel (default/light theme).
	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("Smart Folder Scanner")
	w.Resize(fyne.NewSize(720, 620))

	homeDir, _ := os.UserHomeDir()
	defaultFolder := `C:\Repo`
	defaultFile := filepath.Join(homeDir, "Downloads", "scanner.txt")

	header := widget.NewLabel("Smart Folder Scanner")

	// ------- Folder to scan (manual typing + native Windows dialog) -------
	folderEntry := widget.NewEntry()
	folderEntry.SetPlaceHolder(`e.g. C:\Repo`)
	if _, err := os.Stat(defaultFolder); err == nil {
		folderEntry.SetText(defaultFolder)
	}
	btnBrowseFolder := widget.NewButton("Browse…", func() {
		// Native Windows folder picker
		if path, err := dlg.Directory().Title("Select Folder to Scan").Browse(); err == nil && path != "" {
			folderEntry.SetText(path)
		}
	})
	folderRow := container.NewBorder(nil, nil, nil, btnBrowseFolder, folderEntry)

	// ------- Output file (manual typing + native Windows dialog) -------
	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder(`e.g. C:\Users\me\Downloads\scanner.txt`)
	outputEntry.SetText(defaultFile)
	btnSaveAs := widget.NewButton("Save As…", func() {
		startDir := filepath.Dir(outputEntry.Text)
		if startDir == "." {
			startDir = homeDir
		}
		// Native Windows save dialog
		filename, err := dlg.File().
			Title("Choose Output File").
			SetStartDir(startDir).
			Save()
		if err == nil && filename != "" {
			outputEntry.SetText(filename)
		}
	})
	outputRow := container.NewBorder(nil, nil, nil, btnSaveAs, outputEntry)

	// ------- Skip folders -------
	const defaultSkip = ".venv,__pycache__,backup,build,dev,dist,go,log,screenshots,.git,node_modules,dist,build,.idea,.vscode"
	skipEntry := widget.NewEntry()
	skipEntry.SetPlaceHolder(".git,node_modules,dist,build,.idea,.vscode")
	skipEntry.SetText(defaultSkip)

	// ------- File types to scan -------
	extDefaults := []string{
		".go", ".py", ".java", ".js", ".ts", ".tsx", ".jsx",
		".html", ".css", ".json", ".md",
		".yaml", ".yml", ".xml",
		".sh", ".bat", ".ps1", ".sql",
		".cs", ".cpp", ".c", ".h",
		".kt", ".rs", ".swift", ".php", ".rb", ".properties",
	}

	extChecks := make(map[string]*widget.Check)

	// create 5 columns
	cols := []fyne.CanvasObject{
		container.NewVBox(),
		container.NewVBox(),
		container.NewVBox(),
		container.NewVBox(),
		container.NewVBox(),
	}

	for i, ext := range extDefaults {
		c := widget.NewCheck(ext, nil)
		if slices.Contains([]string{".go", ".py", ".java"}, ext) {
			c.SetChecked(true)
		}
		extChecks[ext] = c
		// distribute evenly across 5 columns
		colIndex := i % 5
		cols[colIndex] = container.NewVBox(cols[colIndex], c)
	}

	customExt := widget.NewEntry()
	customExt.SetPlaceHolder(".toml,.gradle,README,LICENSE (comma-separated)")

	btnAll := widget.NewButton("Select All", func() {
		for _, c := range extChecks {
			c.SetChecked(true)
		}
	})
	btnNone := widget.NewButton("Select None", func() {
		for _, c := range extChecks {
			c.SetChecked(false)
		}
	})

	// grid with 5 columns
	extGrid := container.NewGridWithColumns(5, cols...)
	extControls := container.NewHBox(btnAll, btnNone)

	// ------- Status / progress -------
	status := widget.NewLabel("Ready.")
	bar := widget.NewProgressBarInfinite()
	bar.Hide()

	// ------- Include tree structure -------
	includeTree := widget.NewCheck("Include folder tree", nil)
	includeTree.SetChecked(true) // default ON

	// ------- Scan button -------
	var btnScan *widget.Button
	btnScan = widget.NewButton("Scan", func() {
		if strings.TrimSpace(folderEntry.Text) == "" {
			dialog.ShowInformation("Missing folder", "Please enter or choose a folder to scan.", w)
			return
		}
		if strings.TrimSpace(outputEntry.Text) == "" {
			dialog.ShowInformation("Missing output", "Please enter or choose where to save results.", w)
			return
		}

		var selected []string
		for ext, c := range extChecks {
			if c.Checked {
				selected = append(selected, ext)
			}
		}
		if s := strings.TrimSpace(customExt.Text); s != "" {
			selected = append(selected, splitAndTrim(s)...)
		}
		extList := strings.Join(selected, ",")

		bar.Show()
		status.SetText("Scanning…")
		btnScan.Disable()
		start := time.Now()

		err := scanAndSave(folderEntry.Text, outputEntry.Text, skipEntry.Text, extList, includeTree.Checked)

		bar.Hide()
		btnScan.Enable()
		if err != nil {
			status.SetText(fmt.Sprintf("Error: %v", err))
			dialog.ShowError(err, w)
			return
		}
		status.SetText(fmt.Sprintf("Done. Saved to %s (%.1fs)", outputEntry.Text, time.Since(start).Seconds()))
	})

	// ------- Layout (plain form-like) -------
	form := container.NewVBox(
		header,
		widget.NewSeparator(),

		widget.NewLabel("Folder to scan:"),
		folderRow,

		widget.NewLabel("Output file:"),
		outputRow,

		widget.NewLabel("Folders to skip (comma-separated):"),
		skipEntry,

		widget.NewLabel("File types to include:"),
		extGrid,
		widget.NewLabel("Custom extensions or filenames (comma-separated):"),
		customExt,
		extControls,

		widget.NewSeparator(),
		includeTree,
		container.NewHBox(btnScan),
		status,
		bar,
	)

	scroll := container.NewVScroll(form)
	scroll.SetMinSize(fyne.NewSize(700, 520))

	w.SetContent(scroll)
	w.ShowAndRun()
}

// ===== scanning =====

func scanAndSave(root, output, skipList, extList string, includeTree bool) error {
	skipFolders := splitAndTrim(skipList)
	exts := splitAndTrim(extList)

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()
	writer := bufio.NewWriter(out)

	// Write tree first if requested
	if includeTree {
		tree, err := scanner.WriteTree(root)
		if err == nil {
			fmt.Fprintln(writer, tree)
			fmt.Fprintln(writer) // blank line after tree
		}
	}

	// Dump file contents via scanner.ScanTree
	dump, err := scanner.ScanTree(root, skipFolders, exts)
	if err == nil {
		fmt.Fprintln(writer, dump)
	}

	_ = writer.Flush()
	return err
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
