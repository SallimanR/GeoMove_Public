import type { paths } from "./generated/api.geolocation"

export type MovingDriver = paths["/drivers/moving/closest"]["get"]["responses"]["200"]["content"]["application/json"][number]
