// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {wallet} from '../models';
import {context} from '../models';

export function DeleteWallet():void;

export function GetWalletInfo():Promise<wallet.Wallet|Error>;

export function SaveWallet(arg1:string,arg2:string):Promise<wallet.Wallet|Error>;

export function WailsInit(arg1:context.Context):Promise<Error>;
