export interface RealtimeDriver {
	id: number;
	coordinates: [number, number][];
	time: number;
	distance: number;
}

export interface RealtimeDriverPosition {
	id: number,
	position: [number, number]
	bearing: number,
}

export interface DriverAnimationInfo {
	segmentsPercentages: number[];
	segments: [number, number][];
	currentSegment: number;
	segmentPercentageElapsed: number;
}
