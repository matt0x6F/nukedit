<script lang="ts">
	import '../app.pcss';

	import { Icon } from '@steeze-ui/svelte-icon'
	import { Home, Settings, Bomb, FlaskRound } from '@steeze-ui/lucide-icons'
	import { AppShell, AppBar, AppRail, AppRailAnchor, LightSwitch, Toast, Modal, type ModalComponent, initializeStores, getModalStore, type ModalSettings } from '@skeletonlabs/skeleton';
	
	import { page } from '$app/stores';
	import { active } from '../stores/account';
	import { dryRunMode } from '../stores/application';
	import AddAccountForm from '$lib/AddAccountForm.svelte';

	initializeStores();

	let activeAccount = "";

	active.subscribe((value) => {
		activeAccount = value;
	});

	const modalRegistery: Record<string, ModalComponent> = {
		addAccountModal: { ref: AddAccountForm }
	};
</script>

<Toast />

<Modal components={modalRegistery} />

<!-- App Shell -->
<AppShell slotPageContent="">
	<svelte:fragment slot="header">
		<!-- App Bar -->
		<AppBar gridColumns="grid grid-cols-3" gap="gap-4" spacing="space-x-4" padding="p-4" slotDefault="place-self-center m-auto" slotTrail="place-content-end w-full" slotLead="w-full">
			<svelte:fragment slot="lead">
				<strong class="text-xl uppercase">nukedit</strong>
			</svelte:fragment>
			<svelte:fragment slot="trail">
				{#if activeAccount === ""}
				<a
					class="btn btn-sm variant-ghost-surface"
					href="https://discord.gg/EXqV7W8MtY"
					target="_blank"
					rel="noreferrer"
				>
					Discord
				</a>
				<a
					class="btn btn-sm variant-ghost-surface"
					href="https://twitter.com/SkeletonUI"
					target="_blank"
					rel="noreferrer"
				>
					Twitter
				</a>
				<a
					class="btn btn-sm variant-ghost-surface"
					href="https://github.com/skeletonlabs/skeleton"
					target="_blank"
					rel="noreferrer"
				>
					GitHub
				</a>
				{:else}
				<span class="font-semibold">{activeAccount}</span>
				{/if}
				<LightSwitch />
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>
	<svelte:fragment slot="sidebarLeft">
		<AppRail>
			<!-- --- -->
			<AppRailAnchor href="/" name="Home" selected={$page.url.pathname === '/'} title="Home">
				<svelte:fragment slot="lead"><Icon src="{Home}" size="25px" class="color-gray-900 mx-auto" /></svelte:fragment>
				<span>Home</span>
			</AppRailAnchor>
			{#if activeAccount !== ""}
			<AppRailAnchor href="/nuke" name="Nuke" selected={$page.url.pathname === '/nuke'} title="Nuke">
				<svelte:fragment slot="lead"><Icon src="{Bomb}" size="25px" class="color-gray-900 mx-auto" /></svelte:fragment>
				<span>Nuke</span>
			</AppRailAnchor>
			{/if}
			
			<!-- --- -->
			<svelte:fragment slot="trail">
				{#if dryRunMode}
				<AppRailAnchor title="Dry run mode enabled">
					<Icon src="{FlaskRound}" size="25px" class="color-blue-500 mx-auto" />
				</AppRailAnchor>
				{/if}
				{#if activeAccount !== ""}
				<AppRailAnchor href="/settings" title="Settings">
					<Icon src="{Settings}" size="25px" class="color-gray-900 mx-auto"  />
				</AppRailAnchor>
				{/if}
			</svelte:fragment>
		</AppRail>
	</svelte:fragment>
	<!-- Page Route Content -->
	<slot />
</AppShell>
