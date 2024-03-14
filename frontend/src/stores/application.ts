import { writable, type Writable } from "svelte/store";

export const dryRunMode: Writable<boolean> = writable(false);
export const showSummary: Writable<boolean> = writable(false);