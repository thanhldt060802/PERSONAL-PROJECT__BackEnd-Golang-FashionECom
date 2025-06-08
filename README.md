Ngôn ngữ lập trình: Go (Golang).<br>
Framework: Gin + Huma v2 (tạo OpenAPI docs, chuẩn hóa RESTful API).<br>
ORM: Bun ORM (truy xuất dữ liệu hiệu quả từ PostgreSQL).<br>
Cơ sở dữ liệu: PostgreSQL (chạy trong Docker container).<br>
Middleware: JWT (JSON Web Token) dùng để xác thực người dùng.<br>
Message Queue: Redis (Pub/Sub) dùng để giao tiếp bất đồng bộ.<br>
Giao tiếp đồng bộ: gRPC giữa các service (tốc độ cao, tiết kiệm băng thông).<br>
Công cụ tìm kiếm: Elasticsearch (tối ưu hóa tìm kiếm sản phẩm).<br>
Triển khai: Docker Compose (chạy tất cả các thành phần dịch vụ).<br>