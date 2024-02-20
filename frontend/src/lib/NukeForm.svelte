<script lang="ts">
    import * as runtime from '$lib/wailsjs/runtime'

    let scheduled: boolean = false;
    let cron: string = "";
    let posts: boolean = true;
    let comments: boolean = true;
    let replacementLength: number = 125;
    let maxAge: number = 0;
    let useAge: boolean = false;
    let maxScore: number = 0;
    let useScore: boolean = false;

    $: buttonText = scheduled ? "Schedule" : "Nuke";

    const OpenCronMaker = (event: any) => {
        runtime.BrowserOpenURL("http://www.cronmaker.com")
    }
</script>

<label class="label my-2">
    <div class="font-semibold mr-2">Schedule <input class="checkbox" bind:value="{scheduled}" type="checkbox" /></div>
    <p class="text-sm my-2">The schedule expressed as a <a href="http://www.cronmaker.com" class="text-sm text-primary-600 drop-shadow-md font-semibold" on:click|preventDefault={OpenCronMaker}>cron</a>. Use the linked website if you are not familiar with crons.</p>
    <input class="input my-2 max-w-fit" title="Schedule" type="text" bind:value={cron} disabled={!scheduled} placeholder="cron expression" />
</label>
<div class="my-2">
    <div class="font-semibold">Content</div>
    <p class="prose-sm">Content to delete; note: selecting both will match both against the same criteria.</p>
    <label class="flex items-center space-x-2 max-w-fit">
        <input class="checkbox" bind:value="{posts}" type="checkbox" checked />
        <p>Posts</p>
    </label>
    <label class="flex items-center space-x-2 max-w-fit">
        <input class="checkbox" bind:value="{comments}" type="checkbox" checked />
        <p>Comments</p>
    </label>
</div>
<label class="label my-2">
    <div class="font-semibold mr-2">Max age <input class="checkbox" bind:value="{useAge}" type="checkbox" /></div>
    <p class="prose-sm">Maximum age in days of content for it not to be deleted. Anything older than this age will be removed.</p>
    <input class="input my-2 max-w-fit" title="Max age" type="number" disabled={!useAge} bind:value={maxAge}/>
</label>
<label class="label my-2">
    <div class="font-semibold mr-2 max-w-fit">Max score <input class="checkbox" bind:value="{useScore}" type="checkbox" /></div>
    <p class="prose-sm">Maximum score of content for it not to be deleted. Anything less than this score will be removed.</p>
    <input class="input my-2 max-w-fit" title="Max age" type="number" disabled={!useScore} bind:value={maxScore}/>
</label>
<label class="label my-2">
    <div class="font-semibold mr-2">Replacement text length</div>
    <p class="prose-sm">Nukedit performs safe deletes on Reddit content. It overwrites the content and then deletes it. While this does not guarantee someone hasn't cached the comment or post, anyone following Reddit's API Terms of Service will receive randomly generated text and then a delete.</p>
    <input class="input my-2 max-w-fit" title="Replacement text length" type="number" bind:value={replacementLength} />
</label>
<button type="button" class="btn variant-filled-primary">
    {buttonText}
</button>