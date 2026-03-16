# 🛡️ **GuardPot SSH**

GuardPot SSH, kullanıcıların gerçek sunuculara güvenli şekilde bağlanabildiği,
tarayıcı üzerinden canlı SSH terminali sunan bir full-stack SSH management platformudur.

- ❗ Simülasyon yok

- ❗ Fake output yok

- ✅ Girilen komutlar gerçek sunucuya iletilir

- ✅ Çıktılar gerçek zamanlı olarak tarayıcıda gösterilir

## 🚀 Temel Özellikler

**🔐 Kimlik Doğrulama**

- Email + Password ile kullanıcı kaydı

- JWT tabanlı stateless authentication

- Token expiry & middleware koruması

- (Frontend hazır) Google Sign-In desteği (backend finalize ediliyor)

**🔑 SSH Bağlantı Yönetimi**

- SSH bağlantısı ekleme

- Host

- Port

- Kullanıcı adı

- Password / Private Key

- SSH credential’lar backend’de AES-GCM ile şifrelenerek saklanır

- Bağlantılar kullanıcıya özeldir

**🖥️ Canlı Terminal (Gerçek SSH)**

- Tarayıcıda gerçek zamanlı terminal

- WebSocket üzerinden canlı bağlantı

- Komutlar birebir SSH server’a gider

- Çıktılar anlık olarak ekranda görünür

- Resize, reconnect ve copy/paste desteği

**🔎 Güvenlik & Audit**

- SSH host key scan & trust mekanizması

- known_hosts yönetimi

- SSH session audit log (kullanıcı, IP, event bazlı)

- WebSocket connection limitleri


***🛠️ Kullanılan Teknolojiler***

> **Backend**

- Go

- chi router

- SQLite (WAL mode, busy_timeout)

- JWT (HS256)

- bcrypt

- AES-GCM

- native golang.org/x/crypto/ssh

- WebSocket (gorilla / native upgrader)

> **Frontend**

- React

- Vite

- xterm.js

- WebSocket API

- Google Identity Services (frontend hazır)


<img width="567" height="557" alt="image" src="https://github.com/user-attachments/assets/08cf86c6-9272-497c-9d14-d8c926dc145e" />

<img width="713" height="636" alt="image" src="https://github.com/user-attachments/assets/d7d71436-9b97-41f7-9aed-dd8e6cd9f3a7" />

<img width="1304" height="544" alt="image" src="https://github.com/user-attachments/assets/ded3c6dc-8f10-49c3-aad9-3df980137997" />




<img width="2130" height="860" alt="image" src="https://github.com/user-attachments/assets/ba722e60-d9b3-4d7b-b8db-491fb90f8923" />



# **🤖 Yapay Zeka Kullanımı**

Bu proje geliştirilirken, yapay zekâ destekli araçlar; mimari tasarım, hata ayıklama, güvenlik senaryolarının değerlendirilmesi ve kod kalitesinin artırılması amacıyla yardımcı bir rehber olarak kullandım.

***Yapay zeka:***

- Nihai karar verici olarak değil,

- Geliştiricinin teknik değerlendirmelerini destekleyen bir araç olarak,

- Üretilen kodların tamamı geliştirici tarafından incelenmiş, test edilmiş ve projeye özel olarak uyarlanmıştır.

## https://www.alierenygt.com.tr/Projects/Detail/7
