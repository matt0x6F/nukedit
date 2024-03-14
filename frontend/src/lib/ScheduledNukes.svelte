<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
    import { Trash2 as Trash, Pencil as Edit } from '@steeze-ui/lucide-icons';
    import type { models } from './wailsjs/go/models';
	import { DeleteSchedule, GetSchedules } from './wailsjs/go/api/Scheduler';
	import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';

    const toastStore = getToastStore();

    onMount(async () => {
        await refreshSchedules();
    });

    export let scheduled: models.Schedule[] = [];

    const deleteSchedule = async (id: number) => {
        await DeleteSchedule(id).then((response) => {
            const t: ToastSettings = {
                message: "Successfully deleted schedule",
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-success'
            };

            toastStore.trigger(t);
        }).catch((error) => {
            const t: ToastSettings = {
                message: "Error deleting message: " + error,
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-error'
            };

            toastStore.trigger(t);
        });

        await refreshSchedules();
    }

    const refreshSchedules = async () => {
        // refresh the schedules
        await GetSchedules().then((response) => {
            scheduled = response;
        }).catch((error) => {
            const t: ToastSettings = {
                message: "Error refreshing schedules: " + error,
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-error'
            };

            toastStore.trigger(t);
        });
    }

</script>

<!-- list scheduled nukes the user has already created -->
<div class="my-2">
    <div class="font-semibold">Scheduled nukes</div>
    <div class="table-container">
        <table class="table table-interactive">
            <thead>
                <tr>
                    <th class="table-cell-fit">Cron Expression</th>
                    <!-- make rows for all fields in a schedule -->
                    <th class="table-cell-fit">Posts</th>
                    <th class="table-cell-fit">Comments</th>
                    <th class="table-cell-fit">Max Age</th>
                    <th class="table-cell-fit">Min Score</th>
                    <th class="table-cell-fit">Replacement Length</th>
                    <th class="table-cell-fit"></th>
                </tr>
            </thead>
            <tbody>
            {#each scheduled as schedule}
                <tr>
                    <td class="table-cell-fit">{schedule.cronExpression}</td>
                    <td class="table-cell-fit items-center justify-center justify-items-center"><input class="checkbox" disabled={true} type="checkbox" checked={schedule.posts} /></td>
                    <td class="table-cell-fit"><input class="checkbox" disabled={true} type="checkbox" checked={schedule.comments} /></td>
                    <td class="table-cell-fit">
                        {#if schedule.useMaxAge}
                            {schedule.maxAge}
                        {:else}
                            -
                        {/if}
                    </td>
                    <td class="table-cell-fit">
                        {#if schedule.useMinScore}
                            {schedule.minScore}
                        {:else}
                            -
                        {/if}
                    </td>
                    <td class="table-cell-fit">{schedule.replacementTextLength}</td>
                    <td class="table-cell-fit">
                        <div class="flex space-x-2">
                            <button type="button" on:click={() => deleteSchedule(schedule.id)} class="btn-icon btn-icon-sm variant-filled">
                                <Icon src={Trash} size="16px" class="color-gray" />
                            </button>
                            <button type="button" class="btn-icon btn-icon-sm variant-filled"><Icon src={Edit} size="16px" class="color-gray" /></button>
                        </div>
                    </td>
                </tr>
            {/each}
            </tbody>

            {#if scheduled.length === 0}
                <tr>
                    <td>No scheduled nukes</td>
                </tr>
            {/if}
        </table>
    </div>
</div>