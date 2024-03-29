<template>
  <Form>
    <div class="stake-box p-[20px]">
      <div class="flex">
        <img :src="addressAvatar" class="bg-[#D8D8D8] rounded-[50%] w-[60px] h-[60px]" />
        <div class="ml-[16px]">
          <div>{{ t('applications.see.stake') }}</div>
          <div class="text-[20px] font-bold">{{ shortAddress }}</div>
        </div>
      </div>
      <div class="grid grid-cols-2 text-center mt-[20px]">
        <div class="right-line">
          <div>
            <SvgIcon color="#63A0FA" size="16" name="grt" />
            <label class="text-[#7B8082] ml-[6px]">{{ t('applications.see.yStake') }}</label>
          </div>
          <div>
            <label class="text-[18px] font-bold">{{ stakeAmount }}</label>
            <label class="text-[12px]">{{ t('applications.see.grt') }}</label>
          </div>
        </div>
        <div>
          <div>
            <SvgIcon color="#63A0FA" size="16" name="balance" />
            <label class="text-[#7B8082] ml-[6px]">{{ t('applications.see.wBalance') }}</label>
          </div>
          <div>
            <label class="text-[18px] font-bold">{{ addressBalance }}</label>
            <label class="text-[12px]">{{ t('applications.see.grt') }}</label>
          </div>
        </div>
      </div>
    </div>
    <div class="px-[20px] pb-[60px]">
      <div class="font-bold my-[10px]">{{ t('applications.reward.stakeAmount') }}</div>
      <Input
        class="border !border-[#043CC1] rounded-[8px] h-[60px] px-[10px]"
        v-model:value="inputStakeAmount"
        @change="inputStakeChange"
      >
        <template #suffix>
          <div>
            <label class="text-[#7B8082] mr-[10px]">{{ t('applications.see.grt') }}</label>
            <label
              class="bg-[#63A0FA] px-[20px] py-[8px] rounded-[4px] text-white"
              @click="maxClick"
            >
              {{ t('applications.see.max') }}
            </label>
          </div>
        </template>
      </Input>
      <div class="flex justify-between mt-[10px]">
        <div>{{ t('applications.see.minAmount') }}</div>
        <div>
          <label class="text-[18px] font-bold mr-[3px]">100000</label
          >{{ t('applications.see.grt') }}
        </div>
      </div>
      <div class="flex justify-between mt-[10px]">
        <div>{{ t('applications.see.thawPeriod') }}</div>
        <div>
          <label class="text-[18px] font-bold mr-[3px]">28</label>{{ t('applications.see.days') }}
        </div>
      </div>
      <div v-if="approveButtonVisible">
        <div class="flex my-[10px]">
          <div class="seq-div !bg-[#043CC1]">1</div>{{ t('applications.see.gtrStak') }}
        </div>
        <Button type="primary" size="large" @click="approve" :loading="approveLoading">{{
          t('applications.see.grtAccess')
        }}</Button>
      </div>
      <div class="flex my-[10px]">
        <div :class="stakeDisabled ? 'seq-div' : 'seq-div !bg-[#043CC1]'">2</div
        ><label class="text-[#7B8082]">{{ t('applications.see.gtrStak') }}</label>
      </div>
      <Button
        type="primary"
        size="large"
        :disabled="stakeDisabled"
        @click="stake"
        :loading="stakingLoading"
        >{{ t('applications.see.stake') }}</Button
      >
    </div>
  </Form>
</template>
<script lang="ts" setup>
  import { useI18n } from '/@/hooks/web/useI18n';
  import { SvgIcon } from '/@/components/Icon';
  import { Button, Input, Form } from 'ant-design-vue';
  import { computed, ref } from 'vue';
  import { buildContract, createWeb3Api, runContractMethod, web3Abi } from '/@/utils/web3Util';
  import { useMessage } from '/@/hooks/web/useMessage';
  // defines
  const props = defineProps({
    stakeAmount: String,
    addressBalance: String,
    addressAvatar: String,
    shortAddress: String,
    deployInfo: Object as PropType<Recordable>,
  });
  const { t } = useI18n();
  const inputStakeAmount = ref();
  const approveLoading = ref(false);
  const stakeDisabled = ref(true);
  const stakingLoading = ref(false);
  const approveButtonVisible = ref(true);
  const { createErrorModal } = useMessage();
  const emits = defineEmits(['close-drawer', 'query-stake', 'get-balance']);
  const maxClick = () => {
    inputStakeAmount.value = props.addressBalance;
  };
  const web3Api = computed(() => {
    const { initialization, staking } = props.deployInfo;
    const accountMnemonic = initialization.accountMnemonic;
    const networkUrl = staking.networkUrl;
    if (accountMnemonic && networkUrl) {
      return createWeb3Api(networkUrl, accountMnemonic);
    }
    return undefined;
  });
  const inputStakeChange = () => {
    inputStakeAmount.value = inputStakeAmount.value.replace(/[^\d.]/g, '');
    judgeApprove();
  };
  const judgeApprove = async () => {
    approveLoading.value = true;
    const api = web3Api.value;
    if (api && api.__config) {
      const { erc20ContractAddress } = api.__config;
      const address = props.deployInfo?.staking.agentAddress;
      const owner = props.deployInfo?.deployment.indexerAddress;
      const contract = buildContract(api, web3Abi.ecr20Abi, erc20ContractAddress);
      try {
        let data = await runContractMethod({
          api,
          contract,
          method: 'allowance',
          methodArgs: [owner, address],
          type: 'call',
        });
        data = api.utils.fromWei(data);
        if (inputStakeAmount.value * 1 <= data * 1) {
          approveButtonVisible.value = false;
          stakeDisabled.value = false;
        } else {
          approveButtonVisible.value = true;
          stakeDisabled.value = true;
        }
        approveLoading.value = false;
      } catch (e: any) {
        approveButtonVisible.value = true;
        approveLoading.value = false;
      }
    }
  };
  const approve = async () => {
    const api = web3Api.value;
    if (api && api.__config) {
      approveLoading.value = true;
      const { erc20ContractAddress } = api.__config;
      const address = props.deployInfo?.staking.agentAddress;
      const contract = buildContract(api, web3Abi.ecr20Abi, erc20ContractAddress);
      inputStakeAmount.value = inputStakeAmount.value || 1;
      try {
        await runContractMethod({
          api,
          contract,
          method: 'approve',
          methodArgs: [address, api.utils.toWei(inputStakeAmount.value.toString())],
          type: 'send',
        });
        approveLoading.value = false;
        stakeDisabled.value = false;
        approveButtonVisible.value = false;
      } catch (e: any) {
        approveButtonVisible.value = true;
        approveLoading.value = false;
        createErrorModal({
          title: t('common.errorTip'),
          content: e.message,
        });
        console.info('Approve Staking Proxy Contract Error', e);
      }
    }
  };
  const stake = async () => {
    const api = web3Api.value;
    if (api && api.__config) {
      stakingLoading.value = true;
      const address = props.deployInfo?.staking.agentAddress;
      const contract = buildContract(api, web3Abi.stakeDistributionProxyAbi, address);
      try {
        await runContractMethod({
          api,
          contract,
          method: 'rePledge',
          methodArgs: [api.utils.toWei(inputStakeAmount.value.toString())],
          type: 'send',
        });
        inputStakeAmount.value = '';
        emits('close-drawer');
        emits('query-stake');
        emits('get-balance');
        stakingLoading.value = false;
      } catch (e: any) {
        stakingLoading.value = false;
        createErrorModal({
          title: t('common.errorTip'),
          content: e.message,
        });
        console.info('Approve Staking Proxy Contract Error', e);
      }
    }
  };
  // const closeStakeDrawer = () => {};
</script>
<style lang="less" scoped>
  .stake-box {
    box-shadow: 0px 0px 5px 0px rgba(150, 150, 150, 0.2);
  }

  .right-line::after {
    border-color: rgba(216, 216, 216, 0.4);
    top: 30%;
  }

  .seq-div {
    @apply text-white text-center;
    background: rgba(4, 60, 193, 0.5);
    border-radius: 50%;
    height: 20px;
    width: 20px;
    font-size: 12px;
    margin-right: 5px;
  }

  :deep(.ant-btn-primary) {
    @apply w-full;
    border-radius: 24px;
  }

  :deep(.ant-btn[disabled]) {
    @apply text-white border-none;
    background: rgba(0, 0, 0, 0.1);
  }
</style>
