# Configuration

Prometheus được cấu hình thông qua command-line flags và configure file :

- command-line flags cấu hình các tham số hệ thống bất biến (như vị trí lưu trữ, lượng dữ liệu lưu trên disk và trong bộ nhớ, v.v.)
- configure file viết bằng YAML xác định mọi thứ liên quan đến scraping jobs  and instances , rule files để load .

Prometheus có thể reload cấu hình của nó khi chạy. Nếu cấu hình mới không được định dạng tốt, các thay đổi sẽ không được áp dụng. Reload cấu hình được kích hoạt bằng cách gửi `SIGHUP` tới Prometheus process hoặc gửi  HTTP POST request đến `/-/reload` endpoint (khi cờ `--web.enable-lifecycle`  được bật )... đồng thời cũng sẽ reload rule files được configure .

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

Prometheus hỗ trợ 2 kiểu rules : Recording rule và Alerting rule. Các statement được chứa trong tệp YAML `rule_files` và được Prometheus load  vào.

*_Recording rules_*

Cho phép ta tính toán trước các biểu thức cần thiết hoặc tính toán đắt tiền và lưu kết quả của chúng dưới dạng một time series mới.

Truy vấn kết quả được tính toán trước thường sẽ nhanh hơn nhiều so với thực hiện biểu thức gốc mỗi khi cần. Điều này đặc biệt hữu ích cho bảng điều khiển, cần truy vấn cùng một biểu thức liên tục mỗi lần chúng làm mới.

```yml
#rule_files
groups:
  - name: example # tên của group
    rules: #Tên của time series mới để lưu kết quả
    - record: node_memory_MemFree_percent  
      expr: 100 - (100 * node_memory_MemFree_bytes / node_memory_MemTotal_bytes)  # biểu thức tính metric

```

Check rules sử dụng promtool : sudo apt install promtool :

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus5.png?raw=true)
  
```yml
#prometheus.yml
rule_files:  
- "rule_files"
```

*_Alerting rules_*

Cho phép ta xác định các điều kiện cảnh báo dựa trên biểu thức Prometheus và gửi thông báo kích hoạt cảnh báo tới dịch vụ bên ngoài.

```yml
#rule_files
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
- "rule_files"
```

```yml
#docker-compose.yml
 volumes:  
   - ...
   - ./rule_files:/etc/prometheus/rule_files
```

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus6.png?raw=true)

***Templates***

Prometheus hỗ trợ tạo khuôn mẫu trong các annotations và labels của cảnh báo, cũng như trong các trang điều khiển được phục vụ.
Templates có khả năng chạy các truy vấn đối với cơ sở dữ liệu cục bộ, lặp dữ liệu, sử dụng các điều kiện, định dạng dữ liệu, v.v.
Templates của Prometheus dựa trên templates của Go.

**_Simple alert templates_**

```yml
alert: InstanceDown
expr: up == 0
for: 5m
labels:
  severity: page
annotations:
  summary: "Instance {{$labels.instance}} down"
  description: "{{$labels.instance}} of job {{$labels.job}} has been down for more than 5 minutes."
```

Các mẫu cảnh báo sẽ được thực thi trong mỗi lần lặp cảnh báo kích hoạt  , vì vậy nên giữ mọi query và template nhẹ. Nếu bạn cần mẫu phức tạp hơn để cảnh báo, bạn nên liên kết với bảng điều khiển thay thế.

**_Simple iteration_**

```sh
{{ range query "up" }}
  {{ .Labels.instance }} {{ .Value }}
{{ end }}
```

' . ' -> giá trị của biến

**_Advanced iteration_**

```html
<table>
{{ range printf "node_network_receive_bytes{job='node',instance='%s',device!='lo'}" .Params.instance | query | sortByLabel "device"}}
  <tr><th colspan=2>{{ .Labels.device }}</th></tr>
  <tr>
    <td>Received</td>
    <td>{{ with printf "rate(node_network_receive_bytes{job='node',instance='%s',device='%s'}[5m])" .Labels.instance .Labels.device | query }}{{ . | first | value | humanize }}B/s{{end}}</td>
  </tr>
  <tr>
    <td>Transmitted</td>
    <td>{{ with printf "rate(node_network_transmit_bytes{job='node',instance='%s',device='%s'}[5m])" .Labels.instance .Labels.device | query }}{{ . | first | value | humanize }}B/s{{end}}</td>
  </tr>{{ end }}
</table>
```

**_Display one value_**

 ```sh
{{ with query "some_metric{instance='someinstance'}" }}
  {{ . | first | value | humanize }}
{{ end }}
```

**_Using console URL parameters_**

```sh
{{ with printf "node_memory_MemTotal{job='node',instance='%s'}" .Params.instance | query }}
  {{ . | first | value | humanize1024 }}B
{{ end }}
```

**_Defining reusable templates_**

```sh
{{define "myMultiArgTemplate"}}
  First argument: {{.arg0}}
  Second argument: {{.arg1}}
{{end}}
{{template "myMultiArgTemplate" (args 1 2)}}
```

**_Data Structures_**

```sh
type sample struct {
        Labels map[string]string
        Value  float64
}
```
