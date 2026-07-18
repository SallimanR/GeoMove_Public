import { ReadableAtom } from "nanostores";
import { Component } from "vue";

export interface MovingPosition {
	id: number;
	position: [number, number];
	bearing: number;
}

export interface UseMovingIconLayerOptions {
	paths: ReadableAtom<MovingPath[]>;
	iconUrl: string;
	iconWidth: number;
	iconHeight: number;
	layerId?: string;
	iconSize?: number;
	popupComponent?: Component;
	onClick?: (id: number) => void;
	onHover?: (id: number) => void;
}

export interface MovingPath {
	id: number;
	coordinates: [number, number][];
	time: number;
	distance: number;
}

export interface CacheEntry {
	cumulativeDistances: number[];
	totalDistance: number;
	totalTimeSeconds: number;
}
