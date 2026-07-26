<script setup lang="ts">
import { provide, ref } from "vue";

import Tabs from "primevue/tabs";
import TabList from "primevue/tablist";
import Tab from "primevue/tab";
import TabPanels from "primevue/tabpanels";
import TabPanel from "primevue/tabpanel";

import MapsTab from "./MapsTab/MapsTab.vue";
import ProfileTab from "./ProfileTab.vue";
import { ACTIVE_TAB_KEY } from "src/injectionKeys";

const activeTab = ref("mapsTab");
provide(ACTIVE_TAB_KEY, activeTab);
</script>

<template>
  <Tabs v-model:value="activeTab" class="flex h-full w-full flex-col">
    <TabList class="flex w-full">
      <Tab value="mapsTab" class="flex-1 text-center">Главная</Tab>
      <Tab value="profileTab" class="flex-1 text-center">Профиль</Tab>
    </TabList>
    <TabPanels class="flex flex-1 flex-col overflow-hidden">
      <TabPanel value="mapsTab" class="m-0 flex-1 p-0">
        <MapsTab />
      </TabPanel>
      <TabPanel value="profileTab" class="m-0 flex-1 p-0">
        <ProfileTab />
      </TabPanel>
    </TabPanels>
  </Tabs>
</template>

<style scoped>
:deep(.p-tabs) {
  padding: 0 !important;
  height: 100%;
}
:deep(.p-tabpanels) {
  padding: 0 !important;
  flex: 1 !important;
  overflow: hidden !important;
  display: flex !important;
  flex-direction: column !important;
}
:deep(.p-tabpanel) {
  padding: 0 !important;
  flex: 1 !important;
  height: 100% !important;
}
:deep(.p-tablist) {
  padding: 0 !important;
}
</style>
