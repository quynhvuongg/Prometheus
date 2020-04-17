# Exporters

## Node-exporter

Exporter thu thập các số liệu về phần cứng và OS như :

1. CPU
2. Filesystem
3. Diskstats
4. Netdev
5. Meminfo
6. Hwmon
7. Stat
8. Uname
9. Loadavg
10. Network  
....

 ![ ](https://github.com/quynhvuongg/Picture/blob/master/node1.png?raw=true)

## CAdvisor

Exporter thu thập các số liệu về các container như :

1. CPU
2. Memory
3. Usage :  Total Usage, Usage per Core, Usage Breakdown, Usage  Breakdown
4. Network : Interfaces, Throughput, Errors
5. Filesystem

....

![ ](https://github.com/quynhvuongg/Picture/blob/master/node2.png?raw=true)

## HAproxy-exporter

Exporter thu thập các số liệu thống kê HAproxy bao gồm tình trạng server, tốc độ yêu cầu hiện tại, thời gian phản hồi, v.v.. trên mỗi frontend, backend, server.

```sh
#haproxy.cfg
defaults
  mode http
  timeout server 5s
  timeout connect 5s
  timeout client 5s
  
frontend frontend
  bind *:1234
  use_backend backend
  
backend backend
  server node_exporter 127.0.0.1:9100
  
frontend monitoring
bind *:1235
no log
stats uri /
stats enable

```

```yml
#prometheus.yml
scrape_configs:
  - job_name: haproxy
    static_configs:
      - targets:
        - localhost:9101
```

![ ](https://www.robustperception.io/wp-content/uploads/2015/11/Screenshot-111115-205944-640x431.png)

![ ](https://cdn.haproxy.com/wp-content/uploads/2019/04/prometheus_output.png)

## Blackbox-exporter

Cho phép blackbox thăm dò các endpoint qua HTTP, HTTPS, DNS, TCP và ICMP.

```yml
#prometheus.yml
scrape_configs:
- job_name: blackbox
  metrics_path: /probe
  params:
    module: [http_2xx]
static_configs:
  - targets:
      - http://www.prometheus.io
      - http://www.robustperception.io
      - http://demo.robustperception.io
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: localhost:9115
# Blackbox exporter

```

* HTTP prober

```yml
#blackbox.yml
modules:  
  http_post_2xx:
    prober: http
    timeout: 5s
    http:
      method: POST
      headers:
        Content-Type: application/json
      body: '{}'
```

* DNS prober

```yml
#blackbox.yml
dns_tcp:
  prober: dns
  dns:
    transport_protocol: "tcp"
    query_name: "www.prometheus.io"

```

* TCP prober

```yml
#blackbox.yml
tcp_connect_example:
  prober: tcp
  timeout: 5s

```

* ICMP prober

```yml
#blackbox.yml
icmp_ipv4:
  prober: icmp
  icmp:
    preferred_ip_protocol: ip4

```

## SNMP-exporter

SNMP - Simple Network Management Protocol là tập hợp các giao thức được sử dụng để quản lý các thiết bị mạng như router, switch hay các server đang vận hành.

![ ](https://miro.medium.com/max/1400/0*9q6WjdLZ3ui03fnr)

SNMP-exporter thu thập thông tin được từ SNMP bằng cách đọc tệp cấu hình mặc định `snmp.yml`  chứa  thông tin OIDs( Object Identifiers) để `walk/get` từ thiết bị và thông tin đăng nhập sử dụng trong trường hợp nếu đó là SNMP v2 hoặc SNMP v3, sau đó thực hiện thu thập số liệu.

`snmp.yml` không sử dụng viết tay mà sử dụng một generator để tạo ra nó.

Ex: snmp.yml cho một Cisco router

```yml
Cisco:  
  version: 3  
  auth:  
  username: snmp  
  password: password  
  auth_protocol: SHA  
  priv_protocol: DES  
  security_level: authPriv  
  priv_password: private
walk:  
 - 1.3.6.1.2.1.1 # sysInfo  
 - 1.3.6.1.2.1.2.2 # ifTable  
 - 1.3.6.1.2.1.31.1.1 # ifXTable  
metrics:  
#sysInfo  
- name: sysUpTime  
  oid: 1.3.6.1.2.1.1.3  
  type: counter  
  lookups:  
- labels:  
  labelname: sysDescr  
  oid: 1.3.6.1.2.1.1.1.0  
  type: DisplayString  
- labels:  
  labelname: sysName  
  oid: 1.3.6.1.2.1.1.5.0  
  type: DisplayString  
- labels:  
  labelname: sysLocation  
  oid: 1.3.6.1.2.1.1.6.0  
  type: DisplayString  
- labels:  
  labelname: sysContact  
  oid: 1.3.6.1.2.1.1.4.0  
  type: DisplayString  
#Interfaces  
#Interface ifIndex  
- name: ifIndex  
  oid: 1.3.6.1.2.1.2.2.1.1  
  type: gauge  
  indexes:  
  - labelname: ifIndex  
    type: Integer  
  lookups:  
  - labels:  
    - ifIndex  
    labelname: ifDescr  
    oid: 1.3.6.1.2.1.2.2.1.2  
    type: DisplayString  
  - labels:  
    - ifIndex  
    labelname: ifName  
    oid: 1.3.6.1.2.1.31.1.1.1.1  
    type: DisplayString  
  - labels:  
    - ifIndex  
    labelname: ifAlias  
    oid: 1.3.6.1.2.1.31.1.1.1.18  
    type: DisplayString  
#Interface Type  
- name: ifType  
  oid: 1.3.6.1.2.1.2.2.1.3  
  type: gauge  
  indexes:  
  - labelname: ifIndex  
    type: Integer  
  lookups:  
  - labels:  
    - ifIndex  
    labelname: ifDescr  
    oid: 1.3.6.1.2.1.2.2.1.2  
    type: DisplayString  
  - labels:  
    - ifIndex  
    labelname: ifName  
    oid: 1.3.6.1.2.1.31.1.1.1.1  
    type: DisplayString  
  - labels:  
    - ifIndex  
    labelname: ifAlias  
    oid: 1.3.6.1.2.1.31.1.1.1.18  
    type: DisplayString

```

```yml
#prometheus.yml
scrape_configs:
  - job_name: 'snmp'
    static_configs:
      - targets:
        - 192.168.1.2 # SNMP device.
    metrics_path: /snmp
    params:
      module: [if_mib]
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: localhost:9116
        # SNMP exporter

```
