<script lang="ts" setup>
import {
	NConfigProvider,
	NDialogProvider,
	darkTheme,
} from "naive-ui";

import { themeOverrides } from "@/lib/theme";

</script>

<template>
	<n-config-provider :theme="darkTheme" :theme-overrides="themeOverrides">
		<n-dialog-provider>
			<main>
				<router-view />
			</main>
		</n-dialog-provider>
	</n-config-provider>
</template>

<style lang="scss" scoped>
main {
	max-width: 600px;
	margin-left: auto;
	margin-right: auto;
	margin-top: 20px;
}
</style>