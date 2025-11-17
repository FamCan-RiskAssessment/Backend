# FamCan Backend - Entity Relationship Diagram (ERD)

## Database Schema Overview

### Core Tables & Relationships

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          AUTHENTICATION & AUTHORIZATION                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────────────────┐
│       users          │
├──────────────────────┤
│ id (PK)              │
│ phone (UNIQUE)       │
│ password_hash        │
│ first_name           │
│ last_name            │
│ status               │◄─────────┐
│ created_at           │          │
│ updated_at           │          │
└──────────────────────┘          │
         │                         │
         │ 1:N                     │
         ▼                         │
┌──────────────────────┐           │
│   user_roles         │           │
├──────────────────────┤           │
│ user_id (FK)         │───────────┘
│ role_id (FK)         │─────┐
│ assigned_at          │     │
└──────────────────────┘     │
                             │
                    ┌────────┘
                    │
         ┌──────────▼──────────┐
         │      roles          │
         ├─────────────────────┤
         │ id (PK)             │
         │ name (UNIQUE)       │
         │ description         │
         │ created_at          │
         └──────────┬──────────┘
                    │
                    │ 1:N
                    ▼
         ┌─────────────────────┐
         │ role_permissions    │
         ├─────────────────────┤
         │ role_id (FK)        │
         │ permission_id (FK)  │
         └──────────┬──────────┘
                    │
                    │
         ┌──────────▼──────────┐
         │   permissions       │
         ├─────────────────────┤
         │ id (PK)             │
         │ name (UNIQUE)       │
         │ description         │
         │ category            │
         │ created_at          │
         └─────────────────────┘


┌─────────────────────────────────────────────────────────────────────────────┐
│                            FORM MANAGEMENT                                  │
└─────────────────────────────────────────────────────────────────────────────┘

         ┌──────────────────────┐
         │       forms          │
         ├──────────────────────┤
         │ id (PK)              │
         │ user_id (FK)         │◄──────────┐
         │ filled_by_operator_id(FK)       │
         │ status               │          │
         │ created_at           │          │
         │ submitted_at         │          │
         │ updated_at           │          │
         └──────────┬───────────┘          │
                    │                     │
        ┌───────────┼────────┬────────────┤
        │           │        │            │
        │ 1:1       │ 1:1    │ 1:1        │
        ▼           ▼        ▼            │
┌──────────────┐ ┌──────────────┐ ┌──────────────────┐
│  basic_info  │ │general_health│ │ family_history   │
├──────────────┤ ├──────────────┤ ├──────────────────┤
│ form_id (PK) │ │form_id (PK)  │ │ form_id (PK)     │
│ gender       │ │has_alcohol   │ │ mother_cancer    │
│ birth_date   │ │smoking_status│ │ father_cancer    │
│ created_at   │ │ ...more      │ │ siblings_cancer  │
│ updated_at   │ │ health data  │ │ ...more history  │
└──────────────┘ └──────────────┘ │ created_at       │
        │                          │ updated_at       │
        │                          └──────────────────┘
        │                                  │
        │                                  │
        └──────────────┬───────────────────┘
                       │ (Composite Foreign Key: form_id)
                       │
         ┌─────────────▼──────────────┐
         │   cancers                  │
         ├────────────────────────────┤
         │ id (PK)                    │
         │ form_id (FK)               │
         │ cancer_type                │
         │ diagnosed_date             │
         │ age_at_diagnosis           │
         │ description                │
         │ created_at                 │
         │ updated_at                 │
         └────────────────────────────┘


         ┌──────────────────────────────────┐
         │     mamography_info              │
         ├──────────────────────────────────┤
         │ form_id (FK)                     │
         │ picture_path (S3 URL)            │
         │ uploaded_at                      │
         │ created_at                       │
         │ updated_at                       │
         └──────────────────────────────────┘


┌─────────────────────────────────────────────────────────────────────────────┐
│                         RISK CALCULATIONS                                   │
└─────────────────────────────────────────────────────────────────────────────┘

         ┌──────────────────────┐
         │ premm5_results       │
         ├──────────────────────┤
         │ form_id (FK)         │
         │ mlh1                 │
         │ msh2                 │
         │ msh6                 │
         │ pms2                 │
         │ epcam                │
         │ calculated_at        │
         └──────────────────────┘
                    ▲
                    │
                    │ (Similar pattern for other models)
                    │
         ┌──────────┴──────────┐
         │                     │
    ┌────▼────────────┐  ┌─────▼────────────┐
    │  bcra_results   │  │  gail_results    │
    ├─────────────────┤  ├──────────────────┤
    │ form_id (FK)    │  │ form_id (FK)     │
    │ mutation_risk   │  │ breast_cancer    │
    │ calculated_at   │  │ calculated_at    │
    └─────────────────┘  └──────────────────┘

         ┌──────────────────────┐
         │  plco_results        │
         ├──────────────────────┤
         │ form_id (FK)         │
         │ lung_cancer_risk     │
         │ calculated_at        │
         └──────────────────────┘


┌─────────────────────────────────────────────────────────────────────────────┐
│                         AUDIT & LOGGING                                     │
└─────────────────────────────────────────────────────────────────────────────┘

         ┌────────────────────────┐
         │   action_logs          │
         ├────────────────────────┤
         │ id (PK)                │
         │ user_id (FK)           │◄──────┐
         │ form_id (FK)           │       │
         │ action_type            │       │
         │ description            │       │
         │ ip_address             │       │
         │ user_agent             │       │
         │ created_at             │       │
         └────────────────────────┘       │
                                          │
                            (Points back to users table)
```

---

## Table Definitions

### `users`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| id | UUID | PRIMARY KEY | |
| phone | VARCHAR(20) | UNIQUE, NOT NULL | Iranian phone format |
| password_hash | VARCHAR(255) | | bcrypt/argon2 hashed |
| first_name | VARCHAR(100) | | |
| last_name | VARCHAR(100) | | |
| status | ENUM | NOT NULL, DEFAULT 'ACTIVE' | ACTIVE, INACTIVE, SUSPENDED |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | DEFAULT NOW() | |

### `roles`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| id | UUID | PRIMARY KEY | |
| name | VARCHAR(50) | UNIQUE, NOT NULL | PATIENT, OPERATOR, ADMIN |
| description | TEXT | | |
| created_at | TIMESTAMP | DEFAULT NOW() | |

### `permissions`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| id | UUID | PRIMARY KEY | |
| name | VARCHAR(100) | UNIQUE, NOT NULL | USER_MANAGEMENT_READ, FORM_CREATE, etc. |
| description | TEXT | | |
| category | VARCHAR(50) | | Authentication, FormManagement, Calculation, etc. |
| created_at | TIMESTAMP | DEFAULT NOW() | |

### `user_roles`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| user_id | UUID | FOREIGN KEY (users.id), PRIMARY KEY | |
| role_id | UUID | FOREIGN KEY (roles.id), PRIMARY KEY | |
| assigned_at | TIMESTAMP | DEFAULT NOW() | |

### `role_permissions`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| role_id | UUID | FOREIGN KEY (roles.id), PRIMARY KEY | |
| permission_id | UUID | FOREIGN KEY (permissions.id), PRIMARY KEY | |

### `forms`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| id | UUID | PRIMARY KEY | |
| user_id | UUID | FOREIGN KEY (users.id), NOT NULL | Patient who owns the form |
| filled_by_operator_id | UUID | FOREIGN KEY (users.id) | Null if self-filled, operatorID if filled by operator |
| status | ENUM | NOT NULL, DEFAULT 'DRAFT' | DRAFT, SUBMITTED, ACCEPTED, REJECTED |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| submitted_at | TIMESTAMP | | Null until form is submitted |
| updated_at | TIMESTAMP | DEFAULT NOW() | |

### `basic_info`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| gender | ENUM | NOT NULL | MALE, FEMALE, OTHER |
| birth_date | DATE | NOT NULL | |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | DEFAULT NOW() | |

### `general_health`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| has_alcohol | BOOLEAN | | |
| smoking_status | ENUM | | NEVER, FORMER, CURRENT |
| exercise_frequency | VARCHAR(50) | | |
| other_health_conditions | TEXT | | |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | DEFAULT NOW() | |

### `family_history`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| mother_cancer | TEXT | | Cancer type and age if applicable |
| father_cancer | TEXT | | Cancer type and age if applicable |
| siblings_cancer | TEXT | | Cancer type and age if applicable |
| maternal_relatives_cancer | TEXT | | |
| paternal_relatives_cancer | TEXT | | |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | DEFAULT NOW() | |

### `cancers`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| id | UUID | PRIMARY KEY | |
| form_id | UUID | FOREIGN KEY (forms.id), NOT NULL | 1:N relationship |
| cancer_type | ENUM | NOT NULL | BREAST, OVARIAN, COLON, LUNG, PANCREATIC, etc. |
| diagnosed_date | DATE | | |
| age_at_diagnosis | INTEGER | | |
| description | TEXT | | Additional notes |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | DEFAULT NOW() | |

### `mamography_info`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| picture_path | VARCHAR(500) | | S3 URL (e.g., s3://bucket/uuid.jpg) |
| uploaded_at | TIMESTAMP | | When file was uploaded to S3 |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | DEFAULT NOW() | |

### `premm5_results`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| mlh1 | DECIMAL(5,4) | | Probability (0-1) |
| msh2 | DECIMAL(5,4) | | Probability (0-1) |
| msh6 | DECIMAL(5,4) | | Probability (0-1) |
| pms2 | DECIMAL(5,4) | | Probability (0-1) |
| epcam | DECIMAL(5,4) | | Probability (0-1) |
| calculated_at | TIMESTAMP | | When calculation was performed |

### `bcra_results`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| mutation_risk | DECIMAL(5,4) | | Probability (0-1) |
| calculated_at | TIMESTAMP | | |

### `gail_results`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| breast_cancer_risk | DECIMAL(5,4) | | 5-year risk probability |
| calculated_at | TIMESTAMP | | |

### `plco_results`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| form_id | UUID | FOREIGN KEY (forms.id), PRIMARY KEY | 1:1 relationship |
| lung_cancer_risk | DECIMAL(5,4) | | Probability (0-1) |
| calculated_at | TIMESTAMP | | |

### `action_logs`
| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| id | UUID | PRIMARY KEY | |
| user_id | UUID | FOREIGN KEY (users.id), NOT NULL | User who performed action |
| form_id | UUID | FOREIGN KEY (forms.id) | Null if action not form-related |
| action_type | ENUM | NOT NULL | CREATE_FORM, UPDATE_FORM, ACCEPT_FORM, CALCULATE, etc. |
| description | TEXT | | Human-readable description |
| ip_address | VARCHAR(45) | | IPv4 or IPv6 |
| user_agent | VARCHAR(500) | | Browser/client info |
| created_at | TIMESTAMP | DEFAULT NOW() | |

---

## Key Relationships

### One-to-Many (1:N)
- **users** → **forms**: One patient has many forms
- **users** → **user_roles**: One user has many roles (typically 1 in practice)
- **users** → **action_logs**: One user performs many actions
- **roles** → **role_permissions**: One role has many permissions
- **forms** → **cancers**: One form can contain multiple cancer entries

### One-to-One (1:1)
- **forms** → **basic_info**: Every form has exactly one basic info section
- **forms** → **general_health**: Every form has exactly one health info section
- **forms** → **family_history**: Every form has exactly one family history section
- **forms** → **mamography_info**: Form may have one mamography upload
- **forms** → **premm5_results**: Form may have one calculation result
- **forms** → **bcra_results**: Form may have one calculation result
- **forms** → **gail_results**: Form may have one calculation result
- **forms** → **plco_results**: Form may have one calculation result

### Many-to-Many (M:N)
- **users** ↔ **roles**: Through `user_roles` junction table
- **roles** ↔ **permissions**: Through `role_permissions` junction table

---

## Indexes for Performance

```sql
-- Authentication & Authorization
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Forms & Related Data
CREATE INDEX idx_forms_user_id ON forms(user_id);
CREATE INDEX idx_forms_filled_by_operator_id ON forms(filled_by_operator_id);
CREATE INDEX idx_forms_status ON forms(status);
CREATE INDEX idx_forms_created_at ON forms(created_at);
CREATE INDEX idx_cancers_form_id ON cancers(form_id);

-- Audit Logging
CREATE INDEX idx_action_logs_user_id ON action_logs(user_id);
CREATE INDEX idx_action_logs_form_id ON action_logs(form_id);
CREATE INDEX idx_action_logs_created_at ON action_logs(created_at);
CREATE INDEX idx_action_logs_action_type ON action_logs(action_type);
```

---

## Constraints & Cascade Rules

```
forms → basic_info
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → general_health
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → family_history
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → cancers
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → mamography_info
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → premm5_results
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → bcra_results
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → gail_results
  ON DELETE CASCADE
  ON UPDATE CASCADE

forms → plco_results
  ON DELETE CASCADE
  ON UPDATE CASCADE

users → forms
  ON DELETE RESTRICT  (Prevent deletion of user with forms)
  ON UPDATE CASCADE

users → action_logs
  ON DELETE RESTRICT
  ON UPDATE CASCADE
```

---

## Sample Data Flow Through Tables

### Patient Form Submission Flow
```
1. Patient logs in
   └─> Query: users WHERE phone = "09123456789"

2. Patient creates form
   └─> Insert: forms { user_id, status: DRAFT }

3. Patient fills basic info
   └─> Insert: basic_info { form_id, gender, birth_date }

4. Patient fills general health
   └─> Insert: general_health { form_id, has_alcohol, smoking_status }

5. Patient fills family history
   └─> Insert: family_history { form_id, mother_cancer, ... }

6. Patient adds cancer history
   └─> Insert: cancers { form_id, cancer_type, diagnosed_date }

7. Patient uploads mamography
   └─> Insert: mamography_info { form_id, picture_path }

8. Patient submits form
   └─> Update: forms SET status = SUBMITTED, submitted_at = NOW()
   └─> Insert: action_logs { user_id, form_id, action_type: SUBMIT_FORM }

9. Admin triggers calculation
   └─> Insert: premm5_results { form_id, mlh1, msh2, ... }
   └─> Insert: action_logs { user_id, form_id, action_type: CALCULATE }

10. Admin accepts form
    └─> Update: forms SET status = ACCEPTED
    └─> Insert: action_logs { user_id, form_id, action_type: ACCEPT_FORM }
```

---

## Notes

- **PostgreSQL 17** is used for all persistent data storage
- **UUIDs** are used as primary keys for better distribution and security
- **Timestamps** (created_at, updated_at) are automatically managed by GORM
- **Enums** provide type safety and data consistency
- **Soft deletes** are not used; hard deletes with CASCADE rules are preferred
- **Calculation results** are stored as separate tables to maintain normalization and support multiple model versions
- **Action logs** provide complete audit trail for compliance and debugging
