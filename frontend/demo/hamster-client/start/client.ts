import { createClient } from '@urql/core'; //2.3.6
// @ts-ignore
import fetch from 'isomorphic-fetch'; //3.0.0
import { IndexerManagementClient } from '@graphprotocol/indexer-common'; //^0.19.0

//创建client
export const createIndexerManagementClient = async ({
  url,
}: {
  url: string;
}): Promise<IndexerManagementClient> => {
  return createClient({ url, fetch }) as unknown as IndexerManagementClient;
};
