# xmrhealthd

_Update: From version 1.2.1.0, BTCPay Server serves the XMR sync status over its own API, making this tool obsolete._

xmrhealthd queries the [get_info](https://www.getmonero.org/resources/developer-guides/daemon-rpc.html#get_info) API endpoint of a Monero node every ten seconds. It serves the `synchronized` value through an own HTTP endpoint on port 64325:

```
[{"CryptoCode":"XMR","Synced":true}]
```

## Commands

`xmrhealthd` takes the IP address of the monero node as an argument. If none is given, it connects to 127.0.0.1.

`xmrhealthd-btcpay` gets the IP address by running `docker inspect` on the `btcpayserver_monerod` container. Your user must be a member of the `docker` group, which makes it root equivalent, so be careful.

## Integration

### Server

Just create a domain and reverse proxy it. Example nginx configuration:

```
http {
    server {
        listen 80;
        listen [::]:80;
        server_name xmrhealthd.example.com;
        location / {
            proxy_pass http://10.10.10.10:64325/;
        }
    }
}
```

### Application

Get the JSON data, parse and output it and schedule the next execution:

```
<div id="health">
  <noscript>
    <div class="alert alert-info">Enable JavaScript in order to see the health status.</div>
  </noscript>
</div>

<script>
function updateHealth() {
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {
    if(xhr.readyState == 4) {
      let elem = document.getElementById("health");
      if(xhr.status == 200) {
        elem.innerHTML = "";
        for(c of JSON.parse(xhr.responseText)) {
          elem.innerHTML += `<span class="badge bg-${c.Synced ? 'success' : 'danger'}">${c.CryptoCode}: ${c.Synced ? 'synced' : 'out of sync'}</span> `;
        }
      } else {
        elem.innerHTML = '<span class="badge bg-warning">error</span>';
      }
    }
  }
  xhr.open("GET", "xmrhealthd.example.com", true);
  xhr.send(null);

  setTimeout(updateHealth, 10*1000); // 10 seconds
}
</script>
```
