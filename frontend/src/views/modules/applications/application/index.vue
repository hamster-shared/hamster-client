<template>
  <PageWrapper>
    <Header :showBack="true" />
    <template v-if="!isLoading">
      <template v-if="isAppDeployed">
        <component :applicationId="applicationId" :is="ShowBlockchain" />
      </template>
      <Deployment :applicationId="applicationId" v-else />
    </template>
  </PageWrapper>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref, computed } from 'vue';
  import { useRoute } from 'vue-router';
  import { DictCodeEnum } from '/@/enums/dictCodeEnum';
  import { PageWrapper } from '/@/components/Page';
  import { useMessage } from '/@/hooks/web/useMessage';
  import Deployment from './components/Deployment.vue';
  import ApplicationTabs from './components/ApplicationTabs/components/ApplicationTabs.vue';
  import AptosTabs from './components/ApplicationTabs/aptosComponents/AptosTabs.vue';
  import AvalancheTabs from './components/ApplicationTabs/avalancheComponents/AvalancheTabs.vue';
  import BscTabs from './components/ApplicationTabs/bscComponents/BscTabs.vue';
  import EthereumTabs from './components/ApplicationTabs/ethereumComponents/EthereumTabs.vue';
  import NearTabs from './components/ApplicationTabs/nearComponents/NearTabs.vue';
  import OptimismTabs from './components/ApplicationTabs/optimismComponents/OptimismTabs.vue';
  import PolygonTabs from './components/ApplicationTabs/polygonComponents/PolygonTabs.vue';
  import StarkwareTabs from './components/ApplicationTabs/starkwareComponents/StarkwareTabs.vue';
  import SuiTabs from './components/ApplicationTabs/suiComponents/SuiTabs.vue';
  import ZksyncTabs from './components/ApplicationTabs/zksyncComponents/ZksyncTabs.vue';
  import Header from '../index/components/Header.vue';
  import { QueryApplicationById } from '/@wails/go/app/Application';
  import { useI18n } from '/@/hooks/web/useI18n';

  const { t } = useI18n();
  const { createErrorModal } = useMessage();
  const { params } = useRoute();
  const applicationId = Number(params.id);

  const appInfo = reactive({});
  const isLoading = ref(true);

  const ShowBlockchain = ref('ApplicationTabs');

  const isAppDeployed = computed(
    () =>
      DictCodeEnum.ApplicationDeployStatus_Running.is(appInfo.status) ||
      DictCodeEnum.ApplicationDeployStatus_Offline.is(appInfo.status),
  );

  const getAppInfo = async () => {
    isLoading.value = true;

    try {
      const result = await QueryApplicationById(applicationId);
      Object.assign(appInfo, result);
      console.log('appInfo', appInfo.serviceType);

      const mappings =
        {
          default: ApplicationTabs,
          bsc: BscTabs,
          avalanche: AvalancheTabs,
          ethereum: EthereumTabs,
          polygon: PolygonTabs,
          optimism: OptimismTabs,
          zksync: ZksyncTabs,
          starkware: StarkwareTabs,
          near: NearTabs,
          aptos: AptosTabs,
          sui: SuiTabs,
        }[appInfo.serviceType] || ApplicationTabs;

      ShowBlockchain.value = mappings;
    } catch (error: any) {
      createErrorModal({
        title: t('common.errorTip'),
        content: t('common.operateFailText'),
      });
    } finally {
      isLoading.value = false;
    }
  };

  onMounted(async () => {
    getAppInfo();
  });
</script>
