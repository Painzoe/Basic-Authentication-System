package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log" // slog yerine düz log, daha hızlı ve okunaklı
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)

// Ayarlar - env'den çekiyoruz, yoksa default'a düşüyor
var (
	db_yolu    = getEnv("DB_PATH", "./auth.db")
	port       = getEnv("PORT", "8080")
	secret_key = []byte(getEnv("SECRET_KEY", "bunu-gercek-sistemde-degistirin-123"))
)

// argon2id ayarları. owasp değerlerine yakın tuttum.
type argonParams struct {
	m uint32
	t uint32
	p uint8
	l uint32
}

var params = argonParams{
	m: 64 * 1024,
	t: 3,
	p: 2,
	l: 32,
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", db_yolu)
	if err != nil {
		log.Fatalf("Veritabanına bağlanamadık: %v", err)
	}
	defer db.Close()

	// Tabloyu oluştur (eğer yoksa)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Tablo oluşurken patladık: %v", err)
	}

	mux := http.NewServeMux()

	// Ana rotalar
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up"))
	})
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/dashboard", auth_korumasi(dashboardHandler))
	mux.HandleFunc("/logout", logoutHandler)

	// API tarafı (CLI scripti için)
	mux.HandleFunc("/api/login", apiLoginHandler)

	// Middleware'leri giydirelim
	handler := log_middleware(guvenlik_headers(mux))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Sinyal yakalama (ctrl+c vs kapatırken db patlamasın)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server %s portunda başladı...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server kapandı: %v", err)
		}
	}()

	<-stop
	log.Println("Kapatılıyor...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
}

// Gelen istekleri terminalde görmek için basit bir log
func log_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// XSS ve Clickjacking önlemi
func guvenlik_headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline';")
		next.ServeHTTP(w, r)
	})
}

// Session kontrolü. Cookie yoksa logine atar.
func auth_korumasi(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}

		user, err := token_dogrula(c.Value)
		if err != nil {
			log.Printf("Hatalı session denemesi: %v", err)
			http.Redirect(w, r, "/login", 302)
			return
		}

		r.Header.Set("User", user)
		next.ServeHTTP(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render_tmpl(w, "login", nil)
		return
	}

	u := r.FormValue("username")
	p := r.FormValue("password")

	var hash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", u).Scan(&hash)
	if err != nil {
		render_tmpl(w, "login", "Kullanıcı bulunamadı veya şifre yanlış")
		return
	}

	if !sifre_kontrol(p, hash) {
		render_tmpl(w, "login", "Kullanıcı bulunamadı veya şifre yanlış")
		return
	}

	// Token yapıp cookie'ye gömüyoruz
	token := token_yap(u)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600, // 1 saatlik oturum
	})

	http.Redirect(w, r, "/dashboard", 302)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render_tmpl(w, "register", nil)
		return
	}

	u := r.FormValue("username")
	p := r.FormValue("password")

	if len(p) < 8 {
		render_tmpl(w, "register", "Şifre en az 8 karakter olmalı")
		return
	}

	hash := sifre_hashle(p)
	_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", u, hash)
	if err != nil {
		render_tmpl(w, "register", "Bu kullanıcı adı zaten alınmış")
		return
	}

	log.Printf("Yeni kullanıcı kaydoldu: %s", u)
	http.Redirect(w, r, "/login", 302)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	render_tmpl(w, "dashboard", r.Header.Get("User"))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/login", 302)
}

func apiLoginHandler(w http.ResponseWriter, r *http.Request) {
	// CLI login işleri buraya... (TODO: json response dönülecek)
	w.WriteHeader(501)
	w.Write([]byte("API login henüz hazır değil"))
}

// --- Yardımcı Fonksiyonlar ---

func sifre_hashle(p string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	h := argon2.IDKey([]byte(p), salt, params.t, params.m, params.p, params.l)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, params.m, params.t, params.p,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(h))
}

func sifre_kontrol(p, h string) bool {
	pts := strings.Split(h, "$")
	if len(pts) != 6 {
		return false
	}
	var m, t uint32
	var p_val uint8
	fmt.Sscanf(pts[3], "m=%d,t=%d,p=%d", &m, &t, &p_val)
	salt, _ := base64.RawStdEncoding.DecodeString(pts[4])
	originalHash, _ := base64.RawStdEncoding.DecodeString(pts[5])
	testHash := argon2.IDKey([]byte(p), salt, t, m, p_val, uint32(len(originalHash)))
	return subtle.ConstantTimeCompare(originalHash, testHash) == 1
}

func token_yap(u string) string {
	h := hmac.New(sha256.New, secret_key)
	h.Write([]byte(u))
	sig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(u)) + "." + sig
}

func token_dogrula(t string) (string, error) {
	pts := strings.Split(t, ".")
	if len(pts) != 2 {
		return "", fmt.Errorf("token bozuk")
	}
	u_bytes, _ := base64.RawURLEncoding.DecodeString(pts[0])
	u := string(u_bytes)
	if token_yap(u) != t {
		return "", fmt.Errorf("imza tutmuyor")
	}
	return u, nil
}

func render_tmpl(w http.ResponseWriter, name string, data interface{}) {
	t, err := template.ParseFiles("templates/" + name + ".html")
	if err != nil {
		log.Printf("Template hatası: %v", err)
		http.Error(w, "Görünüm yüklenemedi", 500)
		return
	}
	t.Execute(w, data)
}

func getEnv(k, f string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return f
}
