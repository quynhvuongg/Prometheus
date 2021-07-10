# Overview
<!-- TOC -->

- [Overview](#overview)
  - [Monitoring - Logging - Tracing](#monitoring---logging---tracing)
    - [Logging](#logging)
    - [Tracing](#tracing)
    - [Monitoring](#monitoring)
  - [Prometheus](#prometheus)
    - [What is Prometheus](#what-is-prometheus)
    - [Feature](#feature)
    - [Architecture](#architecture)
    - [GLOSSARY](#glossary)
    - [When does it fit](#when-does-it-fit)
    - [When does it not fit](#when-does-it-not-fit)

<!-- /TOC -->

## Monitoring - Logging - Tracing

### Logging

Được sử dụng với mục đích ghi lại các  bản ghi chứa các mốc thời gian về những sự kiện và dữ liệu liên quan đến nó .

Logging có 3 dạng bản ghi chính :

- Plaintext : bản ghi là dạng văn bản tự do  (dạng bản ghi được sử dụng phổ biến nhất).
- Structured : thông thường có định dạng JSON.
- Binary: dạng nhị phân.

-> Logging cung cấp cái nhìn sâu sắc có giá trị cùng với context rộng lớn đặc biệt hữu ích cho việc phát hiện các hành vi mới và không thể đoán trước. Nhưng vấn đề: lấy, chuyển, lưu trữ và phân tích logs khá tốn kém -> nó chỉ log những gì cần thiết, có thể lưu được.

### Tracing

Đại diện cho các bản ghi, cấu trúc dữ liệu của trace gần giống như sự kiện log. Một trace có thể cung cấp khả năng hiển thị vào cả đường dẫn (hành trình) đi qua của một yêu cầu và cấu trúc của yêu cầu.

### Monitoring

Việc giám sát bao gồm các hành động: thu thập, tổng hợp và phân tích các số liệu đưa ra các cảnh báo cho các nhà phát triển khi một hệ thống không hoạt động như bình thường.
 
Metric là một biểu diễn số của dữ liệu được đo theo các khoảng thời gian. Các số liệu có thể khai thác sức mạnh của mô hình toán học và dự đoán để rút ra hiểu biết về hành vi của một hệ thống trong các khoảng thời gian trong hiện tại và tương lai.

## [Prometheus](https://prometheus.io/docs/introduction/overview/)

### What is Prometheus

Prometheus là bộ công cụ giám sát và cảnh báo hệ thống 100% mã nguồn mở, ban đầu được xây dựng tại SoundCloud. Kể từ khi thành lập vào năm 2012, nhiều công ty và tổ chức đã áp dụng Prometheus và dự án có cộng đồng người dùng và nhà phát triển rất tích cực. Hiện tại nó là một dự án mã nguồn mở độc lập và được duy trì độc lập với bất kỳ công ty nào.

Kể từ khi bắt đầu với không ít nhà phát triển làm việc trong SoundCloud vào năm 2012, một cộng đồng và hệ sinh thái đã phát triển xung quanh Prometheus. Prometheus chủ yếu được viết bằng ngôn ngữ Golang và được cấp phép theo giấy phép Apache 2.0. Ước tính đến năm 2018, đã có hàng chục nghìn tổ chức đang sử dụng Prometheus. Năm 2016, Prometheus tham gia Cloud Native Computing Foundation (CNCF) với tư cách là dự án thứ hai, sau Kubernetes.

Prometheus là hệ thống giám sát chủ động dựa trên số liệu (metric) được thiết kế để theo dõi liên tục tình trạng, các hoạt động và hiệu suất tổng thể của đối tượng mà nó giám sát. Prometheus có khả năng thu thập số liệu từ các mục tiêu được cấu hình theo các khoảng thời gian nhất định, đánh giá các số liệu đó dựa trên biểu thức quy tắc, hiển thị kết quả và có thể kích hoạt cảnh báo nếu một số điều kiện được thỏa mãn yêu cầu.

### Feature

- Prometheus có mô hình dữ liệu đa chiều với dữ liệu chuỗi thời gian được xác định bởi tên của số liệu và cặp khóa – giá trị.
- PromQL một ngôn ngữ truy vấn linh hoạt mà Prometheus sử dụng để lấy dữ liệu từ Agent phục vụ cho việc giám sát.
- Việc thu thập các số liệu chuỗi thời gian thông qua một mô hình kéo qua HTTP và đẩy chuỗi thời gian được hỗ trợ thông qua một cổng trung gian gọi là Pushgateway.
- Các mục tiêu giám sát được phát hiện thông qua dịch vụ khám phá hoặc cấu hình tĩnh.
- Prometheus cũng hỗ trợ nhiều chế độ biểu đồ khác nhau hỗ trợ vẽ đồ thị và bảng điều khiển, cùng với các chương trình tích hợp và hỗ trợ bởi bên thứ ba. Hoạt động cảnh báo của Prometheus linh động và dễ cấu hình.

### Architecture

![ ](https://prometheus.io/assets/architecture.png)

Prometheus là một hệ sinh thái gồm nhiều thành phần:

- Prometheus server: Thu thập lưu trữ dữ liệu dạng chuỗi thời gian.
- Client library: Được sử dụng như công cụ có thể tích hợp trong mã nguồn ứng dụng. Thư viện sử dụng để định nghĩa thông số hiển thị thông qua HTTP endpoint, và có các ngôn ngữ hỗ trợ như: Go, Java, Python, Ruby.
- Push Gateway: hỗ trợ các nhiệm vụ có thời gian sống ngắn. Các tác vụ công việc này không đủ lâu để Prometheus chủ động lấy dữ liệu vì vậy các số liệu sẽ được đẩy về PushGateway rồi đẩy về máy chủ Prometheus.
- Exporter: Dịch vụ hỗ trợ giám sát các dịch vụ trên hệ thống bằng cách thu thập các số liệu từ các đối tượng giám sát sau đó chuyển đổi về kiểu số liệu Prometheus. Exporter được phát triển bởi Prometheus và cộng đồng sử dụng Prometheus.
- Service discovery: Prometheus đã tích hợp với nhiều cơ chế phát hiện dịch vụ phổ biến như Kubernetes, Consul hoặc được khai báo trong các file. Prometheus cho phép cấu hình siêu dữ liệu từ phát hiện dịch vụ được ánh xạ tới các mục tiêu giám sát và nhãn của chúng bằng cách chỉnh sửa nhãn (relabeling).
- Recording and alerting rules: Các quy tắc ghi cho phép các biểu thức PromQL được đánh giá một cách thường xuyên và kết quả của chúng được đưa vào nơi lưu trữ. Quy tắc cảnh báo là một hình thức khác của quy tắc ghi. Chúng cũng đánh giá các biểu thức PromQL thường xuyên và bất kỳ kết quả nào từ các biểu thức đó đều có thể trở thành cảnh báo. Cảnh báo được gửi đến Alertmanager
- Alertmanager: Dịch vụ xử lý các cảnh báo nhận từ Prometheus rồi gửi thông báo đi tới quản trị viên qua email hay tin nhắn,…
- Dashboard: Prometheus có một số API HTTP cho phép truy vấn tổng hợp số liệu, kết quả có thể được đưa ra dưới dạng bảng hoặc đồ thị. Tuy nhiên nó không phải là một hệ thống trực quan nên tích hợp với Grafana cho ra bảng điều khiển trực quan tốt để dễ dàng theo dõi hơn.
- Storage: Prometheus ngoài lưu trữ cục bộ các dữ liệu được thu thập trên các đĩa (SSD/HDD) mà còn hỗ trợ lưu trữ dữ liệu trên các kho dữ liệu từ xa.

Prometheus thực hiện lấy các số liệu từ các mục tiêu được chỉ định qua thành phần phát hiện dịch vụ. Cụ thể ở đây Prometheus định kỳ thu thập các số liệu đã được chuyển đổi về định dạng Prometheus có thể hiểu được từ các Exporter, sau đó sẽ lưu trữ các dữ liệu thu thập được ở máy chủ cục bộ. Tiếp đến Prometheus sẽ chạy các quy tắc để đánh giá số liệu được thu thập về, nếu nó vi phạm một quy tắc nào đó sẽ thực hiện đẩy cảnh báo tới Alertmanager để xử lý đưa tới người quản trị.

### [GLOSSARY](https://prometheus.io/docs/introduction/glossary/)

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

### When does it fit

Prometheus hoạt động tốt để ghi lại bất kỳ chuỗi thời gian hoàn toàn bằng số. Nó phù hợp với cả giám sát tập trung vào máy cũng như giám sát các kiến trúc hướng dịch vụ rất năng động. Trong một thế giới của microservice, sự hỗ trợ của nó cho việc thu thập và truy vấn dữ liệu đa chiều là một thế mạnh đặc biệt.

Prometheus được thiết kế để đảm bảo độ tin cậy, là hệ thống bạn sử dụng trong thời gian ngừng hoạt động để cho phép bạn chẩn đoán nhanh các sự cố. Mỗi máy chủ Prometheus là độc lập, không phụ thuộc vào lưu trữ mạng hoặc các dịch vụ từ xa khác. Bạn có thể dựa vào nó khi các phần khác trong cơ sở hạ tầng của bạn bị hỏng và bạn không cần thiết lập cơ sở hạ tầng mở rộng để sử dụng nó.

### When does it not fit

Prometheus values reliability. Bạn luôn có thể xem những số liệu thống kê có sẵn về hệ thống của bạn, ngay cả trong điều kiện thất bại. Nếu bạn cần độ chính xác 100%, chẳng hạn như thanh toán theo yêu cầu, Prometheus không phải là một lựa chọn tốt vì dữ liệu được thu thập có thể sẽ không được chi tiết và đầy đủ. Trong trường hợp như vậy, tốt nhất bạn nên sử dụng một số hệ thống khác để thu thập và phân tích dữ liệu để thanh toán và Prometheus cho phần còn lại của việc theo dõi của bạn.
