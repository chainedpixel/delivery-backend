-- Users Domain
CREATE TABLE users (
                       id CHAR(36) PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       full_name VARCHAR(255) NOT NULL,
                       phone VARCHAR(20),
                       is_active BOOLEAN DEFAULT true,
                       email_verified_at TIMESTAMP NULL,
                       phone_verified_at TIMESTAMP NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL
);

CREATE TABLE roles (
                       id CHAR(36) PRIMARY KEY,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       description TEXT,
                       is_active BOOLEAN DEFAULT true,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE permissions (
                             id CHAR(36) PRIMARY KEY,
                             name VARCHAR(100) UNIQUE NOT NULL,
                             description TEXT,
                             resource VARCHAR(50) NOT NULL,
                             action VARCHAR(50) NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE user_roles (
                            user_id CHAR(36) NOT NULL,
                            role_id CHAR(36) NOT NULL,
                            assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            assigned_by CHAR(36) NOT NULL,
                            is_active BOOLEAN DEFAULT true,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            PRIMARY KEY (user_id, role_id)
);

CREATE TABLE role_permissions (
                                  role_id CHAR(36) NOT NULL,
                                  permission_id CHAR(36) NOT NULL,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_profiles (
                               user_id CHAR(36) PRIMARY KEY,
                               document_type VARCHAR(20),
                               document_number VARCHAR(30),
                               birth_date DATE,
                               profile_picture_url VARCHAR(255),
                               emergency_contact_name VARCHAR(255),
                               emergency_contact_phone VARCHAR(20),
                               additional_info JSON,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE user_sessions (
                               id CHAR(36) PRIMARY KEY,
                               user_id CHAR(36) NOT NULL,
                               token VARCHAR(255) NOT NULL,
                               device_info JSON,
                               ip_address VARCHAR(45),
                               last_activity TIMESTAMP,
                               expires_at TIMESTAMP NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Companies Domain
CREATE TABLE companies (
                           id CHAR(36) PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           legal_name VARCHAR(255) NOT NULL,
                           tax_id VARCHAR(50) UNIQUE NOT NULL,
                           contact_email VARCHAR(255) NOT NULL,
                           contact_phone VARCHAR(20) NOT NULL,
                           website VARCHAR(255),
                           is_active BOOLEAN DEFAULT true,
                           contract_details JSON,
                           delivery_rate DECIMAL(10,2) NOT NULL,
                           logo_url VARCHAR(255),
                           contract_start_date TIMESTAMP NOT NULL,
                           contract_end_date TIMESTAMP,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE company_addresses (
                                   company_id CHAR(36) NOT NULL,
                                   address_line1 VARCHAR(255) NOT NULL,
                                   address_line2 VARCHAR(255),
                                   city VARCHAR(100) NOT NULL,
                                   state VARCHAR(100) NOT NULL,
                                   postal_code VARCHAR(20),
                                   location POINT NOT NULL,
                                   is_main BOOLEAN DEFAULT false,
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                   PRIMARY KEY (company_id)
);

-- Orders Domain
CREATE TABLE orders (
                        id CHAR(36) PRIMARY KEY,
                        company_id CHAR(36) NOT NULL,
                        branch_id CHAR(36) NOT NULL,
                        client_id CHAR(36) NOT NULL,
                        driver_id CHAR(36),
                        tracking_number VARCHAR(50) UNIQUE NOT NULL,
                        status VARCHAR(20) NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE order_details (
                               order_id CHAR(36) PRIMARY KEY,
                               price DECIMAL(10,2) NOT NULL,
                               distance DECIMAL(10,2) NOT NULL,
                               pickup_time TIMESTAMP NOT NULL,
                               delivery_deadline TIMESTAMP NOT NULL,
                               delivered_at TIMESTAMP NULL,
                               requires_signature BOOLEAN DEFAULT false,
                               delivery_notes JSON,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE package_details (
                                 order_id CHAR(36) PRIMARY KEY,
                                 is_fragile BOOLEAN DEFAULT false,
                                 is_urgent BOOLEAN DEFAULT false,
                                 weight DECIMAL(10,2),
                                 dimensions JSON,
                                 special_instructions TEXT,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE delivery_addresses (
                                    order_id CHAR(36) PRIMARY KEY,
                                    recipient_name VARCHAR(255) NOT NULL,
                                    recipient_phone VARCHAR(20) NOT NULL,
                                    address_line1 VARCHAR(255) NOT NULL,
                                    address_line2 VARCHAR(255),
                                    city VARCHAR(100) NOT NULL,
                                    state VARCHAR(100) NOT NULL,
                                    postal_code VARCHAR(20),
                                    location POINT NOT NULL,
                                    address_notes JSON,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pickup_addresses (
                                  order_id CHAR(36) PRIMARY KEY,
                                  contact_name VARCHAR(255) NOT NULL,
                                  contact_phone VARCHAR(20) NOT NULL,
                                  address_line1 VARCHAR(255) NOT NULL,
                                  address_line2 VARCHAR(255),
                                  city VARCHAR(100) NOT NULL,
                                  state VARCHAR(100) NOT NULL,
                                  postal_code VARCHAR(20),
                                  location POINT NOT NULL,
                                  address_notes JSON,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_status_history (
                                      id CHAR(36) PRIMARY KEY,
                                      order_id CHAR(36) NOT NULL,
                                      status VARCHAR(20) NOT NULL,
                                      description TEXT,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_tracking (
                                order_id CHAR(36) PRIMARY KEY,
                                current_location POINT,
                                current_status VARCHAR(20) NOT NULL,
                                last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE qr_codes (
                          order_id CHAR(36) PRIMARY KEY,
                          qr_data TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE company_branches (
                                  id CHAR(36) PRIMARY KEY,
                                  company_id CHAR(36) NOT NULL,
                                  name VARCHAR(255) NOT NULL,
                                  code VARCHAR(50) NOT NULL,
                                  contact_name VARCHAR(255) NOT NULL,
                                  contact_phone VARCHAR(20) NOT NULL,
                                  contact_email VARCHAR(255) NOT NULL,
                                  is_active BOOLEAN DEFAULT true,
                                  zone_id CHAR(36) NOT NULL,
                                  operating_hours JSON,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                  UNIQUE KEY uk_company_branch_code (company_id, code)
);

CREATE TABLE company_users (
                               user_id CHAR(36) NOT NULL,
                               company_id CHAR(36) NOT NULL,
                               position VARCHAR(100) NOT NULL,
                               department VARCHAR(100),
                               permissions JSON,
                               can_create_orders BOOLEAN DEFAULT false,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                               PRIMARY KEY (user_id, company_id)
);

CREATE TABLE company_billing (
                                 id CHAR(36) PRIMARY KEY,
                                 company_id CHAR(36) NOT NULL,
                                 billing_period_start DATE NOT NULL,
                                 billing_period_end DATE NOT NULL,
                                 total_deliveries INT NOT NULL DEFAULT 0,
                                 total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
                                 status VARCHAR(20) NOT NULL,
                                 billing_details JSON,
                                 paid_at TIMESTAMP NULL,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Zones Domain
CREATE TABLE zones (
                       id CHAR(36) PRIMARY KEY,
                       name VARCHAR(100) NOT NULL UNIQUE,
                       code VARCHAR(20) NOT NULL UNIQUE,
                       boundaries POLYGON NOT NULL,
                       center_point POINT NOT NULL,
                       base_rate DECIMAL(10,2) NOT NULL,
                       max_delivery_time INT NOT NULL,
                       is_active BOOLEAN DEFAULT true,
                       priority_level INT NOT NULL DEFAULT 1,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE adjacent_zones (
                                zone_id CHAR(36) NOT NULL,
                                adjacent_zone_id CHAR(36) NOT NULL,
                                distance DECIMAL(10,2) NOT NULL,
                                travel_time INT NOT NULL,
                                coverage_overlap DECIMAL(5,2),
                                is_active BOOLEAN DEFAULT true,
                                PRIMARY KEY (zone_id, adjacent_zone_id)
);

CREATE TABLE zone_coverage (
                               zone_id CHAR(36) NOT NULL,
                               coverage_area POLYGON NOT NULL,
                               operating_hours JSON NOT NULL,
                               max_concurrent_orders INT NOT NULL DEFAULT 10,
                               surge_multiplier DECIMAL(3,2) DEFAULT 1.00,
                               coverage_rules JSON,
                               PRIMARY KEY (zone_id)
);

-- Drivers Domain
CREATE TABLE drivers (
                         user_id CHAR(36) PRIMARY KEY,
                         license_number VARCHAR(50) UNIQUE NOT NULL,
                         license_expiry DATE NOT NULL,
                         vehicle_type VARCHAR(50) NOT NULL,
                         vehicle_plate VARCHAR(20) NOT NULL,
                         vehicle_model VARCHAR(100) NOT NULL,
                         vehicle_color VARCHAR(50) NOT NULL,
                         is_active BOOLEAN DEFAULT true,
                         vehicle_details JSON,
                         documentation JSON,
                         rating DECIMAL(3,2) DEFAULT 5.00,
                         completed_deliveries INT DEFAULT 0,
                         last_delivery TIMESTAMP NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE driver_zones (
                              driver_id CHAR(36) NOT NULL,
                              zone_id CHAR(36) NOT NULL,
                              is_primary BOOLEAN DEFAULT false,
                              efficiency_rating DECIMAL(3,2) DEFAULT 5.00,
                              deliveries_completed INT DEFAULT 0,
                              last_delivery TIMESTAMP NULL,
                              is_active BOOLEAN DEFAULT true,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              PRIMARY KEY (driver_id, zone_id)
);

CREATE TABLE driver_availability (
                                     driver_id CHAR(36) NOT NULL,
                                     current_zone_id CHAR(36) NOT NULL,
                                     current_location POINT NOT NULL,
                                     status VARCHAR(20) NOT NULL,
                                     last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     active_orders INT DEFAULT 0,
                                     can_take_orders BOOLEAN DEFAULT true,
                                     shift_start TIMESTAMP NOT NULL,
                                     shift_end TIMESTAMP NOT NULL,
                                     PRIMARY KEY (driver_id)
);

-- Payments Domain
CREATE TABLE payment_methods (
                                 id CHAR(36) PRIMARY KEY,
                                 user_id CHAR(36) NOT NULL,
                                 method_type VARCHAR(50) NOT NULL,
                                 provider VARCHAR(50) NOT NULL,
                                 details JSON NOT NULL,
                                 is_default BOOLEAN DEFAULT false,
                                 is_active BOOLEAN DEFAULT true,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE payment_transactions (
                                      id CHAR(36) PRIMARY KEY,
                                      order_id CHAR(36) NOT NULL,
                                      payment_method_id CHAR(36) NOT NULL,
                                      amount DECIMAL(10,2) NOT NULL,
                                      currency VARCHAR(3) DEFAULT 'USD',
                                      status VARCHAR(20) NOT NULL,
                                      transaction_reference VARCHAR(100),
                                      payment_details JSON,
                                      processed_at TIMESTAMP NULL,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE payment_status_history (
                                        id CHAR(36) PRIMARY KEY,
                                        payment_transaction_id CHAR(36) NOT NULL,
                                        status VARCHAR(20) NOT NULL,
                                        details JSON,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Billing Domain
CREATE TABLE invoices (
                          id CHAR(36) PRIMARY KEY,
                          company_id CHAR(36) NOT NULL,
                          invoice_number VARCHAR(50) UNIQUE NOT NULL,
                          period_start DATE NOT NULL,
                          period_end DATE NOT NULL,
                          subtotal DECIMAL(10,2) NOT NULL,
                          tax DECIMAL(10,2) NOT NULL,
                          total_amount DECIMAL(10,2) NOT NULL,
                          status VARCHAR(20) NOT NULL,
                          due_date DATE NOT NULL,
                          paid_at TIMESTAMP NULL,
                          payment_reference VARCHAR(100),
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE invoice_items (
                               id CHAR(36) PRIMARY KEY,
                               invoice_id CHAR(36) NOT NULL,
                               order_id CHAR(36) NULL,
                               description VARCHAR(255) NOT NULL,
                               quantity INT NOT NULL DEFAULT 1,
                               unit_price DECIMAL(10,2) NOT NULL,
                               subtotal DECIMAL(10,2) NOT NULL,
                               tax_rate DECIMAL(5,2) NOT NULL DEFAULT 0.00,
                               tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
                               total_amount DECIMAL(10,2) NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE billing_cycles (
                                company_id CHAR(36) PRIMARY KEY,
                                cycle_day INT NOT NULL,
                                last_billed_date DATE,
                                next_billing_date DATE,
                                billing_frequency VARCHAR(20) DEFAULT 'MONTHLY',
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Events Domain
CREATE TABLE system_events (
                               id CHAR(36) PRIMARY KEY,
                               event_type VARCHAR(50) NOT NULL,
                               source VARCHAR(50) NOT NULL,
                               source_id CHAR(36) NOT NULL,
                               event_data JSON NOT NULL,
                               severity VARCHAR(20) NOT NULL DEFAULT 'INFO',
                               occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE event_logs (
                            id CHAR(36) PRIMARY KEY,
                            event_id CHAR(36) NOT NULL,
                            log_level VARCHAR(20) NOT NULL,
                            description TEXT NOT NULL,
                            metadata JSON,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE audit_logs (
                            id CHAR(36) PRIMARY KEY,
                            user_id CHAR(36) NOT NULL,
                            action VARCHAR(50) NOT NULL,
                            entity_type VARCHAR(50) NOT NULL,
                            entity_id CHAR(36) NOT NULL,
                            old_values JSON,
                            new_values JSON,
                            ip_address VARCHAR(45),
                            user_agent TEXT,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notifications Domain
CREATE TABLE notification_templates (
                                        id CHAR(36) PRIMARY KEY,
                                        name VARCHAR(100) NOT NULL,
                                        type VARCHAR(50) NOT NULL,
                                        title_template TEXT NOT NULL,
                                        content_template TEXT NOT NULL,
                                        variables JSON,
                                        is_active BOOLEAN DEFAULT true,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
                               id CHAR(36) PRIMARY KEY,
                               user_id CHAR(36) NOT NULL,
                               template_id CHAR(36) NOT NULL,
                               title VARCHAR(255) NOT NULL,
                               content TEXT NOT NULL,
                               type VARCHAR(50) NOT NULL,
                               metadata JSON,
                               is_read BOOLEAN DEFAULT false,
                               read_at TIMESTAMP NULL,
                               sent_at TIMESTAMP NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notification_preferences (
                                          user_id CHAR(36) NOT NULL,
                                          notification_type VARCHAR(50) NOT NULL,
                                          email_enabled BOOLEAN DEFAULT true,
                                          push_enabled BOOLEAN DEFAULT true,
                                          sms_enabled BOOLEAN DEFAULT false,
                                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                          PRIMARY KEY (user_id, notification_type)
);

CREATE TABLE notification_devices (
                                      id CHAR(36) PRIMARY KEY,
                                      user_id CHAR(36) NOT NULL,
                                      device_token TEXT NOT NULL,
                                      device_type VARCHAR(50) NOT NULL,
                                      is_active BOOLEAN DEFAULT true,
                                      last_used_at TIMESTAMP NULL,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


-- Dominio de control de almacen
CREATE TABLE warehouse (
                           id CHAR(36) PRIMARY KEY,
                           zone_id CHAR(36) NOT NULL,
                           name VARCHAR(100) NOT NULL,
                           address VARCHAR(255) NOT NULL,
                           location POINT NOT NULL,
                           is_active BOOLEAN DEFAULT true,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE package_warehouse_tracking (
                                            id CHAR(36) PRIMARY KEY,
                                            order_id CHAR(36) NOT NULL,
                                            warehouse_id CHAR(36) NOT NULL,
                                            status VARCHAR(50) NOT NULL,
                                            collector_id CHAR(36) NOT NULL,
                                            collected_at TIMESTAMP NULL,
                                            notes TEXT,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE warehouse_inventory (
                                     id CHAR(36) PRIMARY KEY,
                                     warehouse_id CHAR(36) NOT NULL,
                                     order_id CHAR(36) NOT NULL,
                                     status VARCHAR(50) NOT NULL,
                                     shelf_location VARCHAR(50),
                                     received_at TIMESTAMP NOT NULL,
                                     dispatched_at TIMESTAMP NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- INDICES

-- Dominio Usuarios
ALTER TABLE users
    ADD INDEX idx_users_email (email),
   ADD INDEX idx_users_phone (phone),
   ADD INDEX idx_users_status (is_active, deleted_at);

ALTER TABLE user_roles
    ADD INDEX idx_user_roles_user (user_id),
   ADD INDEX idx_user_roles_role (role_id),
   ADD INDEX idx_user_roles_active (is_active);

ALTER TABLE user_sessions
    ADD INDEX idx_sessions_user (user_id),
   ADD INDEX idx_sessions_token (token),
   ADD INDEX idx_sessions_expiry (expires_at);

ALTER TABLE user_profiles
    ADD INDEX idx_profiles_document (document_type, document_number);

-- Dominio Empresas
ALTER TABLE companies
    ADD INDEX idx_companies_tax (tax_id),
   ADD INDEX idx_companies_status (is_active),
   ADD INDEX idx_companies_contract (contract_start_date, contract_end_date);

ALTER TABLE company_branches
    ADD INDEX idx_branches_company (company_id),
   ADD INDEX idx_branches_zone (zone_id),
   ADD INDEX idx_branches_status (is_active);

ALTER TABLE company_users
    ADD INDEX idx_company_users_company (company_id),
   ADD INDEX idx_company_users_status (can_create_orders);

ALTER TABLE company_addresses
    ADD SPATIAL INDEX idx_company_location (location);

-- Dominio Pedidos
ALTER TABLE orders
    ADD INDEX idx_orders_company (company_id),
   ADD INDEX idx_orders_branch (branch_id),
   ADD INDEX idx_orders_client (client_id),
   ADD INDEX idx_orders_driver (driver_id),
   ADD INDEX idx_orders_tracking (tracking_number),
   ADD INDEX idx_orders_status (status),
   ADD INDEX idx_orders_created (created_at);

ALTER TABLE order_details
    ADD INDEX idx_order_details_delivery (delivery_deadline);

ALTER TABLE delivery_addresses
    ADD SPATIAL INDEX idx_delivery_location (location);

ALTER TABLE pickup_addresses
    ADD SPATIAL INDEX idx_pickup_location (location);

ALTER TABLE order_tracking
    ADD SPATIAL INDEX idx_order_current_location (current_location);

-- Dominio Zonas
ALTER TABLE zones
    ADD SPATIAL INDEX idx_zone_boundaries (boundaries),
   ADD SPATIAL INDEX idx_zone_center (center_point),
   ADD INDEX idx_zone_status (is_active);

ALTER TABLE zone_coverage
    ADD SPATIAL INDEX idx_coverage_area (coverage_area);

ALTER TABLE driver_zones
    ADD INDEX idx_driver_zones_driver (driver_id),
   ADD INDEX idx_driver_zones_zone (zone_id),
   ADD INDEX idx_driver_zones_status (is_active);

-- Dominio Pagos
ALTER TABLE payment_transactions
    ADD INDEX idx_payments_order (order_id),
   ADD INDEX idx_payments_method (payment_method_id),
   ADD INDEX idx_payments_status (status),
   ADD INDEX idx_payments_created (created_at);

ALTER TABLE invoices
    ADD INDEX idx_invoices_company (company_id),
   ADD INDEX idx_invoices_status (status),
   ADD INDEX idx_invoices_dates (period_start, period_end);

-- Dominio Eventos y Notificaciones
ALTER TABLE system_events
    ADD INDEX idx_events_type (event_type),
   ADD INDEX idx_events_source (source, source_id),
   ADD INDEX idx_events_occurred (occurred_at);

ALTER TABLE notifications
    ADD INDEX idx_notifications_user (user_id),
   ADD INDEX idx_notifications_template (template_id),
   ADD INDEX idx_notifications_status (is_read),
   ADD INDEX idx_notifications_type (type);

-- RELACIONES

-- Dominio Usuarios
ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id),
   ADD CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id);

ALTER TABLE user_profiles
    ADD CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_sessions
    ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id);

-- Dominio Empresas
ALTER TABLE company_branches
    ADD CONSTRAINT fk_branches_company FOREIGN KEY (company_id) REFERENCES companies(id),
   ADD CONSTRAINT fk_branches_zone FOREIGN KEY (zone_id) REFERENCES zones(id);

ALTER TABLE company_users
    ADD CONSTRAINT fk_company_users_user FOREIGN KEY (user_id) REFERENCES users(id),
   ADD CONSTRAINT fk_company_users_company FOREIGN KEY (company_id) REFERENCES companies(id);

ALTER TABLE company_addresses
    ADD CONSTRAINT fk_company_addresses FOREIGN KEY (company_id) REFERENCES companies(id);

-- Dominio Pedidos
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_company FOREIGN KEY (company_id) REFERENCES companies(id),
   ADD CONSTRAINT fk_orders_branch FOREIGN KEY (branch_id) REFERENCES company_branches(id),
   ADD CONSTRAINT fk_orders_client FOREIGN KEY (client_id) REFERENCES users(id),
   ADD CONSTRAINT fk_orders_driver FOREIGN KEY (driver_id) REFERENCES drivers(user_id);

ALTER TABLE order_details
    ADD CONSTRAINT fk_order_details FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE package_details
    ADD CONSTRAINT fk_package_details FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE delivery_addresses
    ADD CONSTRAINT fk_delivery_addresses FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE pickup_addresses
    ADD CONSTRAINT fk_pickup_addresses FOREIGN KEY (order_id) REFERENCES orders(id);

-- Dominio Zonas
ALTER TABLE adjacent_zones
    ADD CONSTRAINT fk_adjacent_zones_zone FOREIGN KEY (zone_id) REFERENCES zones(id),
   ADD CONSTRAINT fk_adjacent_zones_adjacent FOREIGN KEY (adjacent_zone_id) REFERENCES zones(id);

ALTER TABLE driver_zones
    ADD CONSTRAINT fk_driver_zones_driver FOREIGN KEY (driver_id) REFERENCES drivers(user_id),
   ADD CONSTRAINT fk_driver_zones_zone FOREIGN KEY (zone_id) REFERENCES zones(id);

-- Dominio Pagos
ALTER TABLE payment_transactions
    ADD CONSTRAINT fk_payments_order FOREIGN KEY (order_id) REFERENCES orders(id),
   ADD CONSTRAINT fk_payments_method FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id);

ALTER TABLE invoice_items
    ADD CONSTRAINT fk_invoice_items_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id),
   ADD CONSTRAINT fk_invoice_items_order FOREIGN KEY (order_id) REFERENCES orders(id);

-- Dominio Notificaciones
ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id),
   ADD CONSTRAINT fk_notifications_template FOREIGN KEY (template_id) REFERENCES notification_templates(id);

ALTER TABLE notification_preferences
    ADD CONSTRAINT fk_notification_prefs_user FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE notification_devices
    ADD CONSTRAINT fk_notification_devices_user FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE driver_availability
    ADD CONSTRAINT fk_driver_availability_driver
        FOREIGN KEY (driver_id) REFERENCES drivers(user_id),
    ADD CONSTRAINT fk_driver_availability_zone
    FOREIGN KEY (current_zone_id) REFERENCES zones(id);

ALTER TABLE event_logs
    ADD CONSTRAINT fk_event_logs_event
        FOREIGN KEY (event_id) REFERENCES system_events(id);

-- Dominio de almacen
ALTER TABLE warehouse
    ADD FOREIGN KEY (zone_id) REFERENCES zones(id);

ALTER TABLE package_warehouse_tracking
    ADD FOREIGN KEY (order_id) REFERENCES orders(id),
    ADD FOREIGN KEY (warehouse_id) REFERENCES warehouse(id),
    ADD FOREIGN KEY (collector_id) REFERENCES users(id);

ALTER TABLE warehouse_inventory
    ADD FOREIGN KEY (warehouse_id) REFERENCES warehouse(id),
    ADD FOREIGN KEY (order_id) REFERENCES orders(id);

-- Inicializacion de roles

INSERT INTO roles (id, name, description) VALUES
                                              (UUID(), 'ADMIN', 'Administrador del sistema'),
                                              (UUID(), 'COMPANY_USER', 'Usuario de empresa cliente'),
                                              (UUID(), 'DRIVER', 'Repartidor'),
                                              (UUID(), 'WAREHOUSE_STAFF', 'Personal de almacén'),
                                              (UUID(), 'COLLECTOR', 'Recolector de paquetes');