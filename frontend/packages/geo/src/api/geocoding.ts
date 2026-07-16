import axios from "axios"

const API_BASE = import.meta.env.VITE_GEO_SEARCH_API

export interface SearchResult {
	properties: Partial<{
		state: string,
		city: string;
		name: string;
		street: string;
		housenumber: string;
	}>;
	geometry: {
		type: string
		coordinates: [number, number];
	};

}

export interface SearchResultList {
	features: SearchResult[];
}

export async function getMapSearch(request: string, lat: number, lon: number): Promise<SearchResultList> {
	const result = await axios.get<SearchResultList>(`${API_BASE}/api/?q=${request}&lon=${lon}&lat=${lat}`, {
		headers: {
			"Accept-Language": "ru"
		}

	})
	return result.data
}

export async function getReverseGeocoding(lat: number, lon: number): Promise<SearchResult> {
	interface ReverseGeocodingResponse {
		features: SearchResult[]
	}

	const result = await axios.get<ReverseGeocodingResponse>(`${API_BASE}/reverse?lon=${lon}&lat=${lat}&limit=1&radius=1`, {
		headers: {
			"Accept-Language": "ru"
		}
	})
	return result.data.features[0] ?? null
}
