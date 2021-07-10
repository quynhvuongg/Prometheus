# Labels
<!-- TOC -->

- [Labels](#labels)
  - [What are labels](#what-are-labels)
  - [What are labels used for](#what-are-labels-used-for)
  - [Instrumentation & Target labels](#instrumentation--target-labels)
  - [Relabeling](#relabeling)
    - [`<relabel_config>`](#relabel_config)
    - [`<metric_relabel_configs>`](#metric_relabel_configs)
    - [`<alert_relabel_configs>`](#alert_relabel_configs)

<!-- /TOC -->

## What are labels

Labels là các cặp key-value kết hợp với name metric để xác định cho time series trong Prometheus.

Ex: `http_requests_total{service="users-directory", instance="1.2.3.4"}`

-> metric trên có 2 labels : `service="users-directory"` và `instance="1.2.3.4"`

- key:  
  - có thể bao gồm: chữ, số hay dấu gạch dưới `_` sao cho match với regex `[a-zA-Z_][a-zA-Z0-9_]*`.
  - key bắt đầu với `_` được sử dụng nội bộ.

- value: có thể chứa bất kỳ Unicode nào.

Thay đổi bất kỳ giá trị label nào, bao gồm thêm hoặc xóa label, sẽ tạo ra time series mới.

Không đặt tên label trong name metric, vì điều này gây ra sự dư thừa và sẽ gây nhầm lẫn nếu các nhãn tương ứng được tổng hợp đi.

CAUTION: Hãy nhớ rằng mọi kết hợp khác nhau của các cặp key-value đại diện cho một chuỗi thời gian khác nhau, có thể làm tăng đáng kể lượng dữ liệu được lưu trữ. Không sử dụng label để lưu trữ  dimensions với số lượng thẻ cao (nhiều giá trị nhãn khác nhau), chẳng hạn như ID người dùng, địa chỉ email hoặc các bộ giá trị không giới hạn khác.

## What are labels used for

- Sử dụng nhãn để phân biệt các đặc điểm của đối tượng đang được đo.

Ex :  Trong `api_http_requests_total` phân biệt các yêu cầu : `operation="create|update|delete"`.

- Lọc metric dựa trên labels
  
```sh
#Only consider the Users Directory service
http_requests_total{service="users-directory}"
#Filter out successful requests
http_requests_total{status!="200"}

```

- Tổng hợp nhãn
  
 ```sh
sum(http_requests_total)
count_values("version", build_version)
topk(5, http_requests_total)

```

## Instrumentation & Target labels

Labels đến từ hai nguồn Instrumenation labels và Target labels. Trong khi bạn làm việc với PromQL, không có sự khác biệt giữa chúng, nhưng quan trọng bạn cần phân biệt chúng để có được lợi ích tốt nhất từ nhãn.

**_Instrumentation labels_**:

Là những nhãn đến từ các instrumentation của bạn. Instrumentation có thể là thư viện, dịch vụ hay hệ thống con hỗ trợ Prometheus lấy metrics mong muốn.

**_Target labels_**

- Là những nhãn xác định danh tính mục tiêu của Prometheus, tối thiểu bạn sẽ có nhãn `job` chỉ ra đối tượng giám sát và nhãn `instance` xác định mục tiêu cụ thể. Ngoài ra, có thể có một số nhãn khác cho environment , datacenter, team,...
- Chúng không thay đổi theo thời gian vì mỗi khi thay đổi sẽ ảnh hưởng đến sự liên tục của biểu đồ và có thể gây ra các sự cố với các quy tắc và cảnh báo.
- Để sử dụng chúng hữu ích nhất có thể ngoài việc giữ chúng không thay đổi thì cần phải giữ số lượng nhãn ở mức tối thiểu nhất tránh gây dư thừa hay vô tình làm việc trên các mục tiêu không liên quan.
- Target labels có thể do người dùng cấu hình hoặc đến từ service discovery kết hợp với relabeling.

```sh
scrape_configs:
- job_name: prometheus
  static_configs:
  - targets:
    localhost:9090
```

```sh
scrape_configs:
- job_name: file
  file_sd_configs:
  - files:
    - '*.json'
  relabel_configs:
  - source_labels: [team]
    regex: monitoring
    replacement: monitor
    target_label: team
    action: replace
```

## Relabeling

Relabeling là công cụ mạnh mẽ để tự động viết lại bộ nhãn của một mục tiêu trước khi nó được scrape. Thông qua một cơ chế này, bất kỳ nhãn nào cũng có thể được loại bỏ, tạo hoặc sửa đổi theo cấp độ mục tiêu.

Relabeling hoạt động như sau :

### `<relabel_config>`

```sh
# Xác định danh sách nhãn nguồn
source_labels: '[' <labelname> [, ...] ']'

# Dấu phân cách được đặt giữa các giá trị của nhãn nguồn
separator: <string> | default = ;

# Nhãn dùng để lưu giá trị kết quả của hành động thay thế
target_label: <labelname>

# Regular expression
regex: <regex> | default = (.*)

# Giá trị thay thế nếu thỏa mãn biểu thức trên
replacement: <string> | default = $1

# Hành động của relabel
action: <relabel_action> | default = replace

```

**_Relabel_action_**:

- `replace`: Xét nhãn trong `source_labels` có giá trị match với regex thì thực hiện set `target_label` với giá trị trong `replacement`, ngược lại sẽ không xảy ra gì cả.

- `keep`: Giữ những mục tiêu mà nó thỏa mãn regex nếu không sẽ bị bỏ đi.

```sh
scrape_configs:
  - files:
  - '*.json'
  relabel_configs:
  - source_labels: [team]
  regex: infa
  action: keep

```

- `drop`: Bỏ những mục tiêu mà nó thỏa mãn biểu thức regex.

```sh
scrape_configs:
- job_name: file
  file_sd_configs:
  - files:
    - '*.json'  
  relabel_configs:
  - source_labels: [job, team]
    regex: prometheus;monitoring
    action: drop

```

- `labelmap`: Nếu tất cả tên nhãn thoả mãn regex thì tên nhãn sẽ được thay thế theo giá trị `replacement`, còn giá trị được giữ lại.

- `labeldrop`: Nếu tất cả tên nhãn trong `source_labels` thoả mãn regex sẽ được loại bỏ khỏi bộ nhãn.

- `labelkeep`: Nếu tất cả tên nhãn trong `source_labels` thoả mãn regex sẽ được giữ lại trong bộ nhãn, còn lại sẽ bị bỏ.

### `<metric_relabel_configs>`

- Áp dụng cho time series sau khi scrape nhưng trước khi chúng được lưu trữ, định dạng cấu hình và hoạt động tương tự như `<relabel_configs>`.

- Sử dụng cho 2 trường hợp: khi bỏ các số liệu đắt và khi sửa các số liệu xấu.

```sh
scrape_configs:
- job_name: prometheus
  static_configs:
  - targets:
    - localhost:9090
  metric_relabel_configs:
  - source_labels: [__name__]
    regex: http_request_size_bytes
    action: drop

```

```sh
scrape_configs:
- job_name: mishaving
  static_configs:
  - targets:
     - localhost:1234
  metric_relabel_configs:
  - regex: 'node_.*'
  action: labeldrop

```

### `<alert_relabel_configs>`

Ngoài ra, Relabeling còn được sử dụng trong cấu hình cảnh báo, áp dụng cho các cảnh báo trước khi chúng được gửi đến Alertmanager, định dạng cấu hình và hoạt động tương tư như `<relabel_configs >`.

```sh
alerting:
  alert_relabel_config:
    - source_labels: [time_window]
      regex: never
      action: drop
```
