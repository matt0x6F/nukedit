<script lang="ts">
	import { Icon } from "@steeze-ui/svelte-icon";
    import { available as accounts, active, available } from "../../stores/account";
	import { ArrowLeftRight, Delete } from "@steeze-ui/lucide-icons";
	import { getModalStore, getToastStore, type ModalSettings, type ToastSettings } from "@skeletonlabs/skeleton";
	import { DeleteAccount, GetAccounts, SetActive } from "$lib/wailsjs/go/api/Accounts";

    const toastStore = getToastStore();
    const modalStore = getModalStore();

    const modal: ModalSettings = {
        type: 'component',
        component: 'addAccountModal',
    };

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

    const switchToAccount = async (username: string) => {
        active.set(username);
        await SetActive(username).then(() => {
            const t: ToastSettings = {
                message: 'Switched to account ' + username,
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-success'
            };
            toastStore.trigger(t);
        }).catch((err: string) => {
            const t: ToastSettings = {
                message: 'Error switching account: ' + err,
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-error'
            };
            toastStore.trigger(t);
        });
    }

    const deleteAccount = async (username:string) => {
        await DeleteAccount(username).then(() => {
            const t: ToastSettings = {
                message: 'Deleted account for ' + username,
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-success'
            };
            toastStore.trigger(t);
            refreshAccounts();
        }).catch((err: string) => {
            const t: ToastSettings = {
                message: 'Error deleting account: ' + err,
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-error'
            };
            toastStore.trigger(t);
        });
    }
</script>

<div class="my-2 mx-2">
    <p class="text-sm my-2">Here you can change various settings about nukedit</p>
    <h2 class="text-lg font-bold">Accounts</h2>

    {#each $accounts as account}
        <div class="flex items-center justify-between my-2">
            <div class="flex items-center">
                <p class="ml-2">{account.username} 
                    <button on:click={() => deleteAccount(account.username)}><Icon src={Delete} size="20px" class="color-gray-900 ml-2 mb-1 inline-block" /></button> 
                    {#if account.username !== $active}
                    <button on:click={() => switchToAccount(account.username)}><Icon src={ArrowLeftRight} size="20px" class="color-gray-900 ml-2 mb-1 inline-block" /></button>
                    {/if}
                </p>
            </div>
        </div>
    {/each}

    <button type="button" class="btn btn-sm variant-filled" on:click={() => modalStore.trigger(modal)}>Add account</button>

</div>