-- Inicializacion de roles

INSERT INTO roles (id, name, description) VALUES
                                              (UUID(), 'ADMIN', 'Administrador del sistema'),
                                              (UUID(), 'COMPANY_USER', 'Usuario de empresa cliente'),
                                              (UUID(), 'DRIVER', 'Repartidor'),
                                              (UUID(), 'WAREHOUSE_STAFF', 'Personal de almacén'),
                                              (UUID(), 'COLLECTOR', 'Recolector de paquetes');

-- 1. Usuario Administrador
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at)
VALUES (
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           'admin@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Admin System',
           '+1234567890',
           true,
           CURRENT_TIMESTAMP
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date)
VALUES (
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           'DNI',
           '12345678',
           '1990-01-01'
       );

-- Asignar rol ADMIN
INSERT INTO user_roles (user_id, role_id, assigned_by, is_active)
SELECT
    'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
    id,
    'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
    true
FROM roles WHERE name = 'ADMIN';

-- 2. Usuario Empresa
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at)
VALUES (
           'b2c3d4e5-f6g7-8901-b2c3-d4e5f6g7h8i9',
           'company@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Company User',
           '+1234567891',
           true,
           CURRENT_TIMESTAMP
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date)
VALUES (
           'b2c3d4e5-f6g7-8901-b2c3-d4e5f6g7h8i9',
           'DNI',
           '23456789',
           '1991-02-02'
       );

-- Asignar rol COMPANY_USER
INSERT INTO user_roles (user_id, role_id, assigned_by, is_active)
SELECT
    'b2c3d4e5-f6g7-8901-b2c3-d4e5f6g7h8i9',
    id,
    'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
    true
FROM roles WHERE name = 'COMPANY_USER';

-- 3. Usuario Repartidor
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at)
VALUES (
           'c3d4e5f6-g7h8-9012-c3d4-e5f6g7h8i9j0',
           'driver@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Driver User',
           '+1234567892',
           true,
           CURRENT_TIMESTAMP
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date)
VALUES (
           'c3d4e5f6-g7h8-9012-c3d4-e5f6g7h8i9j0',
           'DNI',
           '34567890',
           '1992-03-03'
       );

-- Asignar rol DRIVER
INSERT INTO user_roles (user_id, role_id, assigned_by, is_active)
SELECT
    'c3d4e5f6-g7h8-9012-c3d4-e5f6g7h8i9j0',
    id,
    'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
    true
FROM roles WHERE name = 'DRIVER';

-- 4. Usuario Personal de Almacén
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at)
VALUES (
           'd4e5f6g7-h8i9-0123-d4e5-f6g7h8i9j0k1',
           'warehouse@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Warehouse Staff',
           '+1234567893',
           true,
           CURRENT_TIMESTAMP
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date)
VALUES (
           'd4e5f6g7-h8i9-0123-d4e5-f6g7h8i9j0k1',
           'DNI',
           '45678901',
           '1993-04-04'
       );

-- Asignar rol WAREHOUSE_STAFF
INSERT INTO user_roles (user_id, role_id, assigned_by, is_active)
SELECT
    'd4e5f6g7-h8i9-0123-d4e5-f6g7h8i9j0k1',
    id,
    'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
    true
FROM roles WHERE name = 'WAREHOUSE_STAFF';

-- 5. Usuario Recolector
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at)
VALUES (
           'e5f6g7h8-i9j0-1234-e5f6-g7h8i9j0k1l2',
           'collector@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Collector User',
           '+1234567894',
           true,
           CURRENT_TIMESTAMP
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date)
VALUES (
           'e5f6g7h8-i9j0-1234-e5f6-g7h8i9j0k1l2',
           'DNI',
           '56789012',
           '1994-05-05'
       );

-- Asignar rol COLLECTOR
INSERT INTO user_roles (user_id, role_id, assigned_by, is_active)
SELECT
    'e5f6g7h8-i9j0-1234-e5f6-g7h8i9j0k1l2',
    id,
    'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
    true
FROM roles WHERE name = 'COLLECTOR';

-- Insertar zonas
INSERT INTO zones
(id, name, code, boundaries, center_point, base_rate, max_delivery_time, is_active, priority_level, created_at)
VALUES
    ('f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f', 'Zona Norte', 'ZNORTE',
     ST_GeomFromText('POLYGON((-74.03 4.70, -74.02 4.70, -74.02 4.72, -74.03 4.72, -74.03 4.70))'),
     ST_GeomFromText('POINT(-74.025 4.71)'),
     25.00, 60, true, 1, NOW()),

    ('e7d6c5b4-a3f2-4e1d-8c9b-7a6b5c4d3e2f', 'Zona Centro', 'ZCENTRO',
     ST_GeomFromText('POLYGON((-74.08 4.60, -74.06 4.60, -74.06 4.62, -74.08 4.62, -74.08 4.60))'),
     ST_GeomFromText('POINT(-74.07 4.61)'),
     30.00, 45, true, 2, NOW()),

    ('d6e5f4c3-b2a1-4d0e-9f8c-7b6a5d4c3e2f', 'Zona Sur', 'ZSUR',
     ST_GeomFromText('POLYGON((-74.12 4.50, -74.10 4.50, -74.10 4.52, -74.12 4.52, -74.12 4.50))'),
     ST_GeomFromText('POINT(-74.11 4.51)'),
     35.00, 75, true, 1, NOW());

-- Insertar cobertura de zonas
INSERT INTO zone_coverage
(zone_id, coverage_area, operating_hours, max_concurrent_orders, surge_multiplier)
VALUES
    ('f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f',
     ST_GeomFromText('POLYGON((-74.04 4.69, -74.01 4.69, -74.01 4.73, -74.04 4.73, -74.04 4.69))'),
     '{"weekdays":{"start":"08:00","end":"20:00"},"weekends":{"start":"09:00","end":"17:00"}}',
     15, 1.5),

    ('e7d6c5b4-a3f2-4e1d-8c9b-7a6b5c4d3e2f',
     ST_GeomFromText('POLYGON((-74.09 4.59, -74.05 4.59, -74.05 4.63, -74.09 4.63, -74.09 4.59))'),
     '{"weekdays":{"start":"07:00","end":"22:00"},"weekends":{"start":"08:00","end":"20:00"}}',
     25, 1.8),

    ('d6e5f4c3-b2a1-4d0e-9f8c-7b6a5d4c3e2f',
     ST_GeomFromText('POLYGON((-74.13 4.49, -74.09 4.49, -74.09 4.53, -74.13 4.53, -74.13 4.49))'),
     '{"weekdays":{"start":"08:00","end":"21:00"},"weekends":{"start":"09:00","end":"18:00"}}',
     20, 1.3);

-- Usuario de empresa
INSERT INTO users
(id, email, password_hash, full_name, phone, is_active, created_at, updated_at)
VALUES
    ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'empresauser@empresa.com',
     '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq', 'Gerente Empresa', '+573007654321', true, NOW(), NOW());

INSERT INTO user_profiles
(user_id, document_type, document_number, birth_date, emergency_contact_name, emergency_contact_phone)
VALUES ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'DUI', '9876543210', '1990-07-20', 'Contacto Empresa', '+573001112233');

INSERT INTO user_roles
(user_id, role_id, assigned_at, assigned_by, is_active, created_at)
VALUES ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', '991dfbd6-f89b-11ef-a120-0242ac120003', NOW(), 'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8', true, NOW());

-- Insertar compañías
INSERT INTO companies
(id, name, legal_name, tax_id, contact_email, contact_phone, is_active, delivery_rate, contract_start_date, created_at, updated_at)
VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Express Delivery Co.', 'Express Delivery S.A.S',
     '900123456-7', 'contacto@expressdelivery.com', '+573001234567', true, 20.50,
     NOW(), NOW(), NOW()),

    ('b1ffc99-8d0a-4be8-aa6c-7aa8ce481a22', 'Rapid Logistics Inc.', 'Rapid Logistics International Inc.',
     '900654321-8', 'contacto@rapidlogistics.com', '+573009876543', true, 25.75,
     NOW(), NOW(), NOW());

-- Insertar direcciones de compañía
INSERT INTO company_addresses
(id, company_id, address_line1, address_line2, city, state, postal_code, location, is_main, created_at)
VALUES
    ('e1b09d38-e71f-415f-b3eb-ffeb8dd3b493', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Calle 100 #15-20', 'Edificio Centro Empresarial', 'Bogotá', 'Cundinamarca', '110121',
     ST_GeomFromText('POINT(-74.05 4.68)'), true, NOW()),

    ('f2c18d47-f81f-416f-c4fc-00fc9ee4c594', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Calle 80 #20-30', 'Torre Norte', 'Bogotá', 'Cundinamarca', '110111',
     ST_GeomFromText('POINT(-74.07 4.67)'), false, NOW()),

    ('03d29d56-091f-417f-d5fd-11fd0ff5d605', 'b1ffc99-8d0a-4be8-aa6c-7aa8ce481a22',
     'Carrera 15 #93-60', 'Piso 3', 'Bogotá', 'Cundinamarca', '110221',
     ST_GeomFromText('POINT(-74.04 4.66)'), true, NOW());

-- Insertar sucursales
INSERT INTO company_branches
(id, company_id, name, code, contact_name, contact_phone, contact_email, is_active, zone_id, created_at, updated_at)
VALUES
    ('b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Sucursal Norte', 'SUC-NORTE-001', 'Gerente Norte', '+573001112233', 'norte@expressdelivery.com',
     true, 'f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f', NOW(), NOW()),

    ('c6f9d4e2-3f6a-5d7c-b9f0-f6f4d3c2b1a6', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Sucursal Centro', 'SUC-CENTRO-001', 'Gerente Centro', '+573004445566', 'centro@expressdelivery.com',
     true, 'e7d6c5b4-a3f2-4e1d-8c9b-7a6b5c4d3e2f', NOW(), NOW()),

    ('d7a0e5f3-4a7b-6e8d-ca01-a7a5e4d3c2b1', 'b1ffc99-8d0a-4be8-aa6c-7aa8ce481a22',
     'Sucursal Principal', 'SUC-PPAL-001', 'Gerente Principal', '+573007778899', 'principal@rapidlogistics.com',
     true, 'd6e5f4c3-b2a1-4d0e-9f8c-7b6a5d4c3e2f', NOW(), NOW());

-- Insertar usuario de compañía
INSERT INTO company_users
(user_id, company_id, position, department, can_create_orders, created_at, updated_at)
VALUES
    ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Gerente de Operaciones', 'Operaciones', true, NOW(), NOW());
