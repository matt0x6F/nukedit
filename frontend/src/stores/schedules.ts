import type { models } from "$lib/wailsjs/go/models";
import { writable, type Writable } from "svelte/store";

export const schedules: Writable<models.Schedule[]> = writable([]);