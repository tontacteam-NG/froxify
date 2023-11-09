# froxy
Upstream Proxy Burp for fuzzing

Thanks [proxify](https://github.com/projectdiscovery/proxify) help me for the development of this repo. Thanks [Sub](https://github.com/n00b-bot), [Catmandx](https://github.com/catmandx) for ideas and support me. 

### Ý Tưởng
![](image/Screenshot%202023-11-03%20110833.png)

Sẽ cần 1 con Mitmproxy mạnh để tải trọng 1 số tools như:

[ ] Fuzz Param

[ ] Command-inj-header

[ ] Eny test obb

[ ] spiderhog

[ ] nuclei

[x] Save log

### Triển khai
[x] Sử dụng elasticsearch để save log request

[x] Sử dụng redis để lọc request trùng nhau

[x] Lọc gói tin để lưu trong elastic lọc từ bộ TLS trong config

[x] Data pipeline qua kafka

[ ] Chuẩn bị 1 số tool fuzzing (Process)

+ Để chuẩn bị cho fuzzing cần parser request trước