# Tax Calculator API

REST API สำหรับคำนวณภาษีเงินได้บุคคลธรรมดา

## เทคโนโลยีที่ใช้
- Go 1.22+
- net/http (Standard Library)
- JSON Encoder/Decoder
- Postman สำหรับทดสอบ API

## วิธีการติดตั้ง
```bash
 สร้าง go.mod  
```bash
go mod init dd-tax-trainee

```

## วิธีการรัน
```bash
go run main.go

```

API จะรันที่ Server running at http://localhost:8080


## API Endpoints

### POST /tax/calculations
คำนวณภาษีเงินได้

**Request Body:**
```json
{
  "totalIncome": 750000,
  "wht": 0,
  "allowances": [
    {
      "allowanceType": "donation",
      "amount": 0
    }
  ]
}

```

**Response:**
```json
{
  "tax": 63500,
  "taxLevel": [
    {
      "level": "0-150,000",
      "tax": 0
    },
    {
      "level": "150,001-500,000",
      "tax": 35000
    },
    {
      "level": "500,001-1,000,000",
      "tax": 28500
    }
  ]
}

```

## ตัวอย่างการใช้งาน

### คำนวณภาษีพื้นฐาน
```bash
curl -X POST http://localhost:8080/tax/calculations \
  -H "Content-Type: application/json" \
  -d '{
    "totalIncome": 750000,
    "wht": 0,
    "allowances": []
  }'

```

### คำนวณภาษีพร้อม WHT
```bash
curl -X POST http://localhost:8080/tax/calculations \
  -H "Content-Type: application/json" \
  -d '{
    "totalIncome": 600000,
    "wht": 15000,
    "allowances": []
  }'

```