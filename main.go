package main

import (
    "fmt"
    "net/http"
    "html/template"
)

// Biodata struct untuk data statis
type Biodata struct {
    Nama    string
    Email   string
    Usia    int
    Alamat  string
}

// Data Biodata statis
var biodataStatis = map[string]Biodata{
    "rifqy@example.com": Biodata{"Rifqy", "rifqy@example.com", 20, "Jakarta"},
    "ryan@example.com": Biodata{"Ryan", "ryan@example.com", 25, "Surabaya"},
    "wildan@example.com": Biodata{"Wildan", "wildan@example.com", 25, "Bandung"},
    "andri@example.com": Biodata{"Andri", "andri@example.com", 25, "Bogor"},
}

func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.Handle("/error.html", http.FileServer(http.Dir("./")))

    fmt.Println("Server is running on :8080")
    http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    // Ambil email dari URL query
    email := r.URL.Query().Get("email")

    // Ambil data biodata berdasarkan email
    biodata, found := biodataStatis[email]
    if !found {
        http.Error(w, "Biodata tidak ditemukan", http.StatusNotFound)
        return
    }

    // Parse template HTML
    tmpl, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Kirim data biodata ke template
    err = tmpl.Execute(w, biodata)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func logoutHandler(w http.ResponseWriter, r *http.Request) {
    // Redirect pengguna ke halaman login
    http.Redirect(w, r, "/login", http.StatusFound)
}




func loginHandler(w http.ResponseWriter, r *http.Request) {
    // Mengambil daftar email dari data biodata statis
    var emailList []string
    for email := range biodataStatis {
        emailList = append(emailList, email)
    }

    data := struct {
        BiodataList []string
        Error       bool // Untuk menampilkan pesan kesalahan
    }{
        emailList,
        false,
    }

    if r.Method == http.MethodPost {
        // Ambil email dari form POST
        email := r.FormValue("email")

        // Cek apakah email ada dalam data biodata statis
        _, found := biodataStatis[email]
        if !found {
            // Jika email tidak terdaftar, arahkan ke halaman error.html
            http.Redirect(w, r, "/error.html", http.StatusFound)
            return
        } else {
            // Jika email terdaftar, redirect ke halaman biodata
            http.Redirect(w, r, "/?email="+email, http.StatusFound)
            return
        }
    }

    // Parse template HTML
    tmpl, err := template.ParseFiles("login.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Kirim data ke template
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}