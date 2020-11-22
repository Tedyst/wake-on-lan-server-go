# Basic wake on lan server

Basic format

```
http://host?ip=192.168.10.108&address=AA:BB:CC:DD:EE:FF&redirectUrl=https://stoicatedy.ovh
```

It works by serving a React app that does the following things:

    1. It sends a WoL packet to the address mentioned in the query url
    2. It waits until a ping has been received from that IP
    3. If IP is ready, the webpage is redirected to redirectUrl
