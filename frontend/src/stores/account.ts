import type { models } from "$lib/wailsjs/go/models";
import { writable, type Writable } from "svelte/store";

// username
export const active: Writable<string> = writable("");
// account list
export const available: Writable<models.Account[]> = writable([]);