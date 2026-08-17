// Leaflet's default marker icon resolves its image paths relative to the
// Leaflet package itself. Vite doesn't rewrite those paths automatically, so
// every L.marker() renders as a broken image unless the URLs are re-pointed
// at Vite's bundled asset URLs — do it once, globally, before any map mounts.
import L from 'leaflet'
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png'
import markerIcon from 'leaflet/dist/images/marker-icon.png'
import markerShadow from 'leaflet/dist/images/marker-shadow.png'

delete (L.Icon.Default.prototype as unknown as { _getIconUrl?: unknown })._getIconUrl

L.Icon.Default.mergeOptions({
  iconRetinaUrl: markerIcon2x,
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
})
