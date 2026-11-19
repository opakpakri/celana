const cache = new Map();

export function getCache(key) {
	return cache.has(key) ? cache.get(key) : null;
}

export function setCache(key, value) {
	cache.set(key, value);
}

export function clearCache() {
	cache.clear();
}
