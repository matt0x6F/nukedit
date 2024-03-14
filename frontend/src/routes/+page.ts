import type { PageLoad } from "./$types";
import type { models } from '$lib/wailsjs/go/models';
import { GetAccounts } from "$lib/wailsjs/go/api/Accounts";
import { available } from "../stores/account";

export const load: PageLoad = async ({params}) => {
    let accounts: models.Account[] = [];

    await GetAccounts().then((accts) => {
        available.set(accts);
    })
}