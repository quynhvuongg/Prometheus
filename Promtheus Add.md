#

## Prometheus High-availability

High-availability không chỉ quan trọng đối với phần mềm đối mặt với khách hàng, mà còn rất quan trọng đối với cơ sở hạn tầng. Nếu cơ sở hạ tầng giám sát không khả dụng cao, thì có nguy cơ mọi người không được thông báo để cảnh báo hay khi máy chủ down sẽ gây giãn đoạn công việc.Do đó, tính sẵn sàng cao phải được cân nhắc kỹ lưỡng.

![ ](https://github.com/quynhvuongg/Picture/blob/master/high1.png?raw=true)

Để chạy Prometheus theo cách khả dụng cao, hai (hoặc nhiều) máy chủ cần phải chạy với cùng một cấu hình, điều đó có nghĩa là chúng scrape các mục tiêu giống nhau -> chúng sẽ có cùng dữ liệu trong bộ nhớ và trên disk, chúng trả lời các yêu cầu theo cùng một cách. Trong thực tế, điều này có thể không hoàn toàn đúng, vì các chu kỳ scrape có thể hơi khác nhau và do đó dữ liệu được ghi có thể hơi khác nhau.

Ngoài ra, các máy chủ Prometheus sử dụng `external labels` khác nhau, vì vậy dữ liệu của chúng không xung đột nếu sử dụng bộ lưu trữ từ xa. Đối với cảnh báo, việc sử dụng `alert_relabel_configs` để đảm bảo chúng vẫn gửi các cảnh báo có nhãn giống hệt nhau mà Alertmanager sẽ tự động de-duplicate.

## Prometheus federation

Federation cho phép một máy chủ Prometheus scrape time series được chọn từ một máy chủ Prometheus khác.

Thông thường, nó được sử dụng để đạt được các thiết lập giám sát Prometheus có thể mở rộng hoặc để kéo các số liệu liên quan từ Prometheus của một dịch vụ sang một dịch vụ khác.

### Hierarchical federation

Federation phân cấp cho phép Prometheus mở rộng quy mô sang môi trường có hàng chục trung tâm dữ liệu và hàng triệu nút. Trong trường hợp sử dụng này, cấu trúc liên kết giống như một cái cây, với các máy chủ Prometheus cấp cao hơn thu thập dữ liệu timeseries tổng hợp từ các máy chủ trực thuộc.

Ví dụ một thiết lập có thể bao gồm nhiều máy chủ Prometheus trên mỗi trung tâm dữ liệu thu thập dữ liệu với độ chi tiết cao (instance-level drill-down) và một bộ máy chủ Prometheus toàn cầu chỉ thu thập và lưu trữ dữ liệu tổng hợp (job-level drill-down ) từ các máy chủ cục bộ. Điều này cung cấp một cái nhìn toàn cầu tổng hợp và cục bộ chi tiết.

### Cross-service federation

Federation dịch vụ chéo, máy chủ Prometheus của một dịch vụ được cấu hình để quét dữ liệu đã chọn từ máy chủ Prometheus của dịch vụ khác để cho phép cảnh báo và truy vấn đối với cả hai bộ dữ liệu trong một máy chủ.

Ví dụ: bộ lập lịch cụm chạy nhiều dịch vụ có thể tiết lộ thông tin sử dụng tài nguyên (như sử dụng bộ nhớ và CPU) về các dịch vụ đang chạy trên cụm. Mặt khác, một dịch vụ chạy trên cụm đó sẽ chỉ hiển thị các số liệu dịch vụ dành riêng cho ứng dụng. Thông thường, hai bộ số liệu này được loại bỏ bởi các máy chủ Prometheus riêng biệt. Sử dụng federation, máy chủ Prometheus chứa các số liệu cấp độ dịch vụ có thể lấy các số liệu sử dụng tài nguyên cụm về dịch vụ cụ thể của nó từ Prometheus cụm, để cả hai bộ số liệu có thể được sử dụng trong máy chủ đó.

### Configuration

Trên bất kỳ máy chủ Prometheus nào, điểm cuối `/federate` cho phép truy xuất giá trị hiện tại cho một time series đã chọn trong máy chủ đó. Ít nhất một tham số `match[]` URL phải được chỉ định để chọn chuỗi để hiển thị. Mỗi đối số `match[]` cần chỉ định một bộ chọn instant vector như `up`hoặc `{job="api-server"}`. Nếu nhiều tham số `match[]` được cung cấp, sự kết hợp của tất cả các chuỗi khớp được chọn.

Để liên kết các số liệu từ máy chủ này sang máy chủ khác, ta cấu hình máy chủ Prometheus đích của bạn để quét từ  điểm cuối`/federate` của máy chủ nguồn, đồng thời cho phép `honor_labels`tùy chọn scrape (không ghi đè bất kỳ nhãn nào được hiển thị bởi máy chủ nguồn) và chuyển các tham số `match[]` mong muốn.

Ví dụ các `scrape_configs` sau đây liên kết bất kỳ chuỗi nào có nhãn `job="prometheus"` hoặc tên số liệu bắt đầu `job:`từ các máy chủ Prometheus tại  `source-prometheus-{1,2,3}:9090` vào Prometheus.

```yml
scrape_configs:
  - job_name: 'federate'
    scrape_interval: 15s

    honor_labels: true
    metrics_path: '/federate'

    params:
      'match[]':
        - '{job="prometheus"}'
        - '{__name__=~"job:.*"}'

    static_configs:
      - targets:
        - 'source-prometheus-1:9090'
        - 'source-prometheus-2:9090'
        - 'source-prometheus-3:9090'

```
