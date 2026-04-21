# NutriTrack — پلتفرم مدیریت تغذیه

سرویس بک‌اند Go برای مدیریت تغذیه، ویژه کاربران ایرانی.

## پیش‌نیازها

- Docker + Docker Compose
- Go 1.24+

## راه‌اندازی سریع

```bash
cp .env.example .env
# ویرایش .env و تنظیم متغیرها
docker compose up -d
```

## اندپوینت‌های اصلی

| مسیر | توضیح |
|------|--------|
| `GET /health` | بررسی سلامت سرویس |
| `POST /api/v1/auth/login` | ورود با ایمیل/رمز عبور |
| `POST /api/v1/auth/otp/send` | ارسال OTP |
| `POST /api/v1/auth/otp/verify` | تایید OTP |

## دستورات رایج

```bash
make test          # اجرای تست‌ها
make test-race     # اجرای تست‌ها با race detector
make lint          # بررسی کیفیت کد
make build         # ساخت باینری
make docker-build  # ساخت Docker image
make up            # راه‌اندازی سرویس‌ها
make down          # توقف سرویس‌ها
make logs          # مشاهده لاگ‌ها
make migrate-up    # اجرای migration ها
```

## ساختار پروژه

```
cmd/server/        # نقطه ورود برنامه
internal/          # منطق داخلی (DDD)
  domain/          # موجودیت‌ها و قراردادها
  application/     # Use Case ها
  infrastructure/  # پیاده‌سازی‌های زیرساخت
  interface/       # هندلرهای HTTP
migrations/        # فایل‌های migration دیتابیس
```

## متغیرهای محیطی

فایل `.env.example` را مشاهده کنید — تمام متغیرها با توضیحات فارسی مستند شده‌اند.
