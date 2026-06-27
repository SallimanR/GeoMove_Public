export function haversineDistance(coord1: [number, number], coord2: [number, number]): number {
	const R = 6371000; // Earth radius in meters
	const lat1 = degreesToRadians(coord1[1]);
	const lat2 = degreesToRadians(coord2[1]);
	const deltaLat = degreesToRadians(coord2[1] - coord1[1]);
	const deltaLng = degreesToRadians(coord2[0] - coord1[0]);

	const a =
		Math.sin(deltaLat / 2) ** 2 +
		Math.cos(lat1) * Math.cos(lat2) * Math.sin(deltaLng / 2) ** 2;
	const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
	return R * c;
}

export function degreesToRadians(degrees: number): number {
	return (degrees * Math.PI) / 180;
}
