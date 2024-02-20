import { sveltekit } from '@sveltejs/kit/vite';
import path from 'path';
import { defineConfig } from 'vite';

export default defineConfig({
	server: {
		fs: {
			allow: ['.']
		}
	},
	plugins: [sveltekit()],
	resolve: {
		alias: {
			'@': path.resolve(__dirname, './'), 
		}
	}
});
