import type { GlobalThemeOverrides } from "naive-ui";

export const themeOverrides: GlobalThemeOverrides = {
	common: {
    	primaryColor: "#2588F4",
    },
	LoadingBar: {
		colorLoading: "#76B3F5",
	},
	Select: {
		peers: {
			InternalSelection: {
				borderHover: "1px solid #2588F4",
				borderFocus: "1px solid #76B3F5",
			}
		}
	},
	Button: {
		borderHover: "1px solid #2588F4",
		borderFocus: "1px solid #2588F4",
		textColorHover: "#2588F4",
		textColorGhostHover: "#2588F4",
		textColorFocus: "#2588F4",
		textColorGhostFocus: "#2588F4",
		colorHoverPrimary: "#76B3F5",
		colorFocusPrimary: "#76B3F5",
		borderFocusPrimary: "1px solid #76B3F5",
		borderHoverPrimary: "1px solid #2588F4",
		borderPressedPrimary: "1px solid #76B3F5",
		borderPressed: "1px solid #76B3F5",
		colorPressedPrimary: "#2588F4",
		textColorPressed: "#76B3F5",
        textColorGhostHoverPrimary: "#2588F4",
	},
	Input: {
		borderFocus: "1px solid #76B3F5",
		borderHover: "1px solid #2588F4",
	},
    Pagination: {
        itemTextColorHover: "#2588F4",
        itemTextColorActive: "#2588F4",
        itemTextColorPressed: "#76B3F5",
    }
};
