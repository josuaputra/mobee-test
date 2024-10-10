# Mendeploy Aplikasi Go Menggunakan Pipeline CI/CD GitHub Actions

Repositori ini berisi konfigurasi dan kode untuk mendeploy aplikasi Go ke cluster AWS EKS menggunakan pipeline CI/CD dengan GitHub Actions sebagai home-test dari Mobee. Deploymen ini mencakup konfigurasi ingress untuk mengakses aplikasi melalui HTTPS.


## Pre-requirements

Sebelum memulai, pastikan Anda memiliki komponen berikut terinstal dan dikonfigurasi:

### Instalasi PostgreSQL

1. **Helm**: Pastikan Anda telah menginstal Helm. Anda dapat mengikuti petunjuk instalasi [di sini](https://helm.sh/docs/intro/install/).
2. **Akses ke Cluster EKS**: Anda harus memiliki akses ke cluster AWS EKS dan kubectl terkonfigurasi untuk mengaksesnya.
3. **Jalankan Perintah Helm**: Gunakan perintah berikut untuk menginstal PostgreSQL di cluster EKS:

   ```bash
   kubectl create secret generic postgresql-auth \
   --from-literal=postgres-password=$(openssl rand -base64 12) \
   --from-literal=mobee-password=$(openssl rand -base64 12)

   ```bash
   helm repo add bitnami https://charts.bitnami.com/bitnami
   helm install my-postgresql bitnami/postgresql \
     --set postgresqlUsername=mobee \
     --set postgresqlPassword=$(kubectl get secret --namespace default my-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode) \
     --set postgresqlDatabase=mobeedb

### Instalasi Cert-Manager

1. **Kubectl**: Pastikan Anda telah menginstal kubectl dan dapat mengakses cluster Kubernetes.
2. **Helm**: Pastikan Anda telah menginstal Helm. Anda dapat mengikuti petunjuk instalasi di sini.
3. **Jalankan Perintah Helm*: Gunakan perintah berikut untuk menginstal cert-manager di cluster EKS:

   ```bash
   kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.5.4/cert-manager.yaml
   helm repo add jetstack https://charts.jetstack.io 
   helm repo update
   helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.5.4


## Penjelasan Pipeline CI/CD

Pipeline CI/CD diimplementasikan menggunakan GitHub Actions. Langkah-langkah utama dalam workflow adalah sebagai berikut:

1. **Build Aplikasi**: Pipeline membangun aplikasi Go dan membuat image Docker. Image ditandai dengan timestamp saat ini untuk versioning.
2. **Push ke Registry Docker**: Image yang dibangun didorong ke registry Docker (misalnya, Docker Hub atau AWS ECR) untuk penyimpanan.
3. **Deploy ke EKS**: Pipeline mendeploy aplikasi ke cluster AWS EKS. Ini menggunakan Kustomize untuk menerapkan manifest Kubernetes yang terletak di direktori `k8s`. Pipeline juga mengatur image di deployment berdasarkan image Docker yang telah didorong.
4. **Konfigurasi Ingress**: Resource ingress dibuat untuk mengelola akses eksternal ke aplikasi melalui domain tertentu, dengan dukungan HTTPS yang dikelola oleh cert-manager.

File workflow GitHub Actions terletak di `.github/workflows/staging-deployment.yml`. Ini di-trigger pada event push ke cabang `master`, memastikan bahwa kode terbaru dideploy secara otomatis.

## Konfigurasi Ingress

Konfigurasi ingress memungkinkan lalu lintas eksternal mencapai aplikasi Go. Komponen utama dari setup ingress adalah:

- **Resource Ingress**: Resource ini menentukan bagaimana lalu lintas harus diarahkan ke layanan aplikasi. Ini mendengarkan pada domain tertentu (misalnya, `example.com`) dan meneruskan lalu lintas ke layanan backend (misalnya, `my-go-app` yang berjalan di port 8080).
- **Dukungan TLS**: Resource ingress dikonfigurasi dengan TLS, menggunakan sertifikat yang diterbitkan oleh `cert-manager`. Ini memastikan koneksi aman HTTPS ke aplikasi.
- **Mengakses Aplikasi**: Setelah ingress diatur, aplikasi dapat diakses melalui domain yang ditentukan di browser web atau klien API.

## Asumsi dan Keputusan

- **Kontainerisasi**: Aplikasi dibangun sebagai kontainer Docker untuk memastikan konsistensi di seluruh lingkungan. Keputusan ini memungkinkan deployment yang mudah ke Kubernetes.
- **Kubernetes**: Target deployment adalah cluster AWS EKS, memanfaatkan Kubernetes untuk orkestrasi. Pilihan ini memberikan skalabilitas dan fitur manajemen untuk aplikasi.
- **GitHub Actions**: Menggunakan GitHub Actions untuk CI/CD memungkinkan integrasi yang mulus dengan repositori GitHub, memudahkan untuk mengotomatisasi proses deployment.
- **Pengendali Ingress**: Pengendali ingress (misalnya, NGINX Ingress Controller) diasumsikan sudah terinstal di cluster EKS, memungkinkan aplikasi untuk menangani lalu lintas HTTP/S eksternal.
- **cert-manager**: Keputusan untuk menggunakan `cert-manager` menyederhanakan pengelolaan sertifikat TLS, mengotomatiskan proses penerbitan dan pembaruan. Ini mengurangi kompleksitas dalam menangani sertifikat SSL secara manual.
- **AWS EKS**: Kami menggunakan AWS EKS sebagai orkestrasi untuk aplikasi kami. 
- **CI/CD**: Untuk CI/CD, kami memisahkan langkah-langkah untuk membangun dan mendeploy aplikasi.
- **Kustomize**: Kami menggunakan Kustomize untuk mengelola konfigurasi Kubernetes dan mempermudah proses deployment.
