// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {models} from '../models';
import {context} from '../models';

export function DeleteSchedule(arg1:number):Promise<void>;

export function GetScheduleByUsername(arg1:string):Promise<Array<models.Schedule>>;

export function GetSchedules():Promise<Array<models.Schedule>>;

export function SaveSchedule(arg1:models.Schedule):Promise<void>;

export function Startup(arg1:context.Context):Promise<void>;