# Backend Implementation Plan - Data Master HadirYuk

## Scope
Implementasi 3 entitas data master: Shifts, Office Locations, Leave Types

## Files Created (16 files)

### Models (3 files)
| File | Table | Key Fields |
|------|-------|------------|
| `internal/models/shift.go` | shifts | id, name (unique), start_time, end_time, break_duration, color_code, total_hours |
| `internal/models/office_location.go` | office_locations | id, name, address, latitude, longitude, radius_meters, is_active |
| `internal/models/leave_type.go` | leave_types | id, name (unique), description, default_days, is_paid |

### DTOs (1 file)
| File | DTOs |
|------|------|
| `internal/dtos/master_dto.go` | ShiftRequest, ShiftDTO, LocationRequest, LocationDTO, LeaveTypeRequest, LeaveTypeDTO + converter functions |

### Repositories (3 files)
| File | Model |
|------|-------|
| `internal/repositories/shift_repository.go` | Shift |
| `internal/repositories/office_location_repository.go` | OfficeLocation |
| `internal/repositories/leave_type_repository.go` | LeaveType |

### Services (3 files)
| File | Methods |
|------|---------|
| `internal/services/shift_service.go` | ShiftCreate, ShiftGetAllPaginated, ShiftGetAllUnpaginated, ShiftGetByID, ShiftUpdate, ShiftDelete |
| `internal/services/location_service.go` | LocationCreate, LocationGetAllPaginated, LocationGetAllUnpaginated, LocationGetByID, LocationUpdate, LocationDelete |
| `internal/services/leave_type_service.go` | LeaveTypeCreate, LeaveTypeGetAllPaginated, LeaveTypeGetAllUnpaginated, LeaveTypeGetByID, LeaveTypeUpdate, LeaveTypeDelete |

### Handlers (3 files)
| File | Handlers |
|------|----------|
| `internal/handlers/shift_handler.go` | ShiftCreate, ShiftGetAll, ShiftGetByID, ShiftUpdate, ShiftDelete |
| `internal/handlers/location_handler.go` | LocationCreate, LocationGetAll, LocationGetByID, LocationUpdate, LocationDelete |
| `internal/handlers/leave_type_handler.go` | LeaveTypeCreate, LeaveTypeGetAll, LeaveTypeGetByID, LeaveTypeUpdate, LeaveTypeDelete |

### Routes (3 files)
| File | Endpoints | Permissions |
|------|-----------|-------------|
| `internal/routes/shift_route.go` | CRUD /api/shifts | shift.create, shift.index, shift.update, shift.delete |
| `internal/routes/location_route.go` | CRUD /api/locations | location.create, location.index, location.update, location.delete |
| `internal/routes/leave_type_route.go` | CRUD /api/leave/types | leave.manage-types |

## Files Modified (2 files)
| File | Change |
|------|--------|
| `internal/repositories/00_repository.go` | Added Shift, OfficeLocation, LeaveType to Repositories struct |
| `cmd/api/main.go` | Registered 3 new route groups |
| `cmd/migration/main.go` | Added models to AutoMigrate and DropTable |

## API Endpoints

### Shifts
| Method | Endpoint | Permission |
|--------|----------|------------|
| POST | /api/shifts | shift.create |
| GET | /api/shifts | shift.index |
| GET | /api/shifts/:id | shift.index |
| PUT | /api/shifts/:id | shift.update |
| DELETE | /api/shifts/:id | shift.delete |

### Locations
| Method | Endpoint | Permission |
|--------|----------|------------|
| POST | /api/locations | location.create |
| GET | /api/locations | location.index |
| GET | /api/locations/:id | location.index |
| PUT | /api/locations/:id | location.update |
| DELETE | /api/locations/:id | location.delete |

### Leave Types
| Method | Endpoint | Permission |
|--------|----------|------------|
| POST | /api/leave/types | leave.manage-types |
| GET | /api/leave/types | leave.manage-types |
| GET | /api/leave/types/:id | leave.manage-types |
| PUT | /api/leave/types/:id | leave.manage-types |
| DELETE | /api/leave/types/:id | leave.manage-types |
