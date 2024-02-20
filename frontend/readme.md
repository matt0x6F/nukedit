This installation tracks https://wails.io/docs/guides/sveltekit#3-important-notes

## Notes

### Server files will cause build failures.

- `+layout.server.ts`, `+page.server.ts`, `+server.ts` or any file with "server" in the name will fail to build as all routes are prerendered.

### The Wails runtime unloads with full page navigations!

- Anything that causes full page navigations: `window.location.href = '/<some>/<page>'` or Context menu reload when using wails dev. What this means is that you can end up losing the ability to call any runtime breaking the app. There are two ways to work around this.
- Use `import { goto } from '$app/navigation'` then call `goto('/<some>/<page>')` in your `+page.svelte`. This will prevent a full page navigation.
- If full page navigation can't be prevented the Wails runtime can be added to all pages by adding the below into the `<head>` of `myapp/frontend/src/app.html``

```
<head>
...
    <meta name="wails-options" content="noautoinject" />
    <script src="/wails/ipc.js"></script>
    <script src="/wails/runtime.js"></script>
...
</head>
```

See https://wails.io/docs/guides/frontend for more information.

### Inital data can be loaded and refreshed from +page.ts/+page.js to +page.svelte.

- `+page.ts`/`+page.js` works well with `load()` https://kit.svelte.dev/docs/load#page-data
- `invalidateAll()` in `+page.svelte` will call `load()` from `+page.ts`/`+page.js` https://kit.svelte.dev/docs/load#rerunning-load-functions-manual-invalidation.

### Error Handling

- Expected errors using Throw error works in +page.ts/+page.js with a +error.svelte page. https://kit.svelte.dev/docs/errors#expected-errors
- Unexpected errors will cause the application to become unusable. Only recovery option (known so far) from unexpected errors is to reload the app. To do this create a file `myapp/frontend/src/hooks.client.ts` then add the below code to the file.

```
import { WindowReloadApp } from '$lib/wailsjs/runtime/runtime' 
export async function handleError() {
    WindowReloadApp()
}
```

### Using Forms and handling functions

- The simplest way is to call a function from the form is the standard, bind:value your variables and prevent submission `<form method="POST" on:submit|preventDefault={handle}>`
- The more advanced way is to use:enhance (progressive enhancement) which will allow for convenient access to formData, formElement, submitter. The important note is to always `cancel()` the form which prevents server side behavior. https://kit.svelte.dev/docs/form-actions#progressive-enhancement Example:

```
<form method="POST" use:enhance={({cancel, formData, formElement, submitter}) => {
    cancel()
    console.log(Object.fromEntries(formData))
    console.log(formElement)
    console.log(submitter)
    handle()
}}>
```
