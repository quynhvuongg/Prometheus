# Service discovery (SD)

![ ](https://o.quizlet.com/8vqZtNfDtb4wz1cdMJiqKA.png)

## Why do we need Service discovery

Cho đến hiện tại, Prometheus xác định những mục tiêu thông qua cấu hình `static_configs` trong file prometheus.yml do người dùng thiết lập đơn giản và dễ sử dụng. Tuy nhiên, trong môi trường năng động việc cập nhật thủ công khi phải thêm hoặc xóa khá khó chịu, do đó Service discovery xuất hiện để công việc này dễ dàng hơn.

## Basic service discoveries

### File

Thường được gọi là file_sd có định dạng là JSON hoặc YAML do người dùng cung cấp trong hệ thống tệp cục bộ, chứa thông tin về các mục tiêu cho Prometheus.

Ex:

```json
#filesd.json
[
 {
   "targets": [ "localhost:9090" ],
   "labels": {"job": "prometheus"}
  },
 {
  "targets": [ "localhost:9100" ],
  "labels": {"job": "node-exporter"}
  },
  {
  "targets": [ "localhost:8080" ],
  "labels": {"job": "cadvisor"}
  }
]
```

```yml
#prometheus.yml
scrape_configs:
- job_name: file
  file_sd_configs:
  - files:
    - '*.json'

```

![ ](https://github.com/quynhvuongg/Picture/blob/master/ser1.png?raw=true)

### Consul

Prometheus truy xuất các mục tiêu từ danh mục API của Consul mà chúng ta cung cấp thông qua Consul CLI , API hoặc sử dụng file cấu hình trong thư mục cấu hình Consul.
Consul sẽ trả về thông tin các mục tiêu bằng metadata, do đó ta cần sử dụng Relabeling để chọn mục tiêu mong muốn và chuyển đổi metadata thành target label.

Ex: Chọn các dịch vụ có thẻ `prod` và sử dụng tên dịch vụ cho nhãn `job`.

```yml
scrape_configs:
  - job_name: consul
    consul_sd_configs:
      - server: 'localhost:8500'
    relabel_configs:
      - source_labels: [__meta_consul_tags]
        regex: .*,prod,.*
        action: keep
      - source_labels: [__meta_consul_service]
        target_label: job

```

### OpenStack

Tương tự như Consul, Prometheus truy xuất thông tin mục tiêu từ OpenStack Nova API  và sử dụng Relabeling để chuyển metadata thành target label.

Ex:

```yml
scrape_configs:
  - job_name: 'openstack'  
    openstack_sd_configs:
      - source_labels:  
        - __address__  
        - __meta_openstack_tag_prometheus_io_port  
        action: replace  
        regex: ([^:]+)(?::\d+)?;(\d+)  
        replacement: $1:$2  
        target_label: job
      - source_labels: [__meta_openstack_instance_name]  
        target_label: instance

```
