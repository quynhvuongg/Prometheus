# Install

Prometheus hỗ trợ 3 hình thức cài đặt các thành phần hệ thống gồm : Docker Image, cài đặt từ source với Go và file chương trình chạy sẵn đã được biên dịch sẵn.

## Install Prometheus with Docker

Step 1: Viết Dockerfile (nếu cần )

```code
FROM prom/prometheus
ADD prometheus.yml /etc/prometheus/
```

Step 2: Viết file prometheus.yml để tùy chỉnh config (nếu cần )

```yml
# prometheus.yml
global:  
  # Set thời gian lấy metrics sau mỗi 15s (defaul = 1m)
  scrape_interval:  15s
  #Đính kèm labels này vào bất kỳ time series hoặc alert nào khi liên lạc với hệ thống bên ngoài(remote storage , Alertmanager )
  external_labels:  
    monitor:  'my-monitor'  
#Một cấu hình scrape chứa chính xác endpoint để scrape
scrape_configs:  
  - job_name:  'prometheus'  
    static_configs:
         - targets: ['localhost:9090']
  
```

Step 3: Viết file docker-compose.yml cho docker compose:

```yml
# docker-compose.yml
 version: '3.7'
 volumes:
   prometheus_data: {}
 services:
  
   prometheus:  
     image: prom/prometheus:v2.12.0
     # mount prometheus.yml vào container
     volumes:  
       - ./prometheus.yml:/etc/prometheus/prometheus.yml
       - prometheus_data:/prometheus
    # chỉ định configure file vào đường dẫn
     command:
       - "--config.file=/etc/prometheus/prometheus.yml"  
     ports:  
       - "9090:9090"
```

Step 4: docker-compose up

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus1.png?raw=true)

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus2.png?raw=true)

## Install node-exporter

Thu thập các số liệu hệ thống như sử dụng cpu / bộ nhớ / lưu trữ và sau đó nó xuất chúng cho Prometheus . Nó có thể được chạy như một container docker đồng thời báo cáo các số liệu thống kê cho hệ thống máy chủ.

docker-compose.yml

```yml
version: '3.7'

volumes:
  prometheus_data: {}

services:

  prometheus:
    image: prom/prometheus:v2.12.0
    volumes:  
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
    ports:
      - "9090:9090"

  node-exporter:
     image: prom/node-exporter
     ports:
       - '9100:9100'
```

prometheus.yml

```yml
global:  
  scrape_interval:  15s
  external_labels:
      monitor:  'my-monitor'  

scrape_configs:  
  - job_name: 'prometheus'  

    static_configs:
         - targets: ['localhost:9090']

  - job_name: 'node-exporter'

    static_configs:
         - targets: ['node-exporter:9100']
```

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus3.png?raw=true)

## Install Grafana

Grafana là một bộ mã nguồn mở sử dụng trong việc phân tích các dữ liệu thu thập được từ server và hiện thị một các trực quan dữ liệu thu thập được ở nhiều dạng khác nhau.

docker-compose.yml

```yml
version: '3.7'

volumes:
  prometheus_data: {}
  grafana_data: {}

services:

  prometheus:
    image: prom/prometheus:v2.12.0
    volumes:  
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
    ports:
      - "9090:9090"

  node-exporter:
    image: prom/node-exporter
    ports:
       - "9100:9100"

  grafana:
    image: grafana/grafana
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=pass
    depends_on:
      - prometheus
    ports:
      - "3000:3000"
```

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus4.png?raw=true)

**_Install CAdvisor_**

Container Advisor phân tích và hiển thị dữ liệu hiệu suất và sử dụng tài nguyên từ các container đang chạy .

```yml
#docker-compose.yml
version: '3.7'

volumes:
  prometheus_data: {}
  grafana_data: {}

services:

  prometheus:
    image: prom/prometheus:v2.12.0
    volumes:  
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
      - ./rule_files:/etc/prometheus/rule_files
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
    ports:
      - "9090:9090"
    depends_on:
    - cadvisor

  cadvisor:
    image: google/cadvisor
    ports:
    - 8080:8080
    volumes:
    - /:/rootfs:ro
    - /var/run:/var/run:rw
    - /sys:/sys:ro
    - /var/lib/docker/:/var/lib/docker:ro

  node-exporter:
    image: prom/node-exporter
    ports:
       - "9100:9100"

  grafana:
    image: grafana/grafana
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=pass
    depends_on:
      - prometheus
    ports:
      - "3000:3000"
```

```yml
#prometheus.yml
global:  
  scrape_interval:  15s
  external_labels:
      monitor:  'my-monitor'  

scrape_configs:  
  - job_name: 'prometheus'  

    static_configs:
         - targets: ['localhost:9090']

  - job_name: 'node-exporter'

    static_configs:
         - targets: ['node-exporter:9100']

  - job_name: cadvisor
    scrape_interval: 5s
    static_configs:
    - targets:
      - cadvisor:8080

rule_files:  
  - "rule_files"
```

docker-compose up

![ ](https://github.com/quynhvuongg/Picture/blob/master/CAdvisor1.png?raw=true)

![ ](https://github.com/quynhvuongg/Picture/blob/master/CAdvisor4.png?raw=true)

Gafana

![  ](https://github.com/quynhvuongg/Picture/blob/master/CAdvisor3.png?raw=true)
