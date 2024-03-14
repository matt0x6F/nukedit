<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon'
	import { Save } from '@steeze-ui/lucide-icons'
	import { GetAccounts, SaveAccount, SetActive } from '$lib/wailsjs/go/api/Accounts';
	import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
	import { active } from '../stores/account';
	import { available } from '../stores/account';

	const toastStore = getToastStore();

	let clientID: string = '';
	let clientSecret: string = '';
	let username: string = '';
	let password: string = '';
	let requirePW: boolean = true;
	let selectedAccount: string = '';

	function validate(): boolean {
		if (clientID === '' || clientSecret === '' || username === '' || password === '') {
			const t: ToastSettings = {
				message: 'All fields are required',
				timeout: 5000,
				hoverable: true,
				background: 'variant-filled-error'
			};
			toastStore.trigger(t);
			return false
		}

		return true
	}

	const refreshAccounts = async () => {
		await GetAccounts().then((accts) => {
			available.set(accts);
		}).catch((err: string) => {
			const t: ToastSettings = {
				message: 'Error getting accounts: ' + err,
				timeout: 5000,
				hoverable: true,
				background: 'variant-filled-error'
			};
			toastStore.trigger(t);
		});

	}

	const setActiveAccount = async () => {
		await SetActive(selectedAccount).then(() => {
			const t: ToastSettings = {
				message: 'Set active account to ' + selectedAccount,
				timeout: 25,
				autohide: true,
				background: 'variant-filled-success'
			};
			toastStore.trigger(t);
			active.set(selectedAccount);
		}).catch((err: string) => {
			const t: ToastSettings = {
				message: 'Error setting active account: ' + err,
				timeout: 5000,
				hoverable: true,
				background: 'variant-filled-error'
			};
			toastStore.trigger(t);
		});
	}

	const saveAccount = async () => {
		let ok = validate();
		if (ok) {
			let setActive = false;

			await SaveAccount(
				clientID,
				clientSecret,
				username,
				password,
				requirePW
			).then(() => {
				const t: ToastSettings = {
					message: 'Saved account for ' + username,
					timeout: 5000,
					hoverable: true,
					background: 'variant-filled-success'
				};
				toastStore.trigger(t);

				setActive = true;
			}).catch((err: string) => {
				const t: ToastSettings = {
					message: 'Error saving account information: ' + err,
					timeout: 5000,
					hoverable: true,
					background: 'variant-filled-error'
				};
				toastStore.trigger(t);
			});

			await refreshAccounts();

			if (setActive) {
				selectedAccount = username;
				setActiveAccount();
			}
		}
	}
</script>

<div class="container h-full mx-auto flex justify-center items-center">
	<div class="space-y-10 text-center flex flex-col items-center">
		<h2 class="h2">Welcome to nukedit!</h2>
		
		{#if $available.length === 0}
		<div class="flex justify-center space-x-2 flex-col">
			<p>In order to get started, we need to collect some information.</p>
			<p class="font-semibold">This information remains private and does not leave the computer you're using it on.</p>
		</div>
		<div class="flex flex-col justify-center space-y-2 space-x-2">
			<label class="label">
				<span>Client ID</span>
				<input bind:value="{clientID}" required class="input" type="text" placeholder="Client ID" />
			</label>
			<label class="label">
				<span>Client Secret</span>
				<input bind:value="{clientSecret}" required class="input" type="text" placeholder="Client Secret" />
			</label>
			<label class="label">
				<span>Username</span>
				<input bind:value="{username}" required class="input" type="text" placeholder="Username" />
			</label>
			<label class="label">
				<span>Password</span>
				<input bind:value="{password}" required class="input" type="password" placeholder="Password" />
			</label>
			<label class="flex items-center space-x-2">
				<input class="checkbox" type="checkbox" bind:checked={requirePW} />
				<p>Require password on sign in</p>
			</label>
			<button type="button" class="btn variant-filled" on:click={() => saveAccount()}>
				<span><Icon src="{Save}" size="25px" theme="solid" class="color-gray-900" /></span>
				<span>Save</span>
			</button>
		</div>
		{/if}

		{#if $active === '' && $available.length > 0}
			<p class="font-semibold">Select an account to use</p>
			<select class="select" bind:value="{selectedAccount}">
				{#each $available as account}
					<option value="{account.username}">{account.username}</option>
				{/each}
			</select>
			{#each $available as account }
				{#if account.username === $active && account.require_password}
				<label class="label">
					<span>Password</span>
					<input bind:value="{password}" required class="input" type="password" placeholder="Password" />
				</label>
				{/if}
			{/each}
			<button type="button" class="btn-icon btn-icon-xl variant-filled-primary font-bold italic" on:click={() => {setActiveAccount()}}>Go!</button>
		{/if}

		{#if $active !== ''}
			<p class="font-semibold">Welcome back, {$active}!</p>
		{/if}
	</div>
</div>
