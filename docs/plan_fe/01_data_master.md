# Frontend Implementation Plan - Data Master HadirYuk

## Scope
Implementasi UI CRUD untuk 3 entitas data master: Shifts, Office Locations, Leave Types

## Files Created (9 files)

### Stores (3 files)
| File | Endpoint | Entity |
|------|----------|--------|
| `src/stores/shift.ts` | /shifts | IShift, IShiftPayload |
| `src/stores/location.ts` | /locations | ILocation, ILocationPayload |
| `src/stores/leaveType.ts` | /leave/types | ILeaveType, ILeaveTypePayload |

### Pages - Shifts (2 files)
| File | Purpose |
|------|---------|
| `src/pages/master/shifts/IndexView.vue` | Daftar shift dengan card grid, loading/empty state, pagination |
| `src/pages/master/shifts/FormModal.vue` | Form: name, start_time, end_time, break_duration, color_code |

### Pages - Locations (2 files)
| File | Purpose |
|------|---------|
| `src/pages/master/locations/IndexView.vue` | Daftar lokasi dengan card grid, active badge, koordinat |
| `src/pages/master/locations/FormModal.vue` | Form: name, address, latitude, longitude, radius_meters, is_active |

### Pages - Leave Types (2 files)
| File | Purpose |
|------|---------|
| `src/pages/master/leave-types/IndexView.vue` | Daftar jenis cuti dengan card grid, paid badge, default days |
| `src/pages/master/leave-types/FormModal.vue` | Form: name, description, default_days, is_paid |

## Files Modified (2 files)
| File | Change |
|------|--------|
| `src/router/index.ts` | Added routes: /master/shifts, /master/locations, /master/leave-types |
| `src/stores/index.ts` | Added barrel exports for new stores |

## Routes

| Path | Name | Component |
|------|------|-----------|
| /master/shifts | MasterShifts | master/shifts/IndexView.vue |
| /master/locations | MasterLocations | master/locations/IndexView.vue |
| /master/leave-types | MasterLeaveTypes | master/leave-types/IndexView.vue |

## UI Components Used
- UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton, UiModal
- FormInput (text, number, time, color, textarea, checkbox)
- PhPlus, PhPencil, PhTrash (Phosphor icons)

## Validation
- Vuelidate dengan required, minLength, minValue, maxValue
- Cross-field validation (end_time > start_time untuk shift)
- Range validation (lat: -90 to 90, lng: -180 to 180, radius: 50-500)
