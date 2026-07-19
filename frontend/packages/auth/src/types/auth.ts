import type { paths } from "./generated/api.auth.ts";

export type AuthUser = NonNullable<
  paths["/auth/me"]["get"]["responses"]["200"]["content"]["application/json"]["user"]
>;

export type OAuthCallbackResponse =
  paths["/auth/{provider}/callback"]["post"]["responses"]["200"]["content"]["application/json"];

export type OAuthCallbackRequest =
  paths["/auth/{provider}/callback"]["post"]["requestBody"]["content"]["application/json"];
