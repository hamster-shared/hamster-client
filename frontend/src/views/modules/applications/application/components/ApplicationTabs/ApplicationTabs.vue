<template>
  <div class="mx-[40px] mb-[100px]">
    <Tabs v-model:activeKey="activeKey">
      <TabPane key="1" :tab="t('applications.see.revenueInfo')">
        <RevenueInfo @modal-confirm="modalConfirm" :deployInfo="deployInfo" />
      </TabPane>
      <TabPane key="2" :tab="t('applications.see.subgraph')">
        <Subgraph />
      </TabPane>
      <TabPane key="3" :tab="t('applications.see.serviceDetails')">
        <ServiceDetails @modal-confirm="modalConfirm" />
      </TabPane>
    </Tabs>
  </div>
</template>
<script lang="ts" setup>
  import { onMounted, ref } from 'vue';
  import { useI18n } from '/@/hooks/web/useI18n';
  import RevenueInfo from './components/RevenueInfo.vue';
  import Subgraph from './components/Subgraph.vue';
  import ServiceDetails from './components/ServiceDetails.vue';
  import { Tabs, TabPane, Modal } from 'ant-design-vue';
  import { GetDeployInfo } from '/@wails/go/app/Deploy';

  const props = defineProps({
    applicationId: Number,
  });

  const { t } = useI18n();

  const activeKey = ref('1');
  const deployInfo = ref<{
    initialization: Recordable;
    staking: Recordable;
    deployment: Recordable;
  }>({
    initialization: {},
    staking: {},
    deployment: {},
  });
  // Get saved deployInfo from API
  const getDeployInfo = async () => {
    const data = await GetDeployInfo(props.applicationId);
    if (data) deployInfo.value = data;
    console.log(deployInfo.value);
    console.log(9999);
  };
  onMounted(() => {
    getDeployInfo();
  });
  const modalConfirm = () => {
    Modal.confirm({
      title: t('applications.see.receiveBenefitsInfo'),
      icon: '',
      okText: t('common.okText'),
      cancelText: t('common.cancelText'),
      onOk() {
        console.log('OK');
      },
      onCancel() {
        console.log('Cancel');
      },
    });
  };
</script>
