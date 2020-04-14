# Alertmanager

Alertmanager nhận tất cả các cảnh báo từ máy chủ Prometheus và chuyển đổi chúng thành các thông báo như emails, messages và pages.

![ ](https://www.oreilly.com/library/view/prometheus-up/9781492034131/assets/prur_1801.png)

## Notification Pipeline

Alertmanager làm nhiều việc cho ta hơn là chỉ chuyển đổi luôn cảnh báo thành thông báo. Alertmanager cung cấp một pipeline kiểm soát để biết cách xử lý cảnh báo trước khi chúng trở thành thông báo. Giống như labels là cốt lõi của chính Prometheus, labels cũng là chìa khóa cho Alertmanager:

![ ](https://user-images.githubusercontent.com/2590559/40104670-cf36c10a-58a5-11e8-909f-133837dd57ac.png)

**_Grouping_**
  
Nhóm các cảnh báo rất hữu ích khi nhiều hệ thống bị lỗi cùng một lúc và hàng trăm đến hàng nghìn cảnh báo có thể được kích hoạt đồng thời.

Ví dụ: Ta cấu hình 100 servers ,khi bị failed thì sẽ gửi cảnh báo đến sysadmin. Khi đó, sysadmin sẽ lập tức nhận 100 notifications một lúc. Thay vì vậy, ta gom nhóm 100 servers này vào 1 group, và sysadmin sẽ chỉ nhận được 1 notification mà thôi. Notification này sẽ chứa đầy đủ thông báo của 100 servers này.

```sh
group_by: ['alertname', 'cluster', 'service']

```

Nhóm các cảnh báo, thời gian cho các thông báo được nhóm và receivers nhận các thông báo đó được cấu hình bởi route tree trong tệp cấu hình.

**_Inhibition_**

Ngăn cản một số cảnh báo thành thông báo khi có một cảnh báo khác nghiêm trọng hơn đang được kích hoạt .

**_Silences_**

Cho phép bỏ qua một số cảnh báo nhất định trong một thời gian. Silences được cấu hình trong giao diện web Alertmanager.

**_Routing_**

Alertmanager sẽ không chuyển tất cả thông báo đến cùng một nơi, mà nó cho phép cấu hình routing tree để xác định điểm cần đến cho thông báo .
Nếu notification trùng với match của route nào đó, thì sẽ được gửi đi theo đường đó. Còn không match với route nào, nó sẽ được gửi theo đường đi mặc định.

**_Deduplication_**

Alertmanager sẽ điều tiết thông báo cho một nhóm nhất định để bạn không bị spam bởi các thông báo giống nhau bằng cách sao chép dữ liệu, loại bỏ những bản sao lặp lại.

**_Retry_**

Trong điều kiện lý tưởng các thông báo được xử lý kịp thời, tuy nhiên thực tế có thể có những ảnh hưởng đến hệ thống mà bị lạc mất. Do đó Alertmanager sẽ lặp lại thông báo để chúng không bị lạc quá lâu.

**_Receiver_**

Cấu hình thông tin các nơi nhận. Ví dụ như tên đăng nhập, mật khẩu,...

**_Notification_**

Cuối cùng chúng đến giai đoạn được gửi đi dưới dạng thông báo qua receiver. Thông báo được tạo theo khuôn mẫu, cho phép tùy chỉnh nội dung của chúng và nhấn mạnh các chi tiết quan trọng .

## Configuration

Alertmanager được cấu hình thông qua command-line flags và configuration file. Tương tự như Prometheus, Alertmanager có thể reload cấu hình khi chạy. Nếu cấu hình mới không được định dạng tốt, các thay đổi sẽ không được áp dụng và lỗi được ghi lại.  
Sử dụng : "--config.file" để chỉ định configuration file.

**_Install Alertmanager_**

```yml
#docker-compose.yml
version : '3.7'
services:  
  
  prometheus:
  ...
  alertmanager:
    image: prom/alertmanager
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml
    command:
      - "--config.file=/etc/alertmanager/alertmanager.yml"
    ports:
      - "9093:9093"

```

```yml
#prometheus.yml
alerting:
  alertmanagers:
  - scheme: http
  static_configs:
  - targets:
    - "alertmanager:9093"

```

**_Config Alertmanager_**

```yml
#alertmanager.yml
global:
  resolve_timeout: 2m

route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 2m
  repeat_interval: 3h
  receiver: 'default'

  routes:
  - match:
      severity: warning
    receiver: 'gmail'
  - match:
      severity: critical
    receiver: 'slack'
  - match_re:
      service: ^(foo1|foo2|baz)$
    receiver: 'slack'


inhibit_rules:
- source_match:
    severity: 'critical'
  target_match:
    severity: 'warning'
  equal: ['alertname']


receivers:
- name: 'default'
  slack_configs:
  - send_resolved: true
    channel: '#default'
    api_url: 'https://hooks.slack.com/services/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'

- name: 'gmail'
  email_configs:
  - to: quynh.vuongg@gmail.com
    from: quynh.vuongg@gmail.com
    smarthost: smtp.gmail.com:587
    auth_username: quynh.vuongg@gmail.com
    auth_identity: quynh.vuongg@gmail.com
    auth_password: <pass>
    send_resolved: true

- name: 'slack'
  slack_configs:
    - send_resolved: true
      channel: '#general'
      api_url: 'https://hooks.slack.com/services/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'

```

**global**: nơi khai báo giá trị các biến toàn cục được sử dụng trong file cấu hình này.

Ví dụ:

```yml
global:
  # The smarthost and SMTP sender used for mail notifications.
  smtp_smarthost: 'localhost:25'
  smtp_from: 'alertmanager@example.org'
  smtp_auth_username: 'alertmanager'
  smtp_auth_password: 'password'
```

-> Những thông tin cần thiết để có thể gửi một cảnh báo đến hòm thư mail.

**route**: nơi cấu hình các thông tin đường đi mặc định. Các cảnh báo sẽ được gưỉ  mặc định theo đường này nếu không match với cấu hình các đường đi con khác.

```yml
#route default
route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 2m
  repeat_interval: 3h
  receiver: 'default'
```

- group_by: Prometheus sẽ gom những thông báo có cùng `alertname` vào 1 thông báo.
- group_wait: Sau khi một cảnh báo được taọ ra. Phải đợi khoảng thời gian này thì cảnh báo mới được gửi đi.
- group_interval: Sau khi cảnh báo đầu tiên gửi đi, phải đợi khoảng thời gian này thì các cảnh báo sau mới được gửi đi.
- repeat_interval: Sau khi cảnh báo được gửi đi thành công. Sau khoảng thời gian này, nếu vấn đề vẫn còn tồn tại, Prometheus sẽ tiếp tục gửi đi cảnh báo sau khoảng thời gian này.

NOTE:

- Khi không khai báo gì trong `group-by: []` nó sẽ gom tất cả vào một thông báo.
- group-by theo labels trong alert rules không phải trong file config.

**routes**: nơi cấu hình các đường đi con. Prometheus sẽ dựa vào labels để chọn ra đường đi. Chúng ta có thể khai báo labels với tền đầy đủ hoặc sử dụng `regular expression`.

```yml
#route-child
  routes:
  - match:
      severity: warning
    receiver: 'gmail'
  - match:
      severity: critical
    receiver: 'slack'
  - match_re:
      service: ^(foo1|foo2|baz)$
    receiver: 'slack'

```

- Nếu thông báo có nhãn `severity` với giá trị  `warning` thì sẽ gửi đến đường đi gmail.

- Nếu thông báo có nhãn `severity` với giá trị `critical`  thì sẽ gửi đến đường đi slack.

- Label là foo1 hoặc foo2 hoặc baz sẽ được gửi đến slack.

**Inhibition**:

```yml
inhibit_rules:
- source_match:
    severity: 'critical'
  target_match:
    severity: 'warning'
  equal: ['alertname']
```

-> Khi cảnh báo có nhãn `critical` được gửi đi, thì các cảnh báo `warning` không cần phải gửi đi nữa, áp dụng với các cảnh báo có cùng `alertname`.

**receivers**:

```yml
receivers:
- name: 'default'
  slack_configs:
  - send_resolved: true
    channel: '#default'
    api_url: 'https://hooks.slack.com/services/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'

- name: 'gmail'
  email_configs:
  - to: quynh.vuongg@gmail.com
    from: quynh.vuongg@gmail.com
    smarthost: smtp.gmail.com:587
    auth_username: quynh.vuongg@gmail.com
    auth_identity: quynh.vuongg@gmail.com
    auth_password: <pass>
    send_resolved: true

- name: 'slack'
  slack_configs:
    - send_resolved: true
      channel: '#general'
      api_url: 'https://hooks.slack.com/services/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'
```

-> Nơi nhận mặc đinh: Thông báo sẽ được gửi đến channel `default` của slack.
Các nơi nhận khác của các route con với các nhãn tương ứng `gmail` và `slack`.

## Practices

Run docker-compose up

![ ](https://github.com/quynhvuongg/Picture/blob/master/noti1.png?raw=true)

Alert 1: Container CPU > 80 % với nhãn "warning"  

![ ](https://github.com/quynhvuongg/Picture/blob/master/noti2.png?raw=true)

Sau đó cảnh báo được gửi tới  Alertmanager để xử lý

![ ](https://github.com/quynhvuongg/Picture/blob/master/noti5.png?raw=true)

Với nhãn "warning" thông báo được gửi tới Gmail
![ ](https://github.com/quynhvuongg/Picture/blob/master/noti3.png?raw=true)

1 Service bị down
![ ](https://github.com/quynhvuongg/Picture/blob/master/noti4.png?raw=true)

 Alert 2: Service down với nhãn "critical"
 ![ ](https://github.com/quynhvuongg/Picture/blob/master/noti6.png?raw=true)

Cảnh báo được đưa đến Alermanager xử lý

![ ](https://github.com/quynhvuongg/Picture/blob/master/noti7.png?raw=true)

Với nhãn "critical" thông báo được gửi tới Slack kênh "#general"

![ ](https://github.com/quynhvuongg/Picture/blob/master/noti8.png?raw=true)

**Silences**: Thiết lập trực tiếp trên giao diện web của Alertmanager

![ ](https://github.com/quynhvuongg/Picture/blob/master/noti9.png?raw=true)

![ ](https://github.com/quynhvuongg/Picture/blob/master/noti10.png?raw=true)
