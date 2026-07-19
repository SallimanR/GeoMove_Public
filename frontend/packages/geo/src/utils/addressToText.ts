import { SearchResult } from "../api/geocoding";

export function addressToText(
	address: SearchResult,
): string {
	if (!address) return ""
	return [
		address.properties.name,
		address.properties.street,
		address.properties.housenumber,
		address.properties.city,
	]
		.filter(Boolean)
		.join(", ");
}
