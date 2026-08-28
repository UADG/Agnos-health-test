-- ใช้ gen_random_uuid() สำหรับสร้าง UUID อัตโนมัติ (รองรับ PostgreSQL 13+)

-- 1. ตารางโรงพยาบาล (Hospitals)
CREATE TABLE hospitals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL, -- เช่น 'HOSPITAL_A'
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. ตารางพนักงาน (Staff)
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- ต้องเก็บเป็นค่า Hash (เช่น bcrypt) ห้ามเก็บ Plain Text
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. ตารางผู้ป่วย (Patients)
-- รองรับโครงสร้างข้อมูลที่ดึงมาจาก Hospital A API
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    
    patient_hn VARCHAR(50) NOT NULL, -- Hospital Number
    national_id VARCHAR(20),
    passport_id VARCHAR(50),
    
    first_name_th VARCHAR(100),
    middle_name_th VARCHAR(100),
    last_name_th VARCHAR(100),
    
    first_name_en VARCHAR(100),
    middle_name_en VARCHAR(100),
    last_name_en VARCHAR(100),
    
    date_of_birth DATE,
    phone_number VARCHAR(20),
    email VARCHAR(255),
    gender CHAR(1) CHECK (gender IN ('M', 'F')), -- บังคับให้กรอกแค่ M หรือ F
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- เงื่อนไข: ใน 1 โรงพยาบาล จะต้องไม่มีรหัสผู้ป่วย (HN) ซ้ำกัน
    UNIQUE (hospital_id, patient_hn) 
);

-- ==========================================
-- สร้าง INDEX สำหรับปรับปรุงความเร็วในการค้นหา (Optimization)
-- สำคัญมากสำหรับ API: /patient/search ที่รับ Optional Parameters
-- ==========================================

-- Index สำหรับค้นหาผู้ป่วยในโรงพยาบาลเดียวกัน
CREATE INDEX idx_patients_hospital_id ON patients(hospital_id);

-- Index สำหรับช่องค้นหาหลักตามโจทย์
CREATE INDEX idx_patients_national_id ON patients(national_id);
CREATE INDEX idx_patients_passport_id ON patients(passport_id);
CREATE INDEX idx_patients_phone ON patients(phone_number);
CREATE INDEX idx_patients_email ON patients(email);

-- Index สำหรับค้นหาชื่อ-นามสกุล (ทำแบบ Composite เผื่อการค้นหาร่วมกัน)
CREATE INDEX idx_patients_name_en ON patients(first_name_en, last_name_en);
CREATE INDEX idx_patients_name_th ON patients(first_name_th, last_name_th);

INSERT INTO hospitals (code, name) VALUES 
('HOSPITAL_A', 'Hospital A'),
('HOSPITAL_B', 'Hospital B (Central Branch)'),
('HOSPITAL_C', 'Hospital C (North Branch)');





INSERT INTO staff (username, password_hash, hospital_id) VALUES 
('admin_a1', '$2a$10$Wq9b...สมมติว่านี่คือรหัสผ่านที่Hashแล้ว', (SELECT id FROM hospitals WHERE code = 'HOSPITAL_A')),
('doctor_a1', '$2a$10$Wq9b...สมมติว่านี่คือรหัสผ่านที่Hashแล้ว', (SELECT id FROM hospitals WHERE code = 'HOSPITAL_A')),
('admin_b1', '$2a$10$Wq9b...สมมติว่านี่คือรหัสผ่านที่Hashแล้ว', (SELECT id FROM hospitals WHERE code = 'HOSPITAL_B')),
('admin_c1', '$2a$10$Wq9b...สมมติว่านี่คือรหัสผ่านที่Hashแล้ว', (SELECT id FROM hospitals WHERE code = 'HOSPITAL_C'));





INSERT INTO patients (
    hospital_id, patient_hn, national_id, passport_id, 
    first_name_th, last_name_th, first_name_en, last_name_en, 
    date_of_birth, phone_number, email, gender
) VALUES 

(
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_A'), 'HN-A-0001', '1100112233445', NULL, 
    'สมชาย', 'ใจดี', 'Somchai', 'Jaidee', 
    '1990-05-15', '0812345678', 'somchai@email.com', 'M'
),
(
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_A'), 'HN-A-0002', '1100555566778', NULL, 
    'สมหญิง', 'รักเรียน', 'Somying', 'Rakrian', 
    '1995-11-20', '0898765432', 'somying@email.com', 'F'
),
(   
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_A'), 'HN-A-0003', NULL, 'A12345678', 
    NULL, NULL, 'John', 'Doe', 
    '1988-01-01', '0911112222', 'john.doe@email.com', 'M'
),


(
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_B'), 'HN-B-0001', '2200111122223', NULL, 
    'มานะ', 'อดทน', 'Mana', 'Aodthon', 
    '1982-08-08', '0822223333', 'mana@email.com', 'M'
),
(   
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_B'), 'HN-B-0002', NULL, 'J98765432', 
    NULL, NULL, 'Yui', 'Aragaki', 
    '1992-06-11', '0833334444', 'yui.a@email.com', 'F'
),


(   
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_C'), 'HN-C-0001', '3300112233445', NULL, 
    'กานต์', 'มั่นคง', 'Karn', 'Mankong', 
    '2000-02-14', '0844445555', NULL, 'M'
);