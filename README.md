# FROXIFY
Burp Upstream Proxy for fuzzing

Thanks [proxify](https://github.com/projectdiscovery/proxify) help me for the development of this repo. Thanks [Catmandx](https://github.com/catmandx) for ideas and support me. 

### Ideas
![](image/Screenshot%202023-11-03%20110833.png)

Not support HTTP/2. Recommend using server 64/32 core and >32Gb RAM:

[ ] Fuzz Param

[ ] nuclei

[x] Save log

### Build
[x] Use elasticsearch to save log request

[x] Use redis to filter same request

[x] Filter same request to save in elastic from config/TLS

[x] Read data pipeline in kafka

[x] Tool fuzzing (sqlmap) (x8 in process)
