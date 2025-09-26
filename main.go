package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Poker Club - Timer Blindes")

	// Exemple de fenêtre "Paramètres"
	openSettings := func() {
		settingsWindow := myApp.NewWindow("Paramètres")
		settingsWindow.SetContent(widget.NewLabel("Ici tu pourras configurer les blindes, la durée, etc."))

		// Parametre pour definir le temps des rounds

		// Parametre pour definir les blindes
		// Parametre pour definir les ante
		// Parametre pour definir les pauses
		// Parametre pour definir le son des alertes

		settingsWindow.Resize(fyne.NewSize(300, 200))
		settingsWindow.Show()
	}

	// Exemple de boîte de dialogue "À propos"
	openAbout := func() {
		dialog.ShowInformation("À propos", "Poker Timer v1.0\nFait avec Fyne", myWindow)
	}

	openUpdate := func() {
		dialog.ShowInformation("Mise à jour", "Vérification des mises à jour...\nVous utilisez la dernière version.", myWindow)
	}

	// Définir le menu
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("Fichier",
			fyne.NewMenuItem("Quitter", func() {
				myApp.Quit()
			}),
		),
		fyne.NewMenu("Paramètres",
			fyne.NewMenuItem("Ouvrir", func() {
				openSettings()
			}),
		),
		fyne.NewMenu("Aide",
			fyne.NewMenuItem("À propos", func() {
				openAbout()
			}),
			fyne.NewMenuItem("Update", func() {
				openUpdate()
			}),
		),
	)

	myWindow.SetMainMenu(mainMenu)

	// Contenu principal
	label := widget.NewLabel("Timer Poker ici")
	myWindow.SetContent(container.NewCenter(label))

	myWindow.Resize(fyne.NewSize(900, 500))
	myWindow.ShowAndRun()
}
