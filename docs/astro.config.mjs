// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'sopsv',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/frostyeti/sopsv' }],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Installation', slug: 'getting-started/installation' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Managing Vaults', slug: 'guides/vaults' },
						{ label: 'Managing Secrets', slug: 'guides/secrets' },
						{ label: 'Executing Commands', slug: 'guides/exec' },
					],
				},
				{
					label: 'Reference',
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
