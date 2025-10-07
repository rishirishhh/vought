# 📺 Vought

A high-performance and scalable video streaming service built with **Golang**.
The system supports video upload, storage, encoding, and adaptive streaming delivery with an architecture designed for extensibility and large-scale use cases.

---

## 🚀 Features

- **Video Uploads** – Secure and resumable uploads.
- **Encoding Pipeline** – Transcodes videos into multiple bitrates and formats (HLS/DASH).
- **Adaptive Streaming** – Optimized playback experience across devices and network conditions.
- **User Authentication & Authorization** – Role-based access using JWT or SSO (pluggable).
- **Metadata Management** – Store and query video metadata (title, description, tags, status).
- **Search & Recommendation Ready** – Extensible APIs for future ML-driven personalization.
- **Monitoring & Observability** – Integrated logging, metrics, and tracing.
- **Horizontal Scalability** – Stateless services with message queues for workload distribution.

---

**Components:**

- **API Gateway / Controller Layer** – Handles incoming requests (upload, playback, metadata).
- **Video Service (Core)** – Business logic for video management.
- **Encoding Worker** – Processes encoding jobs asynchronously.
- **Object Storage** – Stores raw and processed video files (e.g., S3, MinIO, GCS).
- **Metadata Database** – Manages video metadata and state (e.g., PostgreSQL, MySQL).
- **CDN Integration** – Ensures fast, scalable global delivery.

---

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **API Layer**: REST + gRPC (for internal services)
- **Storage**: S3/MinIO (object storage for videos)
- **Database**: PostgreSQL (video metadata)
- **Queue**: Kafka / RabbitMQ (encoding jobs)
- **Transcoding**: FFmpeg-based workers
- **Deployment**: Docker + Kubernetes
- **Monitoring**: Prometheus + Grafana
- **Auth**: JWT / OAuth 2.0

---

## ⚡ Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- FFmpeg installed
- PostgreSQL
- MinIO/S3 bucket

### Run Locally

```bash
# Clone repository
git clone https://github.com/rishirishhh/vought.git
cd vought

# Run server
go run /server/main.go
```

---

## 📈 Scalability Considerations

- **Stateless API layer** for easy horizontal scaling.
- **Async encoding pipeline** for heavy workloads.
- **CDN-backed streaming** for low latency delivery.
- **Sharding/partitioning** in metadata DB for large datasets.
- **Service mesh** for secure inter-service communication.
