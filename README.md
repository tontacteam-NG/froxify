# froxy
Upstream Proxy Burp for fuzzing

### Ý Tưởng
![](image/Screenshot%202023-11-03%20110833.png)

Sẽ cần 1 con Mitmproxy mạnh để tải trọng 1 số tools như:

[ ] Fuzz Param

[ ] Command-inj-header

[ ] Eny test obb

[ ] spiderhog

[ ] nuclei

[ ] Save log (processing 70%)

### Triển khai
[ ] Sử dụng elasticsearch để save log request

[ ] Sử dụng redis để lọc request trùng nhau

[ ] Lọc gói tin để lưu trong elastic lọc từ bộ TLS trong config

[ ] Chuẩn bị 1 số tool fuzzing