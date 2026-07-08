export function displayDistance(distance: number): string {
	distance = Math.floor(distance)

	if (distance < 1000) {
		return `${distance} м`
	}
	return `${distance / 1000} км`
}
