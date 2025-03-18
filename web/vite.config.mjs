import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
	plugins: [react(), tailwindcss()],
	css: {
		postcss: './postcss.config.js',
	},
	build: {
		outDir: 'build', 
		minify: false,
		terserOptions: {
			compress: false,
			mangle: false,
		},
	},
	server: {
		open: true, // Opens browser on dev server start
	},
	envDir: '../',
	// For backwards compatibility
	envPrefix: ['REACT_', 'VITE_'],
});
