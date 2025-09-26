package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	repoOwner      = "TFAURE56"
	repoName       = "poker-app"
	currentVersion = "v0.0.1"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	a := app.New()
	w := a.NewWindow("PokerApp Launcher")

	status := widget.NewLabel("Vérification des mises à jour…")
	progress := widget.NewProgressBar()

	w.SetContent(container.NewVBox(
		status,
		progress,
	))
	w.Resize(fyne.NewSize(400, 120))
	w.Show()

	go func() {
		latest, err := checkLatestRelease()
		if err != nil {
			status.SetText(fmt.Sprintf("Erreur MAJ: %v\nLancement version locale...", err))
			launchApp()
			a.Quit()
			return
		}

		if latest.TagName != currentVersion {
			status.SetText(fmt.Sprintf("Nouvelle version trouvée (%s → %s). Téléchargement…", currentVersion, latest.TagName))
			err := downloadAndReplace(latest, progress)
			if err != nil {
				status.SetText(fmt.Sprintf("Erreur MAJ: %v\nLancement version locale...", err))
				launchApp()
				a.Quit()
				return
			}
			status.SetText("Mise à jour terminée. Lancement…")
		} else {
			status.SetText("Aucune mise à jour trouvée. Lancement…")
		}

		launchApp()
		a.Quit()
	}()

	a.Run()
}

func checkLatestRelease() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func downloadAndReplace(release *Release, progress *widget.ProgressBar) error {
	if len(release.Assets) == 0 {
		return fmt.Errorf("aucun binaire dans la release")
	}

	assetURL := release.Assets[0].BrowserDownloadURL
	resp, err := http.Get(assetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	exePath, _ := os.Executable()
	tmpPath := exePath + ".tmp"

	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer out.Close()

	total := resp.ContentLength
	progress.SetValue(0)

	buf := make([]byte, 32*1024)
	var downloaded int64
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := out.Write(buf[:n]); wErr != nil {
				return wErr
			}
			downloaded += int64(n)
			progress.SetValue(float64(downloaded) / float64(total))
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	// Remplacer l’exe
	return os.Rename(tmpPath, exePath)
}

func launchApp() {
	exePath := "./PokerApp.exe" // <-- ton application poker compilée
	cmd := exec.Command(exePath)
	cmd.Start()
}
