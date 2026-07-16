#!/bin/bash
set -e

# Backend: Go
oapi-codegen -config api/openapi/auth/oapi-codegen.yaml api/openapi/auth/auth.yaml
oapi-codegen -config api/openapi/driver/driver_oapi-codegen.yaml api/openapi/driver/driver.yaml
oapi-codegen -config api/openapi/driver/freely_available_oapi-codegen.yaml api/openapi/driver/freely_available.yaml
oapi-codegen -config api/openapi/geolocation/oapi-codegen.yaml api/openapi/geolocation/geolocation.yaml

# Frontend: TypeScript
cd frontend && pnpm exec openapi-typescript ../api/openapi/driver/driver.yaml \
	-o packages/@domains/driver/src/types/generated/api.driver.ts

pnpm exec openapi-typescript ../api/openapi/driver/freely_available.yaml \
	-o packages/@domains/driver/src/types/generated/api.freely_available.ts

pnpm exec openapi-typescript ../api/openapi/geolocation/geolocation.yaml \
	-o packages/@domains/geolocation/src/types/generated/api.geolocation.ts
