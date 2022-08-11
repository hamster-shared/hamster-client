<template>
  <div class="revenue-bg rounded-[20px] p-[20px] relative">
    <div
      @click="receiveBenefits"
      class="benefit-bg cursor-pointer absolute right-0 top-10 text-white py-[16px] pl-[16px] rounded-l-[100px]"
    >
      <label class="text-[18px] font-bold">{{ t('applications.see.receiveBenefits') }}</label>
      <SvgIcon class="text-white rounded-[50%] mx-[5px]" size="16" name="right" />
    </div>
    <div class="flex">
      <img
        :src="createSvgAvatar(addressAvatar)"
        class="bg-[#D8D8D8] rounded-[50%] h-[120px] w-[120px]"
      />
      <div class="text-[#F4F4F4] ml-[20px]">
        <div class="text-[26px]">{{ shortenAddress(address) }}</div>
        <div>{{ address }}</div>
        <div class="relative text-white"
          ><label class="text-[26px] font-bold">{{ addressBalance }}</label
          ><label class="absolute top-0">GRT</label></div
        >
      </div>
    </div>
    <div class="grid grid-cols-3 text-white text-center mt-[20px]">
      <div class="right-line">
        <div class="text-[22px] font-bold">{{ income }}</div>
        <div>{{ t('applications.see.income') }}</div>
      </div>
      <div class="right-line">
        <div class="text-[22px] font-bold">{{ stakeTotal }}</div>
        <div>{{ t('applications.see.stakAmount') }}</div>
      </div>
      <div>
        <div class="text-[22px] font-bold">{{ unStakeAmount }}</div>
        <div>{{ t('applications.see.unstakAmount') }}</div>
      </div>
    </div>
    <div class="my-[40px] text-center">
      <label
        class="label-btn"
        @click="
          drawerVisible = true;
          stakeVisible = true;
        "
        >{{ t('applications.see.stake') }}</label
      >
      <label
        class="label-btn ml-[30px]"
        @click="
          drawerVisible = true;
          unstakeVisible = true;
        "
        >{{ t('applications.see.unstake') }}</label
      >
      <label
        class="label-btn ml-[30px]"
        @click="
          drawerVisible = true;
          withdrawVisible = true;
        "
        >{{ t('applications.see.withdraw') }}</label
      >
    </div>
  </div>
  <Drawer
    v-model:visible="drawerVisible"
    :closable="false"
    placement="right"
    class="drawer-revenue-info"
    @close="onDrawerClose"
  >
    <StakeDrawer
      :stakeAmount="stakeTotal"
      :addressBalance="addressBalance"
      :addressAvatar="shortenAddress(address)"
      :deployInfo="deployData"
      v-if="stakeVisible"
    />
    <UnstakeDrawer v-if="unstakeVisible" />
    <WithdrawDrawer v-if="withdrawVisible" />
  </Drawer>
</template>
<script lang="ts" setup>
  import { onMounted, ref, watchEffect, computed } from 'vue';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { SvgIcon } from '/@/components/Icon';
  import StakeDrawer from './StakeDrawer.vue';
  import UnstakeDrawer from './UnstakeDrawer.vue';
  import WithdrawDrawer from './WithdrawDrawer.vue';
  import { createSvgAvatar } from '/@/utils/avatar';
  import { shortenAddress } from '/@/utils/thegraphUtil';
  import { buildContract, createWeb3Api, runContractMethod, web3Abi } from '/@/utils/web3Util';
  import { Modal, Drawer } from 'ant-design-vue';
  import { useMessage } from '/@/hooks/web/useMessage';

  // defines
  const props = defineProps({
    deployInfo: Object as PropType<Recordable>,
  });
  const { t } = useI18n();
  // const emits = defineEmits(['modalConfirm']);

  const drawerVisible = ref(false);
  const stakeVisible = ref(false);
  const unstakeVisible = ref(false);
  const withdrawVisible = ref(false);

  // This is an placeholder, address will be get from API later
  const address = ref('');
  // address balance
  const addressBalance = ref('0');
  const addressAvatar = ref();
  const income = ref('0');
  const stakeTotal = ref('0');
  const unStakeAmount = ref('0');
  const { createErrorModal } = useMessage();
  const deployData = ref<{
    initialization: Recordable;
    staking: Recordable;
    deployment: Recordable;
  }>({
    initialization: {},
    staking: {},
    deployment: {},
  });
  async function onDrawerClose() {
    stakeVisible.value = false;
    unstakeVisible.value = false;
    withdrawVisible.value = false;
  }
  const initAddress = () => {
    address.value = props.deployInfo?.deployment.indexerAddress;
    addressAvatar.value = createSvgAvatar(address.value);
  };
  //web3 api
  const web3Api = computed(() => {
    const { initialization, staking } = props.deployInfo;
    const accountMnemonic = initialization.accountMnemonic;
    const networkUrl = staking.networkUrl;
    if (accountMnemonic && networkUrl) {
      return createWeb3Api(networkUrl, accountMnemonic);
    }
    return undefined;
  });
  const getAddressBalance = async () => {
    const api = web3Api.value;
    if (api && api.__config) {
      const erc20 = buildContract(api, web3Abi.ecr20Abi, api.__config.erc20ContractAddress);
      const balance_data = await runContractMethod({
        api,
        contract: erc20,
        method: 'balanceOf',
        methodArgs: ['0x9438BbE4E7AF1ec6b13f75ECd1f53391506A12DF'],
        type: 'call',
      });
      addressBalance.value = api.utils.fromWei(balance_data);
    }
  };
  const getIncome = async () => {
    const api = web3Api.value;
    const address = props.deployInfo?.staking.agentAddress;
    if (api) {
      const contract = buildContract(api, web3Abi.stakeDistributionProxyAbi, address);
      try {
        const data = await runContractMethod({
          api,
          contract,
          method: 'gainIncome',
          methodArgs: [],
          type: 'call',
        });
        income.value = api.utils.fromWei(data.toString());
      } catch (e: any) {
        income.value = '0';
        console.info(e);
      }
    }
  };
  const getStakeAmount = async () => {
    const api = web3Api.value;
    if (api) {
      const address = props.deployInfo?.staking.agentAddress;
      const contract = buildContract(api, web3Abi.stakeDistributionProxyAbi, address);
      try {
        const data = await runContractMethod({
          api,
          contract,
          method: 'getStakingAmount',
          methodArgs: [],
          type: 'call',
        });
        stakeTotal.value = api.utils.fromWei(data.toString());
      } catch (e: any) {
        stakeTotal.value = '0';
        console.info(e);
      }
    }
  };
  const getUnStakeAmount = async () => {
    const api = web3Api.value;
    if (api) {
      const address = props.deployInfo?.staking.agentAddress;
      const contract = buildContract(api, web3Abi.stakeDistributionProxyAbi, address);
      try {
        const data = await runContractMethod({
          api,
          contract,
          method: 'getUnStakingAmount',
          methodArgs: [],
          type: 'call',
        });
        unStakeAmount.value = api.utils.fromWei(data.toString());
      } catch (e: any) {
        unStakeAmount.value = '0';
        console.info(e);
      }
    }
  };
  // const queryDeployInfo = async () => {
  //   const data = await GetDeployInfo(applicationId);
  //   if (data) {
  //     deployInfo.value = data;
  //     address.value = data.deployment.indexerAddress;
  //     addressAvatar.value = createSvgAvatar(address.value);
  //   }
  //   console.log(data);
  //   // if (data) deployInfo.value = data;
  // };

  onMounted(() => {
    // queryDeployInfo();
  });
  watchEffect(() => {
    if (props.deployInfo?.staking.agentAddress) {
      deployData.value = props.deployInfo;
      initAddress();
      getAddressBalance();
      getIncome();
      getStakeAmount();
      getUnStakeAmount();
    }
  });
  const withdrawIncome = async () => {
    const api = web3Api.value;
    const address = props.deployInfo?.staking.agentAddress;
    if (api) {
      const contract = buildContract(api, web3Abi.stakeDistributionProxyAbi, address);
      try {
        await runContractMethod({
          api,
          contract,
          method: 'withdrawIncome',
          methodArgs: [],
          type: 'send',
        });
        await getIncome();
      } catch (e: any) {
        createErrorModal({
          title: t('common.errorTip'),
          content: e.message,
        });
      }
    }
  };
  const receiveBenefits = () => {
    Modal.confirm({
      title: t('applications.see.receiveBenefitsInfo'),
      icon: '',
      okText: t('common.okText'),
      cancelText: t('common.cancelText'),
      onOk() {
        withdrawIncome();
      },
    });
  };
</script>
<style lang="less" scoped>
  .drawer-revenue-info {
    :global(.ant-drawer) {
      display: flex;
      align-items: center;
    }

    :global(.ant-drawer-right .ant-drawer-content-wrapper) {
      height: auto;
    }

    :global(.ant-drawer-content) {
      border-radius: 0 8px 8px 0;
      min-height: 500px;
    }

    :global(.ant-drawer-body) {
      padding: 0px;
    }
  }

  .revenue-bg {
    background: url('src/assets/images/application-bg.png') no-repeat top;
    background-size: 100% 100%;
  }

  .label-btn {
    @apply bg-[#376AED] px-[30px] py-[10px] rounded-[20px] text-[#E2EAFC] cursor-pointer;
  }

  .benefit-bg {
    background: linear-gradient(122deg, #ff893e 0%, #e76a93 100%);
  }
</style>
