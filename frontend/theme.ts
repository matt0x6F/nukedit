
import type { CustomThemeConfig } from '@skeletonlabs/tw-plugin';

export const nukeditTheme: CustomThemeConfig = {
    name: 'nukedit',
    properties: {
		// =~= Theme Properties =~=
		"--theme-font-family-base": `Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'`,
		"--theme-font-family-heading": `Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'`,
		"--theme-font-color-base": "0 0 0",
		"--theme-font-color-dark": "255 255 255",
		"--theme-rounded-base": "9999px",
		"--theme-rounded-container": "8px",
		"--theme-border-base": "1px",
		// =~= Theme On-X Colors =~=
		"--on-primary": "0 0 0",
		"--on-secondary": "255 255 255",
		"--on-tertiary": "0 0 0",
		"--on-success": "0 0 0",
		"--on-warning": "0 0 0",
		"--on-error": "255 255 255",
		"--on-surface": "255 255 255",
		// =~= Theme Colors  =~=
		// primary | #ff7800 
		"--color-primary-50": "255 235 217", // #ffebd9
		"--color-primary-100": "255 228 204", // #ffe4cc
		"--color-primary-200": "255 221 191", // #ffddbf
		"--color-primary-300": "255 201 153", // #ffc999
		"--color-primary-400": "255 161 77", // #ffa14d
		"--color-primary-500": "255 120 0", // #ff7800
		"--color-primary-600": "230 108 0", // #e66c00
		"--color-primary-700": "191 90 0", // #bf5a00
		"--color-primary-800": "153 72 0", // #994800
		"--color-primary-900": "125 59 0", // #7d3b00
		// secondary | #3d3846 
		"--color-secondary-50": "226 225 227", // #e2e1e3
		"--color-secondary-100": "216 215 218", // #d8d7da
		"--color-secondary-200": "207 205 209", // #cfcdd1
		"--color-secondary-300": "177 175 181", // #b1afb5
		"--color-secondary-400": "119 116 126", // #77747e
		"--color-secondary-500": "61 56 70", // #3d3846
		"--color-secondary-600": "55 50 63", // #37323f
		"--color-secondary-700": "46 42 53", // #2e2a35
		"--color-secondary-800": "37 34 42", // #25222a
		"--color-secondary-900": "30 27 34", // #1e1b22
		// tertiary | #0EA5E9 
		"--color-tertiary-50": "219 242 252", // #dbf2fc
		"--color-tertiary-100": "207 237 251", // #cfedfb
		"--color-tertiary-200": "195 233 250", // #c3e9fa
		"--color-tertiary-300": "159 219 246", // #9fdbf6
		"--color-tertiary-400": "86 192 240", // #56c0f0
		"--color-tertiary-500": "14 165 233", // #0EA5E9
		"--color-tertiary-600": "13 149 210", // #0d95d2
		"--color-tertiary-700": "11 124 175", // #0b7caf
		"--color-tertiary-800": "8 99 140", // #08638c
		"--color-tertiary-900": "7 81 114", // #075172
		// success | #84cc16 
		"--color-success-50": "237 247 220", // #edf7dc
		"--color-success-100": "230 245 208", // #e6f5d0
		"--color-success-200": "224 242 197", // #e0f2c5
		"--color-success-300": "206 235 162", // #ceeba2
		"--color-success-400": "169 219 92", // #a9db5c
		"--color-success-500": "132 204 22", // #84cc16
		"--color-success-600": "119 184 20", // #77b814
		"--color-success-700": "99 153 17", // #639911
		"--color-success-800": "79 122 13", // #4f7a0d
		"--color-success-900": "65 100 11", // #41640b
		// warning | #EAB308 
		"--color-warning-50": "252 244 218", // #fcf4da
		"--color-warning-100": "251 240 206", // #fbf0ce
		"--color-warning-200": "250 236 193", // #faecc1
		"--color-warning-300": "247 225 156", // #f7e19c
		"--color-warning-400": "240 202 82", // #f0ca52
		"--color-warning-500": "234 179 8", // #EAB308
		"--color-warning-600": "211 161 7", // #d3a107
		"--color-warning-700": "176 134 6", // #b08606
		"--color-warning-800": "140 107 5", // #8c6b05
		"--color-warning-900": "115 88 4", // #735804
		// error | #c01c28 
		"--color-error-50": "246 221 223", // #f6dddf
		"--color-error-100": "242 210 212", // #f2d2d4
		"--color-error-200": "239 198 201", // #efc6c9
		"--color-error-300": "230 164 169", // #e6a4a9
		"--color-error-400": "211 96 105", // #d36069
		"--color-error-500": "192 28 40", // #c01c28
		"--color-error-600": "173 25 36", // #ad1924
		"--color-error-700": "144 21 30", // #90151e
		"--color-error-800": "115 17 24", // #731118
		"--color-error-900": "94 14 20", // #5e0e14
		// surface | #241f31 
		"--color-surface-50": "222 221 224", // #dedde0
		"--color-surface-100": "211 210 214", // #d3d2d6
		"--color-surface-200": "200 199 204", // #c8c7cc
		"--color-surface-300": "167 165 173", // #a7a5ad
		"--color-surface-400": "102 98 111", // #66626f
		"--color-surface-500": "36 31 49", // #241f31
		"--color-surface-600": "32 28 44", // #201c2c
		"--color-surface-700": "27 23 37", // #1b1725
		"--color-surface-800": "22 19 29", // #16131d
		"--color-surface-900": "18 15 24", // #120f18
	}
}