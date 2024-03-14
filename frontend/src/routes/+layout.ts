import { DryRun } from "$lib/wailsjs/go/main/App";
import { dryRunMode } from "../stores/application";
import type { LayoutLoad } from "./$types";

export const prerender = true
export const ssr = false

export const load: LayoutLoad = async ({params}) => {
    await DryRun().then((value) => {
        dryRunMode.set(value)
    })
}