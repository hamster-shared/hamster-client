// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import { context } from '../models';
import { deploy } from '../models';

export function QueryGraphStatus(arg1: Array<string>): Promise<number | Error>;

export function SaveDeployInfo(arg1: number, arg2: string): Promise<boolean | Error>;

export function WailsInit(arg1: context.Context): Promise<Error>;

export function DeployTheGraph(arg1: number, arg2: string): Promise<boolean | Error>;

export function GetDeployInfo(arg1: number): Promise<deploy.DeployParameter | Error>;
