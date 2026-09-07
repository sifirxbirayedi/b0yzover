


b0yzover Nedir ?

b0yzover, bulut servisleri uzerindeki dangling CNAME ve subdomain takeover (alt alan adi ele gecirme) acilarini hizli ve etkili bir sekilde tespit etmek icin gelistirilmis gelismis bir aracatir. Hedef sistemlerin DNS yapisini ve HTTP yanitlarini analiz ederek guvenlik aciklarini raporlar.

Ne Yapabilir ?
-Hedef domain icin crt.name uzerinden aktif alt alan adlarini toplar ve listeler.
-Toplanan alt alan adlari uzerinde eszamanli (multi-threaded) olarak HTTP istekleri gerceklestirir.
-Bulut servisleri saglayicilarina ait bilinen parmak izi (fingerprint) veritabani ile yanitlari karsilastirir.
-Dangling CNAME (bosta kalan CNAME) kayitlarini tespit ederek potansiyel takeover açiklarini yakalar.
-Wildcard DNS kayitlarini tespit ederek yanlis pozitif (false-positive) sonuçlari filtreler.
-Sonuclari ve açik detaylarini terminale duzenli bir formatta basar.
-Taranan sonuclari ve ayrintili raporlari JSON formatinda dosyalayarak kaydeder.

Kullanim Sekli

Kod:
Tekil hedef domain taramasi icin:
./b0yzover -d hedefdomain.com

Hedefleri iceren liste dosyasi ile tarama icin:
./b0yzover -l hedefler.txt -o sonuc.json

Parametreler

-d, --target: Tekil hedef domain (crt.name kesfi icin)
-l, --list: Taranacak hedefleri iceren TXT dosyasi
-c, --threads: Eszamanli thread/cekirdek sayisi (Varsayilan: 30)
-t, --timeout: HTTP istek zaman asimi saniyesi
-o, --output: Sonuclarin kaydedilecegi JSON dosya adi
-v, --verbose: Detayli log modu

Windows ve geri kalanlar için eğer go yoksa bu https://golang.org/dl/ linkten
indirebilirsiniz sonrası go run main.go eğer bir sorun yaşarsanız bana özelden yazdığınız veya
telegramdan elimden gelenin en iyi şekliyle yardımcı olmaya çalışırım

Yükleme , go ile yapildi basit mantık go run.

GitHub




Bu aracın sorumluluğu size aittir.
Subyz üzerinden geliştirilmiş ve geliştirilmeye devam edecektir.
Program Backend Sahibi : @b0yner 
