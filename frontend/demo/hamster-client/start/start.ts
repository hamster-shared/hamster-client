import { parseIndexingRule, setIndexingRule } from './rules';
import { processIdentifier, IndexingDecisionBasis } from '@graphprotocol/indexer-common'; //^0.19.0
import { createIndexerManagementClient } from './client';

//部署subgraph子图
const startRules = async () => {
  //该值为查询子图列表中currentVersion.subgraphDeployment.ipfsHash
  const deploymentId = 'QmXYZnggrnRpQe5ZGK2VzdhNyTgFqAP4AziQgxEZFk7fJm';
  try {
    const [identifier, identifierType] = await processIdentifier(deploymentId, {
      all: false,
      global: true,
    });

    const inputRule = parseIndexingRule({
      identifier,
      identifierType,
      decisionBasis: IndexingDecisionBasis.ALWAYS,
    });
    //该地址为index的地址(到时候需要前端拼接)
    const url = new URL('http://localhost:8500');
    const client = await createIndexerManagementClient({ url: url.toString() });
    await setIndexingRule(client, inputRule);
  } catch (error: any) {
    console.log(error.toString());
    process.exitCode = 1;
  }
};
