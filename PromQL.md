# PromQL

Prometheus cung cấp một ngôn ngữ truy vấn gọi là PromQL (Prometheus Query Language) cho phép người dùng chọn và tổng hợp dữ liệu time series theo thời gian thực.
Kết quả của một biểu thức có thể được hiển thị dưới dạng biểu đồ, dữ liệu dạng bảng trong trình duyệt biểu thức của Prometheus hoặc được sử dụng bởi các hệ thống bên ngoài thông qua API HTTP.

## Expression language data types

- Instant vector: một set time series chứa một mẫu cho mỗi time series, tất cả đều chia sẻ cùng timestamp .
- Range vector: một set time series chứa một phạm vi các điểm dữ liệu theo thời gian cho mỗi time series .
- Scalar (vô hướng): một giá trị số thực đơn giản
- String: một chuỗi đơn giản ( hiện tại không dùng )

## Literals

### String literals

Nội dung được đóng trong dấu " " , '' , `` .

```sh
"this is a string"
'these are : \n \\ \t'
`these are not unescaped: \n ' " \t`

```

### Float literals

Biểu diễn số thực theo dạng literals theo form: `[-](digits)[.(digits)]`

```sh
-2.43
```

### Time series Selector

**_Instant vector selector_**

```sh
http_requests_total{job="prometheus",group="canary"}

```

- Các toán tử :
  - `=`: Chọn nhãn bằng chuỗi được cung cấp.
  - `!=`: Chọn nhãn không bằng chuỗi được cung cấp.
  - `=~`: Chọn các nhãn match với chuỗi được cung cấp.
  - `!~`: Chọn các nhãn không match với chuỗi được cung cấp.

```sh
http_requests_total{environment=~"staging|testing|development",method!="GET"}
```

-> Chọn những nhãn môi trường là staging, testing, development và sử dụng method khác GET.

- Vector selector phải chỉ định tên cụ thể hoặc ít nhất không match với chuỗi trống.
- Sử dụng internal label  `__name__`
- Metric name không được chưa các ký tự : `bool`, `on`, `ignoring`, `group_left` and `group_right`.

-> Để sử dụng chúng ta có thể dụng nhãn `__name__`

```sh
on{} # Bad!
{__name__="on"} # Good!

```

Regex(regular expressions) Prometheus [RE2 syntax](https://github.com/google/re2/wiki/Syntax).

**_Range Vector Selectors_**

```sh
http_requests_total{job="prometheus"}[5m]

```

-> Chọn tất cả các giá trị được ghi lại trong vòng 5 phút trước so với hiện tại có tên là `http_requests_total` và nhãn `job` là  `prometheus`.

- `s`  - seconds
- `m`  - minutes
- `h`  - hours
- `d`  - days
- `w`  - weeks
- `y`  - years

**_Offset modifier_**

Cho phép thay đổi thời gian bù cho các instant và range vectors trong một truy vấn.

```sh
http_requests_total offset 5m

```

-> trả về giá trị của `http_quests_total` 5 phút trong quá khứ so với thời gian đánh giá truy vấn hiện tại.

```sh
sum(http_requests_total{method="GET"} offset 5m) // GOOD.

sum(http_requests_total{method="GET"}) offset 5m // INVALID.

```

```sh
rate(http_requests_total[5m] offset 1w)

```

-> trả về tốc độ trong 5 phút mà `http_quests_total` đã có một tuần trước.

### Subquery

Cho phép thực hiện một truy vấn tức thời trong một query. Kết quả của truy vấn con là một range vectơ.

Syntax:

```sh
<instant_query> '[' <range> ':' [ <resolution> ] ']' [ offset <duration> ]

```

Ví dụ :
Subquery bên trong hàm `min_over_time` trả về tốc độ 5 phút của metric `http_quests_total` trong 30 phút qua, resolution 1m :

```sh
min_over_time( rate(http_requests_total[5m])[30m:1m] )

```

- `rate (http_quests_total [5m]) [30m: 1m]` là truy vấn con, trong đó `rate(http_quests_total [5m])` là truy vấn đã được thực hiện.
- `rate (http_quests_total [5m])` được thực hiện từ `start=<now> -30m` đến `end=<now>`, với độ chia 1m.  
- Cuối cùng, kết quả của tất cả các đánh giá ở trên được chuyển đến `min_over_time ()`.

## Operators

### Binary operators

 PromQL hỗ trợ các toán tử logic và số học cơ bản .

**_Arithmetic binary operators_**
  
- `+`  (addition)
- `-`  (subtraction)
- `*`  (multiplication)
- `/`  (division)
- `%`  (modulo)
- `^`  (exponentiation)

Được sử dụng giữa các toán hạng:

- scalar-scalar: Kết quả là một scalar khác.
- instant vector -scalar : Kết quả là một instant vector khác có giá trị mẫu là kết quả toán tử số học giữa giá trị mẫu gốc với scalar.
- instant vector-vector value.

**_Comparison binary operators_**

- `==` (equal)
- `!=` (not-equal)
- `>`  (greater-than)
- `<`  (less-than)
- `>=` (greater-or-equal)
- `<=` (less-or-equal)

So sánh :

- scalar-scalar: Trả về giá trị `bool` : `0 (false)`or `1 (true)`
- instant vector-scalar: Thay về trả về `bool`, nó trả về phần tử của vector so sánh với scalar có `bool = 1` và loại bỏ phần tử `bool = 0`
- instant vector-vector : Trả về phần tử vector có labels chung giữa hai vector thay vì trả về giá trị `bool`.

**_Logical/set binary operators_**

- `and` (intersection)
- `or`  (union)
- `unless` (complement)

`vector1 and vector2`: Kết quả là một vector gồm các phần tử chung  nhãn giữa 2 vector, tên và giá trị metric được chuyển từ vectơ bên trái (vector 1).

`vector1 or vector2`: Kết quả là một vector chứa tất cả các phần tử gốc (label sets + value) của vectơ 1 và vectơ 2.

`vector1 unless vector2`: Kết quả là một vectơ bao gồm các phần tử của vectơ 1 mà không có phần tử nào trong vectơ 2.

## Vector matching

**_One-to-one vector matches_**

Hai vector khớp nhau nếu chúng có cùng bộ nhãn và giá trị tương ứng.

`ignoring(<label list>)`: cho phép bỏ qua các nhãn này  
`on (<label list>)`: chỉ cho phép những nhãn này .

```sh
<vector expr> <bin-op> ignoring(<label list>) <vector expr>
<vector expr> <bin-op> on(<label list>) <vector expr>

```

Example input:

```sh
method_code:http_errors:rate5m{method="get", code="500"}  24
method_code:http_errors:rate5m{method="get", code="404"}  30
method_code:http_errors:rate5m{method="put", code="501"}  3
method_code:http_errors:rate5m{method="post", code="500"} 6
method_code:http_errors:rate5m{method="post", code="404"} 21

method:http_requests:rate5m{method="get"}  600
method:http_requests:rate5m{method="del"}  34
method:http_requests:rate5m{method="post"} 120

```

Example query:

```sh
method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m

```

Example result:

```sh
{method="get"}  0.04            //  24 / 600
{method="post"} 0.05            //   6 / 120

```

**_Many-to-one and one-to-many vector matches_**

Đối với matching dạng này sẽ xuất hiện trường hợp phần tử vector trên "one" có thể match với nhiều phần tử trên "many". Do đó, cần phải yêu cầu rõ ràng bằng cách sử dụng `group_left` `group_right` xác định thẻ bên nào cao hơn (sử dụng trong các toán tử số học và so sánh).

```sh
<vector expr> <bin-op> ignoring(<label list>) group_left(<label list>) <vector expr>
<vector expr> <bin-op> ignoring(<label list>) group_right(<label list>) <vector expr>
<vector expr> <bin-op> on(<label list>) group_left(<label list>) <vector expr>
<vector expr> <bin-op> on(<label list>) group_right(<label list>) <vector expr>

```

Example query :

```sh
method_code:http_errors:rate5m / ignoring(code) group_left method:http_requests:rate5m

```

Example result :

```sh
{method="get", code="500"}  0.04            //  24 / 600
{method="get", code="404"}  0.05            //  30 / 600
{method="post", code="500"} 0.05            //   6 / 120
{method="post", code="404"} 0.175           //  21 / 120

```

## Aggregation operators

Prometheus hỗ trợ các toán tử  tích hợp được sử dụng để tổng hợp các phần tử của một instant vector, cho kết quả một vector mới có ít phần tử hơn với các giá trị tổng hợp:

- `sum`  (tổng dimensions)
- `min`  (giá trị nhỏ nhất trong dimensions)
- `max`  (giá trị lớn nhất trong dimensions)
- `avg`  (giá trị trung bình của dimensions)
- `stddev`  (tính độ lệch chuẩn trong dimensions)
- `stdvar`  (tính phương sai chuẩn trong dimensions)
- `count`  (đếm số phần tử trong vector)
- `count_values`  (đếm số phần tử trong vector có cùng giá trị)
- `bottomk`  (k phần tử có giá trị nhỏ nhất)
- `topk`  (k phần tử có giá trị lớn nhất)
- `quantile`  (tính lượng tử với φ (0 ≤ φ ≤ 1) )

Syntax :

```sh
<aggr-op> [without|by (<label list>)] ([parameter,] <vector expression>)

```

or

```sh
<aggr-op> ([parameter,] <vector expression>) [without|by (<label list>)]

```

- `without`: loại bỏ các labels được liệt kê trong vector kết quả .
- `by`: chỉ nhận các labels được liệt kê .
- `parameter`: chỉ yêu cầu đối với `count_values`, `quantile`, `topk`, `bottomk`.

Example :

```sh
sum without (instance) (http_requests_total)

```

```sh
sum by (application, group) (http_requests_total)

```

```sh
sum(http_requests_total)

```

```sh
count_values("version", build_version)

```

```sh
topk(5, http_requests_total)

```

## Binary operator precedence

1. `^`
2. `*`, `/`, `%`
3. `+`, `-`
4. `==`, `!=`, `<=`, `<`, `>=`, `>`
5. `and`, `unless`
6. `or`

---

## FUNCTIONS

- `abs(v instant-vector)`: trả về giá trị tuyệt đối

- `absent(v instant-vector)`:
  - Trả về một vector trống nếu vector truyền vào cho nó có bất kỳ phần từ nào .
  - Trả về một vector có một phần tử có giá trị 1 nếu vector nó truyền vào không có phần tử nào .

```sh
absent(nonexistent{job="myjob"})
# => {job="myjob"}

absent(nonexistent{job="myjob",instance=~".*"})
# => {job="myjob"}

absent(sum(nonexistent{job="myjob"}))
# => {}
```

- `absent_over_time(v range-vector)`:
  - Trả về một vector trống nếu vector truyền vào cho nó cho bất kỳ phần từ nào .
  - Trả về một vector có một phần tử có giá trị 1 nếu vector nó truyền vào không có phần tử nào .
  
```sh
absent_over_time(nonexistent{job="myjob"}[1h])
# => {job="myjob"}

absent_over_time(nonexistent{job="myjob",instance=~".*"}[1h])
# => {job="myjob"}

absent_over_time(sum(nonexistent{job="myjob"})[1h:])
# => {}

```

- `ceil(v instant-vector)`: làm tròn giá trị đến số nguyên gần nhất.
- `changes(v range-vector)`: trả về số lần thay đổi giá trị trong phạm vi thời gian được vector cung cấp.
- `clamp_max(v instant-vector, max scalar)`: kẹp giá trị để có giới hạn trên tối đa.
- `clamp_min(v instant-vector, min scalar)`: kẹp giá trị để có giới hạn dưới tối đa.
- `day_of_month(v=vector(time()) instant-vector)`: trả về ngày trong tháng.
- `day_of_week(v=vector(time()) instant-vector)`: trả về ngày trong tuần.
- `days_in_month(v=vector(time()) instant-vector)`: trả về số ngày trong tháng.
- `delta(v range-vector)`: tính toán sự khác biệt giữa giá trị đầu tiên và cuối cùng của mỗi phần tử chuỗi time series trong một range vector, trả về một instant vector với các delta đã cho và các labels tương đương.

```sh
delta(cpu_temp_celsius{host="zeus"}[2h])

```

- `deriv(v range-vector)`: tính đạo hàm trong phạm vị được vector cung cấp .
- `exp(v instant-vector)`: tính số mũ cho tất cả phần tử .
  - `Exp(+Inf) = +Inf`
  - `Exp(NaN) = NaN`
- `floor(v instant-vector)`: làm tròn xuống số nguyên gần nhất .
- `histogram_quantile(φ float, b instant-vector)`: lượng tử
- `rate(v range-vector)`: tính tốc độ tăng trung bình mỗi giây trong phạm vi .

```sh
rate(http_requests_total{job="api-server"}[5m])

```

- `sort(v instant-vector)`: sắp xếp giá trị tăng dần
- `sort_desc(v instant-vector)`: sắp xếp giá trị giảm dần .
- `timestamp(v instant-vector)`:trả về timestamp của từng mẫu giá trị của vectơ đã cho dưới dạng số giây kể từ ngày 1 tháng 1 năm 1970 UTC.
- `<aggregation>_over_time()`:

  - `avg_over_time (range-vector)`: giá trị trung bình của tất cả các điểm trong khoảng thời gian được chỉ định.
  - `min_over_time (range-vector)`: giá trị tối thiểu của tất cả các điểm trong khoảng thời gian được chỉ định.
  - `max_over_time (range-vector)`: giá trị tối đa của tất cả các điểm trong khoảng thời gian được chỉ định.
  - `sum_over_time (range-vector)`: tổng của tất cả các giá trị trong khoảng thời gian được chỉ định.
  - `Count_over_time (range-vector)`: số lượng của các giá trị trong khoảng thời gian được chỉ định.
  - `quantile_over_time (scalar, range-vector)`: φ-quantile (0 ≤ φ ≤ 1) của các giá trị trong khoảng đã chỉ định.
  - `stddev_over_time (range-vector)`: độ lệch chuẩn của các giá trị trong khoảng thời gian được chỉ định.
  - `stdvar_over_time (range-vector)`: phương sai tiêu chuẩn của các giá trị trong khoảng thời gian được chỉ định.
