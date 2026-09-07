# b0yzover greetz @b0yner
[CENTER][COLOR=rgb(255, 255, 255)][URL='https://hizliresim.com/0g66pb80'][IMG width="594px"]https://i.hizliresim.com/0g66pb80.png[/IMG][/URL]
[URL='https://hizliresim.com/tl3iruqx'][IMG width="587px"]https://i.hizliresim.com/tl3iruqx.png[/IMG][/URL][/COLOR]

[SIZE=5][COLOR=rgb(184, 49, 47)]b0yzover Nedir ?[/COLOR][/SIZE]

[COLOR=rgb(41, 105, 176)]b0yzover[/COLOR][COLOR=rgb(255, 255, 255)], bulut servisleri uzerindeki dangling [/COLOR][COLOR=rgb(184, 49, 47)]CNAME ve subdomain takeover (alt alan adi ele gecirme) [/COLOR][COLOR=rgb(255, 255, 255)]acilarini hizli ve etkili bir sekilde tespit etmek icin gelistirilmis gelismis bir aracatir. Hedef sistemlerin DNS yapisini ve HTTP yanitlarini analiz ederek guvenlik aciklarini raporlar.[/COLOR]

[SIZE=5][COLOR=rgb(184, 49, 47)]Ne Yapabilir ? [/COLOR][/SIZE]
[COLOR=rgb(255, 255, 255)]-Hedef domain icin crt.name uzerinden aktif alt alan adlarini toplar ve listeler.
-Toplanan alt alan adlari uzerinde eszamanli (multi-threaded) olarak HTTP istekleri gerceklestirir.
-Bulut servisleri saglayicilarina ait bilinen parmak izi (fingerprint) veritabani ile yanitlari karsilastirir.
-Dangling CNAME (bosta kalan CNAME) kayitlarini tespit ederek potansiyel takeover açiklarini yakalar.
-Wildcard DNS kayitlarini tespit ederek yanlis pozitif (false-positive) sonuçlari filtreler.
-Sonuclari ve açik detaylarini terminale duzenli bir formatta basar.
-Taranan sonuclari ve ayrintili raporlari JSON formatinda dosyalayarak kaydeder.[/COLOR]

[COLOR=rgb(184, 49, 47)][SIZE=5]Kullanim Sekli[/SIZE][/COLOR]

[CODE]Tekil hedef domain taramasi icin:
./b0yzover -d hedefdomain.com

Hedefleri iceren liste dosyasi ile tarama icin:
./b0yzover -l hedefler.txt -o sonuc.json[/CODE]

[SIZE=5][COLOR=rgb(184, 49, 47)]Parametreler[/COLOR][/SIZE]

[COLOR=rgb(255, 255, 255)]-d, --target: Tekil hedef domain (crt.name kesfi icin)
-l, --list: Taranacak hedefleri iceren TXT dosyasi
-c, --threads: Eszamanli thread/cekirdek sayisi (Varsayilan: 30)
-t, --timeout: HTTP istek zaman asimi saniyesi
-o, --output: Sonuclarin kaydedilecegi JSON dosya adi
-v, --verbose: Detayli log modu[/COLOR]

[SIZE=5][COLOR=rgb(255, 255, 255)]Windows ve geri kalanlar için eğer go yoksa bu [URL]https://golang.org/dl/[/URL] linkten 
indirebilirsiniz sonrası go run main.go eğer bir sorun yaşarsanız bana özelden yazdığınız veya 
telegramdan elimden gelenin en iyi şekliyle yardımcı olmaya çalışırım[/COLOR]

[COLOR=rgb(124, 112, 107)]Yükleme , go ile yapildi basit mantık  go run.[/COLOR][/SIZE]
[CHARGE=10]
[SIZE=5][COLOR=rgb(124, 112, 107)][URL='https://github.com/sifirxbirayedi/b0yzover']GitHu[/URL][/COLOR][/SIZE][URL='https://github.com/sifirxbirayedi/b0yzover']b[/URL]
[/CHARGE]



[COLOR=rgb(184, 49, 47)]Bu aracın sorumluluğu size aittir.
Subyz üzerinden geliştirilmiş ve geliştirilmeye devam edecektir.[/COLOR]
[COLOR=rgb(255, 255, 255)]Program Backend Sahibi : [USER=706355]@b0yner[/USER] [/COLOR][/CENTER]
