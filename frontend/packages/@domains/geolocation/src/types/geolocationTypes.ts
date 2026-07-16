import type { paths } from "./generated/api.geolocation"

export type MovingDriver = paths["/driver/moving/closest"]["get"]["responses"]["200"]["content"]["application/json"][number]
