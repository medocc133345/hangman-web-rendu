package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

type GameState struct {
    Pseudo     string
    Word       string
    Attempts   []string
    Lives      int
    Message    string
}


func lireMotsDepuisFichier(fichier string) ([]string, error) {
    var mots []string

    file, err := os.Open(fichier)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        mot := strings.TrimSpace(scanner.Text())
        if mot != "" {
            mots = append(mots, mot)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return mots, nil
}

func choisirMotAleatoire(mots []string) string {
    rand.Seed(time.Now().UnixNano()) 
    index := rand.Intn(len(mots))     
    return mots[index]                
}

func homePage(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    tmpl.Execute(w, nil)
}

func startGame(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        pseudo := r.FormValue("pseudo")
        difficulty := r.FormValue("difficulty")

        var fichierMots string
        switch difficulty {
        case "facile":
            fichierMots = "mots_faciles.txt"
        case "moyen":
            fichierMots = "mots_moyens.txt"
        case "intermediaire":
            fichierMots = "mots_intermediaires.txt"
        case "difficile":
            fichierMots = "mots_difficiles.txt"
        default:
            fmt.Println("Niveau de difficulté non valide.")
            return
        }

        mots, err := lireMotsDepuisFichier(fichierMots)
        if err != nil {
            fmt.Println("Erreur lors de la lecture du fichier :", err)
            return
        }

        motADeviner := choisirMotAleatoire(mots)
        fmt.Printf("Bienvenue %s, le mot à deviner est : %s\n", pseudo, motADeviner)

        http.Redirect(w, r, "/game", http.StatusSeeOther)
    }
}

func gamePage(w http.ResponseWriter, r *http.Request) {
    gameState := GameState{
        Pseudo:   "PseudoJoueur",          
        Word:     "_ _ _ _",               
        Attempts: []string{"A", "E"},      
        Lives:    3,                       
        Message:  "Bonne chance !",
    }

    tmpl := template.Must(template.ParseFiles("templates/game.html"))
    tmpl.Execute(w, gameState)
}

func main() {

    http.HandleFunc("/", homePage)

    http.HandleFunc("/start", startGame)

    http.HandleFunc("/game", gamePage)

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
