package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
)

// Estrutura para representar um capítulo da Bíblia
type Chapter struct {
    Verses []string `json:"verses"`
}

// Estrutura para representar um livro da Bíblia
type Book struct {
    Abbrev    string    `json:"abbrev"`
    Chapters  [][]string `json:"chapters"` // Atualizado para refletir a estrutura do JSON
    Name      string    `json:"name"`
}

var books []Book

// Carrega os livros da Bíblia a partir do arquivo JSON
func LoadBooks(filename string) error {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }

    err = json.Unmarshal(data, &books)
    if err != nil {
        return err
    }

    return nil
}

// Handler para retornar todos os livros
func GetBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // Extrair apenas abreviaturas dos livros
    names := make([]string, len(books))
    for i, book := range books {
        names[i] = book.Abbrev
    }
    json.NewEncoder(w).Encode(names)
}

// Handler para retornar capítulos de um livro específico
func GetChapters(w http.ResponseWriter, r *http.Request) {
    abbrev := mux.Vars(r)["abbrev"]
    for _, book := range books {
        if book.Abbrev == abbrev {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(book.Chapters)
            return
        }
    }
    http.NotFound(w, r)
}

// Handler para retornar um capítulo específico de um livro
func GetChapter(w http.ResponseWriter, r *http.Request) {
    abbrev := mux.Vars(r)["abbrev"]
    chapterIndex, err := strconv.Atoi(mux.Vars(r)["chapter"])
    if err != nil {
        http.Error(w, "Invalid chapter index", http.StatusBadRequest)
        return
    }
    for _, book := range books {
        if book.Abbrev == abbrev {
            if chapterIndex >= 0 && chapterIndex < len(book.Chapters) {
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(book.Chapters[chapterIndex])
                return
            }
            http.NotFound(w, r)
            return
        }
    }
    http.NotFound(w, r)
}

// Handler para retornar um versículo específico de um capítulo
func GetVerse(w http.ResponseWriter, r *http.Request) {
    abbrev := mux.Vars(r)["abbrev"]
    chapterIndex, err := strconv.Atoi(mux.Vars(r)["chapter"])
    if err != nil {
        http.Error(w, "Invalid chapter index", http.StatusBadRequest)
        return
    }
    verseIndex, err := strconv.Atoi(mux.Vars(r)["verse"])
    if err != nil {
        http.Error(w, "Invalid verse index", http.StatusBadRequest)
        return
    }
    for _, book := range books {
        if book.Abbrev == abbrev {
            if chapterIndex >= 0 && chapterIndex < len(book.Chapters) {
                if verseIndex >= 0 && verseIndex < len(book.Chapters[chapterIndex]) {
                    w.Header().Set("Content-Type", "application/json")
                    json.NewEncoder(w).Encode(book.Chapters[chapterIndex][verseIndex])
                    return
                }
                http.NotFound(w, r)
                return
            }
            http.NotFound(w, r)
            return
        }
    }
    http.NotFound(w, r)
}

// Função principal para inicializar o servidor
func main() {
    // Carrega os livros do arquivo JSON
    err := LoadBooks("acf.json")
    if err != nil {
        log.Fatalf("Error loading books: %v", err)
    }

    // Configura o roteador
    r := mux.NewRouter()
    r.HandleFunc("/books", GetBooks).Methods("GET")
    r.HandleFunc("/books/{abbrev}", GetChapters).Methods("GET")
    r.HandleFunc("/books/{abbrev}/chapters/{chapter:[0-9]+}", GetChapter).Methods("GET")
    r.HandleFunc("/books/{abbrev}/chapters/{chapter:[0-9]+}/verses/{verse:[0-9]+}", GetVerse).Methods("GET")

    // Inicia o servidor
    log.Fatal(http.ListenAndServe(":8080", r))
}

//Lista todos os Liros: http://localhost:8080/books
// Lista todos os Capitulo do livro: http://localhost:8080/books/gn/chapters
// Lista de todos Versiculos: http://localhost:8080/books/gn/chapters/0/verses
// Pega um versiculo: curl http://localhost:8080/books/gn/chapters/0/verses/0

