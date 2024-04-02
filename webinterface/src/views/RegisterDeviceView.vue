<script lang="ts" setup>
import {
	NCard,
	NFormItem,
	NInput,
	NA,
	NSpace,
	NButton,
	NAlert,
} from "naive-ui";
import { computed, ref } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";

import QuitAppButton from "@/components/QuitAppButton.vue";
import type { UserInfoResponse } from "@/lib/responses";

const router = useRouter();

const deviceName = ref("");
const displayError = ref<string|null>(null);

const canSubmit = computed(() => {
	return deviceName.value.length > 0;
});

async function submit() {
	displayError.value = null;

	try {
		const resp = await axios.post("/api/register", {
			devicename: deviceName.value,
		});

		const data = resp.data as UserInfoResponse;
		
		console.log(data);

		if (!data.loggedin) { // shouldn't ever happen
			router.push({ name: "login" });
		} else if (!data.registered) {
			router.push({ name: "register-device" });
		} else {
			router.push({ name: "info" });
		}
	} catch(err: any) {
		displayError.value = err.response.data;
	}
}
</script>

<template>
	<n-card title="Register Device" :segmented="true">
		<template #header-extra>
			<quit-app-button />
		</template>

		<div style="margin-bottom: 20px">
			Enter a name for this device. The name can not include any spaces.
		</div>

		<n-space vertical>
			<n-alert v-if="displayError" type="error">
				{{ displayError }}
			</n-alert>

			<n-form-item
				label="Device Name"
				:show-feedback="false"
				:show-require-mark="true"
				@keypress.enter="submit()"
			>
				<n-input v-model:value="deviceName" placeholder="Device Name" />
			</n-form-item>
		</n-space>

		<template #footer>
			<n-space justify="end">
				<n-button type="primary" @click="submit()" :disabled="!canSubmit">Register</n-button>
			</n-space>
		</template>
	</n-card>
</template>