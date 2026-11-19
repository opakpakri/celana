/** @type {import('tailwindcss').Config} */
export default {
	darkMode: 'class',
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['"Plus Jakarta Sans"', 'Inter', 'system-ui', 'sans-serif']
			}
		}
	},
	plugins: []
};
