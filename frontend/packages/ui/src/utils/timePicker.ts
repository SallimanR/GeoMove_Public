export function buildTimeOptions(): string[] {
	const options: string[] = [];
	for (let h = 0; h < 24; h++) {
		for (let m = 0; m < 60; m += 30) {
			options.push(
				`${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}`,
			);
		}
	}
	return options;
}
