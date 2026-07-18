# GeoMove public repository

## geo APIs:
- https://geomove.online/style/style/style.json - maps style.json
- https://geomove.online/tiles - PMTiles API (tiles for maps)
- https://geomove.online/photon - Search and severse geocoding API
- https://geomove.online/routing - Routing and map matching API

## Use maps in Vue js app
```vue
<script setup lang="ts">
const styleApi = import.meta.env.VITE_STYLE_API; // your API style

</script>

<template>
  <div class="relative flex flex-col h-full">
    <Maps :styleApi="styleApi" />
    <MapsOverlayControls /> <!-- Search and route picker -->
  </div>
</template>
```
