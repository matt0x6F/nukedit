<script lang="ts">
    import { Tab, TabGroup } from '@skeletonlabs/skeleton';
	import { Calendar, Plus } from '@steeze-ui/lucide-icons';
	import { Icon } from '@steeze-ui/svelte-icon';

    import NukeForm from '$lib/NukeForm.svelte';

    // page navigation
    let tabSet: number = 0;
</script>
<div class="my-2 mx-2">
    <p class="text-sm my-2">Nukes are a powerful tool and should be used with caution. They are not reversible and will delete content from Reddit.</p>

    <TabGroup>
        <Tab bind:group={tabSet} name="tab1" value={0}>
            <svelte:fragment slot="lead"><Icon src="{Plus}" size="25px" class="color-gray-900 mx-auto" /></svelte:fragment>
            <span>New</span>
        </Tab>
        <Tab bind:group={tabSet} name="tab2" value={1}>
            <svelte:fragment slot="lead"><Icon src="{Calendar}" size="25px" class="color-gray-900 mx-auto" /></svelte:fragment>
            <span>Scheduled</span>
        </Tab>

        <svelte:fragment slot="panel">
            {#if tabSet === 0}
                <NukeForm />
            {:else if tabSet === 1}
                <!-- list scheduled nukes the user has already created -->
                <div class="my-2">
                    <div class="font-semibold">Scheduled nukes</div>
                    <ul class="list-disc list-inside">
                        <li>None</li>
                    </ul>
                </div>
            {/if}
        </svelte:fragment>
    </TabGroup>
</div>