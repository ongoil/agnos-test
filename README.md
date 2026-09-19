# Hospital Middleware API

ระบบตัวกลางสำหรับค้นหาข้อมูลผู้ป่วยจาก HIS โดยจำกัดสิทธิ์ของเจ้าหน้าที่ตามโรงพยาบาล

## เริ่มต้นใช้งาน

```bash
docker compose up --build
```

Nginx เปิดที่ `http://localhost` และส่งต่อไปยัง Go API ส่วน PostgreSQL เปิดที่พอร์ต `5432`
สำหรับ production ต้องเปลี่ยน `JWT_SECRET` และรหัสผ่านฐานข้อมูลใน compose

## API หลัก

- `POST /staff/create` สร้าง staff โดยรับ `username`, `password`, `hospital` หรือ `hospital_id`
- `POST /staff/login` login โดยรับ `username`, `password`, `hospital` หรือ `hospital_id`
- `GET /patient/search` ต้องมี `Authorization: Bearer <token>` และรับ query parameters ที่เป็น optional:
  `national_id`, `passport_id`, `first_name`, `middle_name`, `last_name`, `date_of_birth`, `phone_number`, `email`

ต้องสร้าง hospital ก่อนผ่าน `POST /backend/api/v1/register/hospital` เช่น `{"hospital_name":"Hospital A"}`

รายละเอียดอยู่ที่ [docs/planning.md](docs/planning.md)
