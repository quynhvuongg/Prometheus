# Writing Exporter

Trong trường hợp không thể thêm Direct instrumentation vào ứng dụng cũng như không tìm thấy một Exporters hiện có đáp ứng được yêu cầu thì phương án giải quyết vấn đề chính là tự viết một Exporter mới đáp ứng mong muốn của bản thân.
Tuy nhiên để Exporter đó tích hợp được với Prometheus thì nó cần tuân thủ theo những quy tắc sau đây.

## Metric

**_Metric naming_**

* Namespaces: Tên ứng dụng,tên dịch vụ hoặc tên xuất phát từ client-library
vd:  
    - **mysql**_up
    - **http**_request_duration_seconds
    - **process**_cpu_seconds_total

* Main part:
  - Nên thể hiện cùng một đối tượng đo logic trên tất cả các nhãn dimension.
  - request duration
  - bytes of data transfer
  - instantaneous resource usage as a percentage
  - Không được bao gồm các nhãn mà chúng được xuất cùng.
  - Theo quy tắc chung, tổng hợp trên tất cả các kích thước nhãn của một số liệu nhất định sẽ có ý nghĩa (mặc dù không nhất thiết phải hữu ích).
  
* Units:
  - Nên sử dụng các đơn vị cơ sở.
  - Sử dụng đơn vị nhất quán.
  - Expose dạng tỷ lệ không nên sử dụng tỷ lệ phần trăm ( Tốt hơn nên sử dụng Counter cho mỗi thành phần tỷ lệ ).

* Suffixes:
  - _total for Counters.
  - ø for Gauges.
  - bucket,_sum,_count for Histograms.
  - _sum,_count for Summaries.
  
**_Metric Labels_**

* Tránh sử dụng các nhãn tổng quát làm tên nhãn như `type` nó quá chung chung và thường vô nghĩa.
* Nếu có thể để tránh tên có khả năng xung đột với target label, chẳng hạn như `region`, `zone`, `cluster`, `availability_zone`, `az`, `datacenter`, `dc`, `owner`, `customer`, `stage`, `service`, `environment` và `env`.
* Đọc / viết và gửi / nhận là tốt nhất dưới dạng metric riêng biệt, thay vì nhãn. Điều này thường là do bạn chỉ quan tâm đến một trong số chúng tại một thời điểm và việc sử dụng chúng theo cách đó sẽ dễ dàng hơn.
* Chuỗi thời gian trong một số liệu không nên chồng chéo. Một tổng hoặc trung bình trên số liệu sẽ có ý nghĩa.
  
  vd: Đừng làm điều này:

```sh
my_metric{label=a} 1
my_metric{label=b} 6
my_metric{label=total} 7
```

**_Metric Types_**

Nên cố gắng match metric types của bạn với metric types của Prometheus (counters, gauges, summaries,histograms).

Thường thì sẽ không rõ loại số liệu là gì, đặc biệt nếu bạn tự động xử lý một bộ số liệu, `UNTYPED` là một mặc định an toàn.

**_Metrics to drop_**

Một số metrics không hữu ích nên cần loại bỏ tránh metrics quá đắt với số lượng thẻ cao hay có số lượng lớn, không hữu ích và thêm lộn xộn.

Vd: - Machine metrics (ví dụ: CPU,Filesystem,memory disk,...) đã được cung cấp bởi node-exporter.
    - Đối với JMX, jmx_exporter đã bao gồm các số liệu thống kê về quy trình và JVM.
    - Quantiles và các vấn đề liên quan có thể chọn bỏ chúng hoặc đưa chúng vào Summary.
    - Tốc độ trung bình kể từ khi bắt đầu ứng dụng, mức tối thiểu, tối đa và độ lệch chuẩn vì Prometheus hoàn toàn có thể tính toán được chúng.
    - ...

## Configuration

* Khi làm việc với các ứng dụng, bạn nên nhắm đến một exporter không yêu cầu cấu hình tùy chỉnh của người dùng ngoài việc cho biết ứng dụng đó ở đâu.

* Bạn cũng có thể cần cung cấp khả năng lọc ra một số metrics nhất định nếu chúng có thể quá chi tiết và tốn kém trên các thiết lập lớn.

* Khi làm việc với các hệ thống giám sát, frameworks và protocols khác, bạn sẽ cần cung cấp cấu hình hoặc tùy chỉnh bổ sung để tạo ra metrics phù hợp với Prometheus. Trong trường hợp tốt nhất, một hệ thống giám sát có mô hình dữ liệu tương tự Prometheus để bạn có thể tự động xác định cách chuyển đổi metrics.Nhiều nhất, chúng tôi cần khả năng cho phép người dùng chọn metrics họ muốn pull.

* Trong các trường hợp khác, metrics từ hệ thống hoàn toàn non-standard, tùy thuộc vào việc sử dụng hệ thống và ứng dụng cơ bản. Trong trường hợp đó, người dùng phải cho biết cách chuyển đổi metrics.

* Khuyến nghị đảm bảo exporter hoạt động ngoài box mà không cần cấu hình và cung cấp lựa chọn các cấu hình ví dụ để chuyển đổi nếu cần.

* YAML là định dạng cấu hình Prometheus tiêu chuẩn, tất cả các cấu hình nên sử dụng YAML theo mặc định.

## Failed scrapes

Hiện tại có hai mẫu cho `failed scrapes` trong đó ứng dụng mà bạn đang nói chuyện không có phản hồi hoặc có vấn đề khác.

* Đầu tiên là trả về lỗi 5xx.

* Thứ hai là có myexporter_up, ví dụ: haproxy_up, biến có giá trị 0 hoặc 1 tùy thuộc vào việc scrape có hoạt động hay không.

Thứ hai tốt hơn khi có vẫn còn một số metrics hữu ích mà bạn có thể nhận được ngay cả với một lần cạo không thành công, chẳng hạn như HAProxy-exporter cung cấp metric thống kê quy trình. Trước đây người dùng dễ dàng xử lý hơn, vì `up` hoạt động theo cách thông thường, mặc dù bạn có thể phân biệt giữa Exporter bị ngừng hoạt động và ứng dụng bị ngừng hoạt động.

## Other tips

* Landing page: Nó tốt hơn cho người dùng nếu truy cập `http://yourexporter/`có một trang HTML đơn giản với tên của Exporters và liên kết đến trang `/metrics`.
* Port numbers: Một người dùng có thể có nhiều Exporters và các thành phần Prometheus trên cùng một máy, vì vậy để dễ dàng hơn, mỗi thành phần có một cổng duy nhất.Đối với Exporters cho các ứng dụng nội bộ nên sử dụng các cổng ngoài phạm vi phân bổ cổng mặc định.
* Không nên đặt timestamp trên metrics bạn đưa ra, hãy để Prometheus thực hiện điều đó.
* Mỗi Exporter nên giám sát chính xác một ứng dụng, tốt nhất là ngồi ngay bên cạnh nó trên cùng một máy.
* Nhiều hệ thống giám sát không có nhãn, thay vào đó làm những việc như `my.class.path.mymetric.labelvalue1.labelvalue2.labelvalue3`.
