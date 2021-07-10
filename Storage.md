# Storage
<!-- TOC -->

- [Storage](#storage)
  - [Local storage](#local-storage)
    - [On-disk layout](#on-disk-layout)
    - [Compaction](#compaction)
    - [Operational aspects (Sizing calculation)](#operational-aspects-sizing-calculation)
    - [Remote storage](#remote-storage)
    - [Compare local and remote storage](#compare-local-and-remote-storage)

<!-- /TOC -->

Prometheus bao gồm cơ sở dữ liệu time series được lưu trữ trên local disk, nhưng cũng tùy ý tích hợp với các hệ thống lưu trữ từ xa.

![ ](https://prometheus.io/assets/architecture.png)

## Local storage

Cơ sở dữ liệu các timeseries được Prometheus lưu trữ tại `local` theo định dạng tùy chỉnh trên các ổ đĩa (HDD/SSD).

### On-disk layout

Các mẫu được nhóm vào thành `blocks` trong hai giờ. Mỗi khối này bao gồm một thư mục chứa một hoặc nhiều tệp `chunk` chứa tất cả các mẫu timeseries cho cửa sổ thời gian đó, cũng như `metadata file` và `index file` (trỏ các `metric names` và `labels` theo chuỗi thời gian trong các tệp `chunk`). Khi chuỗi được xóa thông qua API, các bản ghi xóa sẽ được lưu trữ trong các `tombstone file` riêng biệt (thay vì xóa dữ liệu ngay lập tức khỏi các tệp `chunk`).

Khối cho các mẫu đang đến được giữ trong bộ nhớ và chưa hoàn toàn tồn tại. Nó được bảo vệ chống lại các sự cố bởi bản ghi `write-ahead-log` (WAL), có thể được phát lại khi máy chủ Prometheus khởi động lại sau một sự cố. Các tệp log này được lưu trữ trong thư mục `wal`  có kích thước `128MB` chứa dữ liệu thô chưa được nén, vì vậy chúng có kích thước lớn hơn đáng kể so với các khối thông thường. Prometheus sẽ giữ tối thiểu 3 tệp `wal`, tuy nhiên các máy chủ với lưu lượng truy cập cao, có thể thấy nhiều hơn ba tệp `WAL` vì nó cần giữ dữ liệu thô ít nhất hai giờ.

Cấu trúc của thư mục dữ liệu của máy chủ Prometheus theo TSDB format :

* Index
* Chunks
* Tombstones
* Wal

```sh
./data
├── 01BKGV7JBM69T2G1BGBGM6KB12
│   └── meta.json
├── 01BKGTZQ1SYQJTR4PB43C8PD98
│   ├── chunks
│   │   └── 000001
│   ├── tombstones
│   ├── index
│   └── meta.json
├── 01BKGTZQ1HHWHV8FBJXW1Y3W0K
│   └── meta.json
├── 01BKGV7JC0RY8A6MACW02A2PJD
│   ├── chunks
│   │   └── 000001
│   ├── tombstones
│   ├── index
│   └── meta.json
└── wal
    ├── 00000002
    └── checkpoint.000001

```

Một hạn chế của bộ nhớ cục bộ là nó không được tùy ý gộp hoặc mở rộng do đó nó không đảm bảo được khả năng bền vững khi bị ảnh hưởng bởi sự cố mất disk hay nút.
Sử dụng RAID cho tính khả dụng của disk, `snapshots` để sao lưu, dự tính dung lượng, v.v., được đề xuất để cải thiện độ bền. Với độ bền lưu trữ thích hợp và lập kế hoạch lưu trữ nhiều năm dữ liệu trong `local storage` là có thể.

Ngoài ra, bộ nhớ ngoài có thể được sử dụng thông qua API đọc/ghi từ xa.

### Compaction

* Các `two-hour blocks` ban đầu được nén thành các khối dài hơn trong background.

* Việc nén sẽ tạo ra các khối lớn hơn tới 10% thời gian lưu, hoặc 31 ngày, tùy theo mức nào nhỏ hơn.

### Operational aspects (Sizing calculation)

Prometheus có một số cờ cho phép cấu hình bộ nhớ cục bộ:

* `--storage.tsdb.path`: Xác định nơi Prometheus ghi cơ sở dữ liệu của nó ( default: `data/`).
  
* `--storage.tsdb.retention.time`: Xác định khi nào cần xóa dữ liệu cũ ( default: 15 ngày ).

* `--storage.tsdb.retention.size`: [EXPERIMENTAL], xác định số byte tối đa mà các khối lưu trữ có thể sử dụng (lưu ý rằng điều này không bao gồm kích thước WAL). Dữ liệu cũ nhất sẽ bị xóa đầu tiên. Mặc định là 0 hoặc disabled. Cờ này là thử nghiệm và có thể được thay đổi trong phiên bản tương lai. Các đơn vị được hỗ trợ: KB, MB, GB, PB.
Ex: `512MB`

* `--storage.tsdb.wal-compression`: Cho phép nén `log WAL`. Tùy thuộc vào dữ liệu, bạn có thể để kích thước `WAL` được giảm một nửa. Lưu ý rằng nếu bật cờ này và sau đó hạ cấp Prometheus xuống phiên bản bên dưới 2.11.0, bạn sẽ cần xóa WAL vì nó sẽ không thể đọc được.

Trung bình, Prometheus chỉ sử dụng khoảng 1-2 bytes cho mỗi mẫu. Do đó, để dự tính dung lượng của máy chủ Prometheus, bạn có thể sử dụng công thức:

```sh
needed_disk_space = retention_time_seconds * ingested_samples_per_second * bytes_per_sample

```

Để điều chỉnh tốc độ lấy mẫu trong một giây, có thể giảm số chuỗi thời gian scrape ( ít target hơn hoặc ít chuỗi hơn trên mỗi target ) hoặc bạn có thể tăng khoảng thời gian scrape. Tuy nhiên, việc giảm số lượng series có thể hiệu quả hơn, do nén các mẫu trong một chuỗi.

Nếu `local storage` bị hỏng vì bất kỳ lý do gì, cách tốt nhất là tắt Prometheus và xóa toàn bộ thư mục lưu trữ. Các hệ thống tệp không tuân thủ `POSIX`, không được hỗ trợ bởi `local storage` của Prometheus, các lỗi có thể xảy ra, mà không có khả năng phục hồi. Bạn có thể thử xóa các thư mục khối riêng lẻ để giải quyết vấn đề, điều này có nghĩa là mất một cửa sổ thời gian có giá trị khoảng hai giờ cho mỗi thư mục khối.
Chú ý :  `local storage` của Prometheus không có nghĩa là lưu trữ lâu dài.

Nếu cả chính sách duy trì time và size được chỉ định, bất kỳ chính sách nào được kích hoạt, nó sẽ được sử dụng ngay lúc đó.

Dọn dẹp `expired block` xảy ra trên một nền lịch trình. Có thể mất đến hai giờ để loại bỏ các `expired block`. Chúng phải được hết hạn trước khi chúng được dọn.

### Remote storage

Lưu trữ cục bộ của Prometheus bị giới hạn bởi các nút đơn lẻ về khả năng mở rộng và độ bền. Thay vì cố gắng giải quyết lưu trữ phân cụm trong chính Prometheus, Prometheus có một bộ giao diện cho phép tích hợp với các hệ thống lưu trữ từ xa.

![ ](https://prometheus.io/docs/prometheus/latest/images/remote_integrations.png)

Prometheus tích hợp với các hệ thống lưu trữ từ xa theo hai cách:

* Prometheus có thể ghi các mẫu vào một URL từ xa theo định dạng chuẩn.

* Prometheus có thể đọc (back) dữ liệu mẫu từ một URL từ xa ở định dạng chuẩn.

Cả hai giao thức đọc và ghi đều sử dụng mã hóa bộ đệm giao thức được nén trên HTTP. Các giao thức chưa được coi là API ổn định và có thể thay đổi để sử dụng gRPC qua HTTP / 2 trong tương lai, khi tất cả các bước nhảy giữa Prometheus và bộ lưu trữ từ xa có thể được giả định một cách an toàn để hỗ trợ HTTP / 2.

**_Remote read_**

![ ](https://miro.medium.com/max/1400/0*EE_EK8eHuFQNX9jj.png)

Khi cấu hình, Prometheus lưu trữ các truy vấn ( ví dụ qua HTTP API ) được gửi đến cả `local` và `remote storage` và  kết quả được đồng nhất.

Lưu ý rằng để duy trì độ tin cậy khi đối mặt với các vấn đề lưu trữ từ xa, việc đánh giá những quy tắc alert và record chỉ sử dụng TSDB cục bộ.

**`<remote_read>`**

```yml
# URL của các endpoint để truy vấn.
url: <string>

# Một danh sách tùy chọn các công cụ đối sánh bằng
# phải có trong bộ chọn để truy vấn điểm cuối đọc từ xa.
required_matchers:
   <labelname>: <labelvalue>

# Thời gian chờ cho yêu cầu
  remote_timeout: <duration> | default = 1m

# Việc đọc có nên được thực hiện cho các truy vấn
# trong khoảng thời gian mà bộ nhớ cục bộ nên có dữ liệu đầy đủ.
  read_recent: <boolean> | default = false

# Đặt tiêu đề `Authorization` trên mỗi yêu cầu đọc từ xa với các cách sau:

# 1.Tên người dùng và mật khẩu
basic_auth:
  username: <string>
  password: <string>
  password_file: <string>

# 2. Mã thông báo
  bearer_token: <string>

# 3. Mã thông báo từ file cấu hình
  bearer_token_file: /path/to/bearer/token/file

# Cấu hình cài đặt yêu cầu TLS của remote read.
tls_config:
  <tls_config>

# Optional proxy URL.
proxy_url: <string>

```

**_Remote write_**

![ ](https://miro.medium.com/max/1400/1*tPqV7fAbgb_5Ueixx2Js8A.png)

Ghi từ xa hoạt động bằng cách "theo đuôi" các mẫu chuỗi thời gian được ghi vào bộ lưu trữ cục bộ và xếp hàng chúng để ghi vào bộ lưu trữ từ xa.

Hàng đợi là tập hợp các "shards" được quản lý động: tất cả các mẫu cho bất kỳ chuỗi thời gian cụ thể nào (tức là số liệu duy nhất) sẽ kết thúc trên cùng một shard.

Hàng đợi tự động chia tỷ lệ lên hoặc xuống số lượng phân đoạn ghi vào bộ lưu trữ từ xa để theo kịp tốc độ dữ liệu đến.

Điều này cho phép Prometheus quản lý lưu trữ từ xa trong khi chỉ sử dụng các tài nguyên cần thiết và với cấu hình tối thiểu.

**`<remote_write>`**

```yml
url: <string>

[ remote_timeout: <duration> | default = 30s ]

# Relabeling hoặc hạn chế số lượng số liệu ghi
write_relabel_configs:
  [ - <relabel_config> ... ]

basic_auth:
  [ username: <string> ]
  [ password: <string> ]
  [ password_file: <string> ]

[ bearer_token: <string> ]

[ bearer_token_file: /path/to/bearer/token/file ]

tls_config:
  [ <tls_config> ]

[ proxy_url: <string> ]

# Configures the queue used to write to remote storage.
queue_config:

  # Số lượng mẫu để đệm trên mỗi phân đoạn
  #trước khi bị chặn việc đọc thêm các mẫu từ WAL.
  [ capacity: <int> | default = 500 ]
  
  [ max_shards: <int> | default = 1000 ]
  
  [ min_shards: <int> | default = 1 ]

  [ max_samples_per_send: <int> | default = 100]
  
  [ batch_send_deadline: <duration> | default = 5s ]
  
  [ min_backoff: <duration> | default = 30ms ]
  
  [ max_backoff: <duration> | default = 100ms ]
  
```

### Compare local and remote storage

**_Local Storage_**

* Dữ liệu được lưu theo định dạng TSDB trên disk (SSD/HDD)
* Không có khả năng mở rộng và độ bền kém
* Không lưu trữ lâu dài (mặc định 15d)
* Không có tính sẵn sàng cao
* Độ khả dụng phụ thuộc và disk
* Có thể sử dụng PromQL để truy vấn và đánh giá các rule đối với dữ liệu trên local

**_Remote Storage_**

* Dữ liệu được ghi thông qua các bộ điều hợp lưu trữ và Prometheus không kiểm soát định dạng lưu trữ dữ liệu,được tùy chỉnh
* Có khả năng mở rộng linh hoạt và lưu trữ lâu dài
* Độ khả dụng lớn
* Chỉ có thể sử dụng Prometheus để truy vấn , không được sử dụng dữ liệu trên lưu trữ remote để đánh giá các rule
* Tất cả các đánh giá PromQL trên dữ liệu thô vẫn xảy ra trong chính Prometheus. Điều này có nghĩa là các truy vấn đọc từ xa có một số giới hạn về khả năng mở rộng, vì tất cả dữ liệu cần thiết cần được tải vào máy chủ Prometheus truy vấn trước rồi xử lý ở đó.
