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

const email = ref("");
const password = ref("");
const twoFactorCode = ref("");
const twoFactorRequired = ref(false);
const displayError = ref<string|null>(null);

const canLogin = computed(() => {
	return email.value.length > 0 && password.value.length > 0 && (!twoFactorRequired.value || twoFactorCode.value.length > 0);
});

async function login() {
	displayError.value = null;

	try {
		const resp = await axios.post("/api/login", {
			email: email.value,
			password: password.value,
			twofacode: (twoFactorRequired.value ? twoFactorCode.value : undefined),
		});

		if (resp.data == "2FA_MISSING") {
			twoFactorRequired.value = true;
		} else {
			const data = resp.data as UserInfoResponse;
			
			if (!data.loggedin) { // shouldn't ever happen
				router.push({ name: "login" });
			} else if (!data.registered) {
				router.push({ name: "register-device" });
			} else {
				router.push({ name: "info" });
			}
		}
	} catch(err: any) {
		displayError.value = err.response.data;
	}
}
</script>

<template>
	<n-card title="Login" :segmented="true">
		<template #header-extra>
			<quit-app-button />
		</template>

		<div style="margin-bottom: 20px">
			Login with your Pushover account.
			Create one <n-a href="https://pushover.net/signup">here</n-a> if you haven't already.
		</div>

		<n-space vertical>
			<n-alert v-if="displayError" type="error">
				{{ displayError }}
			</n-alert>

			<n-form-item label="Email" :show-feedback="false" @keypress.enter="login()">
				<n-input v-model:value="email" placeholder="Email" />
			</n-form-item>

			<n-form-item label="Password" :show-feedback="false" @keypress.enter="login()">
				<n-input v-model:value="password" type="password" show-password-on="click" placeholder="Password"/>
			</n-form-item>

			<n-form-item v-if="twoFactorRequired" label="Two-Factor Code" :show-feedback="false">
				<n-input v-model:value="twoFactorCode" placeholder="Two-Factor Code"/>
			</n-form-item>
		</n-space>

		<template #footer>
			<n-space justify="end">
				<n-button type="primary" @click="login()" :disabled="!canLogin">Login</n-button>
			</n-space>
		</template>
	</n-card>
</template>