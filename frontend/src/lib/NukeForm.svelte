<script lang="ts">
    import * as runtime from '$lib/wailsjs/runtime'
	import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
	import { SaveSchedule } from './wailsjs/go/api/Scheduler';
	import { Nuke } from './wailsjs/go/api/Reddit';
	import { active } from '../stores/account';
    import cron from 'cron-validate';
	import type { logger } from './models/log';
	import { showSummary } from '../stores/application';
	import { scrollBottom } from 'svelte-scrolling';
    import type { reddit } from '$lib/wailsjs/go/models';


    const toastStore = getToastStore();

    let username: string = "";
    let scheduled: boolean = false;
    let cronExpression: string = "";
    let posts: boolean = true;
    let comments: boolean = true;
    let replacementLength: number = 125;
    let maxAge: number = 0;
    let useAge: boolean = false;
    let minScore: number = 0;
    let useMinScore: boolean = false;
    

    let nukeResult: reddit.NukeResult | undefined = undefined;

    $: buttonText = scheduled ? "Schedule" : "Nuke";

    active.subscribe((value) => {
        username = value;
    })

    let events: logger.Event[] = [];

    runtime.EventsOn("nuke:log", (event: logger.Event) => {
        events = [...events, event];
        // TODO: auto-scroll isn't working
        scrollBottom();
    });

    let running: boolean = false;

    showSummary.subscribe((value) => {
        running = value;
    });

    const submitRequest = async () => {
        if (cronExpression === "" && scheduled) {
            const t: ToastSettings = {
                message: "Please enter a cron expression",
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-error'
            };

            toastStore.trigger(t);
            return;
        }

        const cronResult = cron(cronExpression);

        if (!cronResult.isValid() && scheduled) {
            const t: ToastSettings = {
                message: "Invalid cron expression",
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-error'
            };

            toastStore.trigger(t);
            return;
        }

        if (!posts && !comments) {
            const t: ToastSettings = {
                message: "Please select at least one content type",
                timeout: 5000,
                hoverable: true,
                background: 'variant-filled-error'
            };

            toastStore.trigger(t);
            return;
        }

        if (scheduled) {
            await SaveSchedule({
                id: -1,
                username: username,
                cronExpression: cronExpression,
                posts: posts,
                comments: comments,
                maxAge: maxAge,
                useMaxAge: useAge,
                minScore: minScore,
                useMinScore: useMinScore,
                replacementTextLength: replacementLength,
            }).then((response) => {
                // reset all field values
                scheduled = false;
                cronExpression = "";
                posts = true;
                comments = true;
                replacementLength = 125;
                maxAge = 0;
                useAge = false;
                minScore = 0;
                useMinScore = false;

                const t: ToastSettings = {
                    message: "Nuke scheduled successfully",
                    timeout: 5000,
                    hoverable: true,
                    background: 'variant-filled-success'
                };

                toastStore.trigger(t);
            }).catch((error) => {
                const t: ToastSettings = {
                    message: 'Error while scheduling nuke: ' + error,
                    timeout: 5000,
                    hoverable: true,
                    background: 'variant-filled-error'
                };

                toastStore.trigger(t);
            });
        } else {
            showSummary.set(true);

            await Nuke({
                scheduled: scheduled,
                cronExpression: cronExpression,
                posts: posts,
                comments: comments,
                maxAge: maxAge,
                useMaxAge: useAge,
                minScore: minScore,
                useMinScore: useMinScore,
                replacementTextLength: replacementLength,
            }).then((response: reddit.NukeResult) => {
                // reset all field values
                scheduled = false;
                cronExpression = "";
                posts = true;
                comments = true;
                replacementLength = 125;
                maxAge = 0;
                useAge = false;
                minScore = 0;
                useMinScore = false;

                nukeResult = response;
                
                const t: ToastSettings = {
                    message: "Nuke executed successfully",
                    autohide: false,
                    hoverable: true,
                    background: 'variant-filled-success'
                };

                toastStore.trigger(t);
            }).catch((error) => {
                const t: ToastSettings = {
                    message: 'Error while executing nuke: ' + error,
                    timeout: 5000,
                    hoverable: true,
                    background: 'variant-filled-error'
                };

                toastStore.trigger(t);
            });
        }
    }

    function resetSummary() {
        nukeResult = undefined;
        showSummary.set(false);
    }

    const OpenCronMaker = (event: any) => {
        runtime.BrowserOpenURL("http://www.cronmaker.com")
    }
</script>

{#if running}
    <!-- Shamelessly borrowed from here: https://tailwindcomponents.com/component/terminal -->
    
    <div class="w-full h-full" >
        {#if nukeResult !== undefined && showSummary}
            <p>Deleted {nukeResult.postsDeleted} posts and {nukeResult.commentsDeleted} comments</p>
            <!-- Button to reset -->
            <button type="button" class="btn variant-filled-primary" on:click={resetSummary}>
                Reset
            </button>
        {/if}

        <div class="coding inverse-toggle px-5 pt-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased 
                    bg-gray-800  pb-6 rounded-lg leading-normal overflow-hidden">
            <div class="top mb-2 flex">
                <div class="h-3 w-3 bg-red-500 rounded-full"></div>
                <div class="ml-2 h-3 w-3 bg-orange-300 rounded-full"></div>
                <div class="ml-2 h-3 w-3 bg-green-500 rounded-full"></div>
            </div>
            <div class="mt-4 flex">
                <span class="text-green-400">nukedit:~$</span>
                <p class="flex-1 typing pl-2">
                    Executing nuke...
                    <br>
                    {#each events as event}
                        <p>{event.message}</p>
                    {/each}
                </p>
            </div>
        </div>
      </div>
{:else}
    <label class="label my-2">
        <div class="font-semibold mr-2">Schedule <input class="checkbox" bind:checked="{scheduled}" type="checkbox" /></div>
        <p class="text-sm my-2">The schedule expressed as a <a href="http://www.cronmaker.com" class="text-sm text-primary-600 drop-shadow-md font-semibold" on:click|preventDefault={OpenCronMaker}>cron</a>. Use the linked website if you are not familiar with crons.</p>
        <input class="input my-2 max-w-fit" title="Schedule" type="text" bind:value={cronExpression} disabled={!scheduled} placeholder="cron expression" />
    </label>
    <div class="my-2">
        <div class="font-semibold">Content</div>
        <p class="prose-sm">Content to delete; note: selecting both will match both against the same criteria.</p>
        <label class="flex items-center space-x-2 max-w-fit">
            <input class="checkbox" type="checkbox" bind:checked={posts} />
            <p>Posts</p>
        </label>
        <label class="flex items-center space-x-2 max-w-fit">
            <input class="checkbox" type="checkbox" bind:checked={comments} />
            <p>Comments</p>
        </label>
    </div>
    <label class="label my-2">
        <div class="font-semibold mr-2">Max age <input class="checkbox" bind:checked="{useAge}" type="checkbox" /></div>
        <p class="prose-sm">Maximum age in days of content for it not to be deleted. Anything older than this age will be removed.</p>
        <input class="input my-2 max-w-fit" title="Max age" type="number" disabled={!useAge} bind:value={maxAge}/>
    </label>
    <label class="label my-2">
        <div class="font-semibold mr-2 max-w-fit">Min score <input class="checkbox" bind:checked="{useMinScore}" type="checkbox" /></div>
        <p class="prose-sm">Minimum score of content for it not to be deleted. Anything less than this score will be removed.</p>
        <input class="input my-2 max-w-fit" title="Max age" type="number" disabled={!useMinScore} bind:value={minScore}/>
    </label>
    <label class="label my-2">
        <div class="font-semibold mr-2">Replacement text length</div>
        <p class="prose-sm">Nukedit performs safe deletes on Reddit content. It overwrites the content and then deletes it. While this does not guarantee someone hasn't cached the comment or post, anyone following Reddit's API Terms of Service will receive randomly generated text and then a delete.</p>
        <input class="input my-2 max-w-fit" title="Replacement text length" type="number" bind:value={replacementLength} />
    </label>
    <button type="button" class="btn variant-filled-primary" on:click={() => {submitRequest()}}>
        {buttonText}
    </button>
{/if}