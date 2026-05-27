# OpenStreetMap Implementation Plan — HadirYuk

> For: Office Location form (lat, lng, radius input)

---

## Overview

Replace manual coordinate input with an interactive map using **Leaflet + OpenStreetMap** (free, no API key required).

## Tech Stack

| Package | Purpose |
|---------|---------|
| `leaflet` | Map rendering engine |
| `@types/leaflet` | TypeScript definitions |
| `leaflet-draw` | Draw circle for radius visualization |

## Files to Create

| File | Purpose |
|------|---------|
| `fe/src/components/utils/FormMap.vue` | Reusable map component for location input |

## Component Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `modelValue` | `{ lat: number, lng: number, radius?: number = 0 }` | - | v-model binding |
| `label` | `string` | `'Lokasi'` | Label text |
| `height` | `number` | `400` | Map container height |

## Component Behavior

1. **Initial State**: Show map centered on Indonesia (lat: -2.5, lng: 118.0)
2. **Click on Map**: Place marker, update lat/lng in form
3. **Radius Slider**: Below map, slider 50-500m, shows circle overlay on map
4. **Current Location** button: Use browser geolocation API

## Integration Points

Used in:
- `fe/src/pages/master/locations/FormModal.vue` — replace lat/lng/radius inputs

## Installation

```bash
cd fe
yarn add leaflet leaflet-draw
yarn add -D @types/leaflet @types/leaflet-draw
```

## Implementation Steps

1. Install dependencies
2. Create `FormMap.vue` component
3. Import Leaflet CSS in component
4. Update `locations/FormModal.vue` to use `FormMap`
5. Test on desktop and mobile

## Estimated Effort

~2 days (1 day BE + FE integration, 1 day testing + polish)
