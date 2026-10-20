<template>
	<svg
		:class="[sizeClass, 'shrink-0 transition-transform duration-200 hover:scale-105']"
		viewBox="0 0 120 120"
		fill="none"
		xmlns="http://www.w3.org/2000/svg"
		role="img"
		aria-label="ZF Hassaan Logo"
	>
		<defs>
			<!-- Ambient Outer Glow -->
			<radialGradient id="zfh-glow" cx="50%" cy="50%" r="50%">
				<stop offset="0%" stop-color="#6366F1" stop-opacity="0.35" />
				<stop offset="100%" stop-color="#4F46E5" stop-opacity="0" />
			</radialGradient>

			<!-- Primary Geometric Gradient -->
			<linearGradient id="zfh-primary-grad" x1="15" y1="10" x2="105" y2="110" gradientUnits="userSpaceOnUse">
				<stop offset="0%" stop-color="#06B6D4" />
				<stop offset="50%" stop-color="#4F46E5" />
				<stop offset="100%" stop-color="#7C3AED" />
			</linearGradient>

			<!-- Accent Core Gradient -->
			<linearGradient id="zfh-accent-grad" x1="30" y1="20" x2="90" y2="100" gradientUnits="userSpaceOnUse">
				<stop offset="0%" stop-color="#FFFFFF" />
				<stop offset="100%" stop-color="#C7D2FE" />
			</linearGradient>

			<!-- Drop Shadow Filter -->
			<filter id="zfh-shadow" x="0" y="0" width="120" height="120" filterUnits="userSpaceOnUse">
				<feDropShadow dx="0" dy="3" stdDeviation="4" flood-color="#312E81" flood-opacity="0.3" />
			</filter>
		</defs>

		<!-- Background Radial Atmosphere -->
		<circle cx="60" cy="60" r="54" fill="url(#zfh-glow)" />

		<!-- Outer Modern Hexagon Base -->
		<path
			d="M60 12L98 34V78L60 100L22 78V34L60 12Z"
			fill="#0F172A"
			stroke="url(#zfh-primary-grad)"
			stroke-width="3.5"
			stroke-linejoin="round"
			filter="url(#zfh-shadow)"
		/>

		<!-- Inner Hexagonal Grid Wireframe -->
		<path
			d="M60 22L89 39V73L60 90L31 73V39L60 22Z"
			stroke="#334155"
			stroke-width="1.2"
			stroke-dasharray="2 3"
			stroke-linejoin="round"
		/>

		<!-- Stylized 'Z' Top & Diagonal -->
		<path
			d="M38 42H82L44 76H72"
			stroke="url(#zfh-accent-grad)"
			stroke-width="6"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>

		<!-- Stylized 'F' Cross-Bar & Precision Diamond -->
		<path
			d="M49 58H71"
			stroke="#06B6D4"
			stroke-width="4.5"
			stroke-linecap="round"
		/>

		<!-- High-Tech Neon Focus Node (Center-Right) -->
		<circle cx="78" cy="76" r="3" fill="#38BDF8" />
	</svg>
</template>

<script setup lang="ts">
withDefaults(
	defineProps<{
		sizeClass?: string
	}>(),
	{
		sizeClass: 'w-7 h-7'
	}
)
</script>
