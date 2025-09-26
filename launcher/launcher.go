package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	repoOwner  = "TFAURE56"
	repoName   = "poker-app"
	appExeName = "PokerApp.exe"
)

// Structure GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Chemin de l'application
func getInstallPath() string {
	local := os.Getenv("LOCALAPPDATA")
	return filepath.Join(local, "PokerApp", appExeName)
}

// Vérifie si l'application existe
func appExists() bool {
	_, err := os.Stat(getInstallPath())
	return err == nil
}

// Récupère la dernière release GitHub
func getLatestRelease() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Télécharge et installe l'application
func downloadAndInstall(release *Release, status *widget.Label) error {
	if len(release.Assets) == 0 {
		return fmt.Errorf("aucun binaire disponible")
	}

	assetURL := release.Assets[0].BrowserDownloadURL
	resp, err := http.Get(assetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	installPath := getInstallPath()
	if err := os.MkdirAll(filepath.Dir(installPath), 0755); err != nil {
		return err
	}

	out, err := os.Create(installPath)
	if err != nil {
		return err
	}
	defer out.Close()

	status.SetText("Téléchargement en cours…")
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	status.SetText("Installation terminée !")
	return nil
}

func main() {
	a := app.NewWithID("com.pokerapp.launcher")
	w := a.NewWindow("PokerApp Launcher")

	status := widget.NewLabel("Vérification de l'application…")
	w.SetContent(container.NewVBox(status))
	w.Resize(fyne.NewSize(400, 120))
	w.Show()

	go func() {
		if !appExists() {
			status.SetText("Application non installée")

			latest, err := getLatestRelease()
			if err != nil {
				status.SetText(fmt.Sprintf("Erreur lors de la vérification de la MAJ: %v", err))
				return
			}

			dialog.ShowConfirm("Installer PokerApp",
				fmt.Sprintf("Voulez-vous installer la version %s ?", latest.TagName),
				func(ok bool) {
					if !ok {
						status.SetText("Installation annulée.")
						return
					}

					go func() {
						if err := downloadAndInstall(latest, status); err != nil {
							status.SetText(fmt.Sprintf("Erreur: %v", err))
							return
						}
					}()
				}, w)

		} else {
			status.SetText("Application déjà installée")
		}
	}()

	a.Run()
}
