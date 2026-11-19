import { get } from 'svelte/store';
import { goto } from '$app/navigation';
import { authStore } from './stores/authStore.js';

const BASE_URL = '/api';

export async function apiFetch(path, options = {}) {
	const { token } = get(authStore);
	const headers = {
		'Content-Type': 'application/json',
		...(options.headers || {})
	};
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const res = await fetch(`${BASE_URL}${path}`, {
		...options,
		headers
	});

	const contentType = res.headers.get('content-type') || '';
	const data = contentType.includes('application/json') ? await res.json() : null;

	if (res.status === 401) {
		authStore.clearSession();
		goto('/login');
		throw new Error(data?.error || 'Sesi berakhir, silakan login kembali');
	}

	if (!res.ok) {
		throw new Error(data?.error || `Permintaan gagal (${res.status})`);
	}

	return data;
}

export function formatRupiah(amount) {
	const num = Number(amount || 0);
	return new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		minimumFractionDigits: 0,
		maximumFractionDigits: 0
	}).format(num);
}

export function formatDate(iso) {
	if (!iso) return '-';
	return new Intl.DateTimeFormat('id-ID', {
		day: '2-digit',
		month: 'short',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	}).format(new Date(iso));
}

export function formatShortDate(iso) {
	if (!iso) return '-';
	return new Intl.DateTimeFormat('id-ID', {
		day: '2-digit',
		month: 'short',
		year: 'numeric'
	}).format(new Date(iso));
}

export function unitLabel(unitType) {
	switch (unitType) {
		case 'PER_KG':
			return 'Per Kg';
		case 'PER_ITEM':
			return 'Per Item';
		case 'PER_M2':
			return 'Per m²';
		default:
			return unitType;
	}
}

export function orderStatusLabel(status) {
	switch (status) {
		case 'PENDING':
			return 'Antre';
		case 'PROCESSING':
			return 'Proses';
		case 'READY':
			return 'Siap Diambil';
		case 'COMPLETED':
			return 'Selesai';
		case 'CANCELLED':
			return 'Dibatalkan';
		default:
			return status;
	}
}

export function paymentStatusLabel(status) {
	switch (status) {
		case 'PAID':
			return 'Lunas';
		case 'PARTIAL':
			return 'DP / Sebagian';
		case 'UNPAID':
			return 'Belum Lunas';
		default:
			return status;
	}
}
