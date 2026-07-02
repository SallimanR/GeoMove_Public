import axios from "axios"

const API_BASE = "http://localhost:2322"

interface ReverseGeocodingResponse {
	features: {
		properties: {
			name: string,
			street: string,
			housenumber: string,
		}
		geometry: {
			type: string,
			coordinates: [number, number]
		}
	}[]
}
export async function getReverseGeocoding(lat: number, lon: number): Promise<ReverseGeocodingResponse> {
	const result = (await axios.get<ReverseGeocodingResponse>(`${API_BASE}/reverse?lon=${lon}&lat=${lat}&limit=1&radius=1`)).data
	console.log(result)
	return result
}

export interface SearchResult {
	properties: Partial<{
		city: string;
		name: string;
		street: string;
		housenumber: string;
	}>;
	geometry: {
		coordinates: [number, number];
	};

}

export interface SearchResults {
	features: SearchResult[];
}

export async function getMapSearch(request: string, lat: number, lon: number): Promise<SearchResults> {
	const result = (await axios.get<SearchResults>(`${API_BASE}/api/?q=${request}&lon=${lon}&lat=${lat}`)).data
	return result
}
