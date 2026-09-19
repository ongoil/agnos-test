# ขั้นตอนการติดตั้ง

1. Clone Project
Clone project จาก GitHub
git clone <repository-url>
cd <project-name>


2. ตั้งค่า Environment
สร้างไฟล์ .env จากตัวอย่างที่โปรเจกต์กำหนด

จากนั้นกำหนดค่าที่จำเป็น เช่น

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hospital_middleware

JWT_SECRET=your-secret-key
```

หมายเหตุ:
ค่าของ Environment สามารถเปลี่ยนได้ตามการตั้งค่าของเครื่องหรือ Docker Compose

3. Run ด้วย Docker Compose
ตรวจสอบว่ามี Docker และ Docker Compose ติดตั้งอยู่ในเครื่องแล้ว

```bash
docker compose up -d --build
```

คำสั่งนี้จะทำการ
* Build Go Application
* Start Go Application
* Start PostgreSQL
* Start Nginx
* เชื่อมต่อ Service ต่าง ๆ ตามที่กำหนดใน `docker-compose.yml`

4. ตรวจสอบ Container
ตรวจสอบว่า Container ทำงานอยู่หรือไม่

```bash
docker compose ps
```

ควรเห็น Service ที่เกี่ยวข้อง เช่น

```text
app
postgres
nginx
```

และสถานะควรเป็น `Up`

5. ตรวจสอบ Log
กรณีต้องการดู Log ของทุก Service

```bash
docker compose logs
```

6. Database

---
PostgreSQL จะถูกสร้างและเริ่มต้นตามค่าที่กำหนดใน `docker-compose.yml`
กรณีมี Migration ให้รันตามคำสั่งที่กำหนดไว้ใน Project
หาก Project มี Seed Data สำหรับ Hospital และ Patient สามารถนำ SQL ที่เตรียมไว้ execute เข้า PostgreSQL ได้

7. ตรวจสอบ API
หลังจาก Container ทำงานแล้ว สามารถเรียก API ผ่าน Nginx ได้ตาม Port ที่กำหนดใน `docker-compose.yml`

ตัวอย่าง:
http://localhost/backend/api/v1


8. Stop Project
หยุด Container

```bash
docker compose down
```

9. Start Project ครั้งถัดไป
หลังจาก Build ครั้งแรกแล้ว สามารถ Start Project ได้ด้วย

```bash
docker compose up -d
```

หากมีการแก้ไข Dockerfile หรือ Dependency และต้องการ Build ใหม่

```bash
docker compose up -d --build
```
