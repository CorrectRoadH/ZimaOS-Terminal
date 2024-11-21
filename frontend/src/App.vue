<script setup lang="ts">
import zimaAuthAxios from './utils/axiosInstance';


const { t } = useI18n();

onMounted(async () => {
  window.document.title = t("title");
  // post /v2/terminal/terminal
  const resp = await zimaAuthAxios.post("/v2/terminal/terminal", {}, {
    headers: {
      "Content-Type": "application/json",
    },
  });

  const data = resp.data;
  const port = data.port;

  const newUrl = `${window.location.protocol}//${window.location.hostname}:${port}/`;
  window.location.href = newUrl;
});
</script>


<template>
  <div class="flex flex-col gap-16 text-center">
    <div class="loader"></div>
  </div>
</template>

<style>
  .loader {
    border: 8px solid #f3f3f3; /* Light grey */
    border-top: 8px solid #3498db; /* Blue */
    border-radius: 50%;
    width: 60px;
    height: 60px;
    animation: spin 1s linear infinite;
    margin: 0 auto;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
</style>
