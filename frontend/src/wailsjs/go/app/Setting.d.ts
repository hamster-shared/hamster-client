// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import { app } from '../models';
import { context } from '../models';

export function GetSetting(): Promise<app.Config | Error>;

export function InitP2pSetting(): Promise<boolean | Error>;

export function Setting(arg1: string, arg2: string): Promise<boolean | Error>;

export function WailsInit(arg1: context.Context): Promise<Error>;