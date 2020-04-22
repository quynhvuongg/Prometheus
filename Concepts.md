
# Concepts

**_Data Model_**

Prometheus về cơ bản lưu trữ tất cả dữ liệu theo time series.
Time series là sự kết hợp của name metric và các cặp key-value được gọi là labels. Kết hợp của các lables cho cùng một name metric khởi tạo nên một dimention cụ thể cho metric đó . Label bắt đầu bằng _ được dành riêng cho sử dụng nội bộ.

![ ](https://image.slidesharecdn.com/copyofprometheusstorage1-160127133731/95/prometheus-storage-4-638.jpg?cb=1453901940)

*_Metric names and labels_*

- Name metric: chỉ ra đặc tính chung của hệ thống cần đo . Ví dụ:  `http_requests_total`  => tổng số HTTP request nhận được.

- PromQL sẽ dựa trên dimension để lọc và tổng hợp dữ liệu. Ví dụ: với metric  `http_requests_total`, chúng ta có thể có 2 label là  `method="POST"`,  `status="200"`

*_Sample_*

Các samples hình thành lên time series thực tế ,mỗi sample bao gồm :

- một  giá trị float64
- một timestamp chính xác đến mili giây.

*_Notation_*

```sh
<metric name>{<label name>=<label value>, ...}
```

**_Metric type_**

Prometheus hỗ trợ 4 metric type là:

1. Counter :
 Metric tích lũy đại diện cho một bộ đếm chỉ có thể tăng dần hoặc được đặt lại về 0 khi khởi động lại. Ví dụ: bạn có thể sử dụng bộ đếm để thể hiện số lượng yêu cầu được phục vụ, nhiệm vụ đã hoàn thành hoặc lỗi.
Không sử dụng bộ đếm với giá trị có thể giảm. Ví dụ: không sử dụng bộ đếm cho số lượng quy trình hiện đang chạy nên sử dụng Gauge.

2. Gauge :
Metric đại diện cho một giá trị số có thể tùy ý lên xuống, thường được sử dụng cho các giá trị đo như nhiệt độ hoặc mức sử dụng bộ nhớ hiện tại, và "counter" có thể tăng và giảm, như số lượng yêu cầu đồng thời.

3. Histogram :

    Loại số liệu biểu đồ đo tần số của các quan sát giá trị rơi vào các buckets cụ thể.
    Ví dụ: Ta có thể đo thời lượng yêu cầu cho một yêu cầu HTTP cụ thể bằng biểu đồ. Thay vì lưu trữ mọi thời lượng cho mỗi yêu cầu, Prometheus sẽ thực hiện xấp xỉ bằng cách lưu trữ tần suất của các yêu cầu rơi vào các bucket cụ thể.
    Bucket : là counter của các quan sát. Nó nên có giới hạn trên và giới hạn dưới.

- `<basename> _bucket {le = "< upper includesive bound >"}` : bộ đếm tích lũy cho các bucket quan sát .
- `<basename> _sum`: tổng cộng của tất cả các giá trị được quan sát.
- `<basename> _count` : số lượng các sự kiện đã được quan sát .( = `<basename>_bucket{le="+Inf"}`)

    Các trường hợp sử dụng Historgram:
- Khi muốn thực hiện nhiều phép đo của một đối tượng, để sau đó tính trung bình hoặc phần trăm .
- Khi không bận tâm về các giá trị chính xác,có thể sử dụng giá trị xấp xỉ .
- Khi biết phạm vi của các giá trị

    Một số trường hợp sử dụng :
- Thời hạn yêu cầu
- Kích thước phản hồi

 Sử dụng hàm histogram_quantile () để tính toán lượng tử từ biểu đồ hoặc thậm chí tổng hợp biểu đồ.

A Histogram looks like:

```sh
    request_latency_seconds_bucket{le="0.025"} 4.0
    request_latency_seconds_bucket{le="0.05"} 7.0
    request_latency_seconds_bucket{le="0.075"} 10.0
    request_latency_seconds_bucket{le="1.0"} 10.0
    request_latency_seconds_bucket{le="+Inf"} 11.0
    request_latency_seconds_count 11.0
    request_latency_seconds_sum 3.3

```

1. Summary :

    Một phần mở rộng của Histogram. Bên cạnh cũng cung cấp tổng và số lượng quan sát , nó tính toán các lượng tử có thể cấu hình qua sliding time window.

    Một Summary hiển thị nhiều chuỗi thời gian trong một lần quét:

- `<basename> {quantile = "<φ>"}`:  φ-quantiles (0 ≤ φ≤ 1) các sự kiện được quan sát.
- `<basename> _sum`: tổng cộng của tất cả các giá trị được quan sát.
- `<basename> _count`: số lượng các sự kiện đã được quan sát .

    Các trường hợp sử dụng Summary cũng giống như Historgram .

    A Summary looks like:

    ```sh
    go_gc_duration_seconds{quantile="0"} 0.000236554
    go_gc_duration_seconds{quantile="0.25"} 0.000474629
    go_gc_duration_seconds{quantile="0.5"} 0.0005691670000000001
    go_gc_duration_seconds{quantile="0.75"} 0.000677597
    go_gc_duration_seconds{quantile="1"} 0.002479919
    go_gc_duration_seconds_sum 12.532527861
    go_gc_duration_seconds_count 24279
    ```

    Một số khác biệt giữa Summary và Historgram:

- Với Historgram lượng tử được tính trên máy chủ Prometheus, với Summary chúng được tính trên máy chủ ứng dụng . Do đó, dữ liệu Summary không thể được tổng hợp từ một số trường hợp ứng dụng.
- Historgram yêu cầu định nghĩa giá trị trước của bucket , vì vậy phù hợp với trường hợp sử dụng khi có ý tưởng tốt về sự tràn dải của các giá trị.
- Summary là một lựa chọn tốt nếu cần tính toán các lượng tử chính xác, nhưng không thể chắc chắn phạm vi của các giá trị sẽ là bao nhiêu.

**_Job và Instance_**

- Instance: một endpoint mà Prometheus có thể thu thập dữ liệu metrics, tương ứng với một single process.
- Job: tập hợp các instance có chung mục đích.

_Ví dụ:_

- job:  `api-server`
  - instance 1:  `1.2.3.4:5670`
  - instance 2:  `1.2.3.4:5671`
  - instance 3:  `5.6.7.8:5670`
  - instance 4:  `5.6.7.8:5671`

# Architecture

![ ](https://prometheus.io/assets/architecture.png)

Prometheus là một hệ sinh thái gồm nhiều thành phần:

- Prometheus server : thu thập và lưu trữ dữ liệu dưới dạng time series.
- Client libraries : đây là API tương tác với Prometheus của các ngôn ngữ lập trình khác như Go, Java/Scala, Ruby, Python, ...
- Push Gateway: hỗ trợ các job có thời gian sống ngắn. Các tác vụ công việc này không đủ lâu để Prometheus chủ động lấy dữ liệu vì vậy các metric sẽ được đẩy về PushGateway rồi đẩy về Prometheus Server.
- Các exporter  như HAProxy, StatsD, Graphite, ... : hỗ trợ giám sát các dịch vụ hệ thống và gửi về Prometheus .
- Alertmanager : xử lý cảnh báo.
- Các công cụ hỗ trợ khác .

-> Prometheus thực hiện lấy các metric từ các job được chỉ định qua kênh trực tiếp hoặc thông qua dịch vụ Gateway chung . Sau đó Promethus sẽ lưu trữ các dữ liệu thu thấp được ở local máy chủ. Tiếp đến sẽ chạy các rules để xử lý dữ liệu theo nhu cầu cũng như kiểm tra thực hiện các cảnh báo mà bạn mong muốn .

## When does it fit

Prometheus hoạt động tốt để ghi lại bất kỳ chuỗi thời gian hoàn toàn bằng số. Nó phù hợp với cả giám sát tập trung vào máy cũng như giám sát các kiến trúc hướng dịch vụ rất năng động. Trong một thế giới của microservice, sự hỗ trợ của nó cho việc thu thập và truy vấn dữ liệu đa chiều là một thế mạnh đặc biệt.

Prometheus được thiết kế để đảm bảo độ tin cậy, là hệ thống bạn sử dụng trong thời gian ngừng hoạt động để cho phép bạn chẩn đoán nhanh các sự cố. Mỗi máy chủ Prometheus là độc lập, không phụ thuộc vào lưu trữ mạng hoặc các dịch vụ từ xa khác. Bạn có thể dựa vào nó khi các phần khác trong cơ sở hạ tầng của bạn bị hỏng và bạn không cần thiết lập cơ sở hạ tầng mở rộng để sử dụng nó.

## When does it not fit

Prometheus values reliability. Bạn luôn có thể xem những số liệu thống kê có sẵn về hệ thống của bạn, ngay cả trong điều kiện thất bại. Nếu bạn cần độ chính xác 100%, chẳng hạn như thanh toán theo yêu cầu, Prometheus không phải là một lựa chọn tốt vì dữ liệu được thu thập có thể sẽ không được chi tiết và đầy đủ. Trong trường hợp như vậy, tốt nhất bạn nên sử dụng một số hệ thống khác để thu thập và phân tích dữ liệu để thanh toán và Prometheus cho phần còn lại của việc theo dõi của bạn.
