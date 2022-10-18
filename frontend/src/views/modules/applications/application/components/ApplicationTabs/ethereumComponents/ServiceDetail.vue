<template>
  <div class="left-line div-box div-card">
    <div class="div-title">{{ t('applications.see.serviceInfo') }}</div>
    <Row :gutter="[10, 10]" class="text-center">
      <Col :span="5" class="font-bold">{{ t('applications.index.nameText') }}</Col>
      <Col :span="5" class="font-bold">{{ t('applications.see.createTime') }}</Col>
      <Col :span="4" class="font-bold">{{ t('applications.index.statusText') }}</Col>
      <Col :span="5" class="font-bold">{{ t('applications.new.leaseTerm') }}</Col>
      <Col :span="5" class="font-bold">{{ t('applications.see.nodeType') }}</Col>
      <Col :span="5">{{ props.application.name }}</Col>
      <Col :span="5">{{ formatToDateTime(application.createdAt) }}</Col>
      <Col :span="4">
        {{ DictCodeEnum.ApplicationDeployStatus.getOptionLabel(application.status) }}
      </Col>
      <Col :span="5">{{ application.leaseTerm }}</Col>
      <Col :span="5">{{ application.selectNodeType }}</Col>
    </Row>
  </div>
</template>

<script lang="ts" setup>
  import { useI18n } from '/@/hooks/web/useI18n';
  import { Row, Col } from 'ant-design-vue';
  import { formatToDateTime } from '/@/utils/dateUtil';
  import { DictCodeEnum } from '/@/enums/dictCodeEnum';

  const { t } = useI18n();

  const props = defineProps({
    application: Object as PropType<Recordable>,
    deployInfo: Object as PropType<Recordable>,
  });
</script>

<style lang="less" scoped>
  .div-box {
    box-shadow: 0px 6px 20px 0px rgb(31 31 35 / 10%);
  }

  .div-card {
    @apply bg-white;
    border-radius: 20px;
    margin-top: 30px;
    padding: 20px 20px 30px 20px;
  }

  .div-title {
    @apply font-bold;
    font-size: 18px;
    color: #072772;
    margin-bottom: 20px;
  }

  .line-color::before {
    border-color: #e24066;
  }

  .line-color2::before {
    border-color: #009eff;
  }
</style>
