# Go Auth Sistemi (Argon2id)

Python/Flask ile yazdığım eski auth projesini Go'ya çevirdim. Daha performanslı ve güvenli oldu. Şifreleme tarafında Argon2id kullanıyor.

### Neler Yaptık?
- **Argon2id Hashing:** OWASP'ın önerdiği ayarlarla şifreleri saklıyoruz. Kırılması çok daha zor.
- **Signed Cookies:** HMAC kullanarak cookie'leri imzalıyoruz. Kullanıcı tarafında değiştirilirse sistem hemen anlıyor.
- **Güvenlik Header'ları:** XSS ve Clickjacking önlemleri için temel header'lar ekli.
- **SQLite:** Kurulumla uğraşmamak için veritabanı olarak SQLite kullandım.

### Nasıl Çalıştırılır?

1. **Bağımlılıkları çek:**
   ```bash
   go mod tidy
   ```

2. **Server'ı ayağa kaldır:**
   ```bash
   go run main.go
   ```
   (Ya da `go build -o auth-server` yapıp binary'yi de çalıştırabilirsin.)

3. **Sayfaya git:**
   `http://localhost:8080/register` adresinden kayıt olup giriş yapabilirsin.

### CLI Tarafı (auth.sh)
API üzerinden işlem yapmak istersen `auth.sh` scriptini kullanabilirsin.

```bash
chmod +x auth.sh
./auth.sh register testuser testpass123
./auth.sh login testuser testpass123
```

### Bazı Notlar
- Production'da `SECRET_KEY` env variable olarak mutlaka verilmeli.
- Loglar terminalde rahat okunması için standart formatta basılıyor (json değil).
- HTTPS (TLS) desteği yerelde yok, üretim ortamında Nginx/Traefik arkasına koymak lazım.

---
*Geliştirici Notu: Proje eğitim amaçlıdır, ufak tefek geliştirmeler yapılabilir (rate limit vs).*
