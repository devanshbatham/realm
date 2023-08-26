<h1 align="center">
    realm
  <br>
</h1>

<h4 align="center">A utility for recursively traversing SSL/TLS certificates for collecting DNS names</h4>


<p align="center">
  <a href="#install">🏗️ Install</a>
  <a href="#usage">⛏️ Usage</a>
  <br>
</p>

![realm](https://github.com/devanshbatham/realm/blob/main/static/banner.png?raw=true)

# Install
```sh
go install github.com/devanshbatham/realm@v0.0.2
```

# Usage

```sh
(~) >>> realm -d example.com -n 2

🔍 Traversing example.com: 8 domains found
example.com
example.org
www.example.com
www.example.edu
www.example.net
www.example.org
example.net
example.edu

🔍 Traversing example.org: 8 domains found
example.com
example.org
www.example.com
www.example.edu
www.example.net
www.example.org
example.net
example.edu

🔍 Traversing www.example.com: 8 domains found
example.org
www.example.com
www.example.edu
www.example.net
www.example.org
example.net
example.edu
example.com

🔍 Traversing www.example.edu: 8 domains found
www.example.org
example.net
example.edu
example.com
example.org
www.example.com
www.example.edu
www.example.net

...
...
```
