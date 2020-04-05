
# Monitoring - Logging - Tracing

## Logging

Được sử dụng với mục đích ghi lại các  bản ghi chứa các mốc thời gian về những sự kiện và dữ liệu liên quan đến nó .
Logging có 3 dạng bản ghi chính :

- Plaintext : bản ghi là dạng văn bản tự do  ( dạng bản ghi được sử dụng phổ biến nhất ).
- Structured : thông thường được định dạng JSON.
- Binary
Logging cung cấp cái nhìn sâu sắc có giá trị cùng với bối cảnh rộng lớn  đặc biệt hữu ích cho việc phát hiện các hành vi mới nổi và không thể đoán trước . Nhưng vấn đề có được, chuyển, lưu trữ và phân tích logs khá tốn kém vì đó nó chỉ log những gì cần thiết, có thể lưu được.

## Tracing

Đại diện cho các bản ghi, cấu trúc dữ liệu của trace gần giống như sự kiện log. Một trace có thể cung cấp khả năng hiển thị vào cả đường dẫn (hành trình ) đi qua của một yêu cầu cũng như cấu trúc của yêu cầu.

## Monitoring

Việc giám sát thu thập, tổng hợp và phân tích các số liệu đưa ra các cảnh báo cho các nhà phát triển khi một hệ thống không hoạt động như bình thường.  
Metric là một biểu diễn số của dữ liệu được đo theo các khoảng thời gian. Các số liệu có thể khai thác sức mạnh của mô hình toán học và dự đoán để rút ra hiểu biết về hành vi của một hệ thống trong các khoảng thời gian trong hiện tại và tương lai.

# Prometheus

## What is Prometheus

Prometheus là bộ công cụ giám sát và cảnh báo hệ thống nguồn mở ,ban đầu được xây dựng tại SoundCloud. Kể từ khi thành lập vào năm 2012, nhiều công ty và tổ chức đã áp dụng Prometheus và dự án có một cộng đồng người dùng và nhà phát triển rất tích cực.

Bây giờ nó là một dự án nguồn mở độc lập và được duy trì độc lập với bất kỳ công ty nào. Prometheus đã tham gia Cloud Native Computing Foundation vào năm 2016 với tư cách là dự án thứ hai, sau Kubernetes.

Prometheus có khả năng thu thâp thông số/số liệu (metric) từ các mục tiêu được cấu hình theo các khoảng thời gian nhất định , đánh giá các biểu thức quy tắc, hiện thị kết quả và có thể kích hoạt cảnh báo nếu một số điều kiện được thỏa mãn yêu cầu.

## Feature

- Mô hình dữ liệu đa chiều với dữ liệu chuỗi thời gian được xác định bởi tên của metric và cặp key -value .
- PromQL, một ngôn ngữ truy vấn linh hoạt .
- Không phụ thuộc vào lưu trữ phân tán; các node máy chủ đơn là tự chủ.
- Việc thu thập chuỗi thời gian xảy ra thông qua một mô hình kéo qua HTTP.
- Đẩy chuỗi thời gian được hỗ trợ thông qua một cổng trung gian.
- Các Target được phát hiện thông qua service discovery hoặc cấu hình tĩnh.
- Nhiều chế độ hỗ trợ vẽ đồ thị và bảng điều khiển.

## Characteristics

- Prometheus 100% là mã nguồn mở (Github: <https://github.com/prometheus/prometheus).>
- Phần lớn các core tính năng của Prometheus được biết bằng ngôn ngữ Go , một số còn lại được viết bằng Java , Python, Ruby .
- Prometheus không phải dùng để lấy dữ liệu logs , thay vì vậy nó là dịch vụ giám sát , thu thập và xử lý dữ liệu dạng metric .
- Prometheus sử dụng cơ chế pull dữ liệu từ remote là chính , chứ không sử dụng cơ chế đợi remote push dữ liệu lên ngoại trừ trường hợp sử dụng PushGateway .
- Prometheus sử dụng chương trình cảnh báo Alertmanager để xử lý và gửi cảnh báo đi .
- Về phần giao diện biểu đồ nó sử dụng mã nguồn Grafana để tích hợp hiển thị .
- Metric của Prometheus sử dụng chuẩn OpenMetrics.

## GLOSSARY

- **Alert** : một cảnh báo là kết quả của việc đạt điều kiện thoả mãn một rule cảnh báo được cấu hình trong Prometheus. Các cảnh báo sẽ được gửi đến dịch vụ Alertmanager.
- **Alertmanager**: chương trình đảm nhận nhiệm vụ tiếp nhận, xử lý các hoạt động cảnh báo.
- **Bridge** : là một thành phần lấy sample từ thư viện khách và đưa chúng ra hệ thống giám sát không phải Prometheus. Ví dụ: các clients Python, Go và Java có thể xuất các số liệu sang Graphite.
- **Client Library**: một số thư viện hỗ trợ người dùng có thể tự tuỳ chỉnh lập trình phương thức riêng để lấy dữ liệu từ hệ thống và đẩy dữ liệu metric về Prometheus.
- **Endpoint**: nguồn dữ liệu của các chỉ số (metric) mà Prometheus sẽ đi lấy thông tin.
- **Exporter**:  là một chương trình được sử dụng với mục đích thu thập, chuyển đổi các metric không ở dạng kiểu dữ liệu chuẩn Prometheus sang chuẩn dữ liệu Prometheus. Sau đấy exporter sẽ expose web service api chứa thông tin các metrics hoặc đẩy về Prometheus.
- **Collector** : là một phần của Exporter đại diện cho một set metric. Nó có thể là một số liệu đơn nếu nó là một phần của thiết bị đo trực tiếp hoặc nhiều số liệu nếu nó lấy số liệu từ hệ thống khác.
- **Instance**: là một nhãn (label) dùng để định danh duy nhất cho một target trong một job .
- **Job**: là một tập hợp các target chung một nhóm mục đích. Ví dụ: giám sát một nhóm các dịch vụ database,… thì ta gọi đó là một job .
- **PromQL**: promql là viết tắt của Prometheus Query Language, ngôn ngữ này cho phép bạn thực hiện các hoạt động liên quan đến dữ liệu metric.
- **Sample**: là một giá trị đơn lẻ tại một thời điểm thời gian trong time series.
- **Target**: một target là định nghĩa một đối tượng sẽ được Prometheus đi lấy dữ liệu (scrape). Ví dụ như: nhãn nào sẽ được sử dụng cho đối tượng, hình thức chứng thực nào sử dụng hoặc các thông tin cần thiết để quá trình đi lấy dữ liệu ở đối tượng được diễn ra.
- **Notification** : thông báo đại diện cho một nhóm  của một hoặc nhiều cảnh báo và được Alertmanager gửi đến email, Pagerduty, Slack, v.v.
- **Remote read** : Đọc từ xa là một tính năng Prometheus cho phép đọc chuỗi thời gian từ các hệ thống khác (chẳng hạn như lưu trữ dài hạn) như một phần của truy vấn.
- **Remote read adapter** : Không phải tất cả các hệ thống hỗ trợ trực tiếp đọc từ xa. Một bộ chuyển đổi đọc từ xa nằm giữa Prometheus và một hệ thống khác, chuyển đổi các yêu cầu và phản hồi chuỗi thời gian giữa chúng.
- **Remote write** : là một tính năng Prometheus cho phép gửi các sample đã   khi đang di chuyển đến các hệ thống khác, chẳng hạn như lưu trữ dài hạn.
- **Remote write adapter** : Không phải tất cả các hệ thống hỗ trợ trực tiếp ghi từ xa. Một bộ chuyển đổi ghi từ xa nằm giữa Prometheus và một hệ thống khác, chuyển đổi các mẫu trong ghi từ xa thành định dạng mà hệ thống khác có thể hiểu được.
