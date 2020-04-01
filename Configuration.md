## Configuration 

Prometheus được cấu hình thông qua command-line flags và configure file :
-  command-line flags cấu hình các tham số hệ thống bất biến (như vị trí lưu trữ, lượng dữ liệu lưu trên đĩa và trong bộ nhớ, v.v.)
-  configure file viết bằng YAML xác định mọi thứ liên quan đến scraping jobs  and instances , rule files to load .

Prometheus có thể reload cấu hình của nó khi chạy. Nếu cấu hình mới không được định dạng tốt, các thay đổi sẽ không được áp dụng. Reload cấu hình được kích hoạt bằng cách gửi `SIGHUP`tới Prometheus process hoặc gửi  HTTP POST request đến `/-/reload` endpoint (khi cờ `--web.enable-lifecycle`  được bật )..  đồng thời cũng sẽ reload rule files được configure .
 Để chỉ định configure file nào được reload ta sử dụng cờ  `--config.file` 
 
 ```yml
#prometheus.yml
global:  
  scrape_interval:  15s  # Set thời gian lấy metrics sau mỗi 15s (defaul = 1m)
  external_labels: #Đính kèm labels này vào bất kỳ time series hoặc alert nào khi liên lạc với hệ thống bên ngoài (remote storage , Alertmanager )
    monitor:  'my-monitor'  
scrape_configs: #Một cấu hình scrape chứa chính xác endpoint để scrape 
  - job_name:  'prometheus'  
    static_configs: 
         - targets: ['localhost:9090']
 
  # metrics_path defaults to '/metrics'
```
 
 ```yml
 # docker-compose.yml
  command: # chỉ định configure file vào đường dẫn chứa file configure
       - "--config.file=/etc/prometheus/prometheus.yml"  
```

 ---
Prometheus hỗ trợ 2 kiểu rules : recording rule và alerting rule. Các statement được chứa trong tệp YAML `rule_file` và được Prometheus load  vào.

**Recording rules**

Cho phép ta tính toán trước các biểu thức cần thiết hoặc tính toán đắt tiền và lưu kết quả của chúng dưới dạng một time series mới. 

Truy vấn kết quả được tính toán trước thường sẽ nhanh hơn nhiều so với thực hiện biểu thức gốc mỗi khi cần. Điều này đặc biệt hữu ích cho bảng điều khiển, cần truy vấn cùng một biểu thức liên tục mỗi lần chúng làm mới. 

```yml
#record_rules.yml
groups: 
  - name: example # tên của group
    rules: #Tên của time series mới để lưu kết quả
    - record: node_memory_MemFree_percent  
      expr: 100 - (100 * node_memory_MemFree_bytes / node_memory_MemTotal_bytes)  # biểu thức tính các metric 

```
Check rules  sử dụng promtool : sudo apt install promtool :

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus5.png?raw=true)
  

```yml
#prometheus.yml
rule_files:  
- "record_rules.yml"
```

**Alerting rules**

Cho phép ta xác định các điều kiện cảnh báo dựa trên biểu thức  Prometheus và gửi thông báo về kích hoạt cảnh báo  tới dịch vụ bên ngoài.

```yml
#alert_rules.yml
groups:
- name: example
  rules:
# Cảnh báo cho bất kì instance không tới được trong 2m
  - alert: service_down # tên của alert
    expr: up == 0 # biểu thức đánh giá 
    for: 2m # thời gian chờ xử lý trước khi kích hoạt cảnh báo
   labels: #label được đính kèm cùng alert ,nhãn đã tồn tại nào xung đột sẽ bị ghi đè
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} down"
      description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 2 minutes."
#annotations chỉ định một bộ nhãn thông tin có thể được sử dụng để lưu trữ thông tin bổ sung dài hơn như mô tả cảnh báo hoặc liên kết runbook. 
```

```yml
#prometheus.yml
rule_files:  
- "alert_rules.yml"
```


