# Katkıda Bulunma Rehberi (Contribution Note)

Bu projeyi forklayıp katkıda bulunmak isterseniz, lütfen aşağıdaki adımları ve kuralları takip edin:

1. **Fork ve Clone:** Projeyi kendi GitHub hesabınıza forklayın ve bilgisayarınıza kopyalayın (`git clone`).
2. **Bağımlılıklar:** Proje Go 1.25.0 gerektirir. `go mod download` ile bağımlılıkları yükleyebilirsiniz.
3. **Geliştirme:** Değişikliklerinizi yeni bir branch üzerinde yapın (`git checkout -b feature/yeniozellik`).
4. **Test:** Yaptığınız değişikliklerin mevcut yapıyı bozmadığından emin olmak için `make test` komutunu çalıştırın.
5. **GitHub Actions:**
    * Her push işleminde GitHub Actions otomatik olarak projenizi derleyecek ve `/health` kontrolü ile çalışabilirliğini test edecektir.
    * **Release Oluşturma:** Eğer yaptığınız değişikliklerin otomatik olarak bir paket (release) olarak yayınlanmasını istiyorsanız, commit mesajınızın içerisine `<release:true>` etiketini ekleyin. Aksi takdirde sadece derleme ve test işlemleri yapılacaktır.
6. **Pull Request:** Değişikliklerinizi tamamladığınızda, ana repoya bir Pull Request (PR) gönderin.

Not: Lütfen commit mesajlarınızın açıklayıcı olmasına özen gösterin. İyi çalışmalar!
