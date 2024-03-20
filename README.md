# FROXIFY
Burp Upstream Proxy for fuzzing

### Ideas
![](image/Screenshot%202023-11-03%20110833.png)

Not support HTTP/2. This project using osmedeus for sub workflow fuzzing.

[x] Fuzz Param (x8)

[x] Sqli (Sqlmap raw request)

[ ] Crawl find path (Katana)

[ ] XSS (XSS Hunter)

[ ] Authen ()

[ ] SSTI (tplmap)

[ ] Os Command ()

[ ] Inj Header ()

[ ] Secret Key (Trufflehog)

[ ] Shellsock (ShellShockHunter)

[ ] nuclei

[ ] LFI/RFI

[x] Save log

### Build
[x] Use elasticsearch to save log request

[x] Use redis to filter same request

[x] Filter same request to save in elastic from config/TLS

[x] Read data pipeline in kafka

[x] Tool fuzzing
