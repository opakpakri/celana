import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { clearCache } from './dataCache.js';

const TOKEN_KEY = 'celana_token';
const USER_KEY = 'celana_user';

function createAuthStore() {
	const initialToken = browser ? localStorage.getItem(TOKEN_KEY) || '' : '';
	let initialUser = null;
	if (browser) {
		try {
			const raw = localStorage.getItem(USER_KEY);
			if (raw) initialUser = JSON.parse(raw);
		} catch {
			initialUser = null;
		}
	}

	const { subscribe, set, update } = writable({
		token: initialToken,
		user: initialUser
	});

	return {
		subscribe,
		setSession(token, user) {
			if (browser) {
				localStorage.setItem(TOKEN_KEY, token);
				localStorage.setItem(USER_KEY, JSON.stringify(user));
			}
			set({ token, user });
		},
		updateUser(user) {
			if (browser) {
				localStorage.setItem(USER_KEY, JSON.stringify(user));
			}
			update((state) => ({ ...state, user }));
		},
		clearSession() {
			clearCache();
			if (browser) {
				localStorage.removeItem(TOKEN_KEY);
				localStorage.removeItem(USER_KEY);
			}
			set({ token: '', user: null });
		}
	};
}

export const authStore = createAuthStore();
