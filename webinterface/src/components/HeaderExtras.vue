<script lang="ts" setup>
import {
	NSpace,
	NButton,
	NIcon,
} from "naive-ui";
import axios from "axios";
import {
	ArrowClockwise16Regular
} from "@vicons/fluent";
import { ref } from "vue";

import QuitAppButton from "@/components/QuitAppButton.vue";

const emit = defineEmits(["reset-connection"]);

const resetLoading = ref(false);

async function resetConnection() {
	resetLoading.value = true;

	try {
		await axios.post("/api/reset-connection");

		emit("reset-connection");
	} finally {
		resetLoading.value = false;
	}
}
</script>

<template>
	<n-space>
		<n-button @click="resetConnection" :loading="resetLoading">
			<template #icon>
				<n-icon><ArrowClockwise16Regular /></n-icon>
			</template>
			Retry Pushover connection
		</n-button>

		<quit-app-button />
	</n-space>
</template>