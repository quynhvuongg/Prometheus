
### How to Prometheus work ?

![ ](https://images.viblo.asia/87ba2c92-9be1-436f-816e-f3be586e4c6a.png)

Prometheus sẽ chủ động  `pull`  (kéo) các  `metrics`  về qua HTTP mỗi 10s hay 30s do chúng ta thiết lập. Bản thân các service thì thường không thể tự export được các metrics cho Prometheus mà cần đến các  `Instructmentation/Exporter`.  

-   Exporter là những app được viết cho mấy cái thông dụng như Database, Server. Chúng ta chỉ cần chạy nó và nó sẽ export các metrics thu thập được cho mình.
-   Instructmentation thì chỉ những client-libraries được cung cấp bởi Prometheus hoặc một bên thứ 3 nào đó, để mình cài vào ứng dụng của mình, giúp tùy biến những metrics riêng của hệ thống.  

Một số Exporter :

-   Prometheus: chính bản thân prometheus cũng có một built-in exporter, export các metrics về service prometheus ra tại URL:  [http://prometheus.lc:9090/metrics](http://prometheus.lc:9090/metrics)
-   cAdvisor  : export các metrics của các docker service, các process trên server.
-   Node Exporter: export các metrics một con node (hiểu là một server) như CPU, RAM của node, dung lượng ổ đĩa, số lượng request tới node đấy, .etc.
-   Postgres Exporter, giúp đọc dữ liệu từ các bảng trong Postgres và export ra cho Prometheus
-   HAProxy Exporter

 
## Install :

Prometheus hỗ trợ 3 hình thức cài đặt các thành phần hệ thống gồm : Docker Image, cài đặt từ source với Go và file chương trình chạy sẵn đã được biên dịch sẵn.

**Install Prometheus with Docker** 

Step 1 : Viết Dockerfile (nếu cần )
```code
FROM prom/prometheus
ADD prometheus.yml /etc/prometheus/
```

Step 2 : Viết file prometheus.yml để tùy chỉnh config (nếu cần )

```yml
# prometheus.yml 
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


Step 2 :   Viết file docker-compose.yml cho docker compose :

```yml
# docker-compose.yml 
 version: '3.7' 
 volumes:
   prometheus_data: {}
    
 services: 
  
   prometheus:  
     image: prom/prometheus:v2.12.0
     volumes: # mount prometheus.yml vào container 
    #   - ./prometheus.yml:/etc/prometheus/prometheus.yml 
       - prometheus_data:/prometheus
     command: # chỉ định configure file vào đường dẫn chứa file configure
       - "--config.file=/etc/prometheus/prometheus.yml"  
     ports:  
       - "9090:9090"
```

Step 3 : docker -compose up 

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus1.png?raw=true)

![ ](https://github.com/quynhvuongg/Picture/blob/master/prometheus2.png?raw=true)

**Install node-exporter**
 
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

**Install Grafana**

Grafana  là một bộ mã nguồn mở sử dụng trong việc phân tích các dữ liệu thu thập được từ server và hiện thị một các trực quan dữ liệu thu thập được ở nhiều dạng khác nhau.

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
