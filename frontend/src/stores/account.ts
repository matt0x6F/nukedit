import type { models } from "$lib/wailsjs/go/models";
import { writable, type Writable } from "svelte/store";

export const active = writable("");
export const available: Writable<models.Account[]> = writable([]);