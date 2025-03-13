-- IMPORTANTE EJECUTAR EL SCRIPT LUEGO DE EJECUTAR EL PROGRAMA

-- Borrar la llave foranea erronea de la tabla company_users
ALTER TABLE company_users DROP FOREIGN KEY fk_company_users_company_branch;

-- Inicializacion de roles
INSERT INTO roles (id, name, description, is_active, created_at, updated_at) VALUES
                                                                                 ('991c53a8-f89b-11ef-a120-0242ac120003', 'ADMIN', 'Administrador del sistema', 1, '2025-03-04 01:54:24', '2025-03-04 01:54:24'),
                                                                                 ('991dfbd6-f89b-11ef-a120-0242ac120003', 'COMPANY_USER', 'Usuario de empresa cliente', 1, '2025-03-04 01:54:24', '2025-03-04 01:54:24'),
                                                                                 ('991e016f-f89b-11ef-a120-0242ac120003', 'DRIVER', 'Repartidor', 1, '2025-03-04 01:54:24', '2025-03-04 01:54:24'),
                                                                                 ('991e01c7-f89b-11ef-a120-0242ac120003', 'WAREHOUSE_STAFF', 'Personal de almacén', 1, '2025-03-04 01:54:24', '2025-03-04 01:54:24'),
                                                                                 ('991e01ed-f89b-11ef-a120-0242ac120003', 'COLLECTOR', 'Recolector de paquetes', 1, '2025-03-04 01:54:24', '2025-03-04 01:54:24'),
                                                                                 ('991a01ed-f89b-11ef-a120-0242ac120003', 'FINAL_USER', 'Usuario final de la aplicacion', 1, '2025-03-04 01:54:24', '2025-03-04 01:54:24');

-- Insertar compañías (deben insertarse primero porque son referenciadas por usuarios)
INSERT INTO companies
(id, name, legal_name, tax_id, contact_email, contact_phone, website, is_active, contract_details, delivery_rate, logo_url, contract_start_date, contract_end_date, created_at, updated_at)
VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Express Delivery Co.', 'Express Delivery S.A.S',
     '900123456-7', 'contacto@expressdelivery.com', '+573001234567', NULL, true, NULL, 20.50, NULL,
     '2025-03-04 03:52:32', NULL, '2025-03-04 03:52:32', '2025-03-04 03:52:32'),

    ('b1ffc99-8d0a-4be8-aa6c-7aa8ce481a22', 'Rapid Logistics Inc.', 'Rapid Logistics International Inc.',
     '900654321-8', 'contacto@rapidlogistics.com', '+573009876543', NULL, true, NULL, 25.75, NULL,
     '2025-03-04 03:52:32', NULL, '2025-03-04 03:52:32', '2025-03-04 03:52:32');

-- 1. Usuario Administrador
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at, phone_verified_at, created_at, updated_at, deleted_at, company_id)
VALUES (
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           'admin@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Admin System',
           '+1234567890',
           true,
           '2025-03-04 01:54:24',
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24',
           NULL,
           'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date, profile_picture_url, emergency_contact_name, emergency_contact_phone, additional_info, created_at, updated_at)
VALUES (
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           'DNI',
           '12345678',
           '1990-01-01',
           NULL,
           NULL,
           NULL,
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24'
       );

-- Asignar rol ADMIN
INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by, is_active, created_at)
VALUES (
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           '991c53a8-f89b-11ef-a120-0242ac120003',
           '2025-03-04 01:54:24',
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           true,
           '2025-03-04 01:54:24'
       );

-- 2. Usuario Empresa
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at, phone_verified_at, created_at, updated_at, deleted_at, company_id)
VALUES (
           'b2c3d4e5-f6g7-8901-b2c3-d4e5f6g7h8i9',
           'company@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Company User',
           '+1234567891',
           true,
           '2025-03-04 01:54:24',
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24',
           NULL,
           'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date, profile_picture_url, emergency_contact_name, emergency_contact_phone, additional_info, created_at, updated_at)
VALUES (
           'b2c3d4e5-f6g7-8901-b2c3-d4e5f6g7h8i9',
           'DNI',
           '23456789',
           '1991-02-02',
           NULL,
           NULL,
           NULL,
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24'
       );

-- Asignar rol COMPANY_USER
INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by, is_active, created_at)
VALUES (
           'b2c3d4e5-f6g7-8901-b2c3-d4e5f6g7h8i9',
           '991dfbd6-f89b-11ef-a120-0242ac120003',
           '2025-03-04 01:54:24',
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           true,
           '2025-03-04 01:54:24'
       );

-- 3. Usuario Repartidor
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at, phone_verified_at, created_at, updated_at, deleted_at, company_id)
VALUES (
           'c3d4e5f6-g7h8-9012-c3d4-e5f6g7h8i9j0',
           'driver@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Driver User',
           '+1234567892',
           true,
           '2025-03-04 01:54:24',
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24',
           NULL,
           'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date, profile_picture_url, emergency_contact_name, emergency_contact_phone, additional_info, created_at, updated_at)
VALUES (
           'c3d4e5f6-g7h8-9012-c3d4-e5f6g7h8i9j0',
           'DNI',
           '34567890',
           '1992-03-03',
           NULL,
           NULL,
           NULL,
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24'
       );

-- Asignar rol DRIVER
INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by, is_active, created_at)
VALUES (
           'c3d4e5f6-g7h8-9012-c3d4-e5f6g7h8i9j0',
           '991e016f-f89b-11ef-a120-0242ac120003',
           '2025-03-04 01:54:24',
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           true,
           '2025-03-04 01:54:24'
       );

-- 4. Usuario Personal de Almacén
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at, phone_verified_at, created_at, updated_at, deleted_at, company_id)
VALUES (
           'd4e5f6g7-h8i9-0123-d4e5-f6g7h8i9j0k1',
           'warehouse@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Warehouse Staff',
           '+1234567893',
           true,
           '2025-03-04 01:54:24',
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24',
           NULL,
           'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date, profile_picture_url, emergency_contact_name, emergency_contact_phone, additional_info, created_at, updated_at)
VALUES (
           'd4e5f6g7-h8i9-0123-d4e5-f6g7h8i9j0k1',
           'DNI',
           '45678901',
           '1993-04-04',
           NULL,
           NULL,
           NULL,
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24'
       );

-- Asignar rol WAREHOUSE_STAFF
INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by, is_active, created_at)
VALUES (
           'd4e5f6g7-h8i9-0123-d4e5-f6g7h8i9j0k1',
           '991e01c7-f89b-11ef-a120-0242ac120003',
           '2025-03-04 01:54:24',
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           true,
           '2025-03-04 01:54:24'
       );

-- 5. Usuario Recolector
INSERT INTO users (id, email, password_hash, full_name, phone, is_active, email_verified_at, phone_verified_at, created_at, updated_at, deleted_at, company_id)
VALUES (
           'e5f6g7h8-i9j0-1234-e5f6-g7h8i9j0k1l2',
           'collector@delivery.com',
           '$2a$10$2o6x9aCZsWM8oRHy/ZJqLuNmDYFzZAbfzUPBLc4pRJrto2VbmlIAq',
           'Collector User',
           '+1234567894',
           true,
           '2025-03-04 01:54:24',
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24',
           NULL,
           'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
       );

INSERT INTO user_profiles (user_id, document_type, document_number, birth_date, profile_picture_url, emergency_contact_name, emergency_contact_phone, additional_info, created_at, updated_at)
VALUES (
           'e5f6g7h8-i9j0-1234-e5f6-g7h8i9j0k1l2',
           'DNI',
           '56789012',
           '1994-05-05',
           NULL,
           NULL,
           NULL,
           NULL,
           '2025-03-04 01:54:24',
           '2025-03-04 01:54:24'
       );

-- Asignar rol COLLECTOR
INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by, is_active, created_at)
VALUES (
           'e5f6g7h8-i9j0-1234-e5f6-g7h8i9j0k1l2',
           '991e01ed-f89b-11ef-a120-0242ac120003',
           '2025-03-04 01:54:24',
           'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8',
           true,
           '2025-03-04 01:54:24'
       );

-- Insertar zonas
INSERT INTO zones
(id, name, code, boundaries, center_point, base_rate, max_delivery_time, is_active, priority_level, created_at, updated_at)
VALUES
    ('f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f', 'Zona Norte', 'ZNORTE',
     ST_GeomFromText('POLYGON((-74.03 4.70, -74.02 4.70, -74.02 4.72, -74.03 4.72, -74.03 4.70))'),
     ST_GeomFromText('POINT(-74.025 4.71)'),
     25.00, 60, true, 1, '2025-03-04 03:55:19', '2025-03-04 03:55:19'),

    ('e7d6c5b4-a3f2-4e1d-8c9b-7a6b5c4d3e2f', 'Zona Centro', 'ZCENTRO',
     ST_GeomFromText('POLYGON((-74.08 4.60, -74.06 4.60, -74.06 4.62, -74.08 4.62, -74.08 4.60))'),
     ST_GeomFromText('POINT(-74.07 4.61)'),
     30.00, 45, true, 2, '2025-03-04 03:55:19', '2025-03-04 03:55:19'),

    ('d6e5f4c3-b2a1-4d0e-9f8c-7b6a5d4c3e2f', 'Zona Sur', 'ZSUR',
     ST_GeomFromText('POLYGON((-74.12 4.50, -74.10 4.50, -74.10 4.52, -74.12 4.52, -74.12 4.50))'),
     ST_GeomFromText('POINT(-74.11 4.51)'),
     35.00, 75, true, 1, '2025-03-04 03:55:19', '2025-03-04 03:55:19');

-- Insertar cobertura de zonas
INSERT INTO zone_coverage
(zone_id, coverage_area, operating_hours, max_concurrent_orders, surge_multiplier, coverage_rules)
VALUES
    ('f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f',
     ST_GeomFromText('POLYGON((-74.04 4.69, -74.01 4.69, -74.01 4.73, -74.04 4.73, -74.04 4.69))'),
     '{"weekdays":{"start":"08:00","end":"20:00"},"weekends":{"start":"09:00","end":"17:00"}}',
     15, 1.5, NULL),

    ('e7d6c5b4-a3f2-4e1d-8c9b-7a6b5c4d3e2f',
     ST_GeomFromText('POLYGON((-74.09 4.59, -74.05 4.59, -74.05 4.63, -74.09 4.63, -74.09 4.59))'),
     '{"weekdays":{"start":"07:00","end":"22:00"},"weekends":{"start":"08:00","end":"20:00"}}',
     25, 1.8, NULL),

    ('d6e5f4c3-b2a1-4d0e-9f8c-7b6a5d4c3e2f',
     ST_GeomFromText('POLYGON((-74.13 4.49, -74.09 4.49, -74.09 4.53, -74.13 4.53, -74.13 4.49))'),
     '{"weekdays":{"start":"08:00","end":"21:00"},"weekends":{"start":"09:00","end":"18:00"}}',
     20, 1.3, NULL);

-- Usuario de empresa (adicional)
INSERT INTO users
(id, email, password_hash, full_name, phone, is_active, email_verified_at, phone_verified_at, created_at, updated_at, deleted_at, company_id)
VALUES
    ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'marlon.hg2003@gmail.com',
     '$2a$10$tg2Nk6phZ8/CRVm9RyMTueUKMF0EK7RB2mjpNjoC6Ld4vL8hbMyyW', 'Marlon Hernandez', '21212828', true, NULL, NULL,
     '2025-03-04 01:54:46', '2025-03-12 01:34:34', NULL, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

INSERT INTO user_profiles
(user_id, document_type, document_number, birth_date, profile_picture_url, emergency_contact_name, emergency_contact_phone, additional_info, created_at, updated_at)
VALUES ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'DUI', '066119477', '2003-11-25', NULL, 'Ana Mayra', '70111813', NULL,
        '2025-03-04 01:54:46', '2025-03-12 01:34:11');

INSERT INTO user_roles
(user_id, role_id, assigned_at, assigned_by, is_active, created_at)
VALUES ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', '991dfbd6-f89b-11ef-a120-0242ac120003', '2025-03-04 01:56:06',
        'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6g7h8', true, '2025-03-04 01:56:06');

-- Insertar direcciones de compañía
INSERT INTO company_addresses
(id, company_id, address_line1, address_line2, city, state, postal_code, location, is_main, created_at)
VALUES
    ('e1b09d38-e71f-415f-b3eb-ffeb8dd3b493', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Calle 100 #15-20', 'Edificio Centro Empresarial', 'Bogotá', 'Cundinamarca', '110121',
     ST_GeomFromText('POINT(-74.05 4.68)'), true, '2025-03-04 03:52:50'),

    ('f2c18d47-f81f-416f-c4fc-00fc9ee4c594', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Calle 80 #20-30', 'Torre Norte', 'Bogotá', 'Cundinamarca', '110111',
     ST_GeomFromText('POINT(-74.07 4.67)'), false, '2025-03-04 03:52:50'),

    ('03d29d56-091f-417f-d5fd-11fd0ff5d605', 'b1ffc99-8d0a-4be8-aa6c-7aa8ce481a22',
     'Carrera 15 #93-60', 'Piso 3', 'Bogotá', 'Cundinamarca', '110221',
     ST_GeomFromText('POINT(-74.04 4.66)'), true, '2025-03-04 03:52:50');

-- Insertar sucursales
INSERT INTO company_branches
(id, company_id, name, code, contact_name, contact_phone, contact_email, is_active, zone_id, operating_hours, created_at, updated_at)
VALUES
    ('b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Sucursal Norte', 'SUC-NORTE-001', 'Gerente Norte', '+573001112233', 'norte@expressdelivery.com',
     true, 'f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f', NULL, '2025-03-04 03:55:55', '2025-03-04 03:55:55'),

    ('c6f9d4e2-3f6a-5d7c-b9f0-f6f4d3c2b1a6', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Sucursal Centro', 'SUC-CENTRO-001', 'Gerente Centro', '+573004445566', 'centro@expressdelivery.com',
     true, 'e7d6c5b4-a3f2-4e1d-8c9b-7a6b5c4d3e2f', NULL, '2025-03-04 03:55:55', '2025-03-04 03:55:55'),

    ('d7a0e5f3-4a7b-6e8d-ca01-a7a5e4d3c2b1', 'b1ffc99-8d0a-4be8-aa6c-7aa8ce481a22',
     'Sucursal Principal', 'SUC-PPAL-001', 'Gerente Principal', '+573007778899', 'principal@rapidlogistics.com',
     true, 'd6e5f4c3-b2a1-4d0e-9f8c-7b6a5d4c3e2f', NULL, '2025-03-04 03:55:55', '2025-03-04 03:55:55');

-- Insertar usuario de compañía
INSERT INTO company_users
(user_id, company_id, position, department, permissions, can_create_orders, created_at, updated_at, branch_id)
VALUES
    ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
     'Gerente de Operaciones', 'Operaciones', NULL, true, '2025-03-04 03:53:32', '2025-03-04 03:53:32',
     'b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5');