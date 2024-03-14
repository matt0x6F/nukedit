import { get } from "svelte/store";
import type { PageLoad } from "./$types";
import { GetScheduleByUsername } from "$lib/wailsjs/go/api/Scheduler";
import { active } from "../../stores/account";
import { schedules } from "../../stores/schedules";

export const load: PageLoad = async ({params}) => {
    const activeUser: string = get(active);

    await GetScheduleByUsername(activeUser).then((storedSchedules) => {
        schedules.set(storedSchedules);
    });
}