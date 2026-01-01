/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        midnight: {
          50: '#f0f4ff',
          100: '#e0e8ff',
          200: '#c0d0ff',
          300: '#8eafff',
          400: '#5680ff',
          500: '#2d4fff',
          600: '#1a2df5',
          700: '#131de0',
          800: '#141ab5',
          900: '#161a8f',
          950: '#0c0d1f',
        },
        accent: {
          emerald: '#10B981',
          coral: '#F43F5E',
          amber: '#F59E0B',
          violet: '#8B5CF6',
        },
      },
      fontFamily: {
        sans: ['Outfit', 'system-ui', 'sans-serif'],
        display: ['Clash Display', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      animation: {
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'shimmer': 'shimmer 2s linear infinite',
        'glow': 'glow 2s ease-in-out infinite alternate',
      },
      keyframes: {
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
        glow: {
          '0%': { boxShadow: '0 0 20px rgba(45, 79, 255, 0.3)' },
          '100%': { boxShadow: '0 0 40px rgba(45, 79, 255, 0.6)' },
        },
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'mesh-gradient': 'linear-gradient(135deg, #0c0d1f 0%, #161a8f 50%, #131de0 100%)',
      },
    },
  },
  plugins: [],
}

