<template>
	<div class="h-screen w-screen overflow-hidden">
		<router-view />
		<ToastContainer />
	</div>
</template>

<script setup lang="ts">
import ToastContainer from '@/components/ToastContainer.vue'
</script>

