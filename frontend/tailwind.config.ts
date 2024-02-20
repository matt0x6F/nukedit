import { join } from 'path';
import type { Config } from 'tailwindcss';

// 1. Import the Skeleton plugin
import { skeleton } from '@skeletonlabs/tw-plugin';
import { nukeditTheme } from './theme';
import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';

const config = {
	// 2. Opt for dark mode to be handled via the class method
	darkMode: 'class',
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		join(require.resolve('@skeletonlabs/skeleton'), '../**/*.{html,js,svelte,ts}')
	],
	theme: {
		extend: {},
	},
	plugins: [
		forms,
		typography,
		// 4. Append the Skeleton plugin (after other plugins)
		skeleton({
			// Requires changing body tag on app.html
			themes: { custom: [nukeditTheme]}
		})
	]
} satisfies Config;

export default config;